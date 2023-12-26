package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

var threaded = true
var renderAheadOfDisplay = false
var threads = max(1, runtime.NumCPU()-2)

var quality = qualityParameters{
	imageWidth:  640,
	imageHeight: 480,

	samples:  500,
	maxDepth: 50,
	dof:      true,
}

func renderScene() *render {
	world, parameters := makeBook1CoverScene()

	camera := makeCamera(parameters, &quality)
	render := camera.render(world)

	if threaded {

		fmt.Printf("Rendering a %d x %d image with %d samples per pixel, using %d threads\n", camera.imageWidth, camera.imageHeight, quality.samples, threads)

		for i := 0; i < threads; i++ {
			go render.run(max(1, quality.samples/threads), i)
		}
	} else {

		fmt.Printf("Rendering a %d x %d image with %d samples per pixel\n", camera.imageWidth, camera.imageHeight, quality.samples)

		if renderAheadOfDisplay {
			render.run(quality.samples, 0)
		} else {
			go render.run(quality.samples, 0)
		}
	}

	return render
}

func ppmMain() {
	render := renderScene()

	// TODO: progress reporting

	// func(percentageComplete float64) {
	// 	fmt.Fprintf(os.Stderr, "\rPercentage complete: %.2f", percentageComplete*100)
	// }

	imagePixels := render.squash()

	writePpm(render.c.imageWidth, render.c.aspectRatio, imagePixels)
}

func ebitenMain() {
	render := renderScene()

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
