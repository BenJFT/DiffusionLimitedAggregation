package processing

import (
	"os"
	"fmt"
	"math/rand"
	"bufio"
	"encoding/gob"

	"github.com/Benjft/DiffusionLimitedAggregation/util/types"
	"github.com/Benjft/DiffusionLimitedAggregation/util/encoding/xyz"
	"github.com/Benjft/DiffusionLimitedAggregation/aggregation"

	"github.com/skratchdot/open-golang/open"
	"github.com/Benjft/DiffusionLimitedAggregation/util/encoding/svg"
)

func init() {
	os.Mkdir("out", os.ModeDir)
	os.Mkdir("out\\plot", os.ModeDir)
	os.Mkdir("out\\saves", os.ModeDir)
}

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

func draw3D(state []types.Point, title string, display bool) {
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
func draw2D(state []types.Point, title string, display bool) {
	name := fmt.Sprintf("out\\plot\\%s.svg", title)
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	str := svg.Format(state)
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
	var tmpRun types.Run
	err = decoder.Decode(&tmpRun)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lastRun = tmpRun
}