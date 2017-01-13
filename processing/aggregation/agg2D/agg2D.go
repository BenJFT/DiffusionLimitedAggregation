package agg2D

import (
	"encoding/gob"
	"math"
	"math/rand"
)

func init() {
	// register the point structure to allow it to be saved and loaded
	gob.Register(Point2D{})
}

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 5
)

// Specialised point structure for a point in 2D. Implements methods required for the Point interface in aggregation
type Point2D struct {
	X, Y int64
}

func (p Point2D) Coordinates() []int64 {
	return []int64{p.X, p.Y}
}

func (p Point2D) SquareDistance(coords []float64) float64 {
	var dx, dy float64 = float64(p.X)-coords[0], float64(p.Y)-coords[1]
	return dx*dx + dy*dy
}

// Structure implemented to remove most, if not all, garbage collector overhead by preallocating all memory to be used
// by the simulation
type cache struct {
	point       Point2D
	pointRadius float64

	rng      *rand.Rand
	lastWalk int64

	state       map[Point2D]int64
	stateRadius float64

	borderRadius    float64
	borderRadiusInt int64

	tempPoint Point2D
	tempFloat float64
}

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

// resets the location of the current point to some location on a sphere touching the border
func (c *cache) pointToBorder() {
	c.tempFloat = 2 * math.Pi * c.rng.Float64()
	c.point.X = int64(math.Sin(c.tempFloat) * c.borderRadius)
	c.point.Y = int64(math.Cos(c.tempFloat) * c.borderRadius)
	c.updateCurrPointRadius()
}

// moves the current point by one, applies periodic boundaries and will not move onto a site which is occupied
func (c *cache) walkPoint() {

	// choose a random from four directions and attempt to walk in that direction
	switch c.rng.Int63n(4) {
	case 0:
		c.point.X++
		// if the point is already occupied return to the starting position
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			c.point.X--
		} else {
			// if the point has stepped over the boundary coordinates wrap it to the other side
			if c.point.X > c.borderRadiusInt {
				c.point.X -= 2 * c.borderRadiusInt
			}
			// save the direction it moved to speed up neighbor checks later
			c.lastWalk = 0
		}
	// repeat as above for other directions
	case 1:
		c.point.X--
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			c.point.X++
		} else {
			if c.point.X < -c.borderRadiusInt {
				c.point.X += 2 * c.borderRadiusInt
			}
			c.lastWalk = 1
		}
	case 2:
		c.point.Y++
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			c.point.Y--
		} else {
			if c.point.Y > c.borderRadiusInt {
				c.point.Y -= 2 * c.borderRadiusInt
			}
			c.lastWalk = 2
		}
	case 3:
		c.point.Y--
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			c.point.Y++
		} else {
			if c.point.Y < -c.borderRadiusInt {
				c.point.Y += 2 * c.borderRadiusInt
			}
			c.lastWalk = 3
		}
	}
	c.updateCurrPointRadius()
}

// tests if the adjacent cite in the +x direction is occupied
func (c *cache) isRightIn() bool {
	// check is if site in only if did not come from that direction. equivalence is faster to check, so results in
	// cumulative gains over many loops
	return c.lastWalk != 1 && c.rightIn()
}
// checks if the point one space to the right is in the set of points
func (c *cache) rightIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.X++
	_, ok = c.state[c.tempPoint]
	return
}

// tests if the adjacent cite in the -x direction is occupied
// following all have same structure as the checks for right
func (c *cache) isLeftIn() bool {
	return c.lastWalk != 0 && c.leftIn()
}
func (c *cache) leftIn() (ok bool) {
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
	// evaluation order allows for speedup as the majority of the particles time may be spent far from the aggregate
	// checking if it is near first allows much calculation to be skipped when it is far, without introducing much
	// overhead while close
	return c.pointRadius <= c.stateRadius+1 && (c.isLeftIn() || c.isUpIn() || c.isRightIn() || c.isDownIn())
}

// runs a new 2d aggregation simulation and returns the finished state
func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[Point2D]int64 {
	// initialize variables
	c := cache{}
	c.rng = rng
	c.state = make(map[Point2D]int64, nPoints)
	c.state[c.point] = 0
	c.updateStateRadius()

	// add points one at a time
	for i := int64(1); i < nPoints; i++ {
		// move the working point to the border
		c.pointToBorder()
		// walk the point until it sticks to a neighbor
		for !c.pointHasNeighbor() || sticking < rng.Float64() {
			c.walkPoint()
		}
		// add the point to the set
		c.state[c.point] = i
		// update the state as necessary
		if c.pointRadius > c.stateRadius {
			c.updateStateRadius()
		}
	}

	return c.state
}
