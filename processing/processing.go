package processing

import (
	//"bufio"
	"encoding/gob"
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/analysis"
	"github.com/Benjft/DiffusionLimitedAggregation/util"
	//"github.com/Benjft/DiffusionLimitedAggregation/util/drawing/svg"
	//"github.com/Benjft/DiffusionLimitedAggregation/util/drawing/xyz"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"

	"github.com/skratchdot/open-golang/open"
)

const (
	SAVE_PATH = "out\\save"
	PLOT_PATH = "out\\plot"
	DRAW_PATH = "out\\draw"
	DATA_PATH = "out\\rawData"
)

func init() {
	os.MkdirAll(SAVE_PATH, os.ModeDir)
	os.MkdirAll(PLOT_PATH, os.ModeDir)
	os.MkdirAll(DRAW_PATH, os.ModeDir)
	os.MkdirAll(DATA_PATH, os.ModeDir)

	gob.Register(RunData{})
}

type RunData struct {
	N, D, Seed int64
	Stick float64
	Points []aggregation.Point
}

func (run *RunData) FileName() string {
	return fmt.Sprintf("N%d-D%d-Se%d-St%b", run.N, run.D, run.Seed, run.Stick)
}

func (run *RunData) Run() {
	run.Points = aggregation.RunNew(run.N, run.Seed, run.D, run.Stick)
}

func (run *RunData) Save() (err error) {
	var (
		path = fmt.Sprintf("%s\\%s.save", SAVE_PATH, run.FileName())

		file *os.File
		encoder *gob.Encoder
	)

	file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder = gob.NewEncoder(file)
	err = encoder.Encode(run)
	if err != nil {
		return err
	}
	return nil
}

func (run *RunData) Load() (err error) {
	var (
		path = fmt.Sprintf("%s\\%s.save", SAVE_PATH, run.FileName())

		file *os.File
		decoder *gob.Decoder
	)

	file, err = os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder = gob.NewDecoder(file)
	err = decoder.Decode(run)
	if err != nil {
		return err
	}
	return nil
}

type Radii struct {
	D int64
	Stick float64
	Radii []float64
}
func (run *RunData) GyrationRadii() *Radii {
	return &Radii{run.D, run.Stick, analysis.GyrationRadii(run.Points)}
}


func SaveAll(runs []*RunData) {
	for _, run := range runs {
		go func(run *RunData) {
			var err = run.Save()
			if err != nil {
				fmt.Println(err)
			}
		} (run)
	}
}

func RunOne(n, d, seed int64, stick float64, tryLoad bool) (run *RunData) {
	run = &RunData{N:n, D:d, Seed:seed, Stick:stick}

	if tryLoad {
		var err = run.Load()
		if err != nil {
			fmt.Println(err, "Load Failed. Running...")
			run.Run()
		}
	} else {
		run.Run()
	}
	return
}

func RunMany(n, d, seed int64, stick float64, tryLoad bool, runs int64) (runData []*RunData) {
	runData = make([]*RunData, runs)
	var rng = rand.New(rand.NewSource(seed))

	var runChannel = make(chan *RunData)

	for range runData {
		go func (seed int64) {
			runChannel <- RunOne(n, d, seed, stick, tryLoad)
		} (rng.Int63())
	}

	for i := range runData {
		runData[i] = <-runChannel
	}

	return
}

func VaryDimension(n, d0, d1, dStep, seed int64, stick float64, tryLoad bool, runsPer int64) (runData []*RunData) {
	var (
		nPoints int64 = 0
		runChannel = make(chan []*RunData)
	)

	for d := d0; d <= d1; d += dStep {
		nPoints++
		go func (d int64) {
			runChannel <- RunMany(n, d, seed, stick, tryLoad, runsPer)
		} (d)
	}

	runData = make([]*RunData, runsPer * nPoints)
	for i := int64(0); i < nPoints; i++ {
		idx := i*runsPer
		for j, data := range <-runChannel {
			runData[idx+int64(j)] = data
		}
	}

	// TODO Remove this once tested
	for _, data := range runData {
		if data == nil {
			panic("Missing Data!")
		}
	}

	return runData
}

