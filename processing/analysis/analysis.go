package analysis

import (
	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation"
	"math"
)

type Ball struct {
	Coords                 []float64
	Radius, SquareDistance float64
}

func extendTo(ball *Ball, p1 aggregation.Point) {
	if d2 := p1.SquareDistance(ball.Coords); d2 <= ball.SquareDistance {
		return
	} else {
		r := math.Sqrt(d2)
		dr := (r - ball.Radius) / 2
		coords := p1.Coordinates()
		dCoords := make([]float64, len(coords))
		for i, x := range coords {
			dCoords[i] = float64(x) - ball.Coords[i]
		}

		for i, x := range dCoords {
			ball.Coords[i] += dr * x / r
		}
		ball.Radius += dr
		ball.SquareDistance = math.Pow(ball.Radius, 2)
	}
}

func ApproxBounding(points []aggregation.Point) (balls []Ball) {
	balls = make([]Ball, len(points))

	var lastBall Ball = Ball{}
	lastBall.Coords = make([]float64, len(points[0].Coordinates()))

	for i, p0 := range points {
		var ball Ball = lastBall
		extendTo(&ball, p0)
		balls[i] = ball
		lastBall = ball
	}
	return balls
}

//TODO radius of gyration http://mathworld.wolfram.com/RadiusofGyration.html

type elem struct {
	N int64
	radius float64
}
func gyrationRadius(points []aggregation.Point, chanel chan elem) {
	var mean []float64 = make([]float64, len(points[0].Coordinates()))
	var N int64 = int64(len(points))
	for _, p := range points {
		for i, x := range p.Coordinates() {
			mean[i] += float64(x)/float64(N)
		}
	}

	var e elem = elem{N:N}
	for _, p := range points {
		e.radius += p.SquareDistance(mean)/float64(N)
	}

	e.radius = math.Sqrt(e.radius)

	chanel <- e
}

func GyrationRadii(points []aggregation.Point) (radii []float64) {
	radii = make([]float64, len(points))

	var (
		chanel chan elem = make(chan elem)
	)

	for i := range points {
		go gyrationRadius(points[:i+1], chanel)
	}

	for range points {
		e := <- chanel
		radii[e.N-1] = e.radius
	}

	return radii
}