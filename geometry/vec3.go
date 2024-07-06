package geometry

import "math"

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func (v Vec3) Add(u Vec3) Vec3 {
	return Vec3{
		X: v.X + u.X,
		Y: v.Y + u.Y,
		Z: v.Z + u.Z,
	}
}

func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{
		X: v.X - u.X,
		Y: v.Y - u.Y,
		Z: v.Z - u.Z,
	}
}

func (v Vec3) Scale(s float64) Vec3 {
	return Vec3{
		X: s * v.X,
		Y: s * v.Y,
		Z: s * v.Z,
	}
}

func (v Vec3) Dot(u Vec3) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{
		v.Y*u.Z - v.Z*u.Y,
		v.Z*u.X - v.X*u.Z,
		v.X*u.Y - v.Y*u.X,
	}
}

func (v Vec3) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Normalize() Vec3 {
	norm := v.Norm()
	return Vec3{v.X / norm, v.Y / norm, v.Z / norm}
}

func (v Vec3) Lerp(u Vec3, t float64) Vec3 {
	return v.Scale(1 - t).Add(u.Scale(t))
}

func (v Vec3) Inverse() Vec3 {
	return Vec3{1 / v.X, 1 / v.Y, 1 / v.Z}
}
