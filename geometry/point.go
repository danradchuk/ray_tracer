package geometry

type Point3 struct {
	X, Y, Z float64
}

func (p Point3) Sub(other Point3) Vec3 {
	return Vec3{
		X: p.X - other.X,
		Y: p.Y - other.Y,
		Z: p.Z - other.Z,
	}
}

func (p Point3) Scale(s float64) Point3 {
	return Point3{
		X: p.X * s,
		Y: p.Y * s,
		Z: p.Z * s,
	}
}

func (p Point3) Add(other Point3) Point3 {
	return Point3{
		X: p.X + other.X,
		Y: p.Y + other.Y,
		Z: p.Z + other.Z,
	}
}

func (p Point3) GetCoordinateByAxis(axis int) float64 {
	if axis == 0 {
		return p.X
	}

	if axis == 1 {
		return p.Y
	}

	return p.Z
}