func VarySticking(n, d, seed int64, s0, s1, sStep float64, tryLoad bool, runsPer int64) (runData []*RunData) {
	var (
		nSteps = int64((s1-s0)/sStep)
		runChannel = make(chan []*RunData)
	)
	runData = make([]*RunData, runsPer * nSteps)

	for s := s0; s < s1; s += sStep {
		go func (stick float64) {
			runChannel <- RunMany(n, d, seed, stick, tryLoad, runsPer)
		} (s)
	}

	for i := int64(0); i < nSteps; i++ {
		idx := i*runsPer
		for j, data := range <-runChannel {
			runData[idx+int64(j)] = data
		}
	}

	// TODO Remove this once tested
	for _, data := range runData {
		if data == nil {
			panic("Missing Data!")
		}
	}

	return runData
}

// MeanAndConf95 returns the mean
// and the magnitude of the 95% confidence
// interval on the mean as low and high
// error values.
//
// MeanAndConf95 may be used as
// the f argument to NewErrorPoints.
func meanAndConf95(vls []float64) (mean, lowerr, higherr float64) {
	n := float64(len(vls))

	sum := 0.0
	for _, v := range vls {
		sum += v
	}
	mean = sum / n

	var stdev float64 = 0.0
	if n > 1 {
		sum = 0.0
		for _, v := range vls {
			diff := v - mean
			sum += diff * diff
		}
		stdev = math.Sqrt(sum / (n - 1.5))
	} else {
		stdev = 0
	}

	conf := 1.96 * stdev / math.Sqrt(n)
	return mean, conf, conf
}
func growth(radii []*Radii, showAllCurves bool, outputRaw bool, title string) (dims, seDims float64) {
	var (
		radiiVals = make([][]float64, len(radii))
	)

	for i, r := range radii {
		radiiVals[i] = r.Radii
	}

	radiiVals = util.Transpose(radiiVals)

	var pts = make([]plotter.XYer, len(radiiVals) - 1)
	for i, r := range radiiVals[1:] {
		xys := make(plotter.XYs, len(r))
		logN := math.Log10(float64(i + 2))
		for j, y := range r {
			xys[j].X = logN
			xys[j].Y = math.Log10(y)
		}
		pts[i] = xys
	}

	mean95, err := plotutil.NewErrorPoints(meanAndConf95, pts...)
	if err != nil {
		fmt.Println(err)
		return
	}
	regressionPoints := plotutil.ErrorPoints{
		XYs:mean95.XYs[len(mean95.XYs)/20:],
		XErrors:mean95.XErrors[len(mean95.XYs)/20:],
		YErrors:mean95.YErrors[len(mean95.XYs)/20:],
	}
	var a, b, eb float64
	if len(radii) > 1 {
		a, b, _, eb = util.WeightedLeastSquares(regressionPoints)
	} else {
		a, b, _, eb = util.LeastSquares(regressionPoints.XYs)
	}
	dims = 1/b
	seDims = eb/(b*b - eb*eb)
	if showAllCurves {
		// Draw the set of data without stopping the program
		go func() {
			plt, err := plot.New()
			if err != nil {
				fmt.Println(err)
				return
			}
			plt.X.Label.Text = "Log(N)"
			plt.Y.Label.Text = "Log(R)"

			scat, err := plotter.NewScatter(mean95)
			if err != nil {
				fmt.Println(err)
				return
			}
			scat.Shape = draw.SquareGlyph{}
			scat.Radius /= 2

			errs, err := plotter.NewYErrorBars(mean95)
			if err != nil {
				fmt.Println(err)
				return
			}
			errs.Width /= 2

			fit := plotter.NewFunction(func (x float64) float64 {return a + b*x})

			plt.Add(scat, errs, fit)
			plt.Legend.Add(fmt.Sprintf("y = %.3f + %.3fx", a, b), fit)

			path := fmt.Sprintf("%s\\%s.pdf", PLOT_PATH, title)
			plt.Save(8*vg.Inch, 5.66*vg.Inch, path)
			err = open.Run(path)
			if err != nil {
				fmt.Println(err)
				return
			}
		} ()
	}
	if outputRaw {
		// Saves the points and errors without stopping the program
		go func() {
			var file, err = os.Create(fmt.Sprintf("%s\\%s.csv", DATA_PATH, title))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			var writer = csv.NewWriter(file)
			for i := 0; i < len(regressionPoints.XYs); i++ {
				err = writer.Write([]string{
					strconv.FormatFloat(regressionPoints.XYs[i].X, 'E', -1, 64),
					strconv.FormatFloat(regressionPoints.XYs[i].Y, 'E', -1, 64),
					strconv.FormatFloat(regressionPoints.YErrors[i].High, 'E', -1, 64),
				})

				if err != nil {
					fmt.Println(err)
					return
				}
			}
		} ()
	}
	return dims, seDims

}

