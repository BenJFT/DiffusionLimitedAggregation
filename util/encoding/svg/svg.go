package svg

import (
	"fmt"
	"math"

	"github.com/Benjft/DiffusionLimitedAggregation/util/types"
)

func DrawAggregate(points []types.Point) string {
	const (
		width int64 = 10
	)

	var (
		strOut string = "<?xml version='1.0' encoding='UTF-8'?>\n"
		strBody string = ""
		hsv types.HSVA = types.HSVA{H: 0, S: 1, V: 0.8, A: 1}
		N int = len(points)

		minX, minY, maxX, maxY int64
	)

	for i, point := range points {
		hsv.H = float64(i*300)/float64(N*360)

		var (
			coords []int64 = point.Coordinates()
			x int64 = coords[0]
			y int64 = coords[1]
			r32, g32, b32, _ uint32 = hsv.RGBA()
			r = uint8(float64(math.MaxUint8) * float64(r32) / float64(math.MaxUint32))
			g = uint8(float64(math.MaxUint8) * float64(g32) / float64(math.MaxUint32))
			b = uint8(float64(math.MaxUint8) * float64(b32) / float64(math.MaxUint32))
		)

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

		line := fmt.Sprintf(
			"<circle cx='%d' cy='%d' r='%d' fill='rgb(%d,%d,%d)' />\n",
			x*width + width/2, y*width + width/2, width/2, r, g, b)
		strBody += line
	}
	X := maxX - minX
	Y := maxY - minY
	strOut += fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' version='1.1' width='%d' height='%d'>\n",
		X*width+width, Y*width+width)
	strOut += fmt.Sprintf("<g transform='translate(%d,%d)'>\n", -minX*width, -minY*width)
	strOut += strBody
	strOut += "</g>\n"
	strOut += "</svg>\n"
	return strOut
}
