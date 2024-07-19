package geometry

import "math"

// Bounds3 represents an Axis-Aligned Bounding Box.
// PMin is the lower left point, PMax is the upper right point.
type Bounds3 struct {
	PMin, PMax Point3
}

// EmptyAABB creates an empty axis-aligned bounding box with extreme values.
func EmptyAABB() Bounds3 {
	return Bounds3{
		PMin: Point3{
			X: math.MaxFloat64,
			Y: math.MaxFloat64,
			Z: math.MaxFloat64,
		},
		PMax: Point3{
			X: math.SmallestNonzeroFloat64,
			Y: math.SmallestNonzeroFloat64,
			Z: math.SmallestNonzeroFloat64,
		},
	}
}

// Diagonal returns the diagonal vector of the bounding box.
func (b Bounds3) Diagonal() Vec3 {
	return b.PMax.Sub(b.PMin)
}

// LongestAxis returns
// 0 - if x is the longest axis
// 1 - if y is the longest axis
// 2 - otherwise
func (b Bounds3) LongestAxis() int {
	d := b.Diagonal()

	if d.X > d.Y && d.X > d.Z {
		return 0
	} else if d.Y > d.Z {
		return 1
	}

	return 2
}

// Union returns the union of two bounding boxes.
func (b Bounds3) Union(other Bounds3) Bounds3 {
	return Bounds3{
		PMin: Point3{
			X: min(b.PMin.X, other.PMin.X),
			Y: min(b.PMin.Y, other.PMin.Y),
			Z: min(b.PMin.Z, other.PMin.Z),
		},
		PMax: Point3{
			X: max(b.PMax.X, other.PMax.X),
			Y: max(b.PMax.Y, other.PMax.Y),
			Z: max(b.PMax.Z, other.PMax.Z),
		},
	}
}

// UnionPoint3 returns the union of the bounding box with a point.
func (b Bounds3) UnionPoint3(p Point3) Bounds3 {
	return Bounds3{
		PMin: Point3{
			X: min(b.PMin.X, p.X),
			Y: min(b.PMin.Y, p.Y),
			Z: min(b.PMin.Z, p.Z),
		},
		PMax: Point3{
			X: max(b.PMax.X, p.X),
			Y: max(b.PMax.Y, p.Y),
			Z: max(b.PMax.Z, p.Z),
		},
	}
}

// GetCoordinatesByAxis returns the minimum and maximum coordinates along the specified axis.
func (b Bounds3) GetCoordinatesByAxis(axis int) (float64, float64) {
	if axis == 0 {
		return b.PMin.X, b.PMax.X
	}
	if axis == 1 {
		return b.PMin.Y, b.PMax.Y
	}

	return b.PMin.Z, b.PMax.Z
}

// Intersect checks if a ray intersects the bounding box.
// Note that we don't need a HitRecord here, because we only care about hit or no hit.
func (b *Bounds3) Intersect(r Ray) bool {
	invDir := r.Direction.Inverse()

	t1 := (b.PMin.X - r.Origin.X) * invDir.X
	t2 := (b.PMax.X - r.Origin.X) * invDir.X

	t3 := (b.PMin.Y - r.Origin.Y) * invDir.Y
	t4 := (b.PMax.Y - r.Origin.Y) * invDir.Y

	t5 := (b.PMin.Z - r.Origin.Z) * invDir.Z
	t6 := (b.PMax.Z - r.Origin.Z) * invDir.Z

	tMin := math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
	tMax := math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

	if tMax >= tMin && tMax >= 0 {
		return true
	}

	return false
}
