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

const samples = 500

func renderScene(imageWidth int, aspectRatio float64) *render {

	material_ground := lambertian{makeVec3(0.8, 0.8, 0.0)}

	material_center := lambertian{makeVec3(0.1, 0.2, 0.5)}
	material_left := dielectric{1.5}
	material_right := metal{makeVec3(0.8, 0.6, 0.2), 0.0}

	world := makeHittableList()
	world.add(makeSphere(makeVec3(0, -100.5, -1), 100, material_ground))
	world.add(makeSphere(makeVec3(0, 0, -1), 0.5, material_center))
	world.add(makeSphere(makeVec3(-1, 0, -1), 0.5, material_left))
	world.add(makeSphere(makeVec3(-1, 0, -1), -0.4, material_left))
	world.add(makeSphere(makeVec3(1, 0, -1), 0.5, material_right))

	camera := makeCamera(imageWidth, aspectRatio, 1000)

	render := camera.render(world)

	if threaded {

		fmt.Printf("Rendering a %d x %d image with %d samples per pixel, using %d threads\n", camera.imageWidth, camera.imageHeight, samples, threads)

		for i := 0; i < threads; i++ {
			go render.run(max(1, samples/threads), i)
		}
	} else {

		fmt.Printf("Rendering a %d x %d image with %d samples per pixel\n", camera.imageWidth, camera.imageHeight, samples)

		if renderAheadOfDisplay {
			render.run(samples, 0)
		} else {
			go render.run(samples, 0)
		}
	}

	return render
}

const imageWidth = 256
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
