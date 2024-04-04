package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const USE_BHV = true

var quality = qualityParameters{

	// TODO: move to render parameters?
	imageWidth:  400,          //  150 * 21 / 9, // 640,
	imageHeight: 400 / 16 * 9, // 150,          // 480,

	samples:  256 * 2 * 2 * 2 * 2,
	maxDepth: 50,
	dof:      false,

	// TODO: support locked time / no motion blur
}

func main() {
	// TODO: move to render parameters
	var threads = flag.Int("threads", max(1, runtime.NumCPU()-2), "number of render threads to use")

	var outputPpm = flag.String("ppm", "", "output image as a PPM file")
	var outputPng = flag.String("png", "", "output image as a PNG file")

	var renderImage = flag.Bool("gui", false, "display the image as it is rendered")

	flag.Parse()

	world, parameters := makeCornellBoxScene()

	var mainHittable hittable

	if USE_BHV {
		mainHittable = makeBvhNodeFromList(world)
	} else {
		mainHittable = world
	}

	startedAt := time.Now()
	render := startRendering(mainHittable, parameters, *threads)

	fmt.Printf("Rendering a %d x %d image with %d samples per pixel, using %d threads\n", render.c.imageWidth, render.c.imageHeight, quality.samples, *threads)

	go func() {
		<-render.finished

		finishedAt := time.Now()

		fmt.Println("Render complete.")

		elapsed := finishedAt.Sub(startedAt)
		samplesTaken := render.getTotalSamplesTaken()

		fmt.Printf("Rendered %d samples in %s (%.2f samples per second)\n", samplesTaken, elapsed, float64(samplesTaken)/elapsed.Seconds())

	}()

	if *renderImage {
		game := makeGame(render)

		ebiten.SetTPS(15)

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

const PIXELS_PER_CHUNK = 128

func startRendering(world hittable, parameters cameraParameters, threads int) *render {

	camera := makeCamera(parameters, &quality)
	render := makeRender(camera, world)

	tasks := make([](chan bool), threads)

	////////////////

	var mu sync.Mutex
	var currentPixelIndex int
	finalPixelIndex := getPixelIndex(camera.imageWidth, camera.imageHeight-1, camera.imageWidth)

	////////////////

	go (func() {
		desiredRandomPixels := camera.imageWidth * camera.imageHeight * 2

		for randomPixelsGenerated := 0; randomPixelsGenerated < desiredRandomPixels; randomPixelsGenerated++ {
			randomIndex := int(randb(0, float64(finalPixelIndex)))
			x, y := getPixelCoordinates(randomIndex, camera.imageWidth)

			render.runSinglePixel(int(x), int(y), 1)
		}
	})()

	////////////////

	for i := 0; i < threads; i++ {
		taskFinished := make(chan bool)
		tasks[i] = taskFinished

		go (func() {

			for {
				mu.Lock()

				currentLoopPixelIndex := currentPixelIndex
				currentPixelIndex += PIXELS_PER_CHUNK

				mu.Unlock()

				for pixelIndex := currentLoopPixelIndex; pixelIndex < min(currentLoopPixelIndex+PIXELS_PER_CHUNK, finalPixelIndex); pixelIndex++ {
					x, y := getPixelCoordinates(pixelIndex, camera.imageWidth)
					render.runSinglePixel(x, y, quality.samples)
				}

				if (currentLoopPixelIndex + PIXELS_PER_CHUNK) >= finalPixelIndex {
					close(taskFinished)
					return
				}
			}
		})()
	}

	////////////////

	// for i := 0; i < threads; i++ {
	// 	taskFinished := make(chan bool)
	// 	tasks[i] = taskFinished

	// 	go render.runWholeImage(max(1, quality.samples/threads), taskFinished)
	// }

	////////////////

	go func() {
		for _, taskFinished := range tasks {
			<-taskFinished
		}

		close(render.finished)
	}()

	return render
}
