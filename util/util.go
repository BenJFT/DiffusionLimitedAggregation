package util

import (
	"bufio"
	"github.com/gonum/plot/plotter"
	"math"
	"os"
	"strings"
)

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

func ReadStrOrEmpty() string {
	b := scanner.Scan()
	if b {
		return scanner.Text()
	} else {
		return ""
	}
}

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
