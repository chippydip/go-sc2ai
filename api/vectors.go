package api

import (
	"math"
)

// VecI is a 2D vector with integer components.
type VecI PointI

// Neg flips a vector to point in the opposite direction.
func (v VecI) Neg() VecI {
	return VecI{-v.X, -v.Y}
}

// Add two vectors and return the result.
func (v VecI) Add(v2 VecI) VecI {
	return VecI{v.X + v2.X, v.Y + v2.Y}
}

// Sub subtracts two vectors and returns the results.
func (v VecI) Sub(v2 VecI) VecI {
	return VecI{v.X - v2.X, v.Y - v2.Y}
}

// Mul scales the vector by a constant.
func (v VecI) Mul(c int32) VecI {
	return VecI{v.X * c, v.Y * c}
}

// Dot computes the dot product with another vector.
func (v VecI) Dot(v2 VecI) int32 {
	return v.X*v2.X + v.Y*v2.Y
}

// Len2 computes the squared length (magnitude) of the vector.
func (v VecI) Len2() int32 {
	return v.Dot(v)
}

// Len computes the length (magnitude) of the vector.
func (v VecI) Len() float64 {
	return math.Sqrt(float64(v.Len2()))
}

// Manhattan computes the manhattan distance represented by this vector.
func (v VecI) Manhattan() int32 {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	return v.X + v.Y
}

// Vec2D is a 2D vector if real components.
type Vec2D Point2D

// Neg flips a vector to point in the opposite direction.
func (v Vec2D) Neg() Vec2D {
	return Vec2D{-v.X, -v.Y}
}

// Add two vectors and return the result.
func (v Vec2D) Add(v2 Vec2D) Vec2D {
	return Vec2D{v.X + v2.X, v.Y + v2.Y}
}

// Sub subtracts two vectors and returns the results.
func (v Vec2D) Sub(v2 Vec2D) Vec2D {
	return Vec2D{v.X - v2.X, v.Y - v2.Y}
}

// Mul scales the vector by a constant.
func (v Vec2D) Mul(c float32) Vec2D {
	return Vec2D{v.X * c, v.Y * c}
}

// Mul64 scales the vector by a 64-bit constant. This involves additional casting so Mul should be preferred when 32-bits are sufficient.
func (v Vec2D) Mul64(c float64) Vec2D {
	return Vec2D{float32(float64(v.X) * c), float32(float64(v.Y) * c)}
}

// Dot computes the dot product with another vector.
func (v Vec2D) Dot(v2 Vec2D) float32 {
	return v.X*v2.X + v.Y*v2.Y
}

// Len2 computes the squared length (magnitude) of the vector.
func (v Vec2D) Len2() float32 {
	return v.Dot(v)
}

// Len computes the length (magnitude) of the vector.
func (v Vec2D) Len() float64 {
	return math.Sqrt(float64(v.Len2()))
}

// Manhattan computes the manhattan distance represented by this vector.
func (v Vec2D) Manhattan() float32 {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	return v.X + v.Y
}

// Norm computes the unit vector pointing is the same direction as v.
func (v Vec2D) Norm() Vec2D {
	return v.Mul64(1.0 / v.Len())
}

// Vec is a 3D vector with real components.
type Vec Point

// Neg flips a vector to point in the opposite direction.
func (v Vec) Neg() Vec {
	return Vec{-v.X, -v.Y, -v.Z}
}

// Add two vectors and return the result.
func (v Vec) Add(v2 Vec) Vec {
	return Vec{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

// Sub subtracts two vectors and returns the results.
func (v Vec) Sub(v2 Vec) Vec {
	return Vec{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

// Mul scales the vector by a constant.
func (v Vec) Mul(c float32) Vec {
	return Vec{v.X * c, v.Y * c, v.Z * c}
}

// Mul64 scales the vector by a 64-bit constant. This involves additional casting so Mul should be preferred when 32-bits are sufficient.
func (v Vec) Mul64(c float64) Vec {
	return Vec{float32(float64(v.X) * c), float32(float64(v.Y) * c), float32(float64(v.Z) * c)}
}

// Dot computes the dot product with another vector.
func (v Vec) Dot(v2 Vec) float32 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

// Len2 computes the squared length (magnitude) of the vector.
func (v Vec) Len2() float32 {
	return v.Dot(v)
}

// Len computes the length (magnitude) of the vector.
func (v Vec) Len() float64 {
	return math.Sqrt(float64(v.Len2()))
}

// Manhattan computes the manhattan distance represented by this vector.
func (v Vec) Manhattan() float32 {
	if v.X < 0 {
		v.X = -v.X
	}
	if v.Y < 0 {
		v.Y = -v.Y
	}
	if v.Z < 0 {
		v.Z = -v.Z
	}
	return v.X + v.Y + v.Z
}

// Norm computes the unit vector pointing is the same direction as v.
func (v Vec) Norm() Vec {
	return v.Mul64(1.0 / v.Len())
}

// Cross computes the cross product of v x v2.
func (v Vec) Cross(v2 Vec) Vec {
	return Vec{v.Y*v2.Z - v.Z*v2.Y, v.Z*v2.X - v.X*v2.Z, v.X*v2.Y - v.Y*v2.X}
}
