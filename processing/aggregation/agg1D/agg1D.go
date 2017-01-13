package agg1D

import (
	"encoding/gob"
	"math/rand"
)

func init() {
	// register the point structure to allow it to be saved and loaded
	gob.Register(Point1D{0})
}

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 5
)

// Specialised point structure for a point in 1D. Implements methods required for the Point interface in aggregation
type Point1D struct{
	X int64
}

func (p Point1D) Coordinates() []int64 {
	return []int64{p.X}
}

func (p Point1D) SquareDistance(coords []float64) float64 {
	var (
		ix int64   = p.X
		fx float64 = float64(ix)
		x  float64 = coords[0]
		dx float64 = fx - x
	)
	return dx*dx
}

// Structure implemented to remove most, if not all, garbage collector overhead by preallocating all memory to be used
// by the simulation
type cache struct {
	point           Point1D
	pointRadius     float64

	rng             *rand.Rand
	lastWalk        int64

	state           map[Point1D]int64
	stateRadius     float64

	borderRadius    float64
	borderRadiusInt int64

	tempPoint       Point1D
	tempFloat       float64
}

// runs a new 2d aggregation simulation and returns the finished state
func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[Point1D]int64 {
	var state = make(map[Point1D]int64, nPoints)
	state[Point1D{0}] = 0

	// add points one at a time
	for i := int64(1); i < nPoints; i++ {
		state[Point1D{i}] = i
	}

	return state
}
