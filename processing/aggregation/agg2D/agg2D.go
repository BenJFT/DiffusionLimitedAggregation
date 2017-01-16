//This is an auto generated file from genAggFiles.py
package agg2D

import (
	"encoding/gob"
	"math"
	"math/rand"
)

func init() {
	gob.Register(Point2D{})
}

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 2
)

type Point2D struct {
	A, B int64
}

func (p Point2D) Coordinates() []int64 {
	return []int64{ p.A, p.B }
}

func (p Point2D) SquareDistance(coords []float64) float64 {
	var dA, dB = float64(p.A)-coords[0], float64(p.B)-coords[1]
	return dA*dA + dB*dB
}


type cache struct {
	point Point2D
	pointRadius float64
	rng *rand.Rand
	lastWalk int64
	state map[Point2D]int64
	stateRadius float64
	borderRadius float64
	borderRadiusInt int64
	tempPoint Point2D
	tempFloatA, tempFloatB float64
}

func (c *cache) updateCurrPointRadius() {
	c.pointRadius = math.Sqrt(float64(c.point.A*c.point.A + c.point.B*c.point.B))
}

func (c *cache) updateStateRadius() {
	c.stateRadius = c.pointRadius
	c.borderRadius = c.stateRadius*BORDER_SCALE+BORDER_CONST
	c.borderRadiusInt = int64(c.borderRadius)
}

func (c *cache) pointIn() (ok bool) {
	_, ok = c.state[c.point]
	return
}


func (c *cache) pointToBorder() {
	c.tempFloatA = 1
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.A = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.borderRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.point.B = int64(c.tempFloatA * c.borderRadius)

	c.updateCurrPointRadius()
}

func (c *cache) walkPoint() {
		switch c.rng.Int63n(4) {
	case 0:
		c.point.A++
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.A--
		} else {
			if c.point.A > c.borderRadiusInt {
				c.point.A -= 2*c.borderRadiusInt
			}
			c.lastWalk = 0
		}
	case 1:
		c.point.A--
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.A++
		} else {
			if c.point.A < -c.borderRadiusInt {
				c.point.A += 2*c.borderRadiusInt
			}
			c.lastWalk = 1
		}
	case 2:
		c.point.B++
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.B--
		} else {
			if c.point.B > c.borderRadiusInt {
				c.point.B -= 2*c.borderRadiusInt
			}
			c.lastWalk = 2
		}
	case 3:
		c.point.B--
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.B++
		} else {
			if c.point.B < -c.borderRadiusInt {
				c.point.B += 2*c.borderRadiusInt
			}
			c.lastWalk = 3
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

func (c *cache) isPlusBIn() bool {
	return c.lastWalk != 3 && c.plusBIn()
}
func (c *cache) plusBIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.B++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusBIn() bool {
	return c.lastWalk != 2 && c.minusBIn()
}
func (c *cache) minusBIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.B--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) pointHasNeighbor() bool {
	return c.pointRadius <= c.stateRadius+1 && (c.isPlusAIn() || c.isMinusAIn() || c.isPlusBIn() || c.isMinusBIn())
}

func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[Point2D]int64 {
	c := cache{}
	c.rng = rng
	c.state = make(map[Point2D]int64, nPoints)
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
