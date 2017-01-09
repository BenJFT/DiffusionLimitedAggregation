package aggND

import (
	"encoding/gob"
	"math"
	"math/rand"
	"strconv"
)

func init() {
	// register the point structure to allow it to be saved and loaded
	gob.Register(PointND{})
}

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 3
)

// Un-specialised point, uses a alice to store coordinates, results in slower read and write speeds, but allows for
// any number of dimensions to be simulated over.  Implements required interface for Point
type PointND []int64

func (p PointND) Coordinates() []int64 {
	return []int64(p)
}

func (p PointND) SquareDistance(coords []float64) float64 {
	var d2 float64
	for i, x := range p {
		var d float64 = float64(x) - coords[i]
		d2 += d * d
	}

	return d2
}

// function required to for lookup in hash table as a slice value has no equality defined. Part of the cause of slowdown
func (p PointND) toString() string {
	var out string = ""
	for _, x := range p {
		out += strconv.FormatInt(x, 36) // display numbers in base 36 for shortest length
		out += "-"
	}
	return out
}

// as a slice is a pointer to an underlying array, pass by value does not create a coppy of the array, only the pointer
// as such the slice must be explicitly copied each time a duplicate is needed
func (p PointND) copyPoint() (q PointND) {
	q = make([]int64, len(p))
	copy(q, p)
	return
}

// as the hash table cannot store the point as it's key (slices have no defined equality, the point must be stored with
// it's end index as the result instead using this structure
type elem struct {
	p   PointND
	idx int64
}

// Structure implemented to remove some garbage collector overhead by pre-allocating all memory to be used
// by the simulation
type cache struct {
	dims int64

	point       PointND
	pointRadius float64

	rng      *rand.Rand
	lastWalk int64

	state       map[string]elem
	stateRadius float64

	borderRadius    float64
	borderRadiusInt int64

	spawnAngle float64
	sinFact    float64
	stepAxis   int64
	stepSign   int64
	tempPoint  PointND
}

func (c *cache) updateCurrPointRadius() {
	c.pointRadius = 0
	for _, x := range c.point {
		c.pointRadius += float64(x * x)
	}
	c.pointRadius = math.Sqrt(c.pointRadius)
}

func (c *cache) updateStateRadius() {
	c.stateRadius = c.pointRadius
	c.borderRadius = c.stateRadius*BORDER_SCALE + BORDER_CONST
	c.borderRadiusInt = int64(c.borderRadius)
}

// returns true if the current point is in the caches state
func (c *cache) pointIn() (ok bool) {
	_, ok = c.state[c.point.toString()]
	return

}

// returns true if the current point is in the caches state
func (c *cache) tempPointIn() (ok bool) {
	_, ok = c.state[c.tempPoint.toString()]
	return

}

// resets the location of the current point to some location on the border
func (c *cache) pointToBorder() {
	// N dimensional generalisation of a point on the surface of a sphere
	c.sinFact = 1

	for i := int64(0); i < c.dims-1; i++ {
		c.spawnAngle = 2 * math.Pi * c.rng.Float64()
		c.point[i] = int64(c.sinFact * math.Cos(c.spawnAngle) * c.borderRadius)
		c.sinFact *= math.Sin(c.spawnAngle)
	}
	c.point[c.dims-1] = int64(c.sinFact * c.borderRadius)

	c.updateCurrPointRadius()
}

// moves the current point by one, applies periodic boundaries and will not move onto a site which is occupied
func (c *cache) walkPoint() {

	// choose a random axis to walk along
	c.stepAxis = c.rng.Int63n(c.dims)
	// choose the direction to walk in that axis
	c.stepSign = c.rng.Int63n(2)*2 - 1

	//walk that way
	c.point[c.stepAxis] += c.stepSign

	// if the site is already occupied return to previous location
	if c.pointRadius < 4+c.stateRadius && c.pointIn() {
		c.point[c.stepAxis] -= c.stepSign
	} else {
		// wrap around at boundary
		if c.point[c.stepAxis] > c.borderRadiusInt {
			c.point[c.stepAxis] -= 2 * c.borderRadiusInt
		} else if c.point[c.stepAxis] < -c.borderRadiusInt {
			c.point[c.stepAxis] += 2 * c.borderRadiusInt
		}
		c.lastWalk = 2 * c.stepAxis
		if c.stepSign < 0 {
			c.lastWalk += 1
		}
	}

	c.updateCurrPointRadius()
}

func (c *cache) neighborIn() bool {
	copy(c.tempPoint, c.point)
	// checks for a neighbor in each direction
	for i := int64(0); i < c.dims; i++ {
		// special case of the direction just walked from, no need to check the spot it came from
		if i == c.stepAxis {
			if c.stepSign > 0 {
				c.tempPoint[i] += 1
				if c.tempPointIn() {
					return true
				}
				c.tempPoint[i] -= 1
			} else {
				c.tempPoint[i] -= 1
				if c.tempPointIn() {
					return true
				}
				c.tempPoint[i] += 1
			}
		} else {
			c.tempPoint[i] += 1
			if c.tempPointIn() {
				return true
			}
			c.tempPoint[i] -= 2
			if c.tempPointIn() {
				return true
			}
			c.tempPoint[i] += 1
		}
	}
	return false
}

// returns true if any adjacent site is occupied
func (c *cache) pointHasNeighbor() bool {
	return c.pointRadius <= c.stateRadius+1 && c.neighborIn()
}

// runs a new 3d aggregation simulation and returns the finished state
func RunNew(nPoints int64, sticking float64, rng *rand.Rand, dimension int64) []PointND {
	// initialize memory
	c := cache{}
	c.rng = rng
	c.state = make(map[string]elem, nPoints)
	c.point = make(PointND, dimension)
	c.dims = dimension

	// set seed point
	c.state[c.point.toString()] = elem{p: c.point.copyPoint(), idx: 0}
	c.tempPoint = c.point.copyPoint()
	c.updateStateRadius()

	// add points until count reached
	for i := int64(1); i < nPoints; i++ {
		// moves the point to the boundary
		c.pointToBorder()
		// walks it until it sticks
		for !c.pointHasNeighbor() || sticking < rng.Float64() {
			c.walkPoint()
		}
		// add it to the set
		c.state[c.point.toString()] = elem{p: c.point.copyPoint(), idx: i}
		if c.pointRadius > c.stateRadius {
			c.updateStateRadius()
		}
	}

	// convert it to an array
	var ret []PointND = make([]PointND, nPoints)
	for _, e := range c.state {
		ret[e.idx] = e.p
	}

	return ret
}
