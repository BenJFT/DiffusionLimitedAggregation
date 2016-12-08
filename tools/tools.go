package tools

import (
	"strings"
	"bufio"
	"os"

	"github.com/gonum/plot/plotter"
)

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)

func SingleSpace(str string) string  {

	for tmpStr := strings.Replace(str, "  ", " ", -1); tmpStr != str; tmpStr = strings.Replace(str, "  ", " ", -1) {
		str = tmpStr
	}

	return str
}

func ReadStrOrEmpty() string {
	b := scanner.Scan()
	if b {
		return scanner.Text()
	} else {
		return ""
	}
}

func LeastSquares(points plotter.XYs) (float64, float64) {
	n := float64(len(points))

	var sumX, sumY, sumXY, sumXX float64

	for _, xy := range points {
		sumX += xy.X
		sumY += xy.Y
		sumXX += xy.X*xy.X
		sumXY += xy.X*xy.Y
	}

	base := (n*sumXX - sumX*sumX)
	a := (n * sumXY - sumX * sumY)/base
	b := (sumXX * sumY - sumXY * sumX)/base

	return a, b
}