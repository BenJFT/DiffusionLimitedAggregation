package agg2D

import (
	"math"
	"math/rand"

	"github.com/Benjft/DiffusionLimitedAggregation/util/types"
)

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 5
)

// caching structure. Exists as an optimisation by preventing memory reassignment, this helps reduce time spent in
// garbage collection, and helps keep variables in the cpu register or caches. Early implementations of this gave a
// a significant speedup
type cache struct {
	point           types.Point2D
	pointRadius     float64

	rng             *rand.Rand
	lastWalk        int64

	state           map[types.Point2D]int64
	stateRadius     float64

	borderRadius    float64
	borderRadiusInt int64

	tempPoint       types.Point2D
	tempFloat       float64
}

// operations acting on the caching structure. These (mostly) avoid creating any memory within themselves
// by changing values within the cache variable instead of using returns. This has in benchmarking shown to give
// significant speed gains
func (c *cache) updateCurrPointRadius() {
	c.pointRadius = math.Sqrt(float64(c.point.X*c.point.X + c.point.Y*c.point.Y))
}

func (c *cache) updateStateRadius() {
	c.stateRadius = c.pointRadius
	c.borderRadius = c.stateRadius*BORDER_SCALE + BORDER_CONST
	c.borderRadiusInt = int64(c.borderRadius)
}

// returns true if the current point is in the caches state
func (c *cache) pointIn() (ok bool) {
	_, ok = c.state[c.point]
	return

}

// resets the location of the current point to some location on the border
func (c *cache) pointToBorder() {
	c.tempFloat = 2*math.Pi*c.rng.Float64()
	c.point.X = int64(math.Sin(c.tempFloat) * c.borderRadius)
	c.point.Y = int64(math.Cos(c.tempFloat) * c.borderRadius)
	c.updateCurrPointRadius()
}

// moves the current point by one, applies periodic boundaries and will not move onto a site which is occupied
func (c *cache) walkPoint() {

	switch point := &c.point; c.rng.Int63n(4) {
	case 0:
		point.X++
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			point.X--
		} else {
			if point.X > c.borderRadiusInt {
				point.X -= 2*c.borderRadiusInt
			}
			c.lastWalk = 0
		}
	case 1:
		point.X--
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			point.X++
		} else {
			if point.X < -c.borderRadiusInt {
				point.X += 2*c.borderRadiusInt
			}
			c.lastWalk = 1
		}
	case 2:
		point.Y++
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			point.Y--
		} else {
			if point.Y > c.borderRadiusInt {
				point.Y -= 2*c.borderRadiusInt
			}
			c.lastWalk = 2
		}
	case 3:
		point.Y--
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
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
func (c *cache) isLeftIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 1 && c.leftIn()
}
func (c *cache) leftIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.X++
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the -x direction is occupied
func (c *cache) isRightIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 0 && c.rightIn()
}
func (c *cache) rightIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.X--
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the +y direction is occupied
func (c *cache) isUpIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 3 && c.upIn()
}
func (c *cache) upIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.Y++
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the -y direction is occupied
func (c *cache) isDownIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 2 && c.downIn()
}
func (c *cache) downIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.Y--
	_, ok = c.state[c.tempPoint]
	return
}

// returns true if any adjacent site is occupied
func (c *cache) pointHasNeighbor() bool {
	return c.pointRadius <= c.stateRadius+1 && ( c.isLeftIn() || c.isUpIn() || c.isRightIn() || c.isDownIn() )
}

// runs a new 2d aggregation simulation and returns the finished state
func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[types.Point2D]int64 {

	c := cache{}
	c.rng = rng
	c.state = make(map[types.Point2D]int64, nPoints)
	c.state[c.point] = 0
	c.updateStateRadius()

	for i := int64(1); i < nPoints; i++ {
		c.pointToBorder()
		for !c.pointHasNeighbor() || sticking < rng.Float64() {
			c.walkPoint()
		}
		//if _, ok := c.state[c.currPoint]; ok {
		//	panic("Something went wrong here!")
		//}
		c.state[c.point] = i
		if c.pointRadius > c.stateRadius {
			c.updateStateRadius()
		}
	}

	return c.state
}