package processing

import (
	"fmt"
	"math"
	"math/rand"

	agg "github.com/Benjft/DiffusionLimitedAggregation/aggregation"
	"github.com/Benjft/DiffusionLimitedAggregation/tools"
	"github.com/Benjft/DiffusionLimitedAggregation/genagg"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"

	"github.com/skratchdot/open-golang/open"
)

func RunNoCache(seed, n, runs int64, sticking float64) [][]tools.Point {
	rand.Seed(seed)

	chans := make([]chan map[tools.Point]int64, runs)
	for i := int64(0); i < runs; i++ {
		c := make(chan map[tools.Point]int64)
		go func() {
			rng := rand.New(rand.NewSource(rand.Int63()))
			c <- agg.RunNew(n, sticking, rng)
		} ()
		chans[i] = c
	}

	ret := make([][]tools.Point, runs)
	for i, c := range chans {
		arr := make([]tools.Point, n)
		m := <-c
		for k, v := range m {
			arr[v] = k
		}
		ret[i] = arr
		if i64 :=  int64(len(arr)); i64 < n {
			fmt.Println("MISSING", string(n-i64), "POINTS!")
		}
	}

	return ret
}

func Run(seed, n, runs int64, sticking float64) [][]tools.Point {
	rand.Seed(seed)

	chans := make([]chan map[tools.Point]int64, runs)
	for i := int64(0); i < runs; i++ {
		c := make(chan map[tools.Point]int64)
		go func() {
			rng := rand.New(rand.NewSource(rand.Int63()))
			c <- genagg.RunNew(n, sticking, rng)
		} ()
		chans[i] = c
	}

	ret := make([][]tools.Point, runs)
	for i, c := range chans {
		arr := make([]tools.Point, n)
		m := <-c
		for k, v := range m {
			arr[v] = k
		}
		ret[i] = arr
		if i64 :=  int64(len(arr)); i64 < n {
			fmt.Println("MISSING", string(n-i64), "POINTS!")
		}
	}

	return ret
}

func Draw(state []tools.Point, title, format string, display bool) {
	// Convert aggregation points to plot points
	xys := make(plotter.XYs, len(state))
	for i, point := range state {
		xy := &xys[i]
		x, y := point.XY()
		xy.X, xy.Y = float64(x), float64(y)
	}

	// Set up the plot
	plt, _ := plot.New()
	plt.Title.Text = title
	plt.X.Label.Text = "X Displacement"
	plt.Y.Label.Text = "Y Displacement"
	grid := plotter.NewGrid()
	plt.Add(grid)

	// Add data to the plot
	s, _ := plotter.NewScatter(xys)
	s.Shape = draw.PlusGlyph{}
	plt.Add(s)

	// Squaring the axis so point are evenly spaced
	min := math.Min(plt.X.Min, plt.Y.Min)
	max := math.Max(plt.X.Max, plt.Y.Max)

	plt.X.Min = min
	plt.Y.Min = min
	plt.X.Max = max
	plt.Y.Max = max

	// Save the plot
	fileName := fmt.Sprintf("%s.%s", title, format)
	D := vg.Length(max - min)
	plt.Save(vg.Millimeter*2*D, vg.Millimeter*2*D, fileName)

	open.Run(fileName)
}

func Dimension(states [][]tools.Point, title, format string, display bool) (float64, float64) {
	
	var r2Max int64 = 0
	for _, state := range states {
		for _, point := range state {
			x, y := point.XY()
			if r2 := x*x + y*y; r2 > r2Max {
				r2Max = r2
			}
		}
	}

	var r int64 = int64(math.Ceil(math.Sqrt(float64(r2Max))))
	chans := make([]chan plotter.XYs, len(states))
	for i, state := range states {
		c := make(chan plotter.XYs)
		go func() {
			c <- dimension(state, r)
		} ()
		chans[i] = c
	}

	var xys plotter.XYs
	for _, c := range chans {
		xyData := <- c
		for i := range xyData {
			xy := &xyData[i]
			xy.X = math.Log(xy.X)
			xy.Y = math.Log(xy.Y)
			if xy.X != math.Inf(+1) && xy.X != math.Inf(-1) && xy.Y != math.Inf(+1) && xy.Y != math.Inf(-1) {
				xys = append(xys, *xy)
			}
		}
	}

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}
	s, err := plotter.NewScatter(xys)
	if err != nil {
		panic(err)
	}

	plt.Add(s)
	plt.Title.Text = title
	plt.X.Label.Text = "Log10 box width"
	plt.Y.Label.Text = "log10 box count"

	fileName := fmt.Sprintf("%s.%s", title, format)

	plt.Save(vg.Inch*10, vg.Inch*10, fileName)
	open.Run(fileName)

	return tools.LeastSquares(xys)
}

func dimension(state []tools.Point, r int64) plotter.XYs {
	log2 := int64(math.Ceil(math.Log2(float64(r*2))))
	pow2 := int64(math.Pow(2, float64(log2)))

	xys := make(plotter.XYs, log2)
	for i := int64(1); i <= log2; i++ {
		n := int64(math.Pow(2, float64(i)))
		w := pow2 / n

		boxs := make([][]int64, n)
		for i := range boxs {
			boxs[i] = make([]int64, n)
		}

		for _, point := range state {
			x, y := point.XY()

			x = (r + x) / w
			y = (r + x) / w

			boxs[x][y] += 1
		}

		sum := 0
		for _, row := range boxs {
			for _, box := range row {
				if box < w*w && box != 0{
					sum += 1
				}
			}
		}

		xys[i-1].X = float64(n)
		xys[i-1].Y = float64(sum)
	}

	return xys
}

func Density(states [][]tools.Point, title, format string, display bool) (float64, float64) {

	var r2Max int64 = 0
	for _, state := range states {
		for _, point := range state {
			x, y := point.XY()
			if r2 := x*x + y*y; r2 > r2Max {
				r2Max = r2
			}
		}
	}

	var r int64 = int64(math.Ceil(math.Sqrt(float64(r2Max))))
	chans := make([]chan plotter.XYs, len(states))
	for i, state := range states {
		c := make(chan plotter.XYs)
		go func() {
			c <- dimension(state, r)
		} ()
		chans[i] = c
	}

	var xys plotter.XYs
	for _, c := range chans {
		xyData := <- c
		for i := range xyData {
			xy := &xyData[i]
			xy.X = math.Log(xy.X)
			xy.Y = math.Log(xy.Y)
			if xy.X != math.Inf(+1) && xy.X != math.Inf(-1) && xy.Y != math.Inf(+1) && xy.Y != math.Inf(-1) {
				xys = append(xys, *xy)
			}
		}
	}

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}
	s, err := plotter.NewScatter(xys)
	if err != nil {
		panic(err)
	}

	plt.Add(s)
	plt.Title.Text = title
	plt.X.Label.Text = "Log10 r (lattica const.)"
	plt.Y.Label.Text = "log10 C(r)"

	fileName := fmt.Sprintf("%s.%s", title, format)

	plt.Save(vg.Inch*10, vg.Inch*10, fileName)
	open.Run(fileName)

	return tools.LeastSquares(xys)
}

func density(state []tools.Point, r int64) plotter.XYs {

	R := make(map[int64]int64)

	for i, point1 := range state {
		for _, point2 := range state[:i] {
			x1, y1 := point1.XY()
			x2, y2 := point2.XY()
			dx := x1-x2
			dy := y1-y2

			d := math.Sqrt(float64(dx*dx+dy*dy))

			R[int64(math.Ceil(d))] += 1
		}
	}

	//TODO
	return nil
}