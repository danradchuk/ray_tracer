package geometry

import (
	"math"
)

const Fov = 1.04666667

type Ray struct {
	Origin    Vec3
	Direction Vec3
}

func NewPrimaryRay(camera Vec3, width, height float64, x, y float64) Ray {
	aspectRatio := width / height
	angle := math.Tan(Fov * 0.5)

	// camera space coordinates
	viewX := (2.*((x+.5)/width) - 1) * angle * aspectRatio
	viewY := (1 - 2.*((y+.5)/height)) * angle
	viewZ := camera.Z + 1

	// lookFrom := Vec3{0., 0., 0.}
	// lookAt := Vec3{0., 0., -1.} // our eye looks along positive z-axis
	// up := Vec3{0., 1., 0.}
	//
	//    // u,v,w unit basis vectors
	// w := lookFrom.Sub(lookAt).Normalize() // z-axis
	// u := w.Cross(up).Normalize()          // x-axis
	// v := u.Cross(w)                       // y-axis
	//
	// d := w.Scale(-1.)

	return Ray{
		Origin:    camera,
		Direction: Vec3{viewX, viewY, viewZ}.Sub(camera),
	}
}

func NewSecondaryRay(o Vec3, d Vec3) Ray {
	return Ray{
		Origin:    o,
		Direction: d,
	}
}

func (r Ray) At(t float64) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
