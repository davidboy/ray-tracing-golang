package main

import (
	"fmt"
	"os"
)

func rayColor(r Ray) Vec3 {
	unit_direction := r.direction.UnitVector()
	a := 0.5 * (unit_direction.Y() + 1.0)

	white := MakeVec3(1.0, 1.0, 1.0)
	blue := MakeVec3(0.5, 0.7, 1.0)

	return white.MultiplyScalar(1.0 - a).Add(blue.MultiplyScalar(a))
}

func main() {

	// Image

	aspect_ratio := 16.0 / 9.0

	image_width := 400
	image_height := int(float64(image_width) / aspect_ratio)

	// Camera

	focal_length := 1.0
	viewport_height := 2.0
	viewport_width := viewport_height * (float64(image_width) / float64(image_height))
	camera_center := MakeVec3(0, 0, 0)

	viewport_u := MakeVec3(viewport_width, 0, 0)
	viewport_v := MakeVec3(0, viewport_height, 0)

	pixel_delta_u := viewport_u.DivideScalar(float64(image_width))
	pixel_delta_v := viewport_v.DivideScalar(float64(image_height))

	viewport_upper_left := camera_center.
		Subtract(MakeVec3(0, 0, focal_length)).
		Subtract(viewport_u.DivideScalar(2)).
		Subtract(viewport_v.DivideScalar(2))

	pixel00_loc := viewport_upper_left.Add(pixel_delta_u.Add(pixel_delta_v).MultiplyScalar(0.5))

	fmt.Printf("P3\n%d %d\n255\n", image_width, image_height)

	for y := image_height - 1; y >= 0; y-- {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", y)

		for x := 0; x < image_width; x++ {
			pixel_center := pixel00_loc.
				Add(pixel_delta_u.MultiplyScalar(float64(x))).
				Add(pixel_delta_v.MultiplyScalar(float64(y)))

			ray_direction := pixel_center.Subtract(camera_center)

			ray := MakeRay(camera_center, ray_direction)

			pixel_color := rayColor(ray)
			pixel_color.WriteAsColor()
		}
	}

	fmt.Fprintf(os.Stderr, "\rDone.                 \n")

	mv := Vec3{[3]float64{1.0, 2.0, 3.0}}

	mv.Negate()
}
