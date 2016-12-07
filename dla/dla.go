package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	agg "github.com/Benjft/DiffusionLimitedAggregation/aggregation"
	"github.com/Benjft/DiffusionLimitedAggregation/tools"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"

	"github.com/skratchdot/open-golang/open"
)

var lastStates [][]agg.Point

func runAggregation(seed, n, runs int64, sticking float64) {
	rand.Seed(seed)

	chans := make([]chan map[agg.Point]int64, runs)
	for i := int64(0); i < runs; i++ {
		a := &agg.Aggregator{}
		c := make(chan map[agg.Point]int64)
		go func() {
			rng := rand.New(rand.NewSource(rand.Int63()))
			c <- a.Aggregate(n, sticking, rng)
		} ()
		chans[i] = c
	}

	lastStates = make([][]agg.Point, runs)
	for i, c := range chans {
		arr := make([]agg.Point, n)
		m := <-c
		for k, v := range m {
			arr[v] = k
		}
		lastStates[i] = arr
		if i64 :=  int64(len(arr)); i64 < n {
			fmt.Println("MISSING", string(n-i64), "POINTS!")
		}
	}
}

func processRun(args []string) {
	var seed, n, runs int64
	var sticking float64

	if len(args) > 0 {
		var argstr string = args[0]

		for _, arg := range args[1:] {
			argstr += " "
			argstr += arg
		}

		argstr = strings.Replace(argstr, "=", " ", -1)
		args = strings.Split(argstr, " ")

		if len(args) > 6 {
			fmt.Println("Too much input")
			return
		}

		var key string = ""
		for _, arg := range args {
			switch key {
			case "seed":
				fmt.Sscanf(arg, "%d", &seed)
				key = ""
			case "n":
				fmt.Sscanf(arg, "%d", &n)
				key = ""
			case "runs":
				fmt.Sscanf(arg, "%d", &runs)
				key = ""
			case "sticking":
				fmt.Sscanf(arg, "%f", &sticking)
				key = ""
			default:
				switch arg {
				case "seed": key = "seed"
				case "n": key = "n"
				case "runs": key = "runs"
				case "sticking": key = "sticking"
				default:
					var num float64
					_, err := fmt.Sscanf(arg, "%f", &num)
					if err != nil {
						fmt.Println(arg, "Not recognised as a key or value")
						return
					}
					if seed == 0 {
						seed = int64(num)
					} else if n == 0 {
						n = int64(num)
					} else if sticking == 0 {
						sticking = num
					} else if runs == 0 {
						runs = int64(num)
					} else {
						fmt.Println("Too much input")
						return
					}
				}
			}
		}

		if seed < 1 {
			fmt.Println("Seed must be an int value of at least 1")
			return
		}
		if n < 2 {
			fmt.Println("N must be an int value of at least 2")
			return
		}
		if runs < 1 {
			fmt.Println("Runs must be an int value of at least 1")
			return
		}

	} else {
		fmt.Print("Seed = ")
		if _, err := fmt.Scanf("%d\n", &seed); err != nil || seed < 1 {
			fmt.Println("Seed must be an int value of at least 1")
			return
		}

		fmt.Print("N = ")
		if _, err := fmt.Scanf("%d\n", &n); err != nil || n < 2 {
			fmt.Println("N must be an int value of at least 2")
			return
		}

		fmt.Print("Sticking = ")
		if _, err := fmt.Scanf("%f\n", &sticking); err != nil || sticking <= 0 {
			fmt.Print("Sticking shouldbe a float latger than 0")
			return
		}

		fmt.Print("Runs = ")
		if _, err := fmt.Scanf("%d\n", &runs); err != nil || runs < 1 {
			fmt.Println("Runs mis be an int value of at least 1")
			return
		}
	}
	fmt.Println("Running, please wait")
	runAggregation(seed, n, runs, sticking)
	fmt.Println("Done!")
}

func processPlot(args []string) {
	if len(lastStates) > 0 {
		state := lastStates[0]
		xys := make(plotter.XYs, len(state))
		for i, v := range state {
			xys[i].X = float64(v.X)
			xys[i].Y = float64(v.Y)
		}

		plt, _ := plot.New()
		plt.Add(plotter.NewGrid())

		s, _ := plotter.NewScatter(xys)
		s.Shape = draw.PlusGlyph{}
		plt.Add(s)

		min := math.Min(plt.X.Min, plt.Y.Min)
		max := math.Max(plt.X.Max, plt.Y.Max)

		plt.X.Min = min
		plt.Y.Min = min
		plt.X.Max = max
		plt.Y.Max = max

		plt.Save(vg.Inch*10, vg.Inch*10, args[0])
		open.Run(args[0])
	}
}

func handle(str string) {
	str = tools.SingleSpace(str)
	strs := strings.Split(str, " ")
	head := strs[0]
	tail := strs[1:]

	switch head {
	case "run":
		processRun(tail)
	case "plot":
		processPlot(tail)
	default:
		fmt.Println("Command not recognized")
	}
}

func mainLoop() {
	fmt.Print(">>")
	for cmd := strings.ToLower(tools.ReadStrOrEmpty()); cmd != "quit"; cmd = strings.ToLower(tools.ReadStrOrEmpty()) {
		if cmd != "" {
			handle(cmd)
		}

		fmt.Print(">>")
	}
}

func main() {
	fmt.Println("Starting")

	args := os.Args[1:]
	if len(args) > 0 {
		// TODO
	} else {
		mainLoop()
	}

	fmt.Println("Stopped")
}
