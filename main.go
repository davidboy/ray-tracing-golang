package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func hitSphere(center Vec3, radius float64, r Ray) bool {
	oc := r.Origin.Subtract(center)
	a := r.Direction.Dot(r.Direction)
	b := 2.0 * oc.Dot(r.Direction)
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c
	return discriminant >= 0
}

func rayColor(r Ray) Vec3 {
	if hitSphere(MakeVec3(0, 0, -1), 0.5, r) {
		return MakeVec3(1, 0, 0)
	}

	unitDirection := r.Direction.UnitVector()
	a := 0.5 * (unitDirection.Y() + 1.0)

	white := MakeVec3(1.0, 1.0, 1.0)
	blue := MakeVec3(0.5, 0.7, 1.0)

	return white.MultiplyScalar(1.0 - a).Add(blue.MultiplyScalar(a))
}

func calculateImageHeight(imageWidth int, aspectRatio float64) int {
	return int(float64(imageWidth) / aspectRatio)
}

func getPixelIndex(x, y, imageWidth int) int {
	return y*imageWidth + x
}

func raycastScene(imageWidth int, aspectRatio float64, updateProgress func(percentageComplete float64)) []Vec3 {
	imageHeight := calculateImageHeight(imageWidth, aspectRatio)
	pixels := make([]Vec3, imageWidth*imageHeight)

	// Camera
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	cameraCenter := MakeVec3(0, 0, 0)

	viewportU := MakeVec3(viewportWidth, 0, 0)
	viewportV := MakeVec3(0, viewportHeight, 0)

	pixelDeltaU := viewportU.DivideScalar(float64(imageWidth))
	pixelDeltaV := viewportV.DivideScalar(float64(imageHeight))

	viewportUpperLeft := cameraCenter.
		Subtract(MakeVec3(0, 0, focalLength)).
		Subtract(viewportU.DivideScalar(2)).
		Subtract(viewportV.DivideScalar(2))

	pixel00Loc := viewportUpperLeft.Add(pixelDeltaU.Add(pixelDeltaV).MultiplyScalar(0.5))

	for x := 0; x < imageWidth; x++ {
		updateProgress(float64(x) / float64(imageWidth))

		for y := 0; y < imageHeight; y++ {
			pixelCenter := pixel00Loc.
				Add(pixelDeltaU.MultiplyScalar(float64(x))).
				Add(pixelDeltaV.MultiplyScalar(float64(y)))

			rayDirection := pixelCenter.Subtract(cameraCenter)
			ray := MakeRay(cameraCenter, rayDirection)

			pixels[getPixelIndex(x, y, imageWidth)] = rayColor(ray)
		}
	}

	return pixels
}

const imageWidth = 400
const aspectRatio = 16.0 / 9.0

func ppmMain() {
	imagePixels := raycastScene(imageWidth, aspectRatio, func(percentageComplete float64) {
		fmt.Fprintf(os.Stderr, "\rPercentage complete: %.2f", percentageComplete*100)
	})

	writePpm(imageWidth, aspectRatio, imagePixels)

	fmt.Fprintf(os.Stderr, "\rDone.                                       \n")
}

func ebitenMain() {
	imagePixels := raycastScene(imageWidth, aspectRatio, func(percentageComplete float64) {
		fmt.Fprintf(os.Stdout, "\rPercentage complete: %.2f", percentageComplete*100)
	})

	fmt.Fprintf(os.Stdout, "\rDone.                                       \n")

	game := &Game{
		imageWidth:  imageWidth,
		imageHeight: calculateImageHeight(imageWidth, aspectRatio),

		imagePixels: imagePixels,
	}

	ebiten.SetWindowSize(game.imageWidth, game.imageHeight)
	ebiten.SetWindowTitle("Raytracer")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// ppmMain()
	ebitenMain()
}
