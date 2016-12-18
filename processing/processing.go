package processing

import (
	"os"
	"fmt"
	"math"
	"math/rand"
	"encoding/gob"

	"github.com/Benjft/DiffusionLimitedAggregation/util/types"
	"github.com/Benjft/DiffusionLimitedAggregation/aggregation"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	vgdraw "github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/plotter"

	"github.com/skratchdot/open-golang/open"
)

var (
	lastRun = types.Run{}
)

func run(nPoints, seed, nDimension int64, sticking float64, chanel chan []types.Point) {
	chanel <- aggregation.RunNew(nPoints, seed, nDimension, sticking)
}
func Run(nPoints, nRuns, seed, nDimension int64, sticking float64) {
	var channels []chan []types.Point = make([]chan []types.Point, nRuns)

	rand.Seed(seed)
	for i := range channels {
		var channel chan []types.Point
		channel = make(chan []types.Point)
		go run(nPoints, rand.Int63(), nDimension, sticking, channel)
		channels[i] = channel
	}

	var (
		runSuccessful bool = true
		points [][]types.Point = make([][]types.Point, nRuns)
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
		lastRun = types.Run {
			NPoints: nPoints,
			NDimension: nDimension,
			NRuns: nRuns,
			Seed: seed,
			Sticking: sticking,
			Points: points,
		}
	}
}

func draw(state []types.Point, title, format string, width int64, display bool) {
	var (
		N int = len(state)
		plt *plot.Plot
		hsva types.HSVA = types.HSVA{H:0, S:1, V:.8, A:1}
	)

	plt, _ = plot.New()

	for i, point := range state {

		xy := point.Coordinates()
		x, y := float64(xy[0]), float64(xy[1])
		xys := make(plotter.XYs, 1)
		xys[0].X = x
		xys[0].Y = y

		s, _ := plotter.NewScatter(xys)

		hsva.H = 300*float64(i)/(360*float64(N))
		s.Color = hsva
		s.Shape = vgdraw.BoxGlyph{}
		plt.Add(s)
	}

	var min, max float64
	min = math.Min(plt.X.Min, plt.Y.Min)
	max = math.Max(plt.X.Max, plt.Y.Max)

	plt.X.Min = min
	plt.Y.Min = min
	plt.X.Max = max
	plt.Y.Max = max

	var name string = fmt.Sprintf("out\\plot\\%s.%s", title, format)
	var w vg.Length
	if width == 0 {
		w = vg.Millimeter*2*vg.Length(max-min)
	} else {
		w = vg.Millimeter*vg.Length(width)
	}

	plt.Save(w, w, name)
	if display {
		open.Run(name)
	}
}
func Draw(title, format string, width int64, display bool) {
	if len(title) == 0 {
		title = fmt.Sprintf("aggregate-n%d-seed%d-dims%d-stick%f",
			lastRun.NPoints,
			lastRun.Seed,
			lastRun.NDimension,
			lastRun.Sticking)
	}
	for run, state := range lastRun.Points {
		runtitle := fmt.Sprintf("%s-run%d", title, run)

		go draw(state, runtitle, format, width, display)
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
	var tmpRun types.Run
	err = decoder.Decode(&tmpRun)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lastRun = tmpRun
}