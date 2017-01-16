//This is an auto generated file from genAggFiles.py
package agg7D

import (
	"encoding/gob"
	"math"
	"math/rand"
)

func init() {
	gob.Register(Point7D{})
}

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 2
)

type Point7D struct {
	A, B, C, D, E, F, G int64
}

func (p Point7D) Coordinates() []int64 {
	return []int64{ p.A, p.B, p.C, p.D, p.E, p.F, p.G }
}

func (p Point7D) SquareDistance(coords []float64) float64 {
	var dA, dB, dC, dD, dE, dF, dG = float64(p.A)-coords[0], float64(p.B)-coords[1], float64(p.C)-coords[2], float64(p.D)-coords[3], float64(p.E)-coords[4], float64(p.F)-coords[5], float64(p.G)-coords[6]
	return dA*dA + dB*dB + dC*dC + dD*dD + dE*dE + dF*dF + dG*dG
}


type cache struct {
	point Point7D
	pointRadius float64
	rng *rand.Rand
	lastWalk int64
	state map[Point7D]int64
	stateRadius float64
	borderRadius float64
	borderRadiusInt int64
	tempPoint Point7D
	tempFloatA, tempFloatB float64
}

func (c *cache) updateCurrPointRadius() {
	c.pointRadius = math.Sqrt(float64(c.point.A*c.point.A + c.point.B*c.point.B + c.point.C*c.point.C + c.point.D*c.point.D + c.point.E*c.point.E + c.point.F*c.point.F + c.point.G*c.point.G))
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
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.B = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.borderRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.C = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.borderRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.D = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.borderRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.E = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.borderRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.F = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.borderRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.point.G = int64(c.tempFloatA * c.borderRadius)

	c.updateCurrPointRadius()
}

func (c *cache) walkPoint() {
		switch c.rng.Int63n(14) {
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
	case 4:
		c.point.C++
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.C--
		} else {
			if c.point.C > c.borderRadiusInt {
				c.point.C -= 2*c.borderRadiusInt
			}
			c.lastWalk = 4
		}
	case 5:
		c.point.C--
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.C++
		} else {
			if c.point.C < -c.borderRadiusInt {
				c.point.C += 2*c.borderRadiusInt
			}
			c.lastWalk = 5
		}
	case 6:
		c.point.D++
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.D--
		} else {
			if c.point.D > c.borderRadiusInt {
				c.point.D -= 2*c.borderRadiusInt
			}
			c.lastWalk = 6
		}
	case 7:
		c.point.D--
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.D++
		} else {
			if c.point.D < -c.borderRadiusInt {
				c.point.D += 2*c.borderRadiusInt
			}
			c.lastWalk = 7
		}
	case 8:
		c.point.E++
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.E--
		} else {
			if c.point.E > c.borderRadiusInt {
				c.point.E -= 2*c.borderRadiusInt
			}
			c.lastWalk = 8
		}
	case 9:
		c.point.E--
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.E++
		} else {
			if c.point.E < -c.borderRadiusInt {
				c.point.E += 2*c.borderRadiusInt
			}
			c.lastWalk = 9
		}
	case 10:
		c.point.F++
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.F--
		} else {
			if c.point.F > c.borderRadiusInt {
				c.point.F -= 2*c.borderRadiusInt
			}
			c.lastWalk = 10
		}
	case 11:
		c.point.F--
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.F++
		} else {
			if c.point.F < -c.borderRadiusInt {
				c.point.F += 2*c.borderRadiusInt
			}
			c.lastWalk = 11
		}
	case 12:
		c.point.G++
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.G--
		} else {
			if c.point.G > c.borderRadiusInt {
				c.point.G -= 2*c.borderRadiusInt
			}
			c.lastWalk = 12
		}
	case 13:
		c.point.G--
		if c.pointRadius < c.stateRadius+BORDER_CONST && c.pointIn() {
			c.point.G++
		} else {
			if c.point.G < -c.borderRadiusInt {
				c.point.G += 2*c.borderRadiusInt
			}
			c.lastWalk = 13
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

func (c *cache) isPlusCIn() bool {
	return c.lastWalk != 5 && c.plusCIn()
}
func (c *cache) plusCIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.C++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusCIn() bool {
	return c.lastWalk != 4 && c.minusCIn()
}
func (c *cache) minusCIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.C--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusDIn() bool {
	return c.lastWalk != 7 && c.plusDIn()
}
func (c *cache) plusDIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.D++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusDIn() bool {
	return c.lastWalk != 6 && c.minusDIn()
}
func (c *cache) minusDIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.D--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusEIn() bool {
	return c.lastWalk != 9 && c.plusEIn()
}
func (c *cache) plusEIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.E++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusEIn() bool {
	return c.lastWalk != 8 && c.minusEIn()
}
func (c *cache) minusEIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.E--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusFIn() bool {
	return c.lastWalk != 11 && c.plusFIn()
}
func (c *cache) plusFIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.F++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusFIn() bool {
	return c.lastWalk != 10 && c.minusFIn()
}
func (c *cache) minusFIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.F--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusGIn() bool {
	return c.lastWalk != 13 && c.plusGIn()
}
func (c *cache) plusGIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.G++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusGIn() bool {
	return c.lastWalk != 12 && c.minusGIn()
}
func (c *cache) minusGIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.G--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) pointHasNeighbor() bool {
	return c.pointRadius <= c.stateRadius+1 && (c.isPlusAIn() || c.isMinusAIn() || c.isPlusBIn() || c.isMinusBIn() || c.isPlusCIn() || c.isMinusCIn() || c.isPlusDIn() || c.isMinusDIn() || c.isPlusEIn() || c.isMinusEIn() || c.isPlusFIn() || c.isMinusFIn() || c.isPlusGIn() || c.isMinusGIn())
}

func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[Point7D]int64 {
	c := cache{}
	c.rng = rng
	c.state = make(map[Point7D]int64, nPoints)
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
