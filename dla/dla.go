package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	agg "github.com/Benjft/DiffusionLimitedAggregation/aggregation"
)

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
var lastStates []map[agg.Point]int64

func readInt(a *int64) (err error) {
	_, err = fmt.Scanf("%d\n", a)
	return
}

func readStr() string {
	b := scanner.Scan()
	if b {
		return scanner.Text()
	} else {
		return ""
	}
}

// func getInput() (seed, n, runs int64) {
// 	fmt.Print("Seed = ")
// 	for readInt(&seed) != nil || seed < 1 {
// 		fmt.Print("Seed must be an int larger than 0")
// 	}
// 	fmt.Print("n = ")
// 	for readInt(&n) != nil || n < 2 {
// 		fmt.Print("n must be an int larger than 1")
// 	}
// 	fmt.Print("runs = ")
// 	for readInt(&runs) != nil || n < 1 {
// 		fmt.Print("runs must be an int larger than 0")
// 	}
// 	return
// }

func runAggregation(seed, n, runs int64) {
	rand.Seed(seed)
	// TODO
}

func processRun(args []string) {
	fmt.Println(args[1:])
}

func handle(str string) {
	for str2:=strings.Replace(str, "  ", " ", -1); str2 != str;

}

func mainLoop() {
	fmt.Print(">>")
	for cmd := strings.ToLower(readStr()); cmd != "quit"; cmd = strings.ToLower(readStr()) {
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
