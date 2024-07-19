package core

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"

	"github.com/danradchuk/raytracer/geometry"
	"github.com/danradchuk/raytracer/shading"
)

const MaxDepth = 3
const Bias = 0.0000001

type Camera struct {
	LookAt      geometry.Vec3
	LookFrom    geometry.Vec3 // eye
	Up          geometry.Vec3
	Right       geometry.Vec3
	AspectRatio float64
	FocalLength float64
}

type ImagePlane struct {
	X1 geometry.Vec3
	X2 geometry.Vec3
	X3 geometry.Vec3
	X4 geometry.Vec3
}

type Light struct {
	Pos               geometry.Vec3
	DiffuseIntensity  shading.Color
	SpecularIntensity shading.Color
}

type Scene struct {
	Background       shading.Color
	Light            *Light
	AmbientIntensity shading.Color
	Camera           geometry.Vec3
	ViewPort         geometry.Vec3
	Primitives       []geometry.Primitive
	AccelBVH         *geometry.BVHNode
}

func (s *Scene) CreatePPM(width, height int, fov int, outputFile string) error {
	// create an image file
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// write the format of the file
	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", width, height)
	if err != nil {
		return err
	}

	// visibility + shading
	frameBuffer := createFrameBuffer(width, height)
	numWorkers := runtime.NumCPU()
	chunk := height / numWorkers
	job := func(wg *sync.WaitGroup, startY, endY int) error {
		defer wg.Done()

		for y := startY; y < endY; y++ {
			for x := 0; x < width; x++ {
				r := geometry.NewPrimaryRay(s.Camera, float64(width), float64(height), float64(x), float64(y), fov)
				color := s.castRay(r, 0).ToImageColor()
				frameBuffer[x][y] = color
			}
		}

		return nil
	}

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		startY := i * chunk
		endY := (i + 1) * chunk
		if i == numWorkers-1 {
			endY = height
		}
		wg.Add(1)
		// fmt.Printf("%d worker has started. width %d, height %d. begin %d, end %d\n", i+1, width, height, startY, endY-1)
		go job(&wg, startY, endY)
	}

	wg.Wait()

	// fill in the .ppm file from the frame buffer
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			color := frameBuffer[x][y]
			_, err := fmt.Fprintf(f, "%d %d %d\n", color.R, color.G, color.B)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Scene) castRay(ray geometry.Ray, depth int) shading.Color {
	// stop recursion
	if depth >= MaxDepth {
		return s.Background
	}

	var closestT float64
	var closestPrimitive geometry.Primitive
	var material shading.Material

	hitRecord := s.AccelBVH.Intersect(ray)
	if hitRecord == nil {
		return s.Background
	}
	closestT = hitRecord.T
	closestPrimitive = hitRecord.Primitive
	material = hitRecord.Material

	// Phong Illumination Model
	hitPoint := ray.At(closestT)
	hitNormal := hitRecord.Normal
	viewDir := s.Camera.Sub(hitPoint).Normalize() // vector from the eye to the hitPoint
	lightDistance := s.Light.Pos.Sub(hitPoint).Norm()
	lightDir := s.Light.Pos.Sub(hitPoint).Normalize()

	// 1. compute the reflection component
	var reflective shading.Color
	var reflectionDir = reflect(ray.Direction.Normalize(), hitNormal).Normalize()
	var reflectionRayOrig geometry.Vec3
	// avoid self-reflecting
	if hitNormal.Dot(reflectionDir) < .0 {
		reflectionRayOrig = hitPoint.Sub(hitNormal.Scale(Bias))
	} else {
		reflectionRayOrig = hitPoint.Add(hitNormal.Scale(Bias))
	}
	reflectionRay := geometry.NewSecondaryRay(reflectionRayOrig, reflectionDir)
	reflective = s.castRay(reflectionRay, depth+1).Mul(material.KReflection)

	// 2. compute the shadow component
	shadowIntensity := 1.
	var shadowRayOrig geometry.Vec3
	// avoid self-shadowing
	if hitNormal.Dot(lightDir) < .0 {
		shadowRayOrig = hitPoint.Sub(hitNormal.Scale(Bias))
	} else {
		shadowRayOrig = hitPoint.Add(hitNormal.Scale(Bias))
	}
	shadowRay := geometry.NewSecondaryRay(shadowRayOrig, lightDir)
	for _, obj := range s.Primitives {
		if obj != closestPrimitive {
			hitRecord := obj.Intersect(shadowRay)
			if hitRecord != nil && hitRecord.T > .0 && hitRecord.T < lightDistance {
				shadowIntensity = .0
				break
			}
		}
	}

	// 3. compute diffuse, specular, and ambient components
	dot := hitNormal.Dot(lightDir) // when dot < .0 then a primitive points away from the light
	r := hitNormal.Scale(2 * math.Max(.0, dot)).Sub(lightDir)

	ambient := s.AmbientIntensity.Mul(material.KAmbient)
	diffuse := s.Light.DiffuseIntensity.Mul(material.KDiffuse).MulByNum(math.Max(.0, dot)).MulByNum(shadowIntensity)
	specular := s.Light.SpecularIntensity.Mul(material.KSpecular).MulByNum(math.Pow(math.Max(.0, viewDir.Dot(r)), material.Alpha)).MulByNum(shadowIntensity)

	return ambient.Add(diffuse).Add(specular).Add(reflective)
}

func createFrameBuffer(w, h int) [][]shading.ImageColor {
	frameBuffer := make([][]shading.ImageColor, w)
	for i := range frameBuffer {
		frameBuffer[i] = make([]shading.ImageColor, h)
	}

	return frameBuffer
}

func reflect(V geometry.Vec3, N geometry.Vec3) geometry.Vec3 {
	return V.Sub(N.Scale(2. * V.Dot(N)))
}
