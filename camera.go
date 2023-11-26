package main

type camera struct {
	imageWidth  int
	imageHeight int
	aspectRatio float64

	center      vec3 // camera center
	pixel00Loc  vec3 // location of pixel (0, 0)
	pixelDeltaU vec3 // offset to pixel to the right
	pixelDeltaV vec3 // offset to pixel below

	onUpdateProgress func(progress float64)
}

func makeCamera(imageWidth int, aspectRatio float64) *camera {
	imageHeight := calculateImageHeight(imageWidth, aspectRatio)

	center := makeVec3(0, 0, 0)

	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))

	viewportU := makeVec3(viewportWidth, 0, 0)
	viewportV := makeVec3(0, viewportHeight, 0)

	pixelDeltaU := viewportU.divideScalar(float64(imageWidth))
	pixelDeltaV := viewportV.divideScalar(float64(imageHeight))

	viewportUpperLeft := center.
		subtract(makeVec3(0, 0, focalLength)).
		subtract(viewportU.divideScalar(2)).
		subtract(viewportV.divideScalar(2))

	pixel00Loc := viewportUpperLeft.add(pixelDeltaU.add(pixelDeltaV).multiplyScalar(0.5))

	return &camera{
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
		aspectRatio: aspectRatio,

		center:      center,
		pixel00Loc:  pixel00Loc,
		pixelDeltaU: pixelDeltaU,
		pixelDeltaV: pixelDeltaV,

		onUpdateProgress: func(progress float64) {},
	}

}

func (c *camera) render(world hittable) []vec3 {
	pixels := make([]vec3, c.imageWidth*c.imageHeight)

	for x := 0; x < imageWidth; x++ {
		c.onUpdateProgress(float64(x) / float64(imageWidth))

		for y := 0; y < c.imageHeight; y++ {
			pixelCenter := c.pixel00Loc.
				add(c.pixelDeltaU.multiplyScalar(float64(x))).
				add(c.pixelDeltaV.multiplyScalar(float64(y)))

			rayDirection := pixelCenter.subtract(c.center)
			ray := MakeRay(c.center, rayDirection)

			pixels[getPixelIndex(x, y, imageWidth)] = rayColor(ray, world)
		}
	}

	return pixels
}

var white = makeVec3(1.0, 1.0, 1.0)
var blue = makeVec3(0.5, 0.7, 1.0)

func rayColor(r Ray, world hittable) vec3 {
	var rec hitRecord

	if world.hit(r, makeInterval(0.0, infinity), &rec) {
		return rec.normal.addScalar(1.0).multiplyScalar(0.5)
	}

	unitDirection := r.Direction.unitVector()
	a := 0.5 * (unitDirection.y() + 1.0)

	return white.multiplyScalar(1.0 - a).add(blue.multiplyScalar(a))
}
