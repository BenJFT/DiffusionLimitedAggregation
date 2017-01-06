package xyz

import (
	"fmt"

	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation"
)

func DrawAggregate(points []aggregation.Point) string {
	var strOut, strFrame string
	for _, point := range points {
		var x, y, z float64
		coords := point.Coordinates()
		x = float64(coords[0])
		y = float64(coords[1])
		if len(coords) > 2 {
			z = float64(coords[2])
		}
		strFrame += fmt.Sprintf("  C\t%f\t%f\t%f\n", x, y, z)

	}
	N := len(points)
	strOut += fmt.Sprintf("%d\nAggregate after %d points\n%s", N, N, strFrame)
	return strOut
}
