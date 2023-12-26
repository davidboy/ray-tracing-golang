package main

import (
	"image"
	"image/png"
	"io"
	"log"
)

func writePng(imageWidth int, imageHeight int, imagePixels []vec3, w io.Writer) {
	image := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			index := getPixelIndex(x, y, imageWidth)
			color := imagePixels[index]

			image.Set(x, imageHeight-y, color)
		}
	}

	err := png.Encode(w, image)
	if err != nil {
		log.Fatal(err)
	}

}
