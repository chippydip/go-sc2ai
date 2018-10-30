package api

import (
	"math"
)

// PointI

// ToPoint2D ...
func (p PointI) ToPoint2D() Point2D {
	return Point2D{float32(p.X), float32(p.Y)}
}

// ToCenteredPoint2D ...
func (p PointI) ToCenteredPoint2D() Point2D {
	return Point2D{float32(p.X) + 0.5, float32(p.Y) + 0.5}
}

// ToPoint ...
func (p PointI) ToPoint() Point {
	return Point{float32(p.X), float32(p.Y), 0}
}

// ToCenteredPoint ...
func (p PointI) ToCenteredPoint() Point {
	return Point{float32(p.X) + 0.5, float32(p.Y) + 0.5, 0}
}

// Add ...
func (p PointI) Add(p2 PointI) PointI {
	return PointI{p.X + p2.X, p.Y + p2.Y}
}

// Sub ...
func (p PointI) Sub(p2 PointI) PointI {
	return PointI{p.X - p2.X, p.Y - p2.Y}
}

// Mul ...
func (p PointI) Mul(c int32) PointI {
	return PointI{p.X * c, p.Y * c}
}

// Dot ...
func (p PointI) Dot(p2 PointI) int64 {
	return int64(p.X)*int64(p2.X) + int64(p.Y)*int64(p2.Y)
}

// LenSqr ...
func (p PointI) LenSqr() int64 {
	return p.Dot(p)
}

// Len ...
func (p PointI) Len() float64 {
	return math.Sqrt(float64(p.LenSqr()))
}

// Point2D

// ToPointI ...
func (p Point2D) ToPointI() PointI {
	return PointI{int32(p.X), int32(p.Y)}
}

// ToPoint ...
func (p Point2D) ToPoint() Point {
	return Point{p.X, p.Y, 0}
}

// Add ...
func (p Point2D) Add(p2 Point2D) Point2D {
	return Point2D{p.X + p2.X, p.Y + p2.Y}
}

// Sub ...
func (p Point2D) Sub(p2 Point2D) Point2D {
	return Point2D{p.X - p2.X, p.Y - p2.Y}
}

// Mul ...
func (p Point2D) Mul(c float32) Point2D {
	return Point2D{p.X * c, p.Y * c}
}

// Mul64 ...
func (p Point2D) Mul64(c float64) Point2D {
	return Point2D{float32(float64(p.X) * c), float32(float64(p.Y) * c)}
}

// Dot ...
func (p Point2D) Dot(p2 Point2D) float64 {
	return float64(p.X)*float64(p2.X) + float64(p.Y)*float64(p2.Y)
}

// LenSqr ...
func (p Point2D) LenSqr() float64 {
	return p.Dot(p)
}

// Len ...
func (p Point2D) Len() float64 {
	return math.Sqrt(p.LenSqr())
}

// Normalize ...
func (p Point2D) Normalize() Point2D {
	return p.Mul64(1.0 / p.Len())
}

// Point

// ToPointI ...
func (p Point) ToPointI() PointI {
	return PointI{int32(p.X), int32(p.Y)}
}

// ToPoint2D ...
func (p Point) ToPoint2D() Point2D {
	return Point2D{p.X, p.Y}
}

// Add ...
func (p Point) Add(p2 Point) Point {
	return Point{p.X + p2.X, p.Y + p2.Y, p.Z + p2.Z}
}

// Sub ...
func (p Point) Sub(p2 Point) Point {
	return Point{p.X - p2.X, p.Y - p2.Y, p.Z - p2.Z}
}

// Mul ...
func (p Point) Mul(c float32) Point {
	return Point{p.X * c, p.Y * c, p.Z * c}
}

// Mul64 ...
func (p Point) Mul64(c float64) Point {
	return Point{float32(float64(p.X) * c), float32(float64(p.Y) * c), float32(float64(p.Z) * c)}
}

// Dot ...
func (p Point) Dot(p2 Point) float64 {
	return float64(p.X)*float64(p2.X) + float64(p.Y)*float64(p2.Y) + float64(p.Z)*float64(p2.Z)
}

// LenSqr ...
func (p Point) LenSqr() float64 {
	return p.Dot(p)
}

// Len ...
func (p Point) Len() float64 {
	return math.Sqrt(p.LenSqr())
}

// Normalize ...
func (p Point) Normalize() Point {
	return p.Mul64(1.0 / p.Len())
}

// Cross ...
func (p Point) Cross(p2 Point) Point {
	return Point{p.Y*p2.Z - p.Z*p2.Y, p.Z*p2.X - p.X*p2.Z, p.X*p2.Y - p.Y*p2.X}
}
