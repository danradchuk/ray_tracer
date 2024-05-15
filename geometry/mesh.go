package geometry

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/danradchuk/raytracer/shading"
)

/*
OBJ File Format:
v 3228/1234/32 - a position of vertex ( {x: 3228, y: 1234, z: 32} )
f 8/8/8 7/7/7 9/9/9 10/10/10  - a face (8,7,9,10)
*/
type IndexedMesh struct {
	TrianglesToIdxs [][]int
	Verts           []Vec3
}

func LoadOBJ(fName string) *IndexedMesh {
	var tInd [][]int
	var verts []Vec3

	f, err := os.Open(fName)
	check(err)

	s := bufio.NewScanner(f)
	for s.Scan() {
		fields := strings.Fields(s.Text())

		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "v":
			x, err := strconv.ParseFloat(fields[1], 64)
			check(err)
			y, err := strconv.ParseFloat(fields[2], 64)
			check(err)
			z, err := strconv.ParseFloat(fields[3], 64)
			check(err)

			verts = append(verts, Vec3{X: x, Y: y, Z: z})
		case "f":
			// compute all vertices for n-2 vertexes
			// where n - number of all vertexes in the current row
			// a row in .obj file is - v 1/1/1 2/2/2/2 3/3/3/3 4/4/4/4
			numVerts := len(fields[1:])
			for i := 2; i <= numVerts-1; i++ {
				var verts []int

				// push the first vertex into the slice
				x, err := strconv.Atoi(strings.Split(fields[1], "/")[0])
				check(err)
				verts = append(verts, x-1)

				// push y
				y, err := strconv.Atoi(strings.Split(fields[i], "/")[0])
				check(err)
				verts = append(verts, y-1)

				// push z
				z, err := strconv.Atoi(strings.Split(fields[i+1], "/")[0])
				check(err)
				verts = append(verts, z-1)

				tInd = append(tInd, verts)
			}
		default:
			continue // skip textures, normals, etc. for now
		}
	}

	return &IndexedMesh{
		TrianglesToIdxs: tInd,
		Verts:           verts,
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (m *IndexedMesh) GetTrianglesFromMesh(material shading.Material) []*Triangle {
	var triangles []*Triangle
	for _, mapping := range m.TrianglesToIdxs {
		t := Triangle{
			V0:       m.Verts[mapping[0]],
			V1:       m.Verts[mapping[1]],
			V2:       m.Verts[mapping[2]],
			Material: material,
		}
		triangles = append(triangles, &t)
	}

	return triangles
}
