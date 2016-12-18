package agg3D

import (
	"math"
	"math/rand"

	"github.com/Benjft/DiffusionLimitedAggregation/util/types"
)

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 5
)

type cache struct {
	point types.Point3D
	pointRadius float64

	rng *rand.Rand
	lastWalk int64

	state map[types.Point3D]int64
	stateRadius float64

	borderRadius float64
	borderRadiusInt int64

	tempPoint types.Point3D
	tempA float64
	tempB float64
}

func (c *cache) updateCurrPointRadius() {
	c.pointRadius = math.Sqrt(float64(c.point.X*c.point.X + c.point.Y*c.point.Y + c.point.Z*c.point.Z))
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
	c.tempA = 2*math.Pi*c.rng.Float64()
	c.tempB = 2*math.Pi*c.rng.Float64()

	c.point.X = int64(math.Sin(c.tempA) * math.Sin(c.tempB) * c.borderRadius)
	c.point.Y = int64(math.Cos(c.tempA) * math.Sin(c.tempB) * c.borderRadius)
	c.point.Z = int64(math.Cos(c.tempB) * c.borderRadius)

	c.updateCurrPointRadius()
}

// moves the current point by one, applies periodic boundaries and will not move onto a site which is occupied
func (c *cache) walkPoint() {

	switch point := &c.point; c.rng.Int63n(6) {
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
	case 4:
		point.Z++
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			point.Z--
		} else {
			if point.Z > c.borderRadiusInt {
				point.Z -= 2*c.borderRadiusInt
			}
			c.lastWalk = 4
		}
	case 5:
		point.Z--
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			point.Z++
		} else {
			if point.Z > c.borderRadiusInt {
				point.Z += 2*c.borderRadiusInt
			}
			c.lastWalk = 5
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
func (c *cache) isFrontIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 3 && c.frontIn()
}
func (c *cache) frontIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.Y++
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the -y direction is occupied
func (c *cache) isBackIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 2 && c.backIn()
}
func (c *cache) backIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.Y--
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the +z direction is occupied
func (c *cache) isUpIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 5 && c.upIn()
}
func (c *cache) upIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.Z++
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the -z direction is occupied
func (c *cache) isDownIn() bool {
	// check is if site in only if did not come from that direction as walk check is much faster
	return c.lastWalk != 4 && c.downIn()
}
func (c *cache) downIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.Z--
	_, ok = c.state[c.tempPoint]
	return
}

// returns true if any adjacent site is occupied
func (c *cache) pointHasNeighbor() bool {
	return c.pointRadius <= c.stateRadius+1 && (
		c.isLeftIn() ||
			c.isUpIn() ||
			c.isFrontIn() ||
			c.isRightIn() ||
			c.isDownIn() ||
			c.isBackIn())
}

// runs a new 2d aggregation simulation and returns the finished state
func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[types.Point3D]int64 {

	c := cache{}
	c.rng = rng
	c.state = make(map[types.Point3D]int64, nPoints)
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