package tools

import (
	"strings"
	"bufio"
	"math"
	"os"

	"github.com/gonum/plot/plotter"
)

type Point interface {
	XY() (int64, int64)
}

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

func LR_2(points plotter.XYs) (m, em, c, ec float64) {
	var xx, yy, xy float64
	var x, y, n float64
	n = float64(len(points))
	for _, point := range points {
		x += point.X/n
		y += point.Y/n
	}
	for _, point := range points {
		dx := point.X - x
		dy := point.Y - y
		xx += math.Pow(dx, 2)
		yy += math.Pow(dy, 2)
		xy += dx*dy
	}

	m = xy/xx
	c = y - m*x

	s := math.Sqrt((yy-m*xy)/(n-2))

	ec = math.Sqrt(1/n + x*x/xx)
	em = math.Sqrt(s/math.Sqrt(xx))
	return
}