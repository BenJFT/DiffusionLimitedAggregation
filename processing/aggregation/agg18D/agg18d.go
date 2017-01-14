package agg18D

import (
	"encoding/gob"
	"math"
	"math/rand"
)

func init() {
	gob.Register(Point18D{})
}

const (
	BORDER_SCALE float64 = 1.5
	BORDER_CONST float64 = 4
)

type Point18D struct {
	A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R int64
}

func (p Point18D) Coordinates() []int64 {
	return []int64{ p.A, p.B, p.C, p.D, p.E, p.F, p.G, p.H, p.I, p.J, p.K, p.L, p.M, p.N, p.O, p.P, p.Q, p.R }
}

func (p Point18D) SquareDistance(coords []float64) float64 {
	var dA, dB, dC, dD, dE, dF, dG, dH, dI, dJ, dK, dL, dM, dN, dO, dP, dQ, dR = float64(p.A)-coords[0], float64(p.B)-coords[1], float64(p.C)-coords[2], float64(p.D)-coords[3], float64(p.E)-coords[4], float64(p.F)-coords[5], float64(p.G)-coords[6], float64(p.H)-coords[7], float64(p.I)-coords[8], float64(p.J)-coords[9], float64(p.K)-coords[10], float64(p.L)-coords[11], float64(p.M)-coords[12], float64(p.N)-coords[13], float64(p.O)-coords[14], float64(p.P)-coords[15], float64(p.Q)-coords[16], float64(p.R)-coords[17]
	return dA*dA + dB*dB + dC*dC + dD*dD + dE*dE + dF*dF + dG*dG + dH*dH + dI*dI + dJ*dJ + dK*dK + dL*dL + dM*dM + dN*dN + dO*dO + dP*dP + dQ*dQ + dR*dR
}


type cache struct {
	point Point18D
	pointRadius float64
	rng *rand.Rand
	lastWalk int64
	state map[Point18D]int64
	stateRadius float64
	startRadius float64
	borderRadius float64
	borderRadiusInt int64
	tempPoint Point18D
	tempFloatA, tempFloatB float64
}

func (c *cache) updateCurrPointRadius() {
	c.pointRadius = math.Sqrt(float64(c.point.A*c.point.A + c.point.B*c.point.B + c.point.C*c.point.C + c.point.D*c.point.D + c.point.E*c.point.E + c.point.F*c.point.F + c.point.G*c.point.G + c.point.H*c.point.H + c.point.I*c.point.I + c.point.J*c.point.J + c.point.K*c.point.K + c.point.L*c.point.L + c.point.M*c.point.M + c.point.N*c.point.N + c.point.O*c.point.O + c.point.P*c.point.P + c.point.Q*c.point.Q + c.point.R*c.point.R))
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
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.A = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.B = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.C = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.D = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.E = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.F = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.G = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.H = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.I = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.J = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.K = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.L = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.M = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.N = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.O = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.P = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.tempFloatB = c.rng.Float64() * 2 * math.Pi
	c.point.Q = int64(math.Cos(c.tempFloatB) * c.tempFloatA * c.startRadius)
	c.tempFloatA *= math.Sin(c.tempFloatB)
	c.point.R = int64(c.tempFloatA * c.startRadius)

	c.updateCurrPointRadius()
}

