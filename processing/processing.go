package processing

import (
	"os"
	"fmt"
	"math"
	"math/rand"
	"bufio"
	"encoding/gob"

	"github.com/Benjft/DiffusionLimitedAggregation/aggregation"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/analysis"
	"github.com/Benjft/DiffusionLimitedAggregation/util"
	"github.com/Benjft/DiffusionLimitedAggregation/util/encoding/svg"
	"github.com/Benjft/DiffusionLimitedAggregation/util/encoding/xyz"

	"github.com/skratchdot/open-golang/open"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/plotutil"

)

func init() {
	os.Mkdir("out", os.ModeDir)
	os.Mkdir("out\\plot", os.ModeDir)
	os.Mkdir("out\\saves", os.ModeDir)

	gob.Register(RunState{})
}


type RunState struct {
	NPoints, NDimension, NRuns, Seed int64
	Sticking float64
	Points [][]aggregation.Point
}
var (
	lastRun = RunState{}
)

func run(nPoints, seed, nDimension int64, sticking float64, chanel chan []aggregation.Point) {
	chanel <- aggregation.RunNew(nPoints, seed, nDimension, sticking)
}
func Run(nPoints, nRuns, seed, nDimension int64, sticking float64) {
	var channels []chan []aggregation.Point = make([]chan []aggregation.Point, nRuns)

	rand.Seed(seed)
	for i := range channels {
		var channel chan []aggregation.Point = make(chan []aggregation.Point)
		go run(nPoints, rand.Int63(), nDimension, sticking, channel)
		channels[i] = channel
	}

	var (
		runSuccessful bool = true
		points [][]aggregation.Point = make([][]aggregation.Point, nRuns)
	)
	for i, channel := range channels {
		state := <-channel
		if state == nil {
			fmt.Printf("Thread %d Failed!\n", i)
			runSuccessful = false
		} else {
			points[i] = state
		}
	}

	if runSuccessful {
		lastRun = RunState{
			NPoints: nPoints,
			NDimension: nDimension,
			NRuns: nRuns,
			Seed: seed,
			Sticking: sticking,
			Points: points,
		}
	}
}

func draw3D(state []aggregation.Point, title string, display bool) {
	var name string = fmt.Sprintf("out\\plot\\%s.%s", title, "xyz")
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	str := xyz.Format(state)
	_, err = writer.Write([]byte(str))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if display {
		err = open.Run(name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
func draw2D(state []aggregation.Point, title string, display bool) {
	name := fmt.Sprintf("out\\plot\\%s.svg", title)
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	str := svg.DrawAggregate(state)
	_, err = writer.Write([]byte(str))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if display {
		err = open.Run(name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
func Draw(title string, display bool) {
	if len(title) == 0 {
		title = fmt.Sprintf("aggregate-n%d-seed%d-dims%d-stick%f",
			lastRun.NPoints,
			lastRun.Seed,
			lastRun.NDimension,
			lastRun.Sticking)
	}

	if n := lastRun.NDimension; n == 2 {
		for run, state := range lastRun.Points {
			runtitle := fmt.Sprintf("%s-run%d", title, run)
			go draw2D(state, runtitle, display)
		}
	} else if n == 3 {
		for run, state := range lastRun.Points {
			runtitle := fmt.Sprintf("%s-run%d", title, run)
			go draw3D(state, runtitle, display)
		}
	} else {
		fmt.Println("Can only draw 2D and 3D lattices")
	}
}

func Save(title string) {
	if title == "" {
		title = fmt.Sprintf("save-n%d-seed%d-dims%d-stick%f-runs%d", lastRun.NPoints, lastRun.Seed,
			lastRun.NDimension, lastRun.Sticking, lastRun.NRuns)
	}
	path := fmt.Sprintf("out\\saves\\%s.save", title)

	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err.Error())

	}

	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(lastRun)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Load(title string) {
	path := fmt.Sprintf("out\\saves\\%s.save", title)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	decoder := gob.NewDecoder(file)
	var tmpRun RunState
	err = decoder.Decode(&tmpRun)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lastRun = tmpRun
}

func radii(run []aggregation.Point, chanel chan []analysis.Ball) {
	chanel <- analysis.ApproxBounding(run)
}
func Radii() {
	var channels []chan []analysis.Ball = make([]chan []analysis.Ball, lastRun.NRuns)
	for i, run := range lastRun.Points {
		var channel chan []analysis.Ball = make(chan []analysis.Ball)
		go radii(run, channel)
		channels[i] = channel
	}

	radii := make([][]float64, lastRun.NRuns)
	for i, channel := range channels {
		radii[i] = make([]float64, lastRun.NPoints)
		runBalls := <-channel
		for j, ball := range runBalls {
			radii[i][j] = ball.Radius
		}
	}

	//radius, stdErr := util.MeanAndErr(radii)

	radii = util.Transpose(radii)

	pts := make([]plotter.XYer, len(radii))
	allXY := plotter.XYs{}
	for i, r := range radii {
		xys := make(plotter.XYs, len(r))
		N := float64(i+1)
		for j, y := range r {
			xys[j].X = math.Log10(y+.5)
			xys[j].Y = math.Log10(N)
		}
		pts[i] = xys

		// Ignore the first few as growth will be non typical
		if i > 10 {
			allXY = append(allXY, xys...)
		}
	}

	plt , err := plot.New()
	if err != nil {
		panic(err)
	}

	mean95, err := plotutil.NewErrorPoints(plotutil.MeanAndConf95, pts...)
	if err != nil {
		panic(err)
	}

	a, b, _, eb := util.LeastSquares(allXY)
	fmt.Printf("D = %.3f \u00B1 %.3f\n", b, eb)
	label := fmt.Sprintf("y = %.3f + %.3fx", a, b)
	fit := plotter.NewFunction(func (x float64) float64 { return a + b*x })

	plotutil.AddScatters(plt, mean95)
	plotutil.AddXErrorBars(plt, mean95)
	plt.Add(fit)
	plt.Legend.Add(label, fit)

	name := "out\\plot\\TEST.svg"
	plt.Save(15*vg.Inch, 10*vg.Inch, name)
	open.Run(name)
}