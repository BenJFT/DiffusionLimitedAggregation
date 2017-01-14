package agg1D

import (
	"encoding/gob"
	"math"
	"math/rand"
)

func init() {
	gob.Register(Point1D{})
}

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 4
)

type Point1D struct {
	A int64
}

func (p Point1D) Coordinates() []int64 {
	return []int64{ p.A }
}

func (p Point1D) SquareDistance(coords []float64) float64 {
	var dA = float64(p.A)-coords[0]
	return dA*dA
}


type cache struct {
	point Point1D
	pointRadius float64
	rng *rand.Rand
	lastWalk int64
	state map[Point1D]int64
	stateRadius float64
	startRadius float64
	borderRadius float64
	borderRadiusInt int64
	tempPoint Point1D
	tempFloatA, tempFloatB float64
}

func (c *cache) updateCurrPointRadius() {
	c.pointRadius = math.Sqrt(float64(c.point.A*c.point.A))
}

func (c *cache) updateStateRadius() {
	c.stateRadius = c.pointRadius
	c.startRadius = c.stateRadius+BORDER_CONST
	c.borderRadius = c.startRadius*BORDER_SCALE
	c.borderRadiusInt = int64(c.borderRadius)
}

func (c *cache) pointIn() (ok bool) {
	_, ok = c.state[c.point]
	return
}


func (c *cache) pointToBorder() {
		c.tempFloatA = 1
	c.point.A = int64(c.tempFloatA * c.startRadius)

	c.updateCurrPointRadius()
}

func (c *cache) walkPoint() {
		switch c.rng.Int63n(2) {
	case 0:
		c.point.A++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.A--
		} else {
			if c.point.A > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 0
		}
	case 1:
		c.point.A--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.A++
		} else {
			if c.point.A < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 1
		}
}

	c.updateCurrPointRadius()
}

func (c *cache) isPlusAIn() bool {
	return c.lastWalk != 1 && c.plusAIn()
}
func (c *cache) plusAIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.A++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusAIn() bool {
	return c.lastWalk != 0 && c.minusAIn()
}
func (c *cache) minusAIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.A--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) pointHasNeighbor() bool {
	return c.pointRadius <= c.stateRadius+1 && (c.isPlusAIn() || c.isMinusAIn())
}

func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[Point1D]int64 {
	c := cache{}
	c.rng = rng
	c.state = make(map[Point1D]int64, nPoints)
	c.state[c.point] = 0
	c.updateStateRadius()
	
	for i := int64(1); i < nPoints; i++ {
		c.pointToBorder()
		for !c.pointHasNeighbor() || sticking < rng.Float64() {
			c.walkPoint()
		}
		c.state[c.point]=i
		if c.pointRadius > c.stateRadius {
			c.updateStateRadius()
		}
	}
	return c.state
}
