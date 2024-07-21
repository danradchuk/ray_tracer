package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/danradchuk/raytracer/dsl"
	"github.com/danradchuk/raytracer/geometry"
	"github.com/danradchuk/raytracer/shading"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func main() {
	var (
		width   = flag.Int("width", 1366, "width of the picture in pixels")
		height  = flag.Int("height", 768, "height of the picture in pixels")
		fov     = flag.Int("fov", 90, "field of view")
		input   = flag.String("input", "teapot.obj", "a mesh of an object to render")
		output  = flag.String("output", "image", "image to render")
		imgType = flag.String("type", "gif", "ppm or gif")
		world   = flag.String("scene-file", "./scenes/empty.scene", "file for constructing the scene")
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
	content, err := os.ReadFile(*world)
	if err != nil {
		log.Fatal(err)
	}

	p := dsl.NewParser(string(content))
	s, err := p.Parse()
	if err != nil {
		log.Fatal(err)
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
	split := strings.Split(*output, ".")

	fileName := ""
	if len(split) == 2 {
		fileName = split[0] + "." + *imgType
	} else {
		fileName = *output + "." + *imgType
	}

	if *imgType == "ppm" {
		err := s.RenderPPM(*width, *height, *fov, fileName)
		if err != nil {
			log.Fatal(err)
		}
	} else if *imgType == "gif" {
		err := s.RenderGIF(*width, *height, *fov, fileName)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Unknown image type")
	}
}
