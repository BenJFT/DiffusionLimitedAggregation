package genagg

import (
	"math"
	"math/rand"
)

// references a given lattice point
type Point struct {
	X, Y int64
}

// Speed ignorant implementations for external use
// Faster internal values change values in the local cache struct to avoid slowdown from garbage collection, and hence
// are called on the cache variable itself
func (p *Point) Radius() float64 {
	return math.Sqrt(float64(p.X* p.X + p.Y* p.Y))
}

func (p1 *Point) Distance(p2 *Point) float64 {
	var dx, dy int64
	dx = p1.X-p2.X
	dy = p1.Y-p2.Y

	return math.Sqrt(float64(dx*dx + dy*dy))
}

// caching structure. Exists as an optimisation by preventing memory reassignment, this helps reduce time spent in
// garbage collection, and helps keep variables in the cpu register or caches. Early implementations of this gave a
// a significant speedup
type cache struct {
	currPoint *Point
	currPointRadius float64

	tempPoint *Point

	rng *rand.Rand
}

func (c *cache) pointRadius() {
	c.currPointRadius = math.Sqrt(float64(c.currPoint.X*c.currPoint.X + c.currPoint.Y*c.currPoint.Y))
}

func (c *cache) walkPoint() bool {
	//TODO
}

