package svg

import (
	"fmt"
	"math"

	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation"
)

type HSVA struct {
	H, S, V, A float64
}

func (hsva HSVA) RGBA() (r, g, b, a uint32) {
	const (
		max32 = float64(math.MaxUint32)
	)
	var (
		H, S, V = hsva.H, hsva.S, hsva.V
		zone    = math.Floor(H * 6)
		part    = H*6 - zone
		low     = V * (1 - S)
		midA    = V * (1 - part*S)
		midB    = V * (1 - (1-part)*S)
	)

	switch uint8(zone) % 6 {
	case 0:
		r = uint32(max32 * V)
		g = uint32(max32 * midB)
		b = uint32(max32 * low)
	case 1:
		r = uint32(max32 * midA)
		g = uint32(max32 * V)
		b = uint32(max32 * low)
	case 2:
		r = uint32(max32 * low)
		g = uint32(max32 * V)
		b = uint32(max32 * midB)
	case 3:
		r = uint32(max32 * low)
		g = uint32(max32 * midA)
		b = uint32(max32 * V)
	case 4:
		r = uint32(max32 * midB)
		g = uint32(max32 * low)
		b = uint32(max32 * V)
	case 5:
		r = uint32(max32 * V)
		g = uint32(max32 * low)
		b = uint32(max32 * midA)
	}

	a = uint32(max32 * hsva.A)
	return r, g, b, a
}

// Draws the aggregate to an svg file. points are coloured according to when they were added to the set. Hue rotates
// from 0 to 300/360 over the range for clear contrast
func DrawAggregate(points []aggregation.Point) string {

	var (
		strOut  string = "<?xml version='1.0'?>\n"
		strBody string = ""
		hsv     HSVA   = HSVA{H: 0, S: 1, V: 0.8, A: 1}
		N       int    = len(points)

		minX, minY, maxX, maxY int64
	)

	// iterates over the set of points and adds them to the file at their known locations
	for i, point := range points {
		// rotate the hue
		hsv.H = float64(i*300) / float64(N*360)

		var (
			coords           []int64 = point.Coordinates()
			x                int64   = coords[0]
			y                int64   = coords[1]
			r32, g32, b32, _ uint32  = hsv.RGBA()
			r                        = uint8(float64(math.MaxUint8) * float64(r32) / float64(math.MaxUint32))
			g                        = uint8(float64(math.MaxUint8) * float64(g32) / float64(math.MaxUint32))
			b                        = uint8(float64(math.MaxUint8) * float64(b32) / float64(math.MaxUint32))
		)

		// update the min/max coords
		if x < minX {
			minX = x
		} else if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		} else if y > maxY {
			maxY = y
		}

		// concatenate the new line to the end of the body
		strBody += fmt.Sprintf(
			"<circle cx='%f' cy='%f' r='%f' fill='rgb(%d,%d,%d)' />\n",
			float64(x)+0.5, float64(y)+0.5, 0.5, r, g, b)
	}

	var (
		wMax int64 = 1600
		hMax int64 = 900
		dX         = maxX - minX
		dY         = maxY - minY

		scale = math.Min(float64(wMax)/float64(dX+1), float64(hMax)/float64(dY+1))
		width, height = scale*float64(dX+1), scale*float64(dY*1)
	)
	// write all info required for the svg file.
	strOut += fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' version='1.1' width='%f' height='%f'>\n",
		width, height)
	strOut += fmt.Sprintf("<g transform='scale(%f) translate(%d,%d)'>\n", scale, -minX, -minY)
	strOut += strBody
	strOut += "</g>\n"
	strOut += "</svg>\n"
	return strOut
}