func growthWithStick(radii []*Radii, showAllCurves, outputRaw, showTrend bool, title string)  {
	type elem struct {
		s, dims, errDims float64
	}

	var data = make(map[float64][]*Radii)
	for _, r := range radii {
		if a, ok := data[r.Stick]; ok {
			data[r.Stick] = append(a, r)
		} else {
			data[r.Stick] = []*Radii{r}
		}
	}

	var elemChannel = make(chan elem)
	for s, R := range data {
		go func (s float64, R []*Radii) {
			var title = fmt.Sprintf("PLT_D%d_S%b", R[0].D, R[0].Stick)
			var dims, seDims = growth(R, showAllCurves, outputRaw, title)
			elemChannel <- elem{dims:dims, errDims:seDims, s:s}
		} (s, R)
	}
	var N = len(data)
	var points = plotutil.ErrorPoints{
		XYs:make(plotter.XYs, N),
		XErrors:make(plotter.XErrors, N),
		YErrors:make(plotter.YErrors, N),
	}
	for i := range points.XYs {
		var d = <- elemChannel
		points.XYs[i].X = d.s
		points.XYs[i].Y = d.dims
		points.YErrors[i].High = d.errDims
		points.YErrors[i].Low = d.errDims
	}
	if showTrend {
		go func () {
			var plt, err = plot.New()
			if err != nil {
				fmt.Println(err)
				return
			}
			plotutil.AddScatters(plt, points)
			plotutil.AddYErrorBars(plt, points)

			plt.X.Label.Text = "Sticking Probability"
			plt.Y.Label.Text = "Hausdorff Dimension"

			var path = fmt.Sprintf("%s\\%s.pdf", PLOT_PATH, title)
			plt.Save(8*vg.Inch, 5.66*vg.Inch, path)
			open.Run(path)
		} ()
	}
	if outputRaw {
		go func() {
			var file, err = os.Create(fmt.Sprintf("%s\\%s.csv", DATA_PATH, title))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			var writer = csv.NewWriter(file)
			for i := 0; i < len(points.XYs); i++ {
				err = writer.Write([]string{
					strconv.FormatFloat(points.XYs[i].X, 'E', -1, 64),
					strconv.FormatFloat(points.XYs[i].Y, 'E', -1, 64),
					strconv.FormatFloat(points.YErrors[i].High, 'E', -1, 64),
				})

				if err != nil {
					fmt.Println(err)
					return
				}
			}
		} ()
	}
}

