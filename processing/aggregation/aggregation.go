//This is an auto generated file from genAggFiles.py
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

	default: return RunND(nPoints, seed, sticking, nDimension)
	}
}
