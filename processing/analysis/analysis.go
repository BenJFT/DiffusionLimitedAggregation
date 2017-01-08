package analysis

import (
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation"
	"math"
)

type elemRadius struct {
	N int64
	radius float64
}
func gyrationRadius(points []aggregation.Point, mean []float64, chanel chan elemRadius) {
	var N int64 = int64(len(points))

	var e elemRadius = elemRadius{N:N}
	for _, p := range points {
		e.radius += p.SquareDistance(mean)/float64(N)
	}

	e.radius = math.Sqrt(e.radius)

	chanel <- e
}

type elemMean struct {
	N int64
	mean []float64
}
func findMeans(points []aggregation.Point, means chan elemMean) {
	var mean []float64
	for i, p := range points {
		coords := p.Coordinates()
		ret := make([]float64, len(coords))
		if mean == nil {
			mean = make([]float64, len(coords))
		}
		for j, x := range coords {
			mean[j] += float64(x)
			ret[j] = mean[j]/float64(i+1)
		}
		means <- elemMean{N:int64(i+1),mean:ret}
	}
}

func GyrationRadii(points []aggregation.Point) (radii []float64) {
	radii = make([]float64, len(points))

	var (
		radii chan elemRadius = make(chan elemRadius)
		means chan elemMean = make(chan elemMean)
	)

	go findMeans(points, means)
	for range points {
		e := <-means
		go gyrationRadius(points[:e.N], e.mean, radii)
	}

	for range points {
		e := <-radii
		radii[e.N-1] = e.radius
	}

	return radii
}