package analysis

import (
	"math"
	"github.com/Benjft/DiffusionLimitedAggregation/aggregation"
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
		dr := (r - ball.Radius)/2
		coords := p1.Coordinates()
		dCoords := make([]float64, len(coords))
		for i, x := range coords {
			dCoords[i] = float64(x) - ball.Coords[i]
		}

		for i, x := range dCoords {
			ball.Coords[i] += dr*x/r
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
