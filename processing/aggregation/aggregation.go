package aggregation

import (
	"math/rand"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg1D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg2D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg3D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg4D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg5D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg6D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg7D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg8D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg9D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg10D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg11D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg12D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg13D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg14D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg15D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg16D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg17D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg18D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg19D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg20D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg21D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg22D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg23D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/agg24D"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation/aggND"
)

type Point interface {
	Coordinates() []int64
	SquareDistance([]float64) float64
}

func Run1D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg1D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run2D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg2D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run3D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg3D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run4D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg4D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run5D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg5D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run6D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg6D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run7D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg7D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run8D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg8D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run9D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg9D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run10D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg10D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run11D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg11D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run12D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg12D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run13D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg13D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run14D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg14D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run15D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg15D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run16D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg16D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run17D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg17D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run18D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg18D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run19D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg19D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run20D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg20D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run21D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg21D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run22D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg22D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run23D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg23D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func Run24D(nPoints, seed int64, sticking float64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = agg24D.RunNew(nPoints, sticking, rng)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for point, index := range state {
		points[index] = point
	}
	return points
}

func RunND(nPoints, seed int64, sticking float64, nDimension int64) (points []Point) {
	var rng = rand.New(rand.NewSource(seed))
	var state = aggND.RunNew(nPoints, sticking, rng, nDimension)
	if int64(len(state)) != nPoints {
		panic("N != nPoints. This should never happen.")
	}
	points = make([]Point, nPoints)
	for i, p := range state {
		points[i] = p
	}
	return points
}
func RunNew(nPoints, seed, nDimension int64, sticking float64) []Point {
	switch nDimension {
	case 1: return Run1D(nPoints, seed, sticking)
	case 2: return Run2D(nPoints, seed, sticking)
	case 3: return Run3D(nPoints, seed, sticking)
	case 4: return Run4D(nPoints, seed, sticking)
	case 5: return Run5D(nPoints, seed, sticking)
	case 6: return Run6D(nPoints, seed, sticking)
	case 7: return Run7D(nPoints, seed, sticking)
	case 8: return Run8D(nPoints, seed, sticking)
	case 9: return Run9D(nPoints, seed, sticking)
	case 10: return Run10D(nPoints, seed, sticking)
	case 11: return Run11D(nPoints, seed, sticking)
	case 12: return Run12D(nPoints, seed, sticking)
	case 13: return Run13D(nPoints, seed, sticking)
	case 14: return Run14D(nPoints, seed, sticking)
	case 15: return Run15D(nPoints, seed, sticking)
	case 16: return Run16D(nPoints, seed, sticking)
	case 17: return Run17D(nPoints, seed, sticking)
	case 18: return Run18D(nPoints, seed, sticking)
	case 19: return Run19D(nPoints, seed, sticking)
	case 20: return Run20D(nPoints, seed, sticking)
	case 21: return Run21D(nPoints, seed, sticking)
	case 22: return Run22D(nPoints, seed, sticking)
	case 23: return Run23D(nPoints, seed, sticking)
	case 24: return Run24D(nPoints, seed, sticking)

	default: return RunND(nPoints, seed, sticking, nDimension)
	}
}
