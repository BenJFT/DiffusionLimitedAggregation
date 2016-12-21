package main

import (
	"os"
	"fmt"
	"flag"
	"strings"

	"github.com/Benjft/DiffusionLimitedAggregation/util"
	"github.com/Benjft/DiffusionLimitedAggregation/processing"
)

var (
	formats map[string]bool = map[string]bool {
		"png": true,
		"jpg": true,
		"jpeg": true,
		"svg": true,
		"tif": true,
		"tiff": true,
	}
)
type Handler func (args []string) (tail []string)

func handleRun(args []string) (tail []string) {
	var flags *flag.FlagSet
	flags = flag.NewFlagSet("run", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	var nPoints, nRuns, nDimension, seed int64
	var sticking float64

	flags.Int64Var(&nPoints, "points", 2000, "the number of points to aggregate (Minimum 1)")
	flags.Int64Var(&nRuns, "runs", 1, "the number of aggregates to run (Minimum 1)")
	flags.Int64Var(&nDimension, "dimension", 2, "the number of dimensions (Minimum 2)")
	flags.Int64Var(&seed, "seed", 1, "the seed to run the set of simulations from")

	flags.Float64Var(&sticking, "sticking", 1, "probability of a point sticking to an adjacent point per time step")

	var err error = flags.Parse(args)
	if err != nil {
		fmt.Println(err.Error())
	} else if nPoints < 1 || nRuns < 1 || nDimension < 2 || sticking <= 0 || sticking > 1 {
		flags.PrintDefaults()
	} else {
		tail = flags.Args()
		processing.Run(nPoints, nRuns, seed, nDimension, sticking)
	}
	return tail
}

func handleDraw(args []string) (tail []string) {
	var flags *flag.FlagSet = flag.NewFlagSet("draw", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	var name string
	var display bool

	flags.StringVar(&name, "name", "", "The name for the file. (Default based on run args)")

	flags.BoolVar(&display, "disp", false, "Auto open the drawn figures")

	var err error = flags.Parse(args)
	if err != nil {
		fmt.Println(err)
	} else {
		tail = flags.Args()
		processing.Draw(name, display)
	}
	return tail
}

func handleSave(args []string) (tail []string) {
	var flags *flag.FlagSet = flag.NewFlagSet("save", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	var name string

	flags.StringVar(&name, "name", "", "The name for the file. (Default based on run args)")

	var err error = flags.Parse(args)
	if err != nil {
		fmt.Println(err)
	} else {
		tail = flags.Args()
		processing.Save(name)
	}
	return tail
}

func handleLoad(args []string) (tail []string) {
	var flags *flag.FlagSet = flag.NewFlagSet("save", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	var name string
	var nPoints, nRuns, nDimension, seed int64
	var sticking float64

	flags.StringVar(&name, "name", "", "The name for the file. (Default based on run args)")

	flags.Int64Var(&nPoints, "points", 2000, "the number of points to aggregate (Minimum 1)")
	flags.Int64Var(&nRuns, "runs", 1, "the number of aggregates to run (Minimum 1)")
	flags.Int64Var(&nDimension, "dimension", 2, "the number of dimensions (Minimum 2)")
	flags.Int64Var(&seed, "seed", 1, "the seed to run the set of simulations from")

	flags.Float64Var(&sticking, "sticking", 1, "probability of a point sticking to an adjacent point per time step")

	var err error = flags.Parse(args)
	if err != nil {
		fmt.Println(err)
	} else {
		if name == "" {
			name = fmt.Sprintf("save-n%d-seed%d-dims%d-stick%f-runs%d", nPoints, seed, nDimension, sticking,
				nRuns)
		}

		tail = flags.Args()
		processing.Load(name)
	}
	return tail
}

var (
	handles map[string]Handler = map[string]Handler {
		"run": Handler(handleRun),
		"draw": Handler(handleDraw),
		"save": Handler(handleSave),
		"load": Handler(handleLoad),
	}
)

func handleArgs(args []string) bool {
	var head string
	var tail []string
	for len(args) > 0 {
		head = strings.ToLower(args[0])
		tail = args[1:]

		if head == "quit" || head == "stop" {
			return false
		} else if f, ok := handles[head]; ok {
			args = f(tail)
		} else {
			fmt.Printf("'%s' not recognised as valid\n", head)
			return true
		}
	}
	return true
}

func handleInstructions(instructions string) bool {
	var args []string = util.StringToArgs(instructions)
	return handleArgs(args)
}

func mainLoop() {
	fmt.Print(">>")

	// Scans the input line for the next instruction.
	for cmd := util.ReadStrOrEmpty(); handleInstructions(cmd); cmd = util.ReadStrOrEmpty() {
		fmt.Print(">>")
	}
}

func main() {
	// Gets the command line arguments
	args := os.Args[1:]

	// If arguments were given run using those, else start the prompt loop
	if len(args) > 0 {
		handleArgs(args)
	} else {
		mainLoop()
	}
}
