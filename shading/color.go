package shading

import "math"

const Max = 1.0

var Black = Color{0.0, 0.0, 0.0}

// A 8-bit RGB representation of a color
type ImageColor struct {
	R uint8
	G uint8
	B uint8
}

// A representation of a color where every
// component is in the range [0.0, 1.0]
type Color struct {
	R float64
	G float64
	B float64
}

func (c Color) Mul(other Color) Color {
	return Color{
		R: c.R * other.R,
		G: c.G * other.G,
		B: c.B * other.B,
	}
}

func (c Color) MulByNum(s float64) Color {
	return Color{
		R: c.R * s,
		G: c.G * s,
		B: c.B * s,
	}
}

func (c Color) Add(other Color) Color {
	return Color{
		R: c.R + other.R,
		G: c.G + other.G,
		B: c.B + other.B,
	}
}

func (c Color) Clamped() Color {
	var r float64
	var g float64
	var b float64

	if c.R < 0 {
		r = 0.0
	} else if c.R > 1 {
		r = 1.0
	} else {
		r = math.Min(Max, c.R)
	}

	if c.G < 0 {
		g = 0.0
	} else if c.G > 1 {
		g = 1.0
	} else {
		g = math.Min(Max, c.G)
	}

	if c.B < 0 {
		b = 0.0
	} else if c.B > 1 {
		b = 1.0
	} else {
		b = math.Min(Max, c.B)
	}

	return Color{r, g, b}
}

func (c Color) ToImageColor() ImageColor {
	return ImageColor{
		uint8(c.R*255 + 0.5),
		uint8(c.G*255 + 0.5),
		uint8(c.B*255 + 0.5),
	}
}
