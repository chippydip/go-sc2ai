package main

import (
	"math"

	"github.com/chippydip/go-sc2ai/api"
)

type Point struct {
	X, Y float64
}

func PointFrom2D(a *api.Point2D) *Point {
	return &Point{X: float64(a.X), Y: float64(a.Y)}
}

func PointFrom3D(a *api.Point) *Point {
	return &Point{X: float64(a.X), Y: float64(a.Y)}
}

func (a *Point) To2D() *api.Point2D {
	return &api.Point2D{X: float32(a.X), Y: float32(a.Y)}
}

func (a *Point) Add(b *Point) *Point {
	return &Point{X: a.X + b.X, Y: a.Y + b.Y}
}

func (a *Point) Sub(b *Point) *Point {
	return &Point{X: a.X - b.X, Y: a.Y - b.Y}
}

func (a *Point) Mul(b *Point) *Point {
	return &Point{X: a.X * b.X, Y: a.Y * b.Y}
}

func (a *Point) Div(b *Point) *Point {
	return &Point{X: a.X / b.X, Y: a.Y / b.Y}
}

func (a *Point) Distance(b *Point) float64 {
	return math.Sqrt(a.DistanceSquared(b))
}

func (a *Point) DistanceSquared(b *Point) float64 {
	return math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2)
}

// User should check that he receives not nil
func (a *Point) ClosestUnit(units []*api.Unit) *api.Unit {
	var closest *api.Unit
	for _, unit := range units {
		if closest == nil ||
			(a.DistanceSquared(PointFrom3D(closest.Pos)) > a.DistanceSquared(PointFrom3D(unit.Pos))) {
			closest = unit
		}
	}
	return closest
}
