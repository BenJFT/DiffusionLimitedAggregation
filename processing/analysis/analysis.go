package analysis

import (
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation"
	"math"
)

type elemRadius struct {
	N      int64
	radius float64
}

// calculates the radius of gyration for a set of points assuming the provided average position
// Calculation is just RMS distance from mean
func gyrationRadius(points []aggregation.Point, mean []float64, chanel chan elemRadius) {
	var N int64 = int64(len(points))

	var e elemRadius = elemRadius{N: N}
	for _, p := range points {
		e.radius += p.SquareDistance(mean) / float64(N)
	}

	e.radius = math.Sqrt(e.radius)

	chanel <- e
}

type elemMean struct {
	N    int64
	mean []float64
}

// cumulatively calculates the mean position for the set of points
func findMeans(points []aggregation.Point, means chan elemMean) {

	// running sum of coordinates in each axis
	var sumPos []float64
	for i, p := range points {
		coords := p.Coordinates()
		if sumPos == nil {
			sumPos = make([]float64, len(coords))
		}
		mean := make([]float64, len(coords))
		// add the coordinates to the sum and divide by the number of points already included to find the mean
		for j, x := range coords {
			sumPos[j] += float64(x)
			mean[j] = sumPos[j] / float64(i+1)
		}
		// pass the mean through the 'means' channel
		means <- elemMean{N: int64(i + 1), mean: mean}
	}
}

//calculates the cumulative radius of gyration for the set of points
func GyrationRadii(points []aggregation.Point) (radii []float64) {
	// allocate the empty array to contain the final radii
	radii = make([]float64, len(points))
	var (
		chRadii chan elemRadius = make(chan elemRadius)
		chMeans chan elemMean   = make(chan elemMean)
	)

	// start a goroutine calculating the means after each point is considered
	go findMeans(points, chMeans)
	go func() {
		for range points {
			// wait for the next mean to be calculated then find the gyration radius to accompany it
			e := <-chMeans
			go gyrationRadius(points[:e.N], e.mean, chRadii)
		}
	} ()

	for range points {
		// wait for each radius to be calculated and add them the set
		e := <-chRadii
		radii[e.N-1] = e.radius
	}

	return radii
}