func growthWithDims(radii []*Radii, showAllCurves, outputRaw, showTrend bool, title string)  {
	type elem struct {
		dims, errDims float64
		D int64
	}

	var data = make(map[int64][]*Radii)
	for _, r := range radii {
		if a, ok := data[r.D]; ok {
			data[r.D] = append(a, r)
		} else {
			data[r.D] = []*Radii{r}
		}
	}

	var elemChannel = make(chan elem)
	for D, R := range data {
		go func (D int64, R []*Radii) {
			var title = fmt.Sprintf("PLT_D%d_S%b", R[0].D, R[0].Stick)
			var dims, seDims = growth(R, showAllCurves, outputRaw, title)
			elemChannel <- elem{dims:dims, errDims:seDims, D:D}
		} (D, R)
	}
	var N = len(data)
	var points = plotutil.ErrorPoints{
		XYs:make(plotter.XYs, N),
		XErrors:make(plotter.XErrors, N),
		YErrors:make(plotter.YErrors, N),
	}
	for i := range points.XYs {
		var d = <- elemChannel
		points.XYs[i].X = float64(d.D)
		points.XYs[i].Y = d.dims
		points.YErrors[i].High = d.errDims
		points.YErrors[i].Low = d.errDims
	}
	if showTrend {
		go func () {
			var plt, err = plot.New()
			if err != nil {
				fmt.Println(err)
				return
			}
			plotutil.AddScatters(plt, points)
			plotutil.AddYErrorBars(plt, points)

			var dEqD = plotter.NewFunction(func (x float64) float64 {return x})
			var dEqDs2 = plotter.NewFunction(func (x float64) float64 {return x-2})
			dEqDs2.Dashes = []vg.Length{4*vg.Millimeter}
			plt.Add(dEqD, dEqDs2)
			plt.Legend.Add("d=D", dEqD)
			plt.Legend.Add("d=D-2", dEqDs2)

			plt.X.Label.Text = "Spatial Dimension"
			plt.Y.Label.Text = "Hausdorff Dimension"

			var path = fmt.Sprintf("%s\\%s.pdf", PLOT_PATH, title)
			plt.Save(8*vg.Inch, 5.66*vg.Inch, path)
			open.Run(path)
		} ()
	}
	if outputRaw {
		go func() {
			var file, err = os.Create(fmt.Sprintf("%s\\%s.csv", DATA_PATH, title))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			var writer = csv.NewWriter(file)
			for i := 0; i < len(points.XYs); i++ {
				err = writer.Write([]string{
					strconv.FormatFloat(points.XYs[i].X, 'E', -1, 64),
					strconv.FormatFloat(points.XYs[i].Y, 'E', -1, 64),
					strconv.FormatFloat(points.YErrors[i].High, 'E', -1, 64),
				})

				if err != nil {
					fmt.Println(err)
					return
				}
			}
		} ()
	}
}

func GrowthRate(runData []*RunData, showAllCurves, outputRaw, showTrend bool) {
	var (
		radiiChannel = make(chan *Radii)
		radii = make([]*Radii, len(runData))

		minD = int64(math.MaxInt64)
		maxD = int64(math.MinInt64)

		minS = math.Inf(1)
		maxS = math.Inf(-1)

		nData = 0
	)

	for _, data := range runData {
		var any = false
		if data.D < minD {
			minD = data.D
			any = true
		}
		if data.D > maxD {
			maxD = data.D
			any = true
		}
		if data.Stick < minS {
			minS = data.Stick
			any = true
		}
		if data.Stick < maxS {
			maxS = data.Stick
			any = true
		}

		if any {
			nData += 1
		}

		go func(data *RunData) {
			radiiChannel <- data.GyrationRadii()
		} (data)
	}

	for i := range radii {
		radii[i] = <-radiiChannel
	}

	if maxD > minD && maxS > minS {
		panic("This should not be possible!")
	} else if maxD > minD {
		var title = fmt.Sprintf("PLT_D%d-to%d-N%d_S%b", minD, maxD, nData, minS)
		growthWithDims(radii, showAllCurves, outputRaw, showTrend, title)
	} else if maxS > minS {
		var title = fmt.Sprintf("PLT_D%d_S%b-to%b-N%b", minD, minS, maxS, nData)
		growthWithStick(radii, showAllCurves, outputRaw, showTrend, title)
	} else {
		var title = fmt.Sprintf("PLT_D%d_S%b", minD, minS)
		var d, errD = growth(radii, showAllCurves, outputRaw, title)
		fmt.Printf("d = %.3f \u00b1 %.3f\n", d, errD)
	}


}