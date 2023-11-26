package main

import (
	"math"
	r "math/rand"
)

var infinity = math.Inf(1)

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func calculateImageHeight(imageWidth int, aspectRatio float64) int {
	result := int(float64(imageWidth) / aspectRatio)

	if result < 1 {
		result = 1
	}

	return result
}

func getPixelIndex(x, y, imageWidth int) int {
	return y*imageWidth + x
}

func rand() float64 {
	return r.Float64()
}

func randb(min, max float64) float64 {
	return min + (max-min)*rand()
}

func linearToGamma(linear float64) float64 {
	return math.Sqrt(linear)
}

func (v vec3) toGammaSpace() vec3 {
	return makeVec3(linearToGamma(v.x()), linearToGamma(v.y()), linearToGamma(v.z()))
}
