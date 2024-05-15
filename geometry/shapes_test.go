package geometry

import (
	"testing"
)

func TestVectorAddition(t *testing.T) {
	v1 := Vec3{0, 0, 0}
	v2 := Vec3{1, 1, 1}

	resV := v1.Add(v2)
	if resV.X != 1 {
		t.Errorf("X: got %f want %f", resV.X, 1.0)
	}

	if resV.Y != 1 {
		t.Errorf("Y: got %f want %f", resV.Y, 1.0)
	}

	if resV.Z != 1 {
		t.Errorf("Z: got %f want %f", resV.Z, 1.0)
	}
}

func TestVectorSubtraction(t *testing.T) {
	v1 := Vec3{0, 0, 0}
	v2 := Vec3{1, 1, 1}

	resV := v1.Sub(v2)
	if resV.X != -1 {
		t.Errorf("X: got %f want %f", resV.X, -1.0)
	}

	if resV.Y != -1 {
		t.Errorf("Y: got %f want %f", resV.Y, -1.0)
	}

	if resV.Z != -1 {
		t.Errorf("Z: got %f want %f", resV.Z, -1.0)
	}
}

func TestVectorDotProduct(t *testing.T) {
	v1 := Vec3{5, 100, 2}
	v2 := Vec3{3, 2, 5}

	res := v1.Dot(v2)
	if res != 225 {
		t.Errorf("Dot Prdouct: got %f want %f", res, 225.0)
	}
}

func TestVectorScale(t *testing.T) {
	v1 := Vec3{5, 100, 2}

	resV := v1.Scale(5)

	if resV.X != 25 {
		t.Errorf("X: got %f want %f", resV.X, 25.0)
	}

	if resV.Y != 500 {
		t.Errorf("Y: got %f want %f", resV.Y, 500.0)
	}

	if resV.Z != 10 {
		t.Errorf("Z: got %f want %f", resV.Z, 10.0)
	}
}

// TODO fix this test
// func TestRayScaling(t *testing.T) {
// 	r := Ray{
// 		Origin:    Vec3{0, 0, -1},
// 		Direction: Vec3{33, 10, 20},
// 	}
//
// 	scaled := r.At(0)
// 	if scaled != r.Origin {
// 		t.Errorf("got %v want %v", r.Origin, scaled.Origin)
// 	}
//
// }
