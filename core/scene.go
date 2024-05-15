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

var frameBuffer [1366][768]shading.ImageColor

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

type SceneObject interface {
	GetMaterial() shading.Material
	Intersection(r geometry.Ray) (bool, float64)
	NormalAt(p geometry.Vec3) geometry.Vec3
}

type Scene struct {
	Background       shading.Color
	Light            *Light
	AmbientIntensity shading.Color
	Camera           geometry.Vec3
	ViewPort         geometry.Vec3
	Objects          []SceneObject
}

func (s *Scene) CreatePPM(width, height int) error {
	// create an image file
	f, err := os.Create("circle.ppm")
	if err != nil {
		return err
	}
	defer f.Close()

	// write the format of the file
	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", width, height)
	if err != nil {
		return err
	}

	// visibility + shading
	numWorkers := runtime.NumCPU()
	chunk := height / numWorkers
	job := func(wg *sync.WaitGroup, startY, endY int) error {
		defer wg.Done()

		for y := startY; y < endY; y++ {
			for x := 0; x < width; x++ {
				r := geometry.NewPrimaryRay(s.Camera, float64(width), float64(height), float64(x), float64(y))
				color := castRay(s, r, 0).ToImageColor()
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

func castRay(scene *Scene, ray geometry.Ray, depth int) shading.Color {
	// stop recursion
	if depth >= MaxDepth {
		return scene.Background
	}

	var closestT = math.Inf(1)
	var closestObj interface{}

	// choose the closest intersection and the closest intersection object
	for _, o := range scene.Objects {
		if ok, t := o.Intersection(ray); ok {
			if t < closestT {
				closestT, closestObj = t, o
			}
		}
	}

	// Phong Illumination Model
	if o, ok := closestObj.(SceneObject); ok {
		hitPoint := ray.At(closestT)
		hitNormal := o.NormalAt(hitPoint)
		viewDir := scene.Camera.Sub(hitPoint).Normalize() // vector from the eye to the hitPoint
		lightDistance := scene.Light.Pos.Sub(hitPoint).Norm()
		lightDir := scene.Light.Pos.Sub(hitPoint).Normalize()

		m := o.GetMaterial()

		// compute the reflection component
		var reflective shading.Color
		reflectionDir := reflect(ray.Direction.Normalize(), hitNormal).Normalize()

		var reflectionRayOrig geometry.Vec3
		// avoid self-reflecting
		if hitNormal.Dot(reflectionDir) < .0 {
			reflectionRayOrig = hitPoint.Sub(hitNormal.Scale(Bias))
		} else {
			reflectionRayOrig = hitPoint.Add(hitNormal.Scale(Bias))
		}

		reflectionRay := geometry.NewSecondaryRay(reflectionRayOrig, reflectionDir)
		reflective = castRay(scene, reflectionRay, depth+1).Mul(m.KReflection)

		// compute the shadow component
		var shadowRayOrig geometry.Vec3
		// avoid self-shadowing
		if hitNormal.Dot(lightDir) < .0 {
			shadowRayOrig = hitPoint.Sub(hitNormal.Scale(Bias))
		} else {
			shadowRayOrig = hitPoint.Add(hitNormal.Scale(Bias))
		}

		shadowIntensity := 1.
		shadowRay := geometry.NewSecondaryRay(shadowRayOrig, lightDir)
		for _, obj := range scene.Objects {
			if obj != closestObj {
				if ok, t := obj.Intersection(shadowRay); ok {
					if t > .0 && t < lightDistance {
						shadowIntensity = .0
						break
					}
				}
			}
		}

		// compute diffuse, specular, and ambient components
		dot := hitNormal.Dot(lightDir)
		r := hitNormal.Scale(2 * math.Max(.0, dot)).Sub(lightDir)

		ambient := scene.AmbientIntensity.Mul(m.KAmbient)
		// when dot < .0 then an object points away from the light
		diffuse := scene.Light.DiffuseIntensity.Mul(m.KDiffuse).MulByNum(math.Max(.0, dot)).MulByNum(shadowIntensity)
		specular := scene.Light.SpecularIntensity.Mul(m.KSpecular).MulByNum(math.Pow(math.Max(.0, viewDir.Dot(r)), m.Alpha)).MulByNum(shadowIntensity)

		return ambient.Add(diffuse).Add(specular).Add(reflective)
	}

	return scene.Background
}

func reflect(V geometry.Vec3, N geometry.Vec3) geometry.Vec3 {
	return V.Sub(N.Scale(2. * V.Dot(N)))
}
