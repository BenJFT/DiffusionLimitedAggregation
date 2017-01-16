package processing

import (
	"bufio"
	"encoding/gob"
	"encoding/csv"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/analysis"
	"github.com/Benjft/DiffusionLimitedAggregation/util"
	"github.com/Benjft/DiffusionLimitedAggregation/util/drawing/svg"

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

// Structure to hold data from a single run
type RunData struct {
	N, D, Seed int64
	Stick float64
	Points []aggregation.Point
	radii []float64
}

// runs the simulation for the aggreagte's conditions
func (run *RunData) Run() {
	run.Points = aggregation.RunNew(run.N, run.Seed, run.D, run.Stick)
}

// loads the radii for the run
func (run *RunData) Radii() []float64 {
	if run.radii == nil {
		run.radii = analysis.GyrationRadii(run.Points)
	}
	return run.radii
}

// the file name specified by the runs properties
func (run *RunData) FileName() string {
	return fmt.Sprintf("N%d-D%d-Se%d-St%b", run.N, run.D, run.Seed, run.Stick)
}

// saves the gob of the aggregate
func (run *RunData) Save() error {
	var path = fmt.Sprintf("%s\\%s.save", SAVE_PATH, run.FileName())

	var file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var encoder = gob.NewEncoder(file)
	err = encoder.Encode(run)
	if err != nil {
		return err
	}
	return nil
}

// loads the aggregate from a gob
func (run *RunData) Load() error {
	var path = fmt.Sprintf("%s\\%s.save", SAVE_PATH, run.FileName())

	var file, err = os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var decoder = gob.NewDecoder(file)
	err = decoder.Decode(run)
	if err != nil {
		return err
	}
	return nil
}

// saves the aggregate as a csv file of it's coordinates
func (run *RunData) CSV() error {
	var path = fmt.Sprintf("%s\\%s.csv", DATA_PATH, run.FileName())

	var file, err = os.Create(path)
	if err != nil {
		return err

	}
	defer file.Close()
	var writer = csv.NewWriter(file)
	for _, point := range run.Points {
		var strs = make([]string, run.D)
		for i, x := range point.Coordinates() {
			strs[i] = strconv.FormatInt(x, 10)
		}
		err = writer.Write(strs)
		if err != nil {
			return err
		}
	}
	return nil
}

// saves the aggregate as a csv and gob
func SaveAll(runs []*RunData) {
	for _, run := range runs {
		go func() {
			var err = run.Save()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = run.CSV()
			if err != nil {
				fmt.Println(err)
				return
			}
		} ()
	}
}

// runs a single simulation with the passed parameters.
func RunOne(n, d, seed int64, stick float64, tryLoad bool) (run *RunData) {
	run = &RunData{N:n, D:d, Seed:seed, Stick:stick}

	// attempts to load the run from storage if it exists. If not then runs a new one
	switch tryLoad {
	case true:
		var err = run.Load()
		if err == nil{
			fmt.Println(run.FileName(), "Loaded.")
			break
		}
		fallthrough
	default:
		fmt.Println(run.FileName(), "Started Running.")
		run.Run()
		fmt.Println(run.FileName(), "Finished Running.")
	}
	run.Save()
	return
}

// runs a set of runs in separate concurrent processes and waits for them all to compete
func RunMany(n, d, seed int64, stick float64, tryLoad bool, runs int64) (runData []*RunData) {
	runData = make([]*RunData, runs)
	var rng = rand.New(rand.NewSource(seed))

	var runChannel = make(chan *RunData)

	// sets each process running
	for range runData {
		go func (seed int64) {
			runChannel <- RunOne(n, d, seed, stick, tryLoad)
		} (rng.Int63())
	}

	// waits for each to finish and collects their results
	for i := range runData {
		runData[i] = <-runChannel
	}
	return runData
}

// runs a set of runs for each specified dimension in a concurrent process
func VaryDimension(n, seed int64, stick float64, tryLoad bool, runsPer int64, steps []int64) (runData []*RunData) {
	var (
		nPoints int64 = 0
		runChannel = make(chan []*RunData)
	)

	// starts each process
	for _, d := range steps {
		nPoints++
		go func (d int64) {
			runChannel <- RunMany(n, d, seed, stick, tryLoad, runsPer)
		} (d)
	}

	runData = make([]*RunData, runsPer * nPoints)
	// collects results
	for i := int64(0); i < nPoints; i++ {
		idx := i*runsPer
		for j, data := range <-runChannel {
			runData[idx+int64(j)] = data
		}
	}

	return runData
}

func VarySticking(n, d, seed int64, tryLoad bool, runsPer int64, steps []float64) (runData []*RunData) {
	var (
		nSteps = int64(len(steps))
		runChannel = make(chan []*RunData)
	)
	runData = make([]*RunData, runsPer * nSteps)

	// starts each process
	for _, s := range steps {
		go func (stick float64) {
			runChannel <- RunMany(n, d, seed, stick, tryLoad, runsPer)
		} (s)
	}

	// waits for results
	for i := int64(0); i < nSteps; i++ {
		idx := i*runsPer
		for j, data := range <-runChannel {
			runData[idx+int64(j)] = data
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
// estimates the hausdorff dimension for the set of data
func growth(radii []*RunData, showAllCurves bool, outputRaw bool, title string) (dims, seDims float64) {
	var (
		radiiVals = make([][]float64, len(radii))
	)

	for i, r := range radii {
		radiiVals[i] = r.Radii()
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

	var n = int(float64(len(mean95.XYs))*0.5) // last 50% off aggregates
	regressionPoints := plotutil.ErrorPoints{
		XYs:mean95.XYs[n:],
		XErrors:mean95.XErrors[n:],
		YErrors:mean95.YErrors[n:],
	}
	var a, b, _, eb = util.WeightedLeastSquares(regressionPoints)
	dims = 1/b
	seDims = eb/(b*b)
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
			scat.Color = color.RGBA{A:255, R:200, G:55, B:55}

			errs, err := plotter.NewYErrorBars(mean95)
			if err != nil {
				fmt.Println(err)
				return
			}
			errs.Width /= 2
			errs.Color = color.RGBA{A:255, R:255, G:55, B:55}

			fit := plotter.NewFunction(func (x float64) float64 {return a + b*x})

			plt.Add(errs, scat, fit)
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

// finds the hausdorff dimension for each different value of sticking probability
func growthWithStick(radii []*RunData, showAllCurves, outputRaw, showTrend bool, title string)  {
	type elem struct {
		s, dims, errDims float64
	}

	var data = make(map[float64][]*RunData)
	for _, r := range radii {
		if a, ok := data[r.Stick]; ok {
			data[r.Stick] = append(a, r)
		} else {
			data[r.Stick] = []*RunData{r}
		}
	}

	var elemChannel = make(chan elem)
	for s, R := range data {
		go func (s float64, R []*RunData) {
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
		points.XYs[i].X = math.Log(d.s)
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

			plt.X.Label.Text = "Log(S)"
			plt.Y.Label.Text = "D"

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

// finds the hausdorff dimension for each different value of D
func growthWithDims(radii []*RunData, showAllCurves, outputRaw, showTrend bool, title string)  {
	type elem struct {
		dims, errDims float64
		D int64
	}

	var data = make(map[int64][]*RunData)
	for _, r := range radii {
		if a, ok := data[r.D]; ok {
			data[r.D] = append(a, r)
		} else {
			data[r.D] = []*RunData{r}
		}
	}

	var elemChannel = make(chan elem)
	for D, R := range data {
		go func (D int64, R []*RunData) {
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
			var dEqDs2 = plotter.NewFunction(func (x float64) float64 {return x-1})
			var TK = plotter.NewFunction(func (x float64) float64 {return (x*x+1)/(x+1)})
			dEqDs2.Dashes = []vg.Length{4*vg.Millimeter}
			TK.Color = color.RGBA{A:255, B:255}
			plt.Add(dEqD, dEqDs2, TK)
			plt.Legend.Add("Upper bound D=d", dEqD)
			plt.Legend.Add("Causality relationsip D=d-1", dEqDs2)
			plt.Legend.Add("TK Model for growth", TK)

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

// finds the hausdorff dimension and sepeates based on varied parameters
func GrowthRate(runData []*RunData, showAllCurves, outputRaw, showTrend bool) {
	var (
		minD = int64(math.MaxInt64)
		maxD = int64(math.MinInt64)
		D = map[int64]bool{}

		minS = math.Inf(1)
		maxS = math.Inf(-1)
		S = map[float64]bool{}
	)

	for _, data := range runData {
		if data.D < minD {
			minD = data.D
		}
		if data.D > maxD {
			maxD = data.D
		}
		if data.Stick < minS {
			minS = data.Stick
		}
		if data.Stick > maxS {
			maxS = data.Stick
		}
		D[data.D]=true
		S[data.Stick]=true
	}

	var nData = len(D)*len(S)

	if maxD > minD && maxS > minS {
		panic("This should not be possible!")
	} else if maxD > minD {
		var title = fmt.Sprintf("PLT_D%d-to%d-N%d_S%b", minD, maxD, nData, minS)
		growthWithDims(runData, showAllCurves, outputRaw, showTrend, title)
	} else if maxS > minS {
		var title = fmt.Sprintf("PLT_D%d_S%b-to%b-N%b", minD, minS, maxS, nData)
		growthWithStick(runData, showAllCurves, outputRaw, showTrend, title)
	} else {
		var title = fmt.Sprintf("PLT_D%d_S%b", minD, minS)
		var d, errD = growth(runData, showAllCurves, outputRaw, title)
		fmt.Printf("d = %.3f \u00b1 %.3f\n", d, errD*1.96)
	}


}

// draws 2D aggregates to svg files
func Draw(runData []*RunData) {
	for _, run := range runData {
		go func(run *RunData) {
			if run.D != 2 {
				return
			}

			var path = fmt.Sprintf("%s\\%s.svg", DRAW_PATH, run.FileName())
			var file, err = os.Create(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			var writer = bufio.NewWriter(file)
			_, err = writer.Write([]byte(svg.DrawAggregate(run.Points)))
			if err != nil {
				fmt.Println(err)
				return
			}
			open.Run(path)
		} (run)
	}
}