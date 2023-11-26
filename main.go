package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func renderScene(imageWidth int, aspectRatio float64) *render {

	world := makeHittableList()
	world.add(makeSphere(makeVec3(0, 0, -1), 0.5))
	world.add(makeSphere(makeVec3(0, -100.5, -1), 100))

	camera := makeCamera(imageWidth, aspectRatio, 1000)

	render := camera.render(world)
	go render.run()

	return render
}

const imageWidth = 800
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
