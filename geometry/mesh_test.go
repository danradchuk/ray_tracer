package geometry

import (
	"math"
	"os"
	"testing"
)

func TestLoadOBJ(t *testing.T) {
	f, err := os.CreateTemp("", "*.obj")
	if err != nil {
		t.Fatalf("can't create a temporary file %s", err.Error())
	}

	// verteces
	f.WriteString("v 1.231 2.321 3.691\n")
	f.WriteString("v 2.676 3.653 7.234\n")
	f.WriteString("v 10.859 6.771 -1.542\n")

	// faces
	f.WriteString("f 1/32/100 2/4/5 3/123/5\n")
	f.WriteString("f 3/32/100 2/4/5 3/123/5\n")
	f.WriteString("f 1//100 3//1 3//5\n")

	var expectedVerteces = make([][]Vec3, 3)
	expectedVerteces[0] = []Vec3{
		{X: 1.231, Y: 2.321, Z: 3.691},
		{X: 2.676, Y: 3.653, Z: 7.234},
		{X: 10.859, Y: 6.771, Z: -1.542},
	}
	expectedVerteces[1] = []Vec3{
		{X: 10.859, Y: 6.771, Z: -1.542},
		{X: 2.676, Y: 3.653, Z: 7.234},
		{X: 10.859, Y: 6.771, Z: -1.542},
	}
	expectedVerteces[2] = []Vec3{
		{X: 1.231, Y: 2.321, Z: 3.691},
		{X: 10.859, Y: 6.771, Z: -1.542},
		{X: 10.859, Y: 6.771, Z: -1.542},
	}

	var expectedIdxs = make([][]int, 3)
	expectedIdxs[0] = []int{0, 1, 2} // idx - 1
	expectedIdxs[1] = []int{2, 1, 2}
	expectedIdxs[2] = []int{0, 2, 2}

	// tInd = 0 -> [1,2,3], 1 -> [3,2,3], 2 -> [1,3,3]
	mesh := LoadOBJ(f.Name())

	for i, m := range mesh.TrianglesToIdxs {
		if equalSlices(expectedIdxs[i], m) == false {
			t.Fatalf("failed: wrong verts expected %v, got %v\n", expectedIdxs[i], m)
		}
	}

	for i, idxs := range mesh.TrianglesToIdxs {
		for j, idx := range idxs {
			v := mesh.Verts[idx]
			ev := expectedVerteces[i][j]

			// comapre v and ev coords
			if !nearlyEqual(v.X, ev.X, 0.00001) {
				t.Fatalf("failed: %d triangle %d vec expected %f, got %f \n", i, j, ev.X, v.X)
			}

			if !nearlyEqual(v.Y, ev.Y, 0.00001) {
				t.Fatalf("failed: %d triangle %d vec expected %f, got %f \n", i, j, ev.Y, v.Y)
			}

			if !nearlyEqual(v.Z, ev.Z, 0.00001) {
				t.Fatalf("failed: %d triangle %d vec expected %f, got %f \n", i, j, ev.Z, v.Z)
			}
		}
	}
}

func equalSlices(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

func nearlyEqual(a, b, epsilon float64) bool {
	absA := math.Abs(float64(a))
	absB := math.Abs(float64(b))
	diff := math.Abs(float64(a - b))

	if a == b {
		return true
	} else if a == 0 || b == 0 || (float64(absA+absB) < math.SmallestNonzeroFloat64) {
		return diff < (epsilon * math.SmallestNonzeroFloat64)
	} else {
		return float64(diff/math.Min(float64(absA+absB), math.MaxFloat64)) < epsilon
	}
}
