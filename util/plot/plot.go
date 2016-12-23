package plot

import (
	"fmt"
	//"math"
)

type ShapeCircle struct {}
func (s ShapeCircle) Draw(x, y, w float64) string {
	return fmt.Sprintf("<circle cx='%f' cy='%f' cr='%f'/>\n", x, y, w/2)
}

type ShapeSquare struct {}
func (s ShapeSquare) Draw(x, y, w float64) string {
	return fmt.Sprintf("<rect x='%f' y='%f' width='%f' height='%f'/>\n", x - w/2, y - w/2, w, w)
}

type ShapePlus struct {}
func (s ShapePlus) Draw(x, y, w float64) string {
	w /= 2
	return fmt.Sprintf("<line x1='%f' y1='%f' x2='%f' y2='%f'/>\n", x - w, y, x + w, y) +
		fmt.Sprintf("<line x1='%f' y1='%f' x2='%f' y2='%f'/>\n", x, y - w, x, y + w)
}

type ShapeCross struct {}
func (s ShapeCross) Draw(x, y, w float64) string {
	w /= 2
	return fmt.Sprintf("<line x1='%f' y1='%f' x2='%f' y2='%f'/>\n", x - w, y - w, x + w, y + w) +
		fmt.Sprintf("<line x1='%f' y1='%f' x2='%f' y2='%f'/>\n", x - w, y + w, x + w, y - w)
}

type ShapeTri struct {}
func (s ShapeTri) Draw(x, y, w float64) string {
	w /= 2
	return fmt.Sprintf("<polygon points='%f,%f %f,%f %f,%f'/>\n", x, y+w, x + w*0.866, y - w/2, x - w*0.866, y - w/2)
}

type Shape interface {
	Draw(x, y, size float64) string
}

type XYs []struct{X, Y float64}

type Scatter struct {
	XY XYs
	Style string
	VertShape Shape
	Label string
}
func (s Scatter) Draw() (out string, minX, minY, maxX, maxY float64){
	out += fmt.Sprintf("<g style='%s'>\n", s.Style)
	for _, xy := range s.XY {
		out += s.VertShape.Draw(xy.X, xy.Y, 1)
		if xy.X > maxX {
			maxX = xy.X
		} else if xy.X < minX {
			minX = xy.X
		}
		if xy.Y > maxY {
			maxY = xy.Y
		} else if xy.Y < minY {
			minY = xy.Y
		}
	}
	out += "</g>"
	return
}

func NewScatter(x, y []float64) Scatter {
	scat := Scatter { XY: make(XYs, len(x)), Style:"fill:#000000; stroke:#000000;", VertShape:ShapeCross{} }
	for i := range x {
		scat.XY[i].X = x[i]
		scat.XY[i].Y = y[i]
	}
	return scat
}

type Plotter interface {
	Draw() (out string, minX, minY, maxX, maxY float64)
}

func Plot(plts ...Plotter) (out string) {
	var (
		width, height float64 = 200, 200
		minX, minY, maxX, maxY float64
		pltStrs []string = make([]string, len(plts))
	)

	for i, plt := range plts {
		o, mx, my, Mx, My := plt.Draw()
		pltStrs[i] = o
		if mx < minX {
			minX = mx
		}
		if my < minY {
			minY = my
		}
		if Mx > maxX {
			maxX = Mx
		}
		if My > maxY {
			maxY = My
		}
	}

	var (
		wx float64 = maxX - minX
		wy float64 = maxY - minY
	)
	out = "<?xml version='1.0' encoding='UTF-8'?>\n"
	out += fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' version='1.1' width='%fpx' height='%fpx'>\n",
		width, height)
	out += fmt.Sprintf("<g transform=' scale(%f,-%f) translate(%fpx, %fpx)'>\n",
		width/wx, height/wy, -minX, -minY)

	for _, str := range pltStrs {
		out += str
	}

	out += "</g>\n"
	out += "</svg>\n"

	return out
}