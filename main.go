package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/danradchuk/raytracer/core"
	"github.com/danradchuk/raytracer/geometry"
	"github.com/danradchuk/raytracer/shading"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func main() {
	var (
		width        = flag.Int("width", 1366, "width of the picture in pixels")
		height       = flag.Int("height", 768, "height of the picture in pixels")
		fov          = flag.Int("fov", 90, "field of view")
		r            = flag.Float64("radius", 15., "radius of a sphere")
		input        = flag.String("input", "teapot.obj", "a mesh of an object to render")
		output       = flag.String("output", "image.ppm", "image to render")
		withGeometry = flag.Bool("geometry", false, "use geometry")
	)

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.Lookup("allocs").WriteTo(f, 0)
		defer runtime.GC()
	}

	// construct the scene
	s := core.Scene{
		Background:       shading.Color{R: 0.1, G: 0.3, B: 0.3},
		Camera:           geometry.Vec3{X: 10., Y: 15., Z: 10.},
		AmbientIntensity: shading.Color{R: 0.1, G: 0.1, B: 0.1},
		Lights: []*core.Light{
			{
				Pos:               geometry.Vec3{X: 0., Y: 30., Z: -10.},
				DiffuseIntensity:  shading.Color{R: 0.8, G: 0.8, B: 0.8},
				SpecularIntensity: shading.Color{R: 0.8, G: 0.8, B: 0.8},
			},
			{
				Pos:               geometry.Vec3{X: 30., Y: 30., Z: -10.},
				DiffuseIntensity:  shading.Color{R: 0.1, G: 0.1, B: 0.1},
				SpecularIntensity: shading.Color{R: 0.8, G: 0.8, B: 0.8},
			},
		},
		Primitives: []geometry.Primitive{
			&geometry.Plane{
				Width:    250.0,
				Point:    geometry.Vec3{X: 0, Y: -50, Z: 75},
				Normal:   geometry.Vec3{X: 0, Y: 1, Z: 0},
				Material: shading.Glass,
			},
		},
	}

	// shapes
	if *withGeometry {
		s.Primitives = append(s.Primitives, &geometry.Sphere{
			Center:   geometry.Vec3{X: 0., Y: 0., Z: 25},
			R:        *r,
			Material: shading.Glass,
		})

		s.Primitives = append(s.Primitives, &geometry.Sphere{
			Center:   geometry.Vec3{X: 70., Y: -40., Z: 10},
			R:        10,
			Material: shading.RedRubber,
		})

		s.Primitives = append(s.Primitives, &geometry.Sphere{
			Center:   geometry.Vec3{X: -70., Y: -40., Z: 10},
			R:        10,
			Material: shading.RedRubber,
		})

		s.Primitives = append(s.Primitives, &geometry.Sphere{
			Center:   geometry.Vec3{X: 0., Y: -40., Z: 10},
			R:        10,
			Material: shading.Ivory,
		})

		s.Primitives = append(s.Primitives, &geometry.Triangle{
			V0:       geometry.Vec3{X: -100.0, Y: 0.0, Z: 0.0},
			V1:       geometry.Vec3{X: -20.0, Y: 75.0, Z: 50.0},
			V2:       geometry.Vec3{X: 50.0, Y: 0.0, Z: 100.0},
			Material: shading.RedRubber,
		})
	}

	//load a triangle mesh
	if *input != "" {
		mesh := geometry.LoadOBJ(*input)
		for _, t := range mesh.GetTrianglesFromMesh(shading.RedRubber) {
			s.Primitives = append(s.Primitives, t)
		}
	}

	// build a BVH
	s.AccelBVH = geometry.BuildBVH(s.Primitives)

	// render image
	err := s.CreatePPM(*width, *height, *fov, *output)
	if err != nil {
		log.Fatal(err)
	}
}
