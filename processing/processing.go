package processing

import (
	"fmt"
	"math"
	"math/rand"
	"image/color"

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
		go func(rng *rand.Rand) {
			c <- genagg.RunNew(n, sticking, rng)
		} (rand.New(rand.NewSource(rand.Int63())))
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

	if display {
		open.Run(fileName)
	}
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
			y = (r + y) / w

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

func inCircle(p tools.Point, cx, cy, cr2 float64) bool {
	x, y := p.XY()
	dx := float64(x) - cx
	dy := float64(y) - cy

	r2 := dx*dx + dy*dy
	return r2 <= cr2
}

func allIn(state []tools.Point, cx, cy, cr2 float64) bool {
	for _, point := range state {
		if !inCircle(point, cx, cy, cr2) {
			return false
		}
	}
	return true
}

func cross(x0, y0, x1, y1, x2, y2 float64) float64 {
	return (x1 - x0)*(y2 - y0) - (y1 - y0)*(x2 - x0)
}

func makeDiameter(p0, p1 tools.Point) (cx, cy, cr2 float64) {
	x0, y0 := p0.XY()
	x1, y1 := p1.XY()

	cx = float64(x0 + x1)/2
	cy = float64(y0 + y1)/2

	dx := float64(x0) - cx
	dy := float64(y0) - cy

	cr2 = dx*dx + dy*dy

	return
}

func makeCircumCircle(p0, p1, p2 tools.Point) (cx, cy, cr2 float64) {
	x0, y0 := p0.XY()
	x1, y1 := p1.XY()
	x2, y2 := p2.XY()

	d := float64(x0*(y1 - y2) + x1*(y2 - y0) + x2*(y0 - y1))*2

	if d == 0 {
		return 0, 0, -1
	}

	cx = float64((x0*x0 + y0*y0)*(y1 - y2) + (x1*x1 + y1*y1)*(y2 - y0) + (x2*x2 + y2*y2)*(y0 - y1)) / d
	cy = float64((x0*x0 + y0*y0)*(x2 - x1) + (x1*x1 + y1*y1)*(x0 - x2) + (x2*x2 + y2*y2)*(x1 - x0)) / d
	dx := float64(x0) - cx
	dy := float64(y0) - cy
	cr2 = dx*dx + dy*dy

	return cx, cy, cr2
}

func makeTwoPoints(state []tools.Point, p0, p1 tools.Point) (cx, cy, cr2 float64) {
	cx, cy, cr2 = makeDiameter(p0, p1)

	if allIn(state, cx, cy, cr2) {
		return cx, cy, cr2
	}

	xa, ya := p0.XY()
	x0, y0 := float64(xa), float64(ya)
	xb, yb := p1.XY()
	x1, y1 := float64(xb), float64(yb)

	var lx, ly, lr2 float64
	lr2 = -1
	var rx, ry, rr2 float64
	rr2 = -1

	for _, p2 := range state {
		xc, yc := p2.XY()
		x2, y2 := float64(xc), float64(yc)
		crs := cross((x0), (y0), (x1), (y1), (x2), (y2))
		cx, cy, cr2 = makeCircumCircle(p0, p1, p2)

		if cr2 == -1 {
			continue
		} else if crs > 0 && (lr2 == -1 || cross(x0, y0, x1, y1, cx, cy) > cross(x0, y0, x1, y1, lx, ly)) {
			lx, ly, lr2 = cx, cy, cr2
		} else if crs < 0 && (rr2 == -1 || cross(x0, y0, x1, y1, cx, cy) < cross(x0, y0, x1, y1, rx, ry)) {
			rx, ry, rr2 = cx, cy, cr2
		}
	}

	if rr2 == -1 || (lr2 != -1 &&lr2 <= rr2) {
		return lx, ly, lr2
	} else {
		return rx, ry, rr2
	}
}

func makeOnePoint(state []tools.Point, p0 tools.Point) (cx, cy, cr2 float64) {
	x, y := p0.XY()
	cx, cy = float64(x), float64(y)
	for i, p1 := range state {
		if !inCircle(p1, cx, cy, cr2) {
			cx,cy, cr2 = makeTwoPoints(state[:i], p0, p1)
		}
	}
	return
}

func radius(state []tools.Point) (radii []float64) {
	var cx, cy, cr2 float64

	radii = make([]float64, len(state))

	for i, p0 := range state {
		if !inCircle(p0, cx, cy, cr2) {
			cx, cy, cr2 = makeOnePoint(state[:i], p0)
		}
		radii[i] = math.Sqrt(cr2)
	}

	return radii
}

func Radius(states [][]tools.Point) {
	runs := len(states)

	chans := make([]chan []float64, runs)
	for j, state := range states {
		c := make(chan []float64)
		go func (s []tools.Point) {
			c <- radius(s)
		} (state)
		chans[j] = c
	}

	xys := make(plotter.XYs, 0)
	for _, c := range chans {
		radii := <-c
		radiiXY := make(plotter.XYs, len(radii))
		for i, r := range radii {
			xy := &radiiXY[i]
			xy.X = float64(i+1)
			xy.Y += float64(r+.5)/float64(runs)
		}
		xys = append(xys, radiiXY...)
	}

	for i := range xys {
		xy := &xys[i]
		xy.X, xy.Y = math.Log(xy.X), math.Log(xy.Y)
	}
	m, c := tools.LeastSquares(xys)
	println(1/m, c)

	plt, _ := plot.New()
	s, _ := plotter.NewScatter(xys)
	l := plotter.NewFunction(func (x float64) float64 {return x*m + c})
	s.Color = color.RGBA{R: 255, A: 255}
	s.Shape = draw.CrossGlyph{}
	plt.Add(s, l)

	plt.Save(vg.Inch*10, vg.Inch*10, "tmp.svg")
	open.Run("tmp.svg")
}