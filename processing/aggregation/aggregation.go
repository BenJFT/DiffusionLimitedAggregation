package aggregation

import (
	"encoding/gob"
	"math/rand"

	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg2D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg3D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/aggND"
)

func init() {
	gob.Register(agg2D.Point2D{})
	gob.Register(agg3D.Point3D{})
}

type Point interface {
	Coordinates() []int64
	SquareDistance(coords []float64) float64
}

func Run2D(nPoints, seed int64, sticking float64) (points []Point) {

	var rng *rand.Rand = rand.New(rand.NewSource(seed))
	var state = agg2D.RunNew(nPoints, sticking, rng)

	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run3D(nPoints, seed int64, sticking float64) (points []Point) {

	var rng *rand.Rand = rand.New(rand.NewSource(seed))
	var state = agg3D.RunNew(nPoints, sticking, rng)

	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func RunND(nPoints, seed int64, sticking float64, dimension int64) (points []Point) {

	var rng *rand.Rand = rand.New(rand.NewSource(seed))
	var state = aggND.RunNew(nPoints, sticking, rng, dimension)

	points = make([]Point, nPoints)

	for i, p := range state {
		points[i] = p
	}

	return points
}

func RunNew(nPoints, seed, nDimension int64, sticking float64) []Point {
	switch nDimension {
	case 2:
		return Run2D(nPoints, seed, sticking)
	case 3:
		return Run3D(nPoints, seed, sticking)
	default:
		return RunND(nPoints, seed, sticking, nDimension)
	}
}
