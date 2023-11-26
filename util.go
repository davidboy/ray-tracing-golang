package main

import "math"

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
