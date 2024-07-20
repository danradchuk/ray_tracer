package dsl

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/danradchuk/raytracer/core"
	"github.com/danradchuk/raytracer/geometry"
	"github.com/danradchuk/raytracer/shading"
)

type Parser struct {
	Words     []string
	currToken string
	peekToken string
	peekPos   int
}

func NewParser(content string) *Parser {
	tokens := strings.Fields(content)

	p := &Parser{Words: tokens}
	p.nextToken()

	return p
}

func (p *Parser) nextToken() error {
	if p.peekPos >= len(p.Words) {
		return io.EOF
	}
	p.currToken = p.peekToken
	p.peekToken = p.Words[p.peekPos]
	p.peekPos++

	return nil
}

func (p *Parser) Parse() (*core.Scene, error) {
	var scene = &core.Scene{}
	for p.nextToken() != io.EOF {
		switch p.currToken {
		case "background":
			tok := p.peekToken
			if !strings.HasPrefix(tok, "#") && len(tok) != 7 && !isASCII(tok) {
				return nil, fmt.Errorf("background: invalid color format: %s", tok)
			}

			p.nextToken() // consume hex number

			red, err := strconv.ParseInt(tok[1:3], 16, 64)
			if err != nil {
				return nil, err
			}

			green, err := strconv.ParseInt(tok[3:5], 16, 64)
			if err != nil {
				return nil, err
			}

			blue, err := strconv.ParseInt(tok[5:7], 16, 64)
			if err != nil {
				return nil, err
			}

			color := shading.Color{
				R: float64(red) / 255.,
				G: float64(green) / 255.,
				B: float64(blue) / 255.,
			}
			scene.Background = color
		case "ambient":
			c, err := parseColor(p.peekToken)
			if err != nil {
				return nil, err
			}
			scene.AmbientIntensity = *c
		case "light":
			tok := p.peekToken
			if tok != "{" {
				return nil, fmt.Errorf("unexpected character: %s", tok)
			}

			p.nextToken()

			var light = &core.Light{}
			for p.peekToken != "}" {
				switch p.peekToken {
				case "pos":
					p.nextToken()
					pos, err := parseVec(p.peekToken)
					if err != nil {
						return nil, err
					}
					light.Pos = *pos
				case "diffuse":
					p.nextToken()
					c, err := parseColor(p.peekToken)
					if err != nil {
						return nil, err
					}
					light.DiffuseIntensity = *c
				case "specular":
					p.nextToken()
					c, err := parseColor(p.peekToken)
					if err != nil {
						return nil, err
					}
					light.SpecularIntensity = *c
				}
				p.nextToken()
			}

			if p.peekToken != "}" {
				return nil, fmt.Errorf("unexpected character: %s", p.peekToken)
			}
			scene.Lights = append(scene.Lights, light)
		case "camera":
			eye, err := parseVec(p.peekToken)
			if err != nil {
				return nil, err
			}
			scene.Camera = *eye
			p.nextToken()
		case "sphere":
			tok := p.peekToken
			if tok != "{" {
				return nil, fmt.Errorf("unexpected character: %s", tok)
			}

			p.nextToken()

			var sphere = geometry.Sphere{}
			for p.peekToken != "}" {
				switch p.peekToken {
				case "radius":
					p.nextToken()
					r, err := strconv.ParseFloat(p.peekToken, 64)
					if err != nil {
						return nil, err
					}
					sphere.R = r
				case "center":
					p.nextToken()
					center, err := parseVec(p.peekToken)
					if err != nil {
						return nil, err
					}
					sphere.Center = *center
				case "material":
					p.nextToken()
					m := parseMaterial(p.peekToken)
					if m != nil {
						sphere.Material = *m
					}
				}
				p.nextToken()
			}

			if p.peekToken != "}" {
				return nil, fmt.Errorf("unexpected character: %s", p.peekToken)
			}
			scene.Primitives = append(scene.Primitives, sphere)
		case "triangle":
			tok := p.peekToken
			if tok != "{" {
				return nil, fmt.Errorf("unexpected character: %s", tok)
			}

			p.nextToken()

			var triangle = &geometry.Triangle{}
			for p.peekToken != "}" {
				switch p.peekToken {
				case "v0":
					p.nextToken()
					coords, err := parseVec(p.peekToken)
					if err != nil {
						return nil, err
					}
					triangle.V0 = *coords
				case "v1":
					p.nextToken()
					coords, err := parseVec(p.peekToken)
					if err != nil {
						return nil, err
					}
					triangle.V1 = *coords
				case "v2":
					p.nextToken()
					coords, err := parseVec(p.peekToken)
					if err != nil {
						return nil, err
					}
					triangle.V2 = *coords
				case "material":
					p.nextToken()
					m := parseMaterial(p.peekToken)
					if m != nil {
						triangle.Material = *m
					}
				}
				p.nextToken()
			}

			if p.peekToken != "}" {
				return nil, fmt.Errorf("unexpected character: %s", p.peekToken)
			}
			scene.Primitives = append(scene.Primitives, triangle)
		case "plane":
			tok := p.peekToken
			if tok != "{" {
				return nil, fmt.Errorf("unexpected character: %s", tok)
			}

			p.nextToken()

			var plane = geometry.Plane{}
			for p.peekToken != "}" {
				switch p.peekToken {
				case "width":
					p.nextToken()
					w, err := strconv.ParseFloat(p.peekToken, 64)
					if err != nil {
						return nil, err
					}
					plane.Width = w
				case "point":
					p.nextToken()
					point, err := parseVec(p.peekToken)
					if err != nil {
						return nil, err
					}
					plane.Point = *point
				case "normal":
					p.nextToken()
					n, err := parseVec(p.peekToken)
					if err != nil {
						return nil, err
					}
					plane.Normal = *n
				case "material":
					p.nextToken()
					m := parseMaterial(p.peekToken)
					if m != nil {
						plane.Material = *m
					}
				}
				p.nextToken()
			}

			if p.peekToken != "}" {
				return nil, fmt.Errorf("unexpected character: %s", p.peekToken)
			}
			scene.Primitives = append(scene.Primitives, plane)
		}
	}

	return scene, nil
}

func parseMaterial(token string) *shading.Material {
	if token == "red" {
		return &shading.RedRubber
	} else if token == "ivory" {
		return &shading.Ivory
	} else if token == "glass" {
		return &shading.Glass
	}

	return nil
}

func parseColor(token string) (*shading.Color, error) {
	vec := strings.Split(token, ",")
	r, err := strconv.ParseFloat(vec[0], 64)
	if err != nil {
		return nil, err
	}

	g, err := strconv.ParseFloat(vec[1], 64)
	if err != nil {
		return nil, err
	}

	b, err := strconv.ParseFloat(vec[2], 64)
	if err != nil {
		return nil, err
	}
	return &shading.Color{
		R: r,
		G: g,
		B: b,
	}, nil
}

func parseVec(token string) (*geometry.Vec3, error) {
	vec := strings.Split(token, ",")
	x, err := strconv.ParseFloat(vec[0], 64)
	if err != nil {
		return nil, err
	}

	y, err := strconv.ParseFloat(vec[1], 64)
	if err != nil {
		return nil, err
	}

	z, err := strconv.ParseFloat(vec[2], 64)
	if err != nil {
		return nil, err
	}
	return &geometry.Vec3{
		X: x,
		Y: y,
		Z: z,
	}, nil
}

func isASCII(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}
