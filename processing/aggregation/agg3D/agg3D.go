package agg3D

import (
	"encoding/gob"
	"math"
	"math/rand"
)

func init() {
	// register the point structure to allow it to be saved and loaded
	gob.Register(Point3D{})
}

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 3
)

// Specialised point structure for a point in 3D. Implements methods required for the Point interface in aggregation
type Point3D struct {
	X, Y, Z int64
}

func (p Point3D) Coordinates() []int64 {
	return []int64{p.X, p.Y, p.Z}
}

func (p Point3D) SquareDistance(coords []float64) float64 {
	var dx, dy, dz float64
	dx = float64(p.X) - coords[0]
	dy = float64(p.Y) - coords[1]
	dz = float64(p.Z) - coords[2]

	return dx*dx + dy*dy + dz*dz
}

// Structure implemented to remove most, if not all, garbage collector overhead by preallocating all memory to be used
// by the simulation
type cache struct {
	point       Point3D
	pointRadius float64

	rng      *rand.Rand
	lastWalk int64

	state       map[Point3D]int64
	stateRadius float64

	borderRadius    float64
	borderRadiusInt int64

	tempPoint Point3D
	tempA     float64
	tempB     float64
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

// resets the location of the current point to some location on a sphere touching the border
func (c *cache) pointToBorder() {
	c.tempA = 2 * math.Pi * c.rng.Float64()
	c.tempB = 2 * math.Pi * c.rng.Float64()

	c.point.Z = int64(math.Cos(c.tempB) * c.borderRadius)
	c.point.Y = int64(math.Cos(c.tempA) * math.Sin(c.tempB) * c.borderRadius)
	c.point.X = int64(		      math.Sin(c.tempA) * math.Sin(c.tempB) * c.borderRadius)

	c.updateCurrPointRadius()
}

// moves the current point by one, applies periodic boundaries and will not move onto a site which is occupied
func (c *cache) walkPoint() {

	// choose a random from six directions and attempt to walk in that direction
	switch c.rng.Int63n(6) {
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
	case 4:
		c.point.Z++
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			c.point.Z--
		} else {
			if c.point.Z > c.borderRadiusInt {
				c.point.Z -= 2 * c.borderRadiusInt
			}
			c.lastWalk = 4
		}
	case 5:
		c.point.Z--
		if c.pointRadius < 4+c.stateRadius && c.pointIn() {
			c.point.Z++
		} else {
			if c.point.Z < -c.borderRadiusInt {
				c.point.Z += 2 * c.borderRadiusInt
			}
			c.lastWalk = 5
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
func (c *cache) isFrontIn() bool {
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
	// evaluation order allows for speedup as the majority of the particles time may be spent far from the aggregate
	// checking if it is near first allows much calculation to be skipped when it is far, without introducing much
	// overhead while close
	return c.pointRadius <= c.stateRadius+1 && (c.isLeftIn() ||
		c.isUpIn() ||
		c.isFrontIn() ||
		c.isRightIn() ||
		c.isDownIn() ||
		c.isBackIn())
}

// runs a new 3d aggregation simulation and returns the finished state
func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[Point3D]int64 {
	// initialize variables
	c := cache{}
	c.rng = rng
	c.state = make(map[Point3D]int64, nPoints)
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
