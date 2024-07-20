package geometry

import "math"

// Ray represents a ray with an origin and direction in 3D space.
type Ray struct {
	Origin    Vec3
	Direction Vec3
}

// NewPrimaryRay creates a primary (camera) ray from the camera for a given screen position (x, y)
// with the specified field of view (fov).
func NewPrimaryRay(eye Vec3, width, height float64, x, y float64, fov int) Ray {
	aspectRatio := width / height
	fovRad := (float64(fov) * math.Pi) / 180
	angle := math.Tan(fovRad * 0.5)

	// camera space coordinates
	alpha := (2.*((x+.5)/width) - 1) * angle * aspectRatio
	beta := (1 - 2.*((y+.5)/height)) * angle

	// for left-handed coordinate system
	lookAt := Vec3{0., 0., 1.}
	up := Vec3{0., 1., 0.}

	// u,v,w unit basis vectors
	w := eye.Sub(lookAt).Normalize() // z-axis
	u := w.Cross(up).Normalize()     // x-axis
	v := u.Cross(w)                  // y-axis

	// our eye looks along positive z-axis; since w looks along negative z-axis, then d looks along positive z-axis
	d := u.Scale(alpha).Add(v.Scale(beta)).Sub(w).Normalize()

	return Ray{
		Origin:    eye,
		Direction: d,
	}
}

// NewSecondaryRay creates a secondary (shadow) ray with a given origin and direction.
func NewSecondaryRay(o Vec3, d Vec3) Ray {
	return Ray{
		Origin:    o,
		Direction: d,
	}
}

// At calculates the position of the ray at distance t.
func (r Ray) At(t float64) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
