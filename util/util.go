package util

import (
	"bufio"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"math"
	"os"
	"strings"
)

// splits a string at the spaces. Any double spacing is treated as single space
func StringToArgs(str string) (args []string) {
	args = make([]string, 0)

	for _, s := range strings.Split(str, " ") {
		if len(s) > 0 {
			args = append(args, s)
		}
	}

	return
}

var scanner = bufio.NewScanner(os.Stdin)
// reads the line. If nothing was written returns an empty string
func ReadStrOrEmpty() string {
	b := scanner.Scan()
	if b {
		return scanner.Text()
	} else {
		return ""
	}
}

// flips a 2D array on it's diagonal. Does not need to be square.
func Transpose(data [][]float64) [][]float64 {
	out := make([][]float64, len(data[0]))
	for i := range out {
		out[i] = make([]float64, len(data))
	}

	for i, row := range data {
		for j, y := range row {
			out[j][i] = y
		}
	}
	return out
}

// calculates the least squares regression as defined at http://mathworld.wolfram.com/LeastSquaresFitting.html
// returns the a & b for y = a + bx, as well as their respective standard error.
func LeastSquares(xys plotter.XYs) (a, b, ea, eb float64) {
	var (
		S2, SSxx, SSyy, SSxy, meanX, meanY float64
		n                                  float64 = float64(len(xys))
	)

	// calculating sum of squares
	for _, xy := range xys {
		SSxx += xy.X * xy.X
		SSyy += xy.Y * xy.Y
		SSxy += xy.X * xy.Y
		meanX += xy.X
		meanY += xy.Y
	}

	meanX /= n
	meanY /= n

	SSxx -= meanX * meanX * n
	SSyy -= meanY * meanY * n
	SSxy -= meanX * meanY * n

	// use sum of squares to find a & b for y = a + bx
	b = SSxy / SSxx
	a = meanY - b*meanX

	// find the standard error on each of these coefficients
	S2 = (SSyy - b*SSxy) / (n - 2)

	ea = math.Sqrt(S2 * meanX * meanX / (SSxx * n))
	eb = math.Sqrt(S2 / SSxx)

	return a, b, ea, eb
}
// according to https://www.che.udel.edu/pdf/FittingData.pdf
func WeightedLeastSquares(errs plotutil.ErrorPoints) (a, b, ea, eb float64) {
	var (
		sumX, sumY, sumXY, sumXX, sumE float64
		chX, chY, chXY, chXX, chE chan float64
	)

	chX = make(chan float64)
	go func() {
		var sum, x, e float64

		for i, xy := range errs.XYs {
			x = xy.X
			e = errs.YErrors[i].High
			sum += x/(e*e)
		}
		chX <- sum
	}()
	chY = make(chan float64)
	go func() {
		var sum, y, e float64

		for i, xy := range errs.XYs {
			y = xy.Y
			e = errs.YErrors[i].High
			sum += y/(e*e)
		}
		chY <- sum
	}()
	chXX = make(chan float64)
	go func() {
		var sum, x, e float64

		for i, xy := range errs.XYs {
			x= xy.X
			e = errs.YErrors[i].High
			sum += (x*x)/(e*e)
		}
		chXX <- sum
	}()
	chXY = make(chan float64)
	go func() {
		var sum, x, y, e float64

		for i, xy := range errs.XYs {
			x, y = xy.X, xy.Y
			e = errs.YErrors[i].High
			sum += (x*y)/(e*e)
		}
		chXY <- sum
	}()
	chE = make(chan float64)
	go func() {
		var sum, e float64

		for i := range errs.XYs {
			e = errs.YErrors[i].High
			sum += 1/(e*e)
		}
		chE <- sum
	}()

	sumX = <-chX
	sumY = <-chY
	sumXX = <-chXX
	sumXY = <-chXY
	sumE = <-chE

	b = ((sumX*sumY) - (sumXY*sumE))/((sumX*sumX) - (sumXX*sumE))
	a = (sumXY - (b*sumXX))/sumX
	eb = math.Sqrt( sumE/((sumXX*sumE) - (sumX*sumX)) )
	ea = math.Sqrt( sumXX/((sumXX*sumE) - (sumX*sumX)) )
	return a, b, ea, eb
}