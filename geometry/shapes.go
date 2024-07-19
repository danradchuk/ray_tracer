package geometry

import (
	"math"

	"github.com/danradchuk/raytracer/shading"
)

// Primitive represents an interface for 3D objects that can be intersected by rays
// and have bounding boxes.
type Primitive interface {
	Intersect(r Ray) *HitRecord
	Bounds() Bounds3
}

// HitRecord stores information about a ray-object intersection.
type HitRecord struct {
	T         float64
	Primitive Primitive
	Material  shading.Material
	Normal    Vec3
}

// Sphere represents a sphere with a center, radius, and material.
type Sphere struct {
	Center   Vec3
	R        float64
	Material shading.Material
}

// Intersect computes the intersection of a ray with the sphere.
func (s Sphere) Intersect(r Ray) *HitRecord {
	co := r.Origin.Sub(s.Center)

	a := r.Direction.Dot(r.Direction)
	b := 2.0 * co.Dot(r.Direction)
	c := co.Dot(co) - s.R*s.R

	d := b*b - 4.0*a*c
	if d < 0 {
		return nil // no intersection
	}

	t1 := (-b + math.Sqrt(d)) / (2.0 * a)
	t2 := (-b - math.Sqrt(d)) / (2.0 * a)

	tMin := math.Min(t1, t2)

	p := r.At(tMin)

	return &HitRecord{tMin, s, s.Material, p.Sub(s.Center).Normalize()}
}

// Bounds returns the bounding box of the sphere.
func (s Sphere) Bounds() Bounds3 {
	center := s.Center
	radius := s.R

	pMin := Point3{center.X - radius, center.Y - radius, center.Z - radius}
	pMax := Point3{center.X + radius, center.Y + radius, center.Z + radius}

	return Bounds3{pMin, pMax}
}

// Plane represents a plane with a point, normal vector, width, and material.
type Plane struct {
	Width    float64
	Point    Vec3
	Normal   Vec3
	Material shading.Material
}

// Intersect computes the intersection of a ray with the plane.
func (p Plane) Intersect(r Ray) *HitRecord {
	// t = ((p0 - l0) * n) / (l * n)
	// p0 - point on the plane
	// l0 - origin of the ray
	// n - normal
	// l - ray direction

	n := p.Normal
	denom := n.Dot(r.Direction) // l * n
	if math.Abs(denom) < 1e-6 {
		return nil
	}

	p0l0 := p.Point.Sub(r.Origin) // p0 - l0
	t := p0l0.Dot(n) / denom
	if t < 0 {
		return nil
	}

	xMin := p.Point.X - (p.Width / 2)
	xMax := p.Point.X + (p.Width / 2)
	zMin := p.Point.Z - (p.Width / 2)
	zMax := p.Point.Z + (p.Width / 2)

	ir := r.At(t)
	if ir.X >= xMin && ir.X <= xMax && ir.Z >= zMin && ir.Z <= zMax {
		return &HitRecord{t, p, p.Material, p.Normal.Normalize()}
	}

	return nil
}

