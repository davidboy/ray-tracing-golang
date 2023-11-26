package main

import (
	"fmt"
	"os"
)

func main() {
	image_width := 256
	image_height := 256

	fmt.Printf("P3\n%d %d\n255\n", image_width, image_height)

	for y := image_height - 1; y >= 0; y-- {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", y)

		for x := 0; x < image_width; x++ {
			r := float64(x) / float64(image_width-1)
			g := float64(y) / float64(image_height-1)
			b := 0.0

			MakeVec3(r, g, b).WriteAsColor()
		}
	}

	fmt.Fprintf(os.Stderr, "\rDone.                 \n")

	mv := Vec3{[3]float64{1.0, 2.0, 3.0}}

	mv.Negate()
}
