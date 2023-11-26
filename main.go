package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type hitRecord struct {
	p         Vec3
	normal    Vec3
	t         float64
	frontFace bool
}

func (rec *hitRecord) setFaceNormal(r Ray, outwardNormal Vec3) {
	rec.frontFace = r.Direction.Dot(outwardNormal) < 0

	if rec.frontFace {
		rec.normal = outwardNormal
	} else {
		rec.normal = outwardNormal.Negate()
	}
}

type hittable interface {
	hit(r Ray, tMin float64, tMax float64, rec *hitRecord) bool
}

type hittableList struct {
	hittables []hittable
}

func makeHittableList() hittableList {
	return hittableList{make([]hittable, 0)}
}

func (l *hittableList) clear() {
	l.hittables = make([]hittable, 0)
}

func (l *hittableList) add(h hittable) {
	l.hittables = append(l.hittables, h)
}

func (l hittableList) hit(r Ray, tMin float64, tMax float64, rec *hitRecord) bool {
	var tempRec hitRecord

	hitAnything := false
	closestSoFar := tMax

	for _, hittable := range l.hittables {
		if hittable.hit(r, tMin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.t
			*rec = tempRec
		}
	}

	return hitAnything
}

func rayColor(r Ray, world hittable) Vec3 {
	var rec hitRecord

	if world.hit(r, 0.0, infinity, &rec) {
		return rec.normal.AddScalar(1.0).MultiplyScalar(0.5)
	}

	unitDirection := r.Direction.UnitVector()
	a := 0.5 * (unitDirection.Y() + 1.0)

	white := MakeVec3(1.0, 1.0, 1.0)
	blue := MakeVec3(0.5, 0.7, 1.0)

	return white.MultiplyScalar(1.0 - a).Add(blue.MultiplyScalar(a))
}

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
	// Image

	imageHeight := calculateImageHeight(imageWidth, aspectRatio)
	pixels := make([]Vec3, imageWidth*imageHeight)

	// World

	world := makeHittableList()
	world.add(makeSphere(MakeVec3(0, 0, -1), 0.5))
	world.add(makeSphere(MakeVec3(0, -100.5, -1), 100))

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

			pixels[getPixelIndex(x, y, imageWidth)] = rayColor(ray, world)
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
