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

func (point *Point) radius() float64 {
	return math.Sqrt(float64(point.X *point.X +point.Y *point.Y)) + 0.5
}

func (point *Point) moveToBorder(borderRadius float64, rng *rand.Rand) {
	angle := 2 * math.Pi * rng.Float64()
	point.X = int64(math.Sin(angle) * borderRadius)
	point.Y = int64(math.Cos(angle) * borderRadius)
}

func (point *Point) walk(state map[Point]int64, rng *rand.Rand, borderRadius float64) {
	intBorder := int64(borderRadius)

	switch rng.Intn(4) {
	case 0:
		point.X++
		if _, ok := state[*point]; ok {
			point.X--
		}
		if point.X > intBorder {
			point.X -= 2 * intBorder
		}

	case 1:
		point.X--
		if _, ok := state[*point]; ok {
			point.X++
		}
		if point.X < -intBorder {
			point.X += 2 * intBorder
		}
	case 2:
		point.Y++
		if _, ok := state[*point]; ok {
			point.Y--
		}
		if point.Y > intBorder {
			point.Y -= 2 * intBorder
		}
	case 3:
		point.Y--
		if _, ok := state[*point]; ok {
			point.Y++
		}
		if point.Y < -intBorder {
			point.Y += 2 * intBorder
		}
	}
}

func (point *Point) upIn(state map[Point]int64) bool {
	tempPoint := *point
	tempPoint.Y++
	_, ok := state[tempPoint]
	return ok
}

func (point *Point) downIn(state map[Point]int64) bool {
	tempPoint := *point
	tempPoint.Y--
	_, ok := state[tempPoint]
	return ok
}

func (point *Point) rightIn(state map[Point]int64) bool {
	tempPoint := *point
	tempPoint.X++
	_, ok := state[tempPoint]
	return ok
}

func (point *Point) leftIn(state map[Point]int64) bool {
	tempPoint := *point
	tempPoint.X--
	_, ok := state[tempPoint]
	return ok
}

func (point *Point) hasNeighborIn(state map[Point]int64, stateRadius float64) bool {
	return point.radius() <= stateRadius+1 && (point.upIn(state) || point.downIn(state) ||
		point.leftIn(state) || point.rightIn(state))
}

// runs an aggregation simulation using the passed args. Faster version exists in genagg.go using caching. This exists
// to show the performance gains from doing so
// for benchmarking comparison run
// "go test -bench=. benchtime=1m github.com/Benjft/DiffusionLimitedAggregation/processing" from the commandline with
// it installed in your GOPATH
func RunNew(n int64, sticking float64, rng *rand.Rand) map[tools.Point]int64 {
	rawState := make(map[Point]int64, n) //prealocate memory

	point := Point{0, 0}

	rawState[point] = 0
	stateRadius := point.radius()
	borderRadius := BORDER_SCALE * stateRadius

	retState := make(map[tools.Point]int64, n)
	retState[point] = 0

	for i := int64(1); i < n; i++ {
		point.moveToBorder(borderRadius, rng)
		for !point.hasNeighborIn(rawState, stateRadius) || sticking < rng.Float64() {
			point.walk(rawState, rng, borderRadius)
		}
		rawState[point] = i
		retState[point] = i
		if r := point.radius(); r > stateRadius {
			stateRadius = r
			borderRadius = stateRadius * BORDER_SCALE
		}
	}

	return retState
}
