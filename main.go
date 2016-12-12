package main

import (
	"os"
	"fmt"
	"flag"
	"strings"

	"github.com/Benjft/DiffusionLimitedAggregation/tools"
)

type Handler func (args []string) (tail []string)

func handleRun(args []string) (tail []string) {
	var flags flag.FlagSet
	//TODO
	return make([]string, 0)
}

var (
	handles map[string]Handler = Handler {"run": handleRun}
)

func handleArgs(args []string) (bool) {
	var head string
	var tail []string
	for len(args) > 0 {
		head = strings.ToLower(args[0])
		tail = args[1:]

		if head == "quit" || head == "stop" {
			return false
		}

		args = handles[head](tail)
	}
	return true
}

func handleInstructions(instructions string) (cont bool) {
	var args []string = tools.StringToArgs(instructions)
	return handleArgs(args)
}

func mainLoop() {
	var instructions string
	var err error
	_, err = fmt.Scanln(&instructions)
	for ; handleInstructions(instructions); _, err = fmt.Scanln(&instructions) {
		if  err != nil {
			panic(err)
		}
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
