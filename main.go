package main

import (
	"fmt"
	"os"
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

func main() {

	// Image

	aspectRatio := 16.0 / 9.0

	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)

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

	fmt.Printf("P3\n%d %d\n255\n", imageWidth, imageHeight)

	for y := imageHeight - 1; y >= 0; y-- {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", y)

		for x := 0; x < imageWidth; x++ {
			pixelCenter := pixel00Loc.
				Add(pixelDeltaU.MultiplyScalar(float64(x))).
				Add(pixelDeltaV.MultiplyScalar(float64(y)))

			rayDirection := pixelCenter.Subtract(cameraCenter)

			ray := MakeRay(cameraCenter, rayDirection)

			pixelColor := rayColor(ray)
			pixelColor.WriteAsColor()
		}
	}

	fmt.Fprintf(os.Stderr, "\rDone.                 \n")

	mv := Vec3{[3]float64{1.0, 2.0, 3.0}}

	mv.Negate()
}
