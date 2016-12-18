package aggregation

import (
	"fmt"
	"math/rand"

	"github.com/Benjft/DiffusionLimitedAggregation/util/types"
	"github.com/Benjft/DiffusionLimitedAggregation/aggregation/agg2D"
	"github.com/Benjft/DiffusionLimitedAggregation/aggregation/agg3D"
)

func Run2D(nPoints, seed int64, sticking float64) (points []types.Point) {

	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed))
	var state map[types.Point2D]int64
	state = agg2D.RunNew(nPoints, sticking, rng)

	points = make([]types.Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run3D(nPoints, seed int64, sticking float64) (points []types.Point) {

	var rng *rand.Rand
	rng = rand.New(rand.NewSource(seed))
	var state map[types.Point3D]int64
	state = agg3D.RunNew(nPoints, sticking, rng)

	points = make([]types.Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func RunNew(nPoints, seed, nDimension int64, sticking float64) []types.Point {
	switch nDimension {
	case 2:	return Run2D(nPoints, seed, sticking)
	case 3: return Run3D(nPoints, seed, sticking)
	default:
		fmt.Printf("%d dimensions not supported\n", nDimension)
		return nil
	}
}