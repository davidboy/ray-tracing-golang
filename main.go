package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func raycastScene(imageWidth int, aspectRatio float64, updateProgress func(percentageComplete float64)) []Vec3 {

	world := makeHittableList()
	world.add(makeSphere(MakeVec3(0, 0, -1), 0.5))
	world.add(makeSphere(MakeVec3(0, -100.5, -1), 100))

	camera := makeCamera(imageWidth, aspectRatio)
	camera.onUpdateProgress = updateProgress // TODO

	return camera.render(world)
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

	game := makeGame()

	game.SetPixels(&imagePixels)

	ebiten.SetWindowSize(game.imageWidth, game.imageHeight)
	ebiten.SetWindowTitle("Raytracer")

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// ppmMain()
	ebitenMain()
}
