package processing

import (
	"fmt"
	"math"
	"math/rand"

	agg "github.com/Benjft/DiffusionLimitedAggregation/aggregation"
	"github.com/Benjft/DiffusionLimitedAggregation/tools"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"

	"github.com/skratchdot/open-golang/open"
)

func Run(seed, n, runs int64, sticking float64) [][]agg.Point {
	rand.Seed(seed)

	chans := make([]chan map[agg.Point]int64, runs)
	for i := int64(0); i < runs; i++ {
		a := &agg.Aggregator{}
		c := make(chan map[agg.Point]int64)
		go func() {
			rng := rand.New(rand.NewSource(rand.Int63()))
			c <- a.Aggregate(n, sticking, rng)
		} ()
		chans[i] = c
	}

	ret := make([][]agg.Point, runs)
	for i, c := range chans {
		arr := make([]agg.Point, n)
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

func Draw(state []agg.Point, title, format string, display bool) {
	// Convert aggregation points to plot points
	xys := make(plotter.XYs, len(state))
	for i, point := range state {
		xy := &xys[i]
		xy.X = float64(point.X)
		xy.Y = float64(point.Y)
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

func Dimension(states [][]agg.Point, title, format string, display bool) (float64, float64) {
	
	var r2Max int64 = 0
	for _, state := range states {
		for _, point := range state {
			if r2 := point.X*point.X + point.Y*point.Y; r2 > r2Max {
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
		//xys = append(xys, xyData...)
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
	plt.Save(vg.Inch*10, vg.Inch*10, "dims.svg")
	open.Run("dims.svg")

	return tools.LeastSquares(xys)
	//return 1, 2
}

func dimension(state []agg.Point, r int64) plotter.XYs {
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
			x := (r + point.X) / w
			y := (r + point.Y) / w
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