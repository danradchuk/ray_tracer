package core

import (
	"fmt"
	"image"
	"image/color"
	palette2 "image/color/palette"
	"image/gif"
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

type Light struct {
	Pos               geometry.Vec3
	DiffuseIntensity  shading.Color
	SpecularIntensity shading.Color
}

type Scene struct {
	Background       shading.Color
	Lights           []*Light
	AmbientIntensity shading.Color
	Camera           geometry.Vec3
	Primitives       []geometry.Primitive
	AccelBVH         *geometry.BVHNode
}

func (s *Scene) RenderGIF(width, height, fov int, output string) error {
	cameraPosRotation := func(theta, y float64) geometry.Vec3 {
		r := 10.
		x := r * math.Cos(theta)
		z := r * math.Sin(theta)
		return geometry.Vec3{X: x, Y: y, Z: z}
	}

	var images []*image.Paletted
	var delays []int

	for i := 0; i < 360; i++ {
		img := image.NewPaletted(image.Rect(0, 0, width, height), palette2.Plan9)

		theta := float64(i) * (2 * math.Pi / 360) // Convert degrees to radians

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r := geometry.NewPrimaryRay(cameraPosRotation(theta, s.Camera.Y),
					float64(width),
					float64(height),
					float64(x),
					float64(y),
					float64(fov),
				)
				c := s.castRay(r, 0).ToImageColor()
				img.Set(x, y, color.RGBA{R: c.R, G: c.G, B: c.B, A: 0xFF})
			}
		}

		images = append(images, img)
		delays = append(delays, 0)
	}

	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = gif.EncodeAll(outFile, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Scene) RenderPPM(width, height int, fov int, outputFile string) error {
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
				r := geometry.NewPrimaryRay(s.Camera, float64(width), float64(height), float64(x), float64(y), float64(fov))
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

	hitRecord := s.AccelBVH.Intersect(ray)
	if hitRecord == nil {
		return s.Background
	}

	closestT := hitRecord.T
	closestPrimitive := hitRecord.Primitive
	material := hitRecord.Material

	hitPoint := ray.At(closestT)
	hitNormal := hitRecord.Normal
	viewDir := s.Camera.Sub(hitPoint).Normalize() // vector from the eye to the hitPoint

	var (
		diffuseComponent    shading.Color
		specularComponent   shading.Color
		reflectionComponent shading.Color
	)

	// 1. compute reflection component
	var reflectionDir = reflect(ray.Direction.Normalize(), hitNormal).Normalize()
	var reflectionRayOrig geometry.Vec3
	// avoid self-reflecting
	if hitNormal.Dot(reflectionDir) < .0 {
		reflectionRayOrig = hitPoint.Sub(hitNormal.Scale(Bias))
	} else {
		reflectionRayOrig = hitPoint.Add(hitNormal.Scale(Bias))
	}
	reflectionRay := geometry.NewSecondaryRay(reflectionRayOrig, reflectionDir)
	reflectionComponent = s.castRay(reflectionRay, depth+1).Mul(material.KReflection)

	for _, light := range s.Lights {
		lightDir := light.Pos.Sub(hitPoint).Normalize()

		// 2. compute shadow component
		var shadowRayOrig geometry.Vec3
		// avoid self-shadowing
		if hitNormal.Dot(lightDir) < .0 {
			shadowRayOrig = hitPoint.Sub(hitNormal.Scale(Bias))
		} else {
			shadowRayOrig = hitPoint.Add(hitNormal.Scale(Bias))
		}
		shadowRay := geometry.NewSecondaryRay(shadowRayOrig, lightDir)
		lightDistance := light.Pos.Sub(hitPoint).Norm()

		shadowIntensity := 1.
		hitRecord := s.AccelBVH.IntersectExclude(shadowRay, closestPrimitive)
		if hitRecord != nil && hitRecord.T > .0 && hitRecord.T < lightDistance {
			shadowIntensity = .0
		}

		// 3. compute diffuse and specular components
		dot := math.Max(.0, hitNormal.Dot(lightDir)) // when dot < .0 then a primitive points away from the light
		r := hitNormal.Scale(2 * dot).Sub(lightDir)

		diffuseComponent = diffuseComponent.Add(light.DiffuseIntensity.Mul(material.KDiffuse).MulByNum(dot).MulByNum(shadowIntensity))
		specularComponent = specularComponent.Add(light.SpecularIntensity.Mul(material.KSpecular).MulByNum(math.Pow(math.Max(.0, viewDir.Dot(r)), material.Alpha)).MulByNum(shadowIntensity))
	}

	return s.AmbientIntensity.Mul(material.KAmbient).Add(diffuseComponent).Add(specularComponent).Add(reflectionComponent)
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