// Bounds returns the bounding box of the plane.
func (p Plane) Bounds() Bounds3 {
	halfWidth := p.Width / 2

	// find two orthogonal vectors on the plane
	var u, v Vec3
	if math.Abs(p.Normal.X) > math.Abs(p.Normal.Y) {
		u = Vec3{-p.Normal.Z, 0, p.Normal.X}
	} else {
		u = Vec3{0, -p.Normal.Z, p.Normal.Y}
	}
	u = u.Normalize()
	v = p.Normal.Cross(u)
	v = v.Normalize()

	// calculate the corners of the plane
	c1 := Vec3{p.Point.X + halfWidth*u.X + halfWidth*v.X, p.Point.Y + halfWidth*u.Y + halfWidth*v.Y, p.Point.Z + halfWidth*u.Z + halfWidth*v.Z}
	c2 := Vec3{p.Point.X + halfWidth*u.X - halfWidth*v.X, p.Point.Y + halfWidth*u.Y - halfWidth*v.Y, p.Point.Z + halfWidth*u.Z - halfWidth*v.Z}
	c3 := Vec3{p.Point.X - halfWidth*u.X + halfWidth*v.X, p.Point.Y - halfWidth*u.Y + halfWidth*v.Y, p.Point.Z - halfWidth*u.Z + halfWidth*v.Z}
	c4 := Vec3{p.Point.X - halfWidth*u.X - halfWidth*v.X, p.Point.Y - halfWidth*u.Y - halfWidth*v.Y, p.Point.Z - halfWidth*u.Z - halfWidth*v.Z}

	// calculate pMin and max coordinates
	pMin := Point3{
		X: math.Min(math.Min(c1.X, c2.X), math.Min(c3.X, c4.X)),
		Y: math.Min(math.Min(c1.Y, c2.Y), math.Min(c3.Y, c4.Y)),
		Z: math.Min(math.Min(c1.Z, c2.Z), math.Min(c3.Z, c4.Z)),
	}
	pMax := Point3{
		X: math.Max(math.Max(c1.X, c2.X), math.Max(c3.X, c4.X)),
		Y: math.Max(math.Max(c1.Y, c2.Y), math.Max(c3.Y, c4.Y)),
		Z: math.Max(math.Max(c1.Z, c2.Z), math.Max(c3.Z, c4.Z)),
	}

	return Bounds3{pMin, pMax}
}

// Triangle represents a triangle with three vertices.
type Triangle struct {
	V0, V1, V2 Vec3
	Material   shading.Material
}

// Intersect computes the intersection of a ray with the triangle.
// It uses barycentric coordinates method.
func (t *Triangle) Intersect(r Ray) *HitRecord {
	epsilon := 0.000001

	// compute vectors for two edges of the triangle
	edge1 := Vec3{
		X: t.V1.X - t.V0.X,
		Y: t.V1.Y - t.V0.Y,
		Z: t.V1.Z - t.V0.Z,
	}
	edge2 := Vec3{
		X: t.V2.X - t.V0.X,
		Y: t.V2.Y - t.V0.Y,
		Z: t.V2.Z - t.V0.Z,
	}

	// compute determinant to check if ray and t are parallel
	h := r.Direction.Cross(edge2)
	det := edge1.Dot(h)
	if math.Abs(det) < epsilon {
		return nil
	}

	// compute inverse determinant and barycentric coordinates
	invDet := 1.0 / det
	s := Vec3{
		X: r.Origin.X - t.V0.X,
		Y: r.Origin.Y - t.V0.Y,
		Z: r.Origin.Z - t.V0.Z,
	}
	u := invDet * s.Dot(h)
	if u < 0 || u > 1 {
		return nil
	}

	q := s.Cross(edge1)
	v := invDet * r.Direction.Dot(q)
	if v < 0 || u+v > 1 {
		return nil
	}

	// compute intersection distance
	tr := invDet * edge2.Dot(q)

	// calculate a normal
	a := t.V1.Sub(t.V0)
	b := t.V2.Sub(t.V0)
	n := a.Cross(b).Normalize()

	if tr > epsilon {
		return &HitRecord{tr, t, t.Material, n}
	}

	return nil
}

// Bounds returns the bounding box of the triangle.
func (t *Triangle) Bounds() Bounds3 {
	xmin := min(t.V0.X, t.V1.X, t.V2.X)
	xmax := max(t.V0.X, t.V1.X, t.V2.X)

	ymin := min(t.V0.Y, t.V1.Y, t.V2.Y)
	ymax := max(t.V0.Y, t.V1.Y, t.V2.Y)

	zmin := min(t.V0.Z, t.V1.Z, t.V2.Z)
	zmax := max(t.V0.Z, t.V1.Z, t.V2.Z)

	return Bounds3{
		PMin: Point3{X: xmin, Y: ymin, Z: zmin},
		PMax: Point3{X: xmax, Y: ymax, Z: zmax},
	}
}
