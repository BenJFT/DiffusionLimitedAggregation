package genagg

import (
	"math"
	"math/rand"

	"github.com/Benjft/DiffusionLimitedAggregation/tools"
)

const (
	BORDER_SCALE float64 = 1.5
)

// references a given lattice point
type Point struct {
	X, Y int64
}

func (point Point) XY() (int64, int64) {
	return point.X, point.Y
}

// caching structure. Exists as an optimisation by preventing memory reassignment, this helps reduce time spent in
// garbage collection, and helps keep variables in the cpu register or caches. Early implementations of this gave a
// a significant speedup
type cache struct {
	currPoint       Point
	currPointRadius float64

	rng             *rand.Rand
	lastWalk	int64

	state           map[Point]int64
	stateRadius     float64

	borderRadius    float64
	borderRadiusInt int64

	tempPoint       Point
	tempFloat	float64
}

// operations acting on the caching structure. These (mostly) avoid creating any memory within themselves
// by changing values within the cache variable instead of using returns. This has in benchmarking shown to give
// significant speed gains


func (c *cache) updateCurrPointRadius() {
	c.currPointRadius = math.Sqrt(float64(c.currPoint.X*c.currPoint.X + c.currPoint.Y*c.currPoint.Y))
}
func (c *cache) updateStateRadius() {
	c.stateRadius = c.currPointRadius
	c.borderRadius = c.stateRadius*BORDER_SCALE + 1
	c.borderRadiusInt = int64(c.borderRadius)
}

// returns true if the current point is in the caches state
func (c *cache) currPointIn() (ok bool) {
	_, ok = c.state[c.currPoint]
	return

}

// resets the location of the current point to some location on the border
func (c *cache) currPointToBorder() {
	c.tempFloat = 2*math.Pi*c.rng.Float64()
	c.currPoint.X = int64(math.Sin(c.tempFloat) * c.borderRadius)
	c.currPoint.Y = int64(math.Cos(c.tempFloat) * c.borderRadius)
	c.updateCurrPointRadius()
}

// moves the current point by one, applies periodic boundaries and will not move onto a site which is occupied
func (c *cache) walkPoint() {

	switch point := &c.currPoint; c.rng.Int63n(4) {
	case 0:
		point.X++
		if c.currPointRadius < 1+c.stateRadius && c.currPointIn() {
			point.X--
		} else {
			if point.X > c.borderRadiusInt {
				point.X -= 2*c.borderRadiusInt
			}
			c.lastWalk = 0
		}
	case 1:
		point.X--
		if c.currPointRadius < 1+c.stateRadius && c.currPointIn() {
			point.X++
		} else {
			if point.X < -c.borderRadiusInt {
				point.X += 2*c.borderRadiusInt
			}
			c.lastWalk = 1
		}
	case 2:
		point.Y++
		if c.currPointRadius < 1+c.stateRadius && c.currPointIn() {
			point.Y--
		} else {
			if point.Y > c.borderRadiusInt {
				point.Y -= 2*c.borderRadiusInt
			}
			c.lastWalk = 2
		}
	case 3:
		point.Y--
		if c.currPointRadius < 1+c.stateRadius && c.currPointIn() {
			point.Y++
		} else {
			if point.Y < -c.borderRadiusInt {
				point.Y += 2*c.borderRadiusInt
			}
			c.lastWalk = 3
		}
	}
	c.updateCurrPointRadius()
}

// tests if the adjacent cite in the +x direction is occupied
func (c *cache) isPxIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 1 && c.px()
}
func (c *cache) px() (ok bool) {
	c.tempPoint = c.currPoint
	c.tempPoint.X++
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the -x direction is occupied
func (c *cache) isMxIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 0 && c.mxIn()
}
func (c *cache) mxIn() (ok bool) {
	c.tempPoint = c.currPoint
	c.tempPoint.X--
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the +y direction is occupied
func (c *cache) isPyIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 3 && c.py()
}
func (c *cache) py() (ok bool) {
	c.tempPoint = c.currPoint
	c.tempPoint.Y++
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the -y direction is occupied
func (c *cache) isMyIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 2 && c.myIn()
}
func (c *cache) myIn() (ok bool) {
	c.tempPoint = c.currPoint
	c.tempPoint.Y--
	_, ok = c.state[c.tempPoint]
	return
}

// returns true if any adjacent site is occupied
func (c *cache) currPointHasNeighbor() bool {
	return c.currPointRadius <= c.stateRadius+1 && ( c.isPxIn() || c.isPyIn() || c.isMxIn() || c.isMyIn() )
}

func RunNew(n int64, pStick float64, rng *rand.Rand) (state map[tools.Point]int64) {

	c := cache{}
	c.rng = rng
	c.state = make(map[Point]int64, n)
	c.state[c.currPoint] = 0
	c.updateStateRadius()

	state = make(map[tools.Point]int64, n)
	state[c.currPoint] = 0

	for i := int64(1); i < n; i++ {
		c.currPointToBorder()
		for !c.currPointHasNeighbor() || pStick < rng.Float64() {
			c.walkPoint()
		}
		c.state[c.currPoint] = i
		state[c.currPoint] = i
		if c.currPointRadius > c.stateRadius {
			c.updateStateRadius()
		}
	}

	//state = make(map[tools.Point]int64, len(c.state))
	//for k, v := range c.state {
	//	state[k] = v
	//}

	return state
}

