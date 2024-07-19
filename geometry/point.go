package geometry

// Point3 represents a point in 3D space.
type Point3 struct {
	X, Y, Z float64
}

// Sub subtracts another point from this point and returns the result as a vector.
func (p Point3) Sub(other Point3) Vec3 {
	return Vec3{
		X: p.X - other.X,
		Y: p.Y - other.Y,
		Z: p.Z - other.Z,
	}
}

// Scale scales the coordinates of this point by a scalar value and returns the result.
func (p Point3) Scale(s float64) Point3 {
	return Point3{
		X: p.X * s,
		Y: p.Y * s,
		Z: p.Z * s,
	}
}

// Add adds another point to this point and returns the result.
func (p Point3) Add(other Point3) Point3 {
	return Point3{
		X: p.X + other.X,
		Y: p.Y + other.Y,
		Z: p.Z + other.Z,
	}
}

// GetCoordinateByAxis returns the coordinate of this point along the specified axis (0 for X, 1 for Y, 2 for Z).
func (p Point3) GetCoordinateByAxis(axis int) float64 {
	if axis == 0 {
		return p.X
	}

	if axis == 1 {
		return p.Y
	}

	return p.Z
}
