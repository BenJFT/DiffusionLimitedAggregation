/*
Handles user input and passes tasks off to be dealt with
 */

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"strconv"

	"github.com/Benjft/DiffusionLimitedAggregation/processing"
	"github.com/Benjft/DiffusionLimitedAggregation/util"
)

var data []*processing.RunData
// runs a set of simulations with the same parameters
func runMany(args []string) []string {
	var flagSet = flag.NewFlagSet("run", flag.ContinueOnError)
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
// saves the aggregates in csv form
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
// runs multiple simulations with a varying value for d
func varyDimension(args []string) []string {
	var flagSet = flag.NewFlagSet("varydimension", flag.ContinueOnError)
	var (
		n, d0, d1, dStep, seed, runsPer int64
		stick float64
		tryLoad bool
	)
	flagSet.Int64Var(&n, "npoints", 5000, "The number of points desired in the final aggregate")
	flagSet.Int64Var(&d0, "start", 2, "The lowest number of dimensions to simulate (inclusive)")
	flagSet.Int64Var(&d1, "stop", 6, "The highest number of dimensions to simulate (inclusive)")
	flagSet.Int64Var(&dStep, "step", 1, "The gap between simulations (0 to read values from tail)")
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

	var steps []int64= []int64{}
	var tail = flagSet.Args()
	if dStep == 0 {
		for _, arg := range tail {
			i, err := strconv.Atoi(arg)
			if err != nil {
				break
			}
			steps = append(steps, int64(i))
		}
		tail = tail[len(steps):]
	} else {
		steps = []int64{}

		if a, err := strconv.Atoi(tail[0]); err == nil {
			d1 = int64(a)
			tail = tail[1:]
			if b, err := strconv.Atoi(tail[0]); err == nil {
				d0 = d1
				d1 = int64(b)
				tail = tail[1:]
			}
		}

		if math.Signbit(float64(dStep)) != math.Signbit(float64(d1-d0)) {
			fmt.Println("Wrong sign on step")
			return []string{}
		}

		for d := d0; d <= d1; d += dStep {
			steps = append(steps, d)
		}
	}

	data = processing.VaryDimension(n, seed, stick, tryLoad, runsPer, steps)
	return flagSet.Args()
}
// runs several sets of simulations with the specified range of
func varySticking(args []string) []string {
	var flagSet = flag.NewFlagSet("varysticking", flag.ContinueOnError)
	var (
		n, d, seed, runsPer int64
		s0, s1 float64 = 0.1, 1
		sStep float64
		tryLoad bool
	)
	flagSet.Int64Var(&n, "npoints", 5000, "The number of points desired in the final aggregate")
	flagSet.Int64Var(&d, "dimensions", 2, "The number of dimensions for the space of the simulation")
	flagSet.Int64Var(&seed, "seed", 1, "The seed to be input to the random number generator. Generates seeds for" +
		" each individual run so each change generates a complete new set.")
	flagSet.Int64Var(&runsPer, "runs", 32, "The number of simulations to be run for each value of d")

	flagSet.Float64Var(&sStep, "step", .1, "The gap between simulations (0 to read a list of values from tail) setting upper and lower bounds specified in the tail")

	flagSet.BoolVar(&tryLoad, "load", false, "Load the result if it has been found before for the specific config")

	var err = flagSet.Parse(args)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var steps []float64
	var tail = flagSet.Args()
	if sStep <= 0 {
		steps = []float64{}
		for _, arg := range tail {
			f, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				break
			}
			steps = append(steps, f)
		}
		tail = tail[len(steps):]
	} else {
		steps = []float64{}

		if a, err := strconv.ParseFloat(tail[0], 64); err == nil {
			s1 = a
			tail = tail[1:]
			if b, err := strconv.ParseFloat(tail[0], 64); err == nil {
				s0 = s1
				s1 = b
				tail = tail[1:]
			}
		}

		if math.Signbit(sStep) != math.Signbit(s1-s0) {
			fmt.Println("Wrong sign on step")
			return []string{}
		}

		for s := s0; s <= s1; s += sStep {
			steps = append(steps, s)
		}
	}

	data = processing.VarySticking(n, d, seed, tryLoad, runsPer, steps)
	return tail
}

// finds the hausdorff dimensions using the growth rate
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

// draws the loaded agreagtes as svg images
func draw(args []string) []string {
	var flagSet = flag.NewFlagSet("draw", flag.ContinueOnError)


	var err = flagSet.Parse(args)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	if data != nil {
		processing.Draw(data)
	}
	return flagSet.Args()
}

// stores the set of valid instructions
var (
	handles = map[string]func([]string) []string{
		"run": runMany,
		"save": save,
		"varydimension": varyDimension,
		"varysticking": varySticking,
		"growth": growth,
		"draw": draw,
	}
)

// handles which instruction is being called
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
			fmt.Println("Use -h flag for specific help")
			return true
		}
	}
	return true
}

func handleInstructions(instructions string) bool {
	var args []string = util.StringToArgs(instructions)
	return handleArgs(args)
}

// waits for input then passes any input off to a handler
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
