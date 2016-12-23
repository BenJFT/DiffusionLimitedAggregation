package types

import (
	"math"
	"encoding/gob"
)

func init() {
	gob.Register(Point2D{})
	gob.Register(Point3D{})
	gob.Register(Run{})
	gob.Register(HSVA{})
}

//--------------Point Types--------------

type Point interface {
	Coordinates() []int64
	SquareDistance(coords ...float64) float64
}

//func (p0 Point) SquareDistancePoint(p1 Point) float64 {
//	tmp := p1.Coordinates()
//	coord := make([]float64, len(tmp))
//	for i, x := range tmp {
//		coord[i] = float64(x)
//	}
//	return p0.SquareDistance(coord...)
//}

type Point2D struct {
	X, Y int64
}

func (p Point2D) Coordinates() []int64 {
	return []int64{p.X, p.Y}
}

func (p Point2D) SquareDistance(coords ...float64) float64 {
	var x, y int64 = p.X, p.Y
	var fx, fy, X, Y = float64(x), float64(y), coords[0], coords[1]
	dx := fx - X
	dy := fy - Y

	return dx*dx + dy*dy
}

type Point3D struct {
	X, Y, Z int64
}

func (p Point3D) Coordinates() []int64 {
	return []int64{p.X, p.Y, p.Z}
}

func (p Point3D) SquareDistance(coords ...float64) float64 {
	var dx, dy, dz float64
	dx = float64(p.X) - coords[0]
	dy = float64(p.Y) - coords[1]
	dz = float64(p.Z) - coords[2]

	return dx*dx + dy*dy + dz*dz
}

//--------------Utility Structures--------------

type Run struct {
	NPoints, NDimension, NRuns, Seed int64
	Sticking float64
	Points [][]Point
}

type HSVA struct {
	H, S, V, A float64
}

func (hsva HSVA) RGBA() (r, g, b, a uint32) {
	const (
		max32 = float64(math.MaxUint32)
	)
	var (
		H, S, V = hsva.H, hsva.S, hsva.V
		zone = math.Floor(H * 6)
		part = H * 6 - zone
		low = V * (1 - S)
		midA = V * (1 - part*S)
		midB = V * (1 - (1-part)*S)
	)

	switch uint8(zone) % 6 {
	case 0:
		r = uint32(max32 * V)
		g = uint32(max32 * midB)
		b = uint32(max32 * low)
	case 1:
		r = uint32(max32 * midA)
		g = uint32(max32 * V)
		b = uint32(max32 * low)
	case 2:
		r = uint32(max32 * low)
		g = uint32(max32 * V)
		b = uint32(max32 * midB)
	case 3:
		r = uint32(max32 * low)
		g = uint32(max32 * midA)
		b = uint32(max32 * V)
	case 4:
		r = uint32(max32 * midB)
		g = uint32(max32 * low)
		b = uint32(max32 * V)
	case 5:
		r = uint32(max32 * V)
		g = uint32(max32 * low)
		b = uint32(max32 * midA)
	}

	a = uint32(max32 * hsva.A)
	return r, g, b, a
}
