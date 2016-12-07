package aggregation

import (
	"math"
	"math/rand"
)

const (
	BORDER_SCALE float64 = 2
)

type Point struct {
	x, y int64
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
	agg.pointRadius = math.Sqrt(float64(point.x*point.x+point.y*point.y) + 0.5)
}

func (point *Point) moveToBorder(agg *Aggregator, rng *rand.Rand) {
	agg.angle = 2 * math.Pi * rng.Float64()
	point.x = int64(math.Sin(agg.angle) * agg.borderRadius)
	point.y = int64(math.Cos(agg.angle) * agg.borderRadius)
	point.updateRadius(agg)
}

func (point *Point) walk(agg *Aggregator, rng *rand.Rand) {
	switch rng.Intn(4) {
	case 0:
		point.x++
		if point.x > agg.intBorder {
			point.x -= 2 * agg.intBorder
		}
	case 1:
		point.x--
		if point.x < -agg.intBorder {
			point.x += 2 * agg.intBorder
		}
	case 2:
		point.y++
		if point.y > agg.intBorder {
			point.y -= 2 * agg.intBorder
		}
	case 3:
		point.y--
		if point.y < -agg.intBorder {
			point.y += 2 * agg.intBorder
		}
	}
	point.updateRadius(agg)
}

func (point *Point) upIn(state map[Point]int64, agg *Aggregator) bool {
	agg.tempPoint = *point
	agg.tempPoint.y++
	_, ok := state[agg.tempPoint]
	return ok
}

func (point *Point) downIn(state map[Point]int64, agg *Aggregator) bool {
	agg.tempPoint = *point
	agg.tempPoint.y--
	_, ok := state[agg.tempPoint]
	return ok
}

func (point *Point) rightIn(state map[Point]int64, agg *Aggregator) bool {
	agg.tempPoint = *point
	agg.tempPoint.y++
	_, ok := state[agg.tempPoint]
	return ok
}

func (point *Point) leftIn(state map[Point]int64, agg *Aggregator) bool {
	agg.tempPoint = *point
	agg.tempPoint.y--
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

func (agg *Aggregator) aggregate(n int64, rng *rand.Rand) map[Point]int64 {
	state := make(map[Point]int64, n) //prealocate memory

	point := &agg.currPoint //reference point in cache to avoid gc

	state[*point] = 0

	for i := int64(0); i < n; i++ {
		point.moveToBorder(agg, rng)
		for !point.hasNeighborIn(state, agg) {
			point.walk(agg, rng)
		}
		if agg.pointRadius > agg.radius {
			agg.setRadius()
		}
		state[*point] = i
	}
	return state
}
