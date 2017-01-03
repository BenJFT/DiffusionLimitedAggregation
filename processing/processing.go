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
	"github.com/Benjft/DiffusionLimitedAggregation/util/drawing/svg"
	"github.com/Benjft/DiffusionLimitedAggregation/util/drawing/xyz"

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
	loadedRun = RunState{}
)

func run(nPoints, seed, nDimension int64, sticking float64, chanel chan []aggregation.Point) {
	chanel <- aggregation.RunNew(nPoints, seed, nDimension, sticking)
}
func Run(nPoints, nRuns, seed, nDimension int64, sticking float64) {
	var (
		ch0, ch1, ch2, ch3 chan []aggregation.Point
		points [][]aggregation.Point = make([][]aggregation.Point, nRuns)
		out []aggregation.Point
	)

	rand.Seed(seed)

	for i := int64(0); i < nRuns; {

		if ch0 == nil {
			ch0 = make(chan []aggregation.Point)
			go run(nPoints, rand.Int63(), nDimension, sticking, ch0)
		} else if ch1 == nil && nRuns > 1 {
			ch1 = make(chan []aggregation.Point)
			go run(nPoints, rand.Int63(), nDimension, sticking, ch1)
		} else if ch2 == nil && nRuns > 2 {
			ch2 = make(chan []aggregation.Point)
			go run(nPoints, rand.Int63(), nDimension, sticking, ch2)
		} else  if ch3 == nil && nRuns > 3 {
			ch3 = make(chan []aggregation.Point)
			go run(nPoints, rand.Int63(), nDimension, sticking, ch3)
		} else {
			select {
			case out = <-ch0:
				points[i] = out
				if nRuns > 4 && i < nRuns-1 {
					go run(nPoints, rand.Int63(), nDimension, sticking, ch0)
				}
			case out = <-ch1:
				points[i] = out
				if nRuns > 4 && i < nRuns-1 {
					go run(nPoints, rand.Int63(), nDimension, sticking, ch1)
				}
			case out = <-ch2:
				points[i] = out
				if nRuns > 4 && i < nRuns-1 {
					go run(nPoints, rand.Int63(), nDimension, sticking, ch2)
				}
			case out = <-ch3:
				points[i] = out
				if nRuns > 4 && i < nRuns-1 {
					go run(nPoints, rand.Int63(), nDimension, sticking, ch3)
				}
			}

			if out == nil {
				fmt.Println("Run failed!")
				return
			} else {
				i++
			}
		}
	}

	loadedRun = RunState {
		NPoints: nPoints,
		NDimension: nDimension,
		NRuns: nRuns,
		Seed: seed,
		Sticking: sticking,
		Points: points,
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
			loadedRun.NPoints,
			loadedRun.Seed,
			loadedRun.NDimension,
			loadedRun.Sticking)
	}

	if n := loadedRun.NDimension; n == 2 {
		for run, state := range loadedRun.Points {
			runtitle := fmt.Sprintf("%s-run%d", title, run)
			go draw2D(state, runtitle, display)
		}
	} else if n == 3 {
		for run, state := range loadedRun.Points {
			runtitle := fmt.Sprintf("%s-run%d", title, run)
			go draw3D(state, runtitle, display)
		}
	} else {
		fmt.Println("Can only draw 2D and 3D lattices")
	}
}

func Save(title string) {
	if title == "" {
		title = fmt.Sprintf("save-n%d-seed%d-dims%d-stick%f-runs%d", loadedRun.NPoints, loadedRun.Seed,
			loadedRun.NDimension, loadedRun.Sticking, loadedRun.NRuns)
	}
	path := fmt.Sprintf("out\\saves\\%s.save", title)

	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err.Error())

	}

	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(loadedRun)
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

	loadedRun = tmpRun
}

func radii(run []aggregation.Point, chanel chan []analysis.Ball) {
	chanel <- analysis.ApproxBounding(run)
}
func Radii() {
	var channels []chan []analysis.Ball = make([]chan []analysis.Ball, loadedRun.NRuns)
	for i, run := range loadedRun.Points {
		var channel chan []analysis.Ball = make(chan []analysis.Ball)
		go radii(run, channel)
		channels[i] = channel
	}

	radii := make([][]float64, loadedRun.NRuns)
	for i, channel := range channels {
		radii[i] = make([]float64, loadedRun.NPoints)
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
		if i > len(radii)/50 {
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