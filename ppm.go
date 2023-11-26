package main

import (
	"fmt"
)

func writePpm(imageWidth int, aspectRatio float64, imagePixels []Vec3) {
	imageHeight := calculateImageHeight(imageWidth, aspectRatio)

	fmt.Printf("P3\n%d %d\n255\n", imageWidth, imageHeight)

	for y := imageHeight - 1; y >= 0; y-- {
		for x := 0; x < imageWidth; x++ {
			index := getPixelIndex(x, y, imageWidth)
			imagePixels[index].writeAsColor()
		}
	}
}

func (color Vec3) writeAsColor() {
	ir := int(255.999 * color.e[0])
	ig := int(255.999 * color.e[1])
	ib := int(255.999 * color.e[2])
	fmt.Printf("%d %d %d\n", ir, ig, ib)
}
