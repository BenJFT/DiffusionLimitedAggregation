package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Benjft/DiffusionLimitedAggregation/processing"
	"github.com/Benjft/DiffusionLimitedAggregation/util"
)

var data []*processing.RunData
func runOnce(args []string) []string {
	var flagSet = flag.NewFlagSet("runonce", flag.ContinueOnError)
	var (
		n, d, seed int64
		stick float64
		tryLoad bool
	)

	flagSet.Int64Var(&n, "npoints", 5000, "The number of points desired in the final aggregate")
	flagSet.Int64Var(&d, "dimensions", 2, "The number of dimensions for the space of the simulation")
	flagSet.Int64Var(&seed, "seed", 1, "The seet to be input to the random number generator")

	flagSet.Float64Var(&stick, "pstick", 1, "The probability a point will stick to a neighboring filled site")

	flagSet.BoolVar(&tryLoad, "load", false, "Load the result if it has been found before for the specific config")

	var err = flagSet.Parse(args)

	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	data = []*processing.RunData {
		processing.RunOne(n, d, seed, stick, tryLoad),
	}

	return flagSet.Args()
}
func runMany(args []string) []string {
	var flagSet = flag.NewFlagSet("runmany", flag.ContinueOnError)
	var (
		n, d, seed, runs int64
		stick float64
		tryLoad bool
	)

	flagSet.Int64Var(&n, "npoints", 5000, "The number of points desired in the final aggregate")
	flagSet.Int64Var(&d, "dimensions", 2, "The number of dimensions for the space of the simulation")
	flagSet.Int64Var(&seed, "seed", 1, "The seed to be input to the random number generator. Generates seeds for" +
		" each individual run so each change generates a complete new set.")
	flagSet.Int64Var(&runs, "runs", 32, "The number of simulations to be run in this set")

	flagSet.Float64Var(&stick, "pstick", 1, "The probability a point will stick to a neighboring filled site")

	flagSet.BoolVar(&tryLoad, "load", false, "Load the result if it has been found before for the specific config")

	var err = flagSet.Parse(args)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	data = processing.RunMany(n, d, seed, stick, tryLoad, runs)

	return flagSet.Args()
}
func save(args []string) []string {
	var flagSet = flag.NewFlagSet("save", flag.ContinueOnError)


	var err = flagSet.Parse(args)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	if data != nil {
		processing.SaveAll(data)
	}
	return flagSet.Args()
}
func varyDimension(args []string) []string {
	var flagSet = flag.NewFlagSet("varydimension", flag.ContinueOnError)
	var (
		n, d0, d1, dStep, seed, runsPer int64
		stick float64
		tryLoad bool
	)
	flagSet.Int64Var(&n, "npoints", 5000, "The number of points desired in the final aggregate")
	flagSet.Int64Var(&d0, "start", 1, "The lowest number of dimensions to simulate (inclusive)")
	flagSet.Int64Var(&d1, "stop", 6, "The highest number of dimensions to simulate (inclusive)")
	flagSet.Int64Var(&dStep, "step", 1, "The gap between simulations")
	flagSet.Int64Var(&seed, "seed", 1, "The seed to be input to the random number generator. Generates seeds for" +
		" each individual run so each change generates a complete new set.")
	flagSet.Int64Var(&runsPer, "runs", 32, "The number of simulations to be run for each value of d")

	flagSet.Float64Var(&stick, "pstick", 1, "The probability a point will stick to a neighboring filled site")

	flagSet.BoolVar(&tryLoad, "load", false, "Load the result if it has been found before for the specific config")

	var err = flagSet.Parse(args)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	data = processing.VaryDimension(n, d0, d1, dStep, seed, stick, tryLoad, runsPer)
	return flagSet.Args()
}
func varySticking(args []string) []string {
	var flagSet = flag.NewFlagSet("varysticking", flag.ContinueOnError)
	var (
		n, d, seed, runsPer int64
		s0, s1, sStep float64
		tryLoad bool
	)
	flagSet.Int64Var(&n, "-npoints", 5000, "The number of points desired in the final aggregate")
	flagSet.Int64Var(&d, "-dimensions", 2, "The number of dimensions for the space of the simulation")
	flagSet.Int64Var(&seed, "-seed", 1, "The seed to be input to the random number generator. Generates seeds for" +
		" each individual run so each change generates a complete new set.")
	flagSet.Int64Var(&runsPer, "-runs", 32, "The number of simulations to be run for each value of d")

	flagSet.Float64Var(&s0, "-start", .1, "The lowest probability a point will stick to a neighboring filled site")
	flagSet.Float64Var(&s1, "-stop", 1, "The highest probability a point will stick to a neighboring filled site")
	flagSet.Float64Var(&sStep, "-step", .1, "The gap between simulations")

	flagSet.BoolVar(&tryLoad, "-load", false, "Load the result if it has been found before for the specific config")

	var err = flagSet.Parse(args)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	data = processing.VarySticking(n, d, seed, s0, s1, sStep, tryLoad, runsPer)
	return flagSet.Args()
}
func growth(args []string) []string {
	var flagSet = flag.NewFlagSet("growth", flag.ContinueOnError)
	var (
		hideCurves, outputRaw, hideTrend bool
	)

	flagSet.BoolVar(&hideCurves, "hidecurves", false, "Enable to prevent drawing a large number of curves for" +
		" multi value data sets")
	flagSet.BoolVar(&outputRaw, "raw", false, "Enable to allow curve data to be saved to a csv file for external" +
		" processing")
	flagSet.BoolVar(&hideTrend, "hidetrend", false, "Enable to prevent the overall trend from being plotted")

	var err = flagSet.Parse(args)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	processing.GrowthRate(data, !hideCurves, outputRaw, !hideTrend)
	return flagSet.Args()
}

var (
	handles = map[string]func([]string) []string{
		"runone": runOnce,
		"runmany": runMany,
		"save": save,
		"varydimension": varyDimension,
		"varysticking": varySticking,
		"growth": growth,
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
		} else if head == "go" {
			fmt.Println(runtime.NumGoroutine())
			args = tail
		} else if head == "gc" {
			runtime.GC()
			args = tail
		} else if f, ok := handles[head]; ok {
			args = f(tail)
		} else {
			if head != "help" {
				fmt.Printf("'%s' not recognised as valid\n", head)
				fmt.Println("Try:")
			}
			for k := range handles {
				fmt.Println("\t" + k)
			}
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
	println("Hello, World!")

	// If arguments were given run using those, else start the prompt loop
	if len(args) > 0 {
		handleArgs(args)
	} else {
		mainLoop()
	}
}
