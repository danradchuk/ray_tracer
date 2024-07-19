package geometry

import "math"

// Vec3 represents a 3-dimensional vector.
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// Add adds the corresponding components of two vectors.
func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{
		X: v.X + u.X,
		Y: v.Y + u.Y,
		Z: v.Z + u.Z,
	}
}

// Sub subtracts the corresponding components of another vector from this vector.
func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{
		X: v.X - u.X,
		Y: v.Y - u.Y,
		Z: v.Z - u.Z,
	}
}

// Scale scales each component of the vector by a scalar.
func (v Vec3) Scale(s float64) Vec3 {
	return Vec3{
		X: s * v.X,
		Y: s * v.Y,
		Z: s * v.Z,
	}
}

// Dot returns the dot product of two vectors.
func (v Vec3) Dot(u Vec3) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

// Cross returns the cross product of two vectors.
func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{
		v.Y*u.Z - v.Z*u.Y,
		v.Z*u.X - v.X*u.Z,
		v.X*u.Y - v.Y*u.X,
	}
}

// Norm returns the Euclidean norm of the vector.
func (v Vec3) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns a unit vector in the same direction as this vector.
func (v Vec3) Normalize() Vec3 {
	norm := v.Norm()
	return Vec3{v.X / norm, v.Y / norm, v.Z / norm}
}

// Lerp performs a linear interpolation between two vectors.
func (v Vec3) Lerp(u Vec3, t float64) Vec3 {
	return v.Scale(1 - t).Add(u.Scale(t))
}

// Inverse returns a vector with each component being the reciprocal of this vector's components.
func (v Vec3) Inverse() Vec3 {
	return Vec3{1 / v.X, 1 / v.Y, 1 / v.Z}
}
