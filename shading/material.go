package shading

var Glass = Material{
	KAmbient:    Color{R: 0.1, G: 0.1, B: 0.1},
	KDiffuse:    Color{R: 0.3, G: 0.3, B: 0.3},
	KSpecular:   Color{R: 0.7, G: 0.7, B: 0.7},
	KReflection: Color{R: 0.9, G: 0.9, B: 0.9},
	Alpha:       12500,
}

var Ivory = Material{
	KAmbient:    Color{R: 0.4, G: 0.4, B: 0.35},
	KDiffuse:    Color{R: 0.6, G: 0.6, B: 0.5},
	KSpecular:   Color{R: 0.7, G: 0.7, B: 0.7},
	KReflection: Color{R: 0.2, G: 0.2, B: 0.2},
	Alpha:       125.0,
}
var RedRubber = Material{
	KAmbient:    Color{R: 0.3, G: 0.0, B: 0.0},
	KDiffuse:    Color{R: 0.9, G: 0.1, B: 0.0},
	KSpecular:   Color{R: 0.3, G: 0.3, B: 0.3},
	KReflection: Color{R: 0.1, G: 0.1, B: 0.1},
	Alpha:       10.0,
}

// Material represents the properties of a material used in rendering.
// It includes ambient, diffuse, specular, and reflection constants,
// as well as an alpha value for the Phong model.
type Material struct {
	KAmbient    Color
	KDiffuse    Color
	KSpecular   Color
	KReflection Color
	Alpha       float64
}
