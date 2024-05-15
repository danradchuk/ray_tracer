package geometry

import (
	"math"

	"github.com/danradchuk/raytracer/shading"
)

// Sphere
type Sphere struct {
	Center   Vec3
	R        float64
	Material shading.Material
}

func (s Sphere) GetMaterial() shading.Material {
	return s.Material
}

func (s Sphere) Intersection(r Ray) (bool, float64) {
	co := r.Origin.Sub(s.Center)

	a := r.Direction.Dot(r.Direction)
	b := 2.0 * co.Dot(r.Direction)
	c := co.Dot(co) - s.R*s.R

	d := b*b - 4.0*a*c
	if d < 0 {
		return false, 0 // no intersection
	}

	t1 := (-b + math.Sqrt(d)) / (2.0 * a)
	t2 := (-b - math.Sqrt(d)) / (2.0 * a)

	return true, math.Min(t1, t2)
}

func (s Sphere) NormalAt(p Vec3) Vec3 {
	return p.Sub(s.Center).Normalize()
}

// Plane
type Plane struct {
	Width    float64
	Point    Vec3
	Normal   Vec3
	Material shading.Material
}

func (p Plane) GetMaterial() shading.Material {
	return p.Material
}

func (p Plane) Intersection(r Ray) (bool, float64) {
	// t = ((p0 - l0) * n) / (l * n)
	// p0 - point on the plane
	// l0 - origin of the ray
	// n - normal
	// l - ray direction

	n := p.Normal
	denom := n.Dot(r.Direction) // l * n
	if math.Abs(denom) < 1e-6 {
		return false, 0.0
	}

	p0l0 := p.Point.Sub(r.Origin) // p0 - l0
	t := p0l0.Dot(n) / denom
	if t < 0 {
		return false, 0.0
	}

	xMin := p.Point.X - (p.Width / 2)
	xMax := p.Point.X + (p.Width / 2)
	zMin := p.Point.Z - (p.Width / 2)
	zMax := p.Point.Z + (p.Width / 2)

	ir := r.At(t)
	if ir.X >= xMin && ir.X <= xMax && ir.Z >= zMin && ir.Z <= zMax {
		return true, t
	}

	return false, 0.0
}

func (p Plane) NormalAt(_ Vec3) Vec3 {
	return p.Normal.Normalize()
}

// Triangle
type Triangle struct {
	V0, V1, V2 Vec3
	Material   shading.Material
}

func (t *Triangle) GetMaterial() shading.Material {
	return t.Material
}

func (t *Triangle) Intersection(r Ray) (bool, float64) {
	epsilon := 0.000001

	// Compute vectors for two edges of the triangle
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

	// Compute determinant to check if ray and t are parallel
	h := r.Direction.Cross(edge2)
	det := edge1.Dot(h)
	if math.Abs(det) < epsilon {
		return false, .0
	}

	// Compute inverse determinant and barycentric coordinates
	invDet := 1.0 / det
	s := Vec3{
		X: r.Origin.X - t.V0.X,
		Y: r.Origin.Y - t.V0.Y,
		Z: r.Origin.Z - t.V0.Z,
	}
	u := invDet * s.Dot(h)
	if u < 0 || u > 1 {
		return false, .0
	}

	q := s.Cross(edge1)
	v := invDet * r.Direction.Dot(q)
	if v < 0 || u+v > 1 {
		return false, .0
	}

	// Compute intersection distance
	tr := invDet * edge2.Dot(q)
	if tr > epsilon {
		return true, tr
	}

	return false, .0
}

func (t *Triangle) NormalAt(_ Vec3) Vec3 {
	a := t.V1.Sub(t.V0)
	b := t.V2.Sub(t.V0)

	return a.Cross(b).Normalize()
}
