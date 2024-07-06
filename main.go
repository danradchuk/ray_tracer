package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/danradchuk/raytracer/core"
	"github.com/danradchuk/raytracer/geometry"
	"github.com/danradchuk/raytracer/shading"
)

const width = 1366
const height = 768
const R = 15

func main() {
	// cpu profile
	file, _ := os.Create("./cpu.pprof")
	err := pprof.StartCPUProfile(file)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// memory profile
	memProf, _ := os.Create("./mem.pprof")
	defer pprof.Lookup("allocs").WriteTo(memProf, 0)
	defer runtime.GC()

	// construct the scene
	s := core.Scene{
		Background:       shading.Color{R: 0.1, G: 0.3, B: 0.3},
		Camera:           geometry.Vec3{X: 0., Y: 0., Z: -10.},
		AmbientIntensity: shading.Color{R: 0.1, G: 0.1, B: 0.1},
		Light: &core.Light{
			Pos:               geometry.Vec3{X: 0., Y: 30., Z: -10.},
			DiffuseIntensity:  shading.Color{R: 0.8, G: 0.8, B: 0.8},
			SpecularIntensity: shading.Color{R: 0.8, G: 0.8, B: 0.8},
		},
		Primitives: []geometry.Primitive{
			/* &geometry.Sphere{
				Center:   geometry.Vec3{X: 0., Y: 0., Z: 25},
				R:        R,
				Material: shading.Glass,
			},
			&geometry.Sphere{
				Center:   geometry.Vec3{X: 70, Y: -40, Z: 10},
				R:        10,
				Material: shading.RedRubber,
			},
			&geometry.Sphere{
				Center:   geometry.Vec3{X: -70, Y: -40, Z: 10},
				R:        10,
				Material: shading.RedRubber,
			},
			&geometry.Sphere{
				Center:   geometry.Vec3{X: 0, Y: -40, Z: 10},
				R:        10,
				Material: shading.Ivory,
			},
			&geometry.Triangle{
				V0:       geometry.Vec3{X: -100.0, Y: 0.0, Z: 0.0},
				V1:       geometry.Vec3{X: -20.0, Y: 75.0, Z: 50.0},
				V2:       geometry.Vec3{X: 50.0, Y: 0.0, Z: 100.0},
				Material: shading.RedRubber,
			},
			&geometry.Plane{
				Width:    250.0,
				Point:    geometry.Vec3{X: 0, Y: -50, Z: 75},
				Normal:   geometry.Vec3{X: 0, Y: 1, Z: 0},
				Material: shading.Glass,
			}, */
		},
	}

	// load a triangle mesh
	mesh := geometry.LoadOBJ("teapot.obj")
	for _, t := range mesh.GetTrianglesFromMesh(shading.RedRubber) {
		s.Primitives = append(s.Primitives, t)
	}

	// build a BVH
	s.AccelBVH = geometry.BuildBVH(s.Primitives)

	// render image
	err = s.CreatePPM(width, height)
	if err != nil {
		log.Fatal(err)
	}
}
