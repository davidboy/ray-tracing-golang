package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

var quality = qualityParameters{

	// TODO: move to render parameters?
	imageWidth:  150 * 21 / 9, // 640,
	imageHeight: 150,          // 480,

	samples:  80,
	maxDepth: 25,
	dof:      false,

	// TODO: support locked time / no motion blur
}

func main() {
	// TODO: move to render parameters
	var threads = flag.Int("threads", max(1, runtime.NumCPU()-2), "number of render threads to use")

	var outputPpm = flag.String("ppm", "", "output image as a PPM file")
	var outputPng = flag.String("png", "", "output image as a PNG file")

	var renderImage = flag.Bool("gui", true, "display the image as it is rendered")

	flag.Parse()

	world, parameters := makeBook1CoverScene()
	world, parameters := makeBook1CoverSceneWithMotion()
	render := createRender(world, parameters, *threads)

	fmt.Printf("Rendering a %d x %d image with %d samples per pixel, using %d threads\n", render.c.imageWidth, render.c.imageHeight, quality.samples, *threads)

	go func() {
		<-render.finished
		fmt.Println("Render complete.")
	}()

	if *renderImage {
		game := makeGame(render)

		ebiten.SetWindowSize(game.render.c.imageWidth, game.render.c.imageHeight)
		ebiten.SetWindowTitle("Raytracer")

		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		ebiten.MaximizeWindow()

		if err := ebiten.RunGame(game); err != nil {
			log.Fatal(err)
		}
	}

	<-render.finished

	imagePixels := render.squash()

	if *outputPpm != "" {
		f, err := os.Create(*outputPpm)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		writePpm(render.c.imageWidth, render.c.imageHeight, imagePixels, f)
	}

	if *outputPng != "" {
		f, err := os.Create(*outputPng)

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		writePng(render.c.imageWidth, render.c.imageHeight, imagePixels, f)
	}
}
