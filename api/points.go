package api

// PointI

// ToPoint2D converts to a Point2D.
func (p PointI) ToPoint2D() Point2D {
	return Point2D{float32(p.X), float32(p.Y)}
}

// ToPoint2DCentered converts to a Point2D and adds 0.5 to X/Y to center inside that map cell.
func (p PointI) ToPoint2DCentered() Point2D {
	return Point2D{float32(p.X) + 0.5, float32(p.Y) + 0.5}
}

// ToPoint converts to a Point with zero Z coordinate.
func (p PointI) ToPoint() Point {
	return Point{float32(p.X), float32(p.Y), 0}
}

// ToPointCentered to a Point with zero Z coordinate and adds 0.5 to X/Y to center inside that map cell.
func (p PointI) ToPointCentered() Point {
	return Point{float32(p.X) + 0.5, float32(p.Y) + 0.5, 0}
}

// VecTo computes the vector from p -> p2.
func (p PointI) VecTo(p2 PointI) VecI {
	return VecI(p2).Sub(VecI(p))
}

// For DirTo and Offset, just convert to Point2D first since the result may not be integers

// Distance computes the absolute distance between two points.
func (p PointI) Distance(p2 PointI) float32 {
	return p.VecTo(p2).Len()
}

// Distance2 computes the squared distance between two points.
func (p PointI) Distance2(p2 PointI) int32 {
	return p.VecTo(p2).Len2()
}

// Manhattan computes the manhattan distance between two points.
func (p PointI) Manhattan(p2 PointI) int32 {
	return p.VecTo(p2).Manhattan()
}

// Add returns the point at the end of v when starting from p.
func (p PointI) Add(v VecI) PointI {
	return PointI(VecI(p).Add(v))
}

// Offset4By returns the Von Neumann neighborhood (or 4-neighborhood) of p.
func (p PointI) Offset4By(by int32) [4]PointI {
	return [...]PointI{
		PointI{p.X, p.Y - by},
		PointI{p.X + by, p.Y},
		PointI{p.X, p.Y + by},
		PointI{p.X - by, p.Y},
	}
}

// Offset8By returns the Moore neighborhood (or 8-neighborhood) of p.
func (p PointI) Offset8By(by int32) [8]PointI {
	return [...]PointI{
		PointI{p.X, p.Y - by},
		PointI{p.X + by, p.Y - by},
		PointI{p.X + by, p.Y},
		PointI{p.X + by, p.Y + by},
		PointI{p.X, p.Y + by},
		PointI{p.X - by, p.Y + by},
		PointI{p.X - by, p.Y},
		PointI{p.X - by, p.Y - by},
	}
}

// Point2D

// ToPointI converts to a PointI by truncating X/Y.
func (p Point2D) ToPointI() PointI {
	return PointI{int32(p.X), int32(p.Y)}
}

// ToPoint converts to a Point by truncating X/Y and setting Z to zero.
func (p Point2D) ToPoint() Point {
	return Point{p.X, p.Y, 0}
}

// VecTo computes the vector from p -> p2.
func (p Point2D) VecTo(p2 Point2D) Vec2D {
	return Vec2D(p2).Sub(Vec2D(p))
}

// DirTo computes the unit vector pointing from p -> p2.
func (p Point2D) DirTo(p2 Point2D) Vec2D {
	return p.VecTo(p2).Norm()
}

// Offset moves a point toward a target by the specified distance.
func (p Point2D) Offset(toward Point2D, by float32) Point2D {
	return p.Add(p.DirTo(toward).Mul(by))
}

// Distance computes the absolute distance between two points.
func (p Point2D) Distance(p2 Point2D) float32 {
	return p.VecTo(p2).Len()
}

// Distance2 computes the squared distance between two points.
func (p Point2D) Distance2(p2 Point2D) float32 {
	return p.VecTo(p2).Len2()
}

// Manhattan computes the manhattan distance between two points.
func (p Point2D) Manhattan(p2 Point2D) float32 {
	return p.VecTo(p2).Manhattan()
}

// Add returns the point at the end of v when starting from p.
func (p Point2D) Add(v Vec2D) Point2D {
	return Point2D(Vec2D(p).Add(v))
}

// Offset4By returns the Von Neumann neighborhood (or 4-neighborhood) of p.
func (p Point2D) Offset4By(by float32) [4]Point2D {
	return [...]Point2D{
		Point2D{p.X, p.Y - by},
		Point2D{p.X + by, p.Y},
		Point2D{p.X, p.Y + by},
		Point2D{p.X - by, p.Y},
	}
}

// Offset8By returns the Moore neighborhood (or 8-neighborhood) of p.
func (p Point2D) Offset8By(by float32) [8]Point2D {
	return [...]Point2D{
		Point2D{p.X, p.Y - by},
		Point2D{p.X + by, p.Y - by},
		Point2D{p.X + by, p.Y},
		Point2D{p.X + by, p.Y + by},
		Point2D{p.X, p.Y + by},
		Point2D{p.X - by, p.Y + by},
		Point2D{p.X - by, p.Y},
		Point2D{p.X - by, p.Y - by},
	}
}

// Point

// ToPointI converts to a PointI by truncating X/Y and dropping Z.
func (p Point) ToPointI() PointI {
	return PointI{int32(p.X), int32(p.Y)}
}

// ToPoint2D converts to a Point2D by dropping Z.
func (p Point) ToPoint2D() Point2D {
	return Point2D{p.X, p.Y}
}

// VecTo computes the vector from p -> p2.
func (p Point) VecTo(p2 Point) Vec {
	return Vec(p2).Sub(Vec(p))
}

// DirTo computes the unit vector pointing from p -> p2.
func (p Point) DirTo(p2 Point) Vec {
	return p.VecTo(p2).Norm()
}

// Offset moves a point toward a target by the specified distance.
func (p Point) Offset(toward Point, by float32) Point {
	return p.Add(p.DirTo(toward).Mul(by))
}

// Distance computes the absolute distance between two points.
func (p Point) Distance(p2 Point) float32 {
	return p.VecTo(p2).Len()
}

// Distance2 computes the squared distance between two points.
func (p Point) Distance2(p2 Point) float32 {
	return p.VecTo(p2).Len2()
}

// Add returns the point at the end of v when starting from p.
func (p Point) Add(v Vec) Point {
	return Point(Vec(p).Add(v))
}
