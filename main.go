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

// image quality settings

const samples = 500  // number of samples per pixel
const max_depth = 50 // maximum number of bounces

func makeTripleSphereWorld() (hittable, cameraParameters) {

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

	parameters := cameraParameters{
		imageWidth:  640, // 256,
		imageHeight: 480, // 144,

		vFov:     30,
		lookFrom: makeVec3(-2, 2, 1),
		lookAt:   makeVec3(0, 0, -1),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 10,
		focusDist:    3.4,
	}

	return world, parameters
}

func makeBook1CoverWorld() (hittable, cameraParameters) {

	world := makeHittableList()

	// ground
	world.add(makeSphere(makeVec3(0, -1000, 0), 1000, lambertian{makeVec3(0.5, 0.5, 0.5)}))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand()
			center := makeVec3(float64(a)+0.9*rand(), 0.2, float64(b)+0.9*rand())

			if center.subtract(makeVec3(4, 0.2, 0)).length() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					world.add(makeSphere(center, 0.2, lambertian{randVec().multiply(randVec())}))
				} else if chooseMat < 0.95 {
					// metal
					world.add(makeSphere(center, 0.2, metal{randVecB(0.5, 1), randb(0, 0.5)}))
				} else {
					// glass
					world.add(makeSphere(center, 0.2, dielectric{1.5}))
				}
			}
		}
	}

	world.add(makeSphere(makeVec3(0, 1, 0), 1.0, dielectric{1.5}))
	world.add(makeSphere(makeVec3(-4, 1, 0), 1.0, lambertian{makeVec3(0.4, 0.2, 0.1)}))
	world.add(makeSphere(makeVec3(4, 1, 0), 1.0, metal{makeVec3(0.7, 0.6, 0.5), 0.0}))

	parameters := cameraParameters{
		imageWidth:  640, // 256,
		imageHeight: 480, // 144,

		vFov:     30,
		lookFrom: makeVec3(13, 2, 3),
		lookAt:   makeVec3(0, 0, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.6,
		focusDist:    10.0,
	}

	return world, parameters
}

func renderScene() *render {

	world, parameters := makeBook1CoverWorld()
	camera := makeCamera(parameters)
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

// const imageWidth = 256
// const aspectRatio = 16.0 / 9.0

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