func (c *cache) walkPoint() {
		switch c.rng.Int63n(36) {
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
	case 2:
		c.point.B++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.B--
		} else {
			if c.point.B > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 2
		}
	case 3:
		c.point.B--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.B++
		} else {
			if c.point.B < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 3
		}
	case 4:
		c.point.C++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.C--
		} else {
			if c.point.C > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 4
		}
	case 5:
		c.point.C--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.C++
		} else {
			if c.point.C < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 5
		}
	case 6:
		c.point.D++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.D--
		} else {
			if c.point.D > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 6
		}
	case 7:
		c.point.D--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.D++
		} else {
			if c.point.D < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 7
		}
	case 8:
		c.point.E++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.E--
		} else {
			if c.point.E > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 8
		}
	case 9:
		c.point.E--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.E++
		} else {
			if c.point.E < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 9
		}
	case 10:
		c.point.F++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.F--
		} else {
			if c.point.F > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 10
		}
	case 11:
		c.point.F--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.F++
		} else {
			if c.point.F < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 11
		}
	case 12:
		c.point.G++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.G--
		} else {
			if c.point.G > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 12
		}
	case 13:
		c.point.G--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.G++
		} else {
			if c.point.G < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 13
		}
	case 14:
		c.point.H++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.H--
		} else {
			if c.point.H > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 14
		}
	case 15:
		c.point.H--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.H++
		} else {
			if c.point.H < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 15
		}
	case 16:
		c.point.I++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.I--
		} else {
			if c.point.I > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 16
		}
	case 17:
		c.point.I--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.I++
		} else {
			if c.point.I < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 17
		}
	case 18:
		c.point.J++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.J--
		} else {
			if c.point.J > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 18
		}
	case 19:
		c.point.J--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.J++
		} else {
			if c.point.J < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 19
		}
	case 20:
		c.point.K++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.K--
		} else {
			if c.point.K > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 20
		}
	case 21:
		c.point.K--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.K++
		} else {
			if c.point.K < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 21
		}
	case 22:
		c.point.L++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.L--
		} else {
			if c.point.L > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 22
		}
	case 23:
		c.point.L--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.L++
		} else {
			if c.point.L < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 23
		}
	case 24:
		c.point.M++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.M--
		} else {
			if c.point.M > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 24
		}
	case 25:
		c.point.M--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.M++
		} else {
			if c.point.M < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 25
		}
	case 26:
		c.point.N++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.N--
		} else {
			if c.point.N > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 26
		}
	case 27:
		c.point.N--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.N++
		} else {
			if c.point.N < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 27
		}
	case 28:
		c.point.O++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.O--
		} else {
			if c.point.O > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 28
		}
	case 29:
		c.point.O--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.O++
		} else {
			if c.point.O < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 29
		}
	case 30:
		c.point.P++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.P--
		} else {
			if c.point.P > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 30
		}
	case 31:
		c.point.P--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.P++
		} else {
			if c.point.P < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 31
		}
	case 32:
		c.point.Q++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.Q--
		} else {
			if c.point.Q > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 32
		}
	case 33:
		c.point.Q--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.Q++
		} else {
			if c.point.Q < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 33
		}
	case 34:
		c.point.R++
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.R--
		} else {
			if c.point.R > c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 34
		}
	case 35:
		c.point.R--
		if c.pointRadius < c.startRadius && c.pointIn() {
			c.point.R++
		} else {
			if c.point.R < -c.borderRadiusInt {
				c.pointToBorder()
			}
			c.lastWalk = 35
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

func (c *cache) isPlusHIn() bool {
	return c.lastWalk != 15 && c.plusHIn()
}
func (c *cache) plusHIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.H++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusHIn() bool {
	return c.lastWalk != 14 && c.minusHIn()
}
func (c *cache) minusHIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.H--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusIIn() bool {
	return c.lastWalk != 17 && c.plusIIn()
}
func (c *cache) plusIIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.I++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusIIn() bool {
	return c.lastWalk != 16 && c.minusIIn()
}
func (c *cache) minusIIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.I--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusJIn() bool {
	return c.lastWalk != 19 && c.plusJIn()
}
func (c *cache) plusJIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.J++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusJIn() bool {
	return c.lastWalk != 18 && c.minusJIn()
}
func (c *cache) minusJIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.J--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusKIn() bool {
	return c.lastWalk != 21 && c.plusKIn()
}
func (c *cache) plusKIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.K++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusKIn() bool {
	return c.lastWalk != 20 && c.minusKIn()
}
func (c *cache) minusKIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.K--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusLIn() bool {
	return c.lastWalk != 23 && c.plusLIn()
}
func (c *cache) plusLIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.L++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusLIn() bool {
	return c.lastWalk != 22 && c.minusLIn()
}
func (c *cache) minusLIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.L--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusMIn() bool {
	return c.lastWalk != 25 && c.plusMIn()
}
func (c *cache) plusMIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.M++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusMIn() bool {
	return c.lastWalk != 24 && c.minusMIn()
}
func (c *cache) minusMIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.M--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusNIn() bool {
	return c.lastWalk != 27 && c.plusNIn()
}
func (c *cache) plusNIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.N++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusNIn() bool {
	return c.lastWalk != 26 && c.minusNIn()
}
func (c *cache) minusNIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.N--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusOIn() bool {
	return c.lastWalk != 29 && c.plusOIn()
}
func (c *cache) plusOIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.O++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusOIn() bool {
	return c.lastWalk != 28 && c.minusOIn()
}
func (c *cache) minusOIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.O--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusPIn() bool {
	return c.lastWalk != 31 && c.plusPIn()
}
func (c *cache) plusPIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.P++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusPIn() bool {
	return c.lastWalk != 30 && c.minusPIn()
}
func (c *cache) minusPIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.P--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusQIn() bool {
	return c.lastWalk != 33 && c.plusQIn()
}
func (c *cache) plusQIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.Q++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusQIn() bool {
	return c.lastWalk != 32 && c.minusQIn()
}
func (c *cache) minusQIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.Q--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) isPlusRIn() bool {
	return c.lastWalk != 35 && c.plusRIn()
}
func (c *cache) plusRIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.R++
	_, ok = c.state[c.tempPoint]
	return
}
func (c *cache) isMinusRIn() bool {
	return c.lastWalk != 34 && c.minusRIn()
}
func (c *cache) minusRIn() (ok bool) {
	c.tempPoint = c.point
	c.tempPoint.R--
	_, ok = c.state[c.tempPoint]
	return
}

func (c *cache) pointHasNeighbor() bool {
	return c.pointRadius <= c.stateRadius+1 && (c.isPlusAIn() || c.isMinusAIn() || c.isPlusBIn() || c.isMinusBIn() || c.isPlusCIn() || c.isMinusCIn() || c.isPlusDIn() || c.isMinusDIn() || c.isPlusEIn() || c.isMinusEIn() || c.isPlusFIn() || c.isMinusFIn() || c.isPlusGIn() || c.isMinusGIn() || c.isPlusHIn() || c.isMinusHIn() || c.isPlusIIn() || c.isMinusIIn() || c.isPlusJIn() || c.isMinusJIn() || c.isPlusKIn() || c.isMinusKIn() || c.isPlusLIn() || c.isMinusLIn() || c.isPlusMIn() || c.isMinusMIn() || c.isPlusNIn() || c.isMinusNIn() || c.isPlusOIn() || c.isMinusOIn() || c.isPlusPIn() || c.isMinusPIn() || c.isPlusQIn() || c.isMinusQIn() || c.isPlusRIn() || c.isMinusRIn())
}

func RunNew(nPoints int64, sticking float64, rng *rand.Rand) map[Point18D]int64 {
	c := cache{}
	c.rng = rng
	c.state = make(map[Point18D]int64, nPoints)
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
