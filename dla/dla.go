package main

import (
	"os"
	"fmt"
	"flag"
	"strings"

	"github.com/Benjft/DiffusionLimitedAggregation/tools"
	proc "github.com/Benjft/DiffusionLimitedAggregation/processing"
)

var (
	runSeed int64
	runN int64
	runRuns int64
	runSticking float64

	runStates [][]tools.Point
)

//var formats := map[string]

func handleRun(args []string) {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)

	flags.Int64Var(&runSeed, "seed", 1, "an integer of at least 1")
	flags.Int64Var(&runN, "num", 1000, "an integer of at least 2")
	flags.Int64Var(&runRuns, "runs", 1, "an integer of at least 1")
	flags.Float64Var(&runSticking, "sticking", 1, "a float satifying 0 < f <= 1")

	err := flags.Parse(args)
	if err == nil {
		fmt.Printf("Seed = %d\n Num = %d\n Runs = %d\n Sticking = %f\n", runSeed, runN, runRuns, runSticking)
		fmt.Println("Running, please wait")
		runStates = proc.Run(runSeed, runN, runRuns, runSticking)
		fmt.Println("Done")
	}
}

func handleDraw(args []string) {
	flags := flag.NewFlagSet("plot", flag.ContinueOnError)

	var title, format string
	var display bool

	title = fmt.Sprintf("aggregate-seed%d-n%d-sticking%f", runSeed, runN, runSticking)

	flags.StringVar(&title, "title", title, "the header for the plot and the name of the file")
	flags.StringVar(&format, "format", "svg", "the file type to output the plot as. (allowed svg, png, jpg, tif")
	flags.BoolVar(&display, "display", true, "open the plot after saving. Opens in the befault web browser")

	err := flags.Parse(args)
	if err != nil {
		return
	} else if format != "svg" && format != "png" && format != "jpg" && format != "tif" {
		flags.Usage()
		return
	}

	for i, state := range runStates {
		rTitle := fmt.Sprintf("%s-run%d", title, i)
		fmt.Printf("title-%s  fmt-%s  disp-%t\n", rTitle, format, display)
		proc.DrawTime(state, rTitle, format, display)
	}
}

func handleDimension(args []string) {
	flags := flag.NewFlagSet("dimesion", flag.ContinueOnError)

	var title, format string
	var display bool

	title = fmt.Sprintf("dimension-seed%d-n%d-runs%d-sticking%f", runSeed, runN, runRuns, runSticking)

	flags.StringVar(&title, "title", title, "the header for the plot and the name of the file")
	flags.StringVar(&format, "format", "svg", "the file type to output the plot as. (allowed svg, png, jpg, tif")
	flags.BoolVar(&display, "display", true, "open the plot after saving. Opens in the befault web browser")

	flags.Parse(args)

	a, b := proc.Dimension(runStates, title, format, true)
	fmt.Printf("y=%.3fx+%.3f\n",a, b)
}

func handleRadii(args []string) {
	flags := flag.NewFlagSet("radii", flag.ContinueOnError)

	var title, format string
	var display bool

	title = fmt.Sprintf("radii-seed%d-n%d-runs%d-sticking%f", runSeed, runN, runRuns, runSticking)

	flags.StringVar(&title, "title", title, "the header for the plot and the name of the file")
	flags.StringVar(&format, "format", "svg", "the file type to output the plot as. (allowed svg, png, jpg, tif")
	flags.BoolVar(&display, "display", true, "open the plot after saving. Opens in the befault web browser")

	flags.Parse(args)
	proc.Radius(runStates, title, format, display)
}

func handle(strs []string) {
	head := strs[0]
	tail := strs[1:]

	switch head {
	case "run":
		handleRun(tail)
	case "draw":
		handleDraw(tail)
	case "dimension":
		handleDimension(tail)
	case "radii":
		handleRadii(tail)
	default:
		fmt.Println("Command not recognized")
	}
}

func mainLoop() {
	fmt.Print(">>")
	for cmd := strings.ToLower(tools.ReadStrOrEmpty()); cmd != "quit"; cmd = strings.ToLower(tools.ReadStrOrEmpty()) {
		if cmd != "" {
			handle(strings.Split(tools.SingleSpace(cmd), " "))
		}

		fmt.Print(">>")
	}
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		handle(args)
	} else {
		mainLoop()
	}
}
