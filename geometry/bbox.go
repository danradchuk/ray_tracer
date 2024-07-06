package geometry

import "math"

// Bounds3 represents an Axis-Aligned Bounding Box
type Bounds3 struct {
	Pmin, Pmax Point3
}

func EmptyAABB() Bounds3 {
	return Bounds3{
		Pmin: Point3{
			X: math.MaxFloat64,
			Y: math.MaxFloat64,
			Z: math.MaxFloat64,
		},
		Pmax: Point3{
			X: math.SmallestNonzeroFloat64,
			Y: math.SmallestNonzeroFloat64,
			Z: math.SmallestNonzeroFloat64,
		},
	}
}

func (b Bounds3) Diagonal() Vec3 {
	return b.Pmax.Sub(b.Pmin)
}

// LongestAxis returns 0 - if x is the logest axis, 1 - if y, and 2 - otherwise
func (b Bounds3) LongestAxis() int {
	d := b.Diagonal()

	if d.X > d.Y && d.X > d.Z {
		return 0
	} else if d.Y > d.Z {
		return 1
	}

	return 2
}

func (b Bounds3) Union(other Bounds3) Bounds3 {
	return Bounds3{
		Pmin: Point3{
			X: min(b.Pmin.X, other.Pmin.X),
			Y: min(b.Pmin.Y, other.Pmin.Y),
			Z: min(b.Pmin.Z, other.Pmin.Z),
		},
		Pmax: Point3{
			X: max(b.Pmax.X, other.Pmax.X),
			Y: max(b.Pmax.Y, other.Pmax.Y),
			Z: max(b.Pmax.Z, other.Pmax.Z),
		},
	}
}

func (b Bounds3) UnionPoint3(p Point3) Bounds3 {
	return Bounds3{
		Pmin: Point3{
			X: min(b.Pmin.X, p.X),
			Y: min(b.Pmin.Y, p.Y),
			Z: min(b.Pmin.Z, p.Z),
		},
		Pmax: Point3{
			X: max(b.Pmax.X, p.X),
			Y: max(b.Pmax.Y, p.Y),
			Z: max(b.Pmax.Z, p.Z),
		},
	}
}

func (b Bounds3) GetCoordinatesByAxis(axis int) (float64, float64) {
	if axis == 0 {
		return b.Pmin.X, b.Pmax.X
	}
	if axis == 1 {
		return b.Pmin.Y, b.Pmax.Y
	}

	return b.Pmin.Z, b.Pmax.Z
}

// We don't need a HitRecord here, because we only care about hit or no hit
func (b *Bounds3) Intersect(r Ray) bool {
	invDir := r.Direction.Inverse()

	t1 := (b.Pmin.X - r.Origin.X) * invDir.X
	t2 := (b.Pmax.X - r.Origin.X) * invDir.X

	t3 := (b.Pmin.Y - r.Origin.Y) * invDir.Y
	t4 := (b.Pmax.Y - r.Origin.Y) * invDir.Y

	t5 := (b.Pmin.Z - r.Origin.Z) * invDir.Z
	t6 := (b.Pmax.Z - r.Origin.Z) * invDir.Z

	tmin := math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
	tmax := math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

	if tmax >= tmin && tmax >= 0 {
		return true
	}

	return false
}
