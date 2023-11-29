package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

func renderScene(imageWidth int, aspectRatio float64) *render {

	samples := 1000
	threads := runtime.NumCPU()

	world := makeHittableList()
	world.add(makeSphere(makeVec3(0, -100.5, -1), 100)) // ground
	world.add(makeSphere(makeVec3(0, 0, -1), 0.5))
	world.add(makeSphere(makeVec3(-1, 0, -1), 0.5))
	world.add(makeSphere(makeVec3(1, 0, -1), 0.5))

	camera := makeCamera(imageWidth, aspectRatio, 1000)

	fmt.Printf("Rendering a %d x %d image with %d samples per pixel, using %d threads\n", camera.imageWidth, camera.imageHeight, samples, threads)

	render := camera.render(world)

	for i := 0; i < threads; i++ {
		go render.run(max(1, samples/threads), i)
	}

	return render
}

const imageWidth = 320
const aspectRatio = 16.0 / 9.0

func ppmMain() {
	render := renderScene(imageWidth, aspectRatio)

	// TODO: progress reporting

	// func(percentageComplete float64) {
	// 	fmt.Fprintf(os.Stderr, "\rPercentage complete: %.2f", percentageComplete*100)
	// }

	imagePixels := render.squash()

	writePpm(imageWidth, aspectRatio, imagePixels)
}

func ebitenMain() {
	render := renderScene(imageWidth, aspectRatio)

	// TODO: progress reporting

	game := makeGame(render)

	ebiten.SetWindowSize(game.render.c.imageWidth, game.render.c.imageHeight)
	ebiten.SetWindowTitle("Raytracer")

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.MaximizeWindow()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// ppmMain()
	ebitenMain()
}
