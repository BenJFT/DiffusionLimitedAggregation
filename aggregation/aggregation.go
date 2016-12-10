package aggregation

import (
	"math"
	"math/rand"
	"github.com/Benjft/DiffusionLimitedAggregation/tools"
)



const (
	BORDER_SCALE float64 = 1.5
)

type Point struct {
	X, Y int64
}

func (point Point) XY() (int64, int64) {
	return point.X, point.Y
}

type Aggregator struct {
	angle        float64
	radius       float64
	borderRadius float64
	pointRadius  float64

	intBorder int64

	currPoint Point
	tempPoint Point
}

func (point *Point) updateRadius(agg *Aggregator) {
	agg.pointRadius = math.Sqrt(float64(point.X *point.X +point.Y *point.Y)) + 0.5
}

func (point *Point) moveToBorder(agg *Aggregator, rng *rand.Rand) {
	agg.angle = 2 * math.Pi * rng.Float64()
	point.X = int64(math.Sin(agg.angle) * agg.borderRadius)
	point.Y = int64(math.Cos(agg.angle) * agg.borderRadius)
	point.updateRadius(agg)
}

func (point *Point) walk(state map[Point]int64, agg *Aggregator, rng *rand.Rand) {
	switch rng.Intn(4) {
	case 0:
		point.X++
		if _, ok := state[*point]; ok {
			point.X--
		}
		if point.X > agg.intBorder {
			point.X -= 2 * agg.intBorder
		}

	case 1:
		point.X--
		if _, ok := state[*point]; ok {
			point.X++
		}
		if point.X < -agg.intBorder {
			point.X += 2 * agg.intBorder
		}
	case 2:
		point.Y++
		if _, ok := state[*point]; ok {
			point.Y--
		}
		if point.Y > agg.intBorder {
			point.Y -= 2 * agg.intBorder
		}
	case 3:
		point.Y--
		if _, ok := state[*point]; ok {
			point.Y++
		}
		if point.Y < -agg.intBorder {
			point.Y += 2 * agg.intBorder
		}
	}
	point.updateRadius(agg)
}

func (point *Point) upIn(state map[Point]int64, agg *Aggregator) bool {
	agg.tempPoint = *point
	agg.tempPoint.Y++
	_, ok := state[agg.tempPoint]
	return ok
}

func (point *Point) downIn(state map[Point]int64, agg *Aggregator) bool {
	agg.tempPoint = *point
	agg.tempPoint.Y--
	_, ok := state[agg.tempPoint]
	return ok
}

func (point *Point) rightIn(state map[Point]int64, agg *Aggregator) bool {
	agg.tempPoint = *point
	agg.tempPoint.X++
	_, ok := state[agg.tempPoint]
	return ok
}

func (point *Point) leftIn(state map[Point]int64, agg *Aggregator) bool {
	agg.tempPoint = *point
	agg.tempPoint.X--
	_, ok := state[agg.tempPoint]
	return ok
}

func (point *Point) hasNeighborIn(state map[Point]int64, agg *Aggregator) bool {
	return agg.pointRadius <= agg.radius+1 && (point.upIn(state, agg) || point.downIn(state, agg) ||
		point.leftIn(state, agg) || point.rightIn(state, agg))
}

func (agg *Aggregator) setRadius() {
	agg.radius = agg.pointRadius
	agg.borderRadius = agg.pointRadius*BORDER_SCALE + 1
	agg.intBorder = int64(agg.borderRadius)
}

func (agg *Aggregator) Aggregate(n int64, sticking float64, rng *rand.Rand) map[tools.Point]int64 {
	state := make(map[Point]int64, n) //prealocate memory

	point := &agg.currPoint //reference point in cache to avoid gc

	state[*point] = 0
	point.updateRadius(agg)
	agg.setRadius()

	for i := int64(1); i < n; i++ {
		point.moveToBorder(agg, rng)
		for !point.hasNeighborIn(state, agg) || sticking < rng.Float64() {
			point.walk(state, agg, rng)
		}
		if agg.pointRadius > agg.radius {
			agg.setRadius()
		}
		state[*point] = i
	}

	ret := make(map[tools.Point]int64, len(state))
	for k, v := range state {
		ret[k] = v
	}
	return ret
}
