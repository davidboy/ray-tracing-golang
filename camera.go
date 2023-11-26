package main

import "sync"

type camera struct {
	imageWidth  int
	imageHeight int
	aspectRatio float64

	samples int

	center      vec3 // camera center
	pixel00Loc  vec3 // location of pixel (0, 0)
	pixelDeltaU vec3 // offset to pixel to the right
	pixelDeltaV vec3 // offset to pixel below

	onUpdateProgress func(progress float64)
}

type render struct {
	samples int    // number of samples that have been taken
	pixels  []vec3 // rendered pixels

	finished bool // whether the render has finished

	mu sync.Mutex
}

func makeRender(imageWidth, imageHeight int) *render {
	return &render{
		pixels: make([]vec3, imageWidth*imageHeight),
	}
}

func (r *render) addSample(pixels []vec3) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, pixel := range pixels {
		r.pixels[i].addMut(pixel)
	}

	r.samples++
}

func (r *render) squash() []vec3 {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]vec3, len(r.pixels))

	for i, pixel := range r.pixels {
		result[i] = pixel.divideScalar(float64(r.samples))
	}

	return result
}

func makeCamera(imageWidth int, aspectRatio float64, samples int) *camera {
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

		samples: samples,

		center:      center,
		pixel00Loc:  pixel00Loc,
		pixelDeltaU: pixelDeltaU,
		pixelDeltaV: pixelDeltaV,

		onUpdateProgress: func(progress float64) {},
	}

}

func (c *camera) render(world hittable) *render {

	result := makeRender(c.imageWidth, c.imageHeight)

	pixels := make([]vec3, c.imageWidth*c.imageHeight)

	intensity := makeInterval(0.000, 0.999)

	for x := 0; x < imageWidth; x++ {
		c.onUpdateProgress(float64(x) / float64(imageWidth))

		for y := 0; y < c.imageHeight; y++ {

			pixelColor := makeVec3(0, 0, 0)

			for sample := 0; sample < c.samples; sample++ {
				r := c.getRay(x, y)
				pixelColor.addMut(rayColor(r, world))
			}

			pixelColor.divideScalarMut(float64(c.samples))
			pixelColor.clampMut(intensity)

			pixels[getPixelIndex(x, y, imageWidth)] = pixelColor

		}
	}

	result.addSample(pixels)

	return result
}

var white = makeVec3(1.0, 1.0, 1.0)
var blue = makeVec3(0.5, 0.7, 1.0)

func (c *camera) getRay(x, y int) Ray {
	pixelCenter := c.pixel00Loc.
		add(c.pixelDeltaU.multiplyScalar(float64(x))).
		add(c.pixelDeltaV.multiplyScalar(float64(y)))

	pixelSample := pixelCenter.add(c.pixelSampleSquare())

	rayOrigin := c.center
	rayDirection := pixelSample.subtract(rayOrigin)

	return makeRay(rayOrigin, rayDirection)
}

func (c *camera) pixelSampleSquare() vec3 {
	px := -0.5 + rand()
	py := -0.5 + rand()
	return c.pixelDeltaU.multiplyScalar(px).add(c.pixelDeltaV.multiplyScalar(py))
}

func rayColor(r Ray, world hittable) vec3 {
	var rec hitRecord

	if world.hit(r, makeInterval(0.0, infinity), &rec) {
		return rec.normal.addScalar(1.0).multiplyScalar(0.5)
	}

	unitDirection := r.direction.unitVector()
	a := 0.5 * (unitDirection.y() + 1.0)

	return white.multiplyScalar(1.0 - a).add(blue.multiplyScalar(a))
}
