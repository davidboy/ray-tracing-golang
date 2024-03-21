package main

import "math"

type texture interface {
	value(u, v float64, p vec3) vec3
}

type solidColorTexture struct {
	color vec3
}

func makeColorTexture(r, g, b float64) solidColorTexture {
	return solidColorTexture{makeVec3(r, g, b)}
}

func (s solidColorTexture) value(u, v float64, p vec3) vec3 {
	return s.color
}

type checkerTexture struct {
	invScale  float64
	even, odd texture
}

func makeCheckerTexture(scale float64, even, odd texture) checkerTexture {
	return checkerTexture{1 / scale, even, odd}
}

func (c checkerTexture) value(u, v float64, p vec3) vec3 {
	xInteger := int(math.Floor(c.invScale * p.x()))
	yInteger := int(math.Floor(c.invScale * p.y()))
	zInteger := int(math.Floor(c.invScale * p.z()))

	isEven := (xInteger+yInteger+zInteger)%2 == 0

	if isEven {
		return c.even.value(u, v, p)
	} else {
		return c.odd.value(u, v, p)
	}
}

func makeNoiseTexture(scale float64) noiseTexture {
	return noiseTexture{makePerlin(), scale}
}

type noiseTexture struct {
	noise perlin
	scale float64
}

func (n noiseTexture) value(u, v float64, p vec3) vec3 {
	return makeVec3(1, 1, 1).multiplyScalar(n.noise.noise(p.multiplyScalar(n.scale)))
}
