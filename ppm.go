package main

import (
	"fmt"
	"io"
)

func writePpm(imageWidth int, imageHeight int, imagePixels []vec3, w io.Writer) {

	fmt.Fprintf(w, "P3\n%d %d\n255\n", imageWidth, imageHeight)

	for y := imageHeight - 1; y >= 0; y-- {
		for x := 0; x < imageWidth; x++ {
			index := getPixelIndex(x, y, imageWidth)
			color := imagePixels[index]

			ir := int(255.999 * color.e[0])
			ig := int(255.999 * color.e[1])
			ib := int(255.999 * color.e[2])

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
	}
}
