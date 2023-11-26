package main

import "sync"

// TODO: consolidate all quality-related options
const max_depth = 10 // maximum number of bounces

type camera struct {
	imageWidth  int
	imageHeight int
	aspectRatio float64

	center      vec3 // camera center
	pixel00Loc  vec3 // location of pixel (0, 0)
	pixelDeltaU vec3 // offset to pixel to the right
	pixelDeltaV vec3 // offset to pixel below
}

type render struct {
	c *camera
	w hittable

	samples int    // number of samples that have been taken
	pixels  []vec3 // rendered pixels

	finished bool // whether the render has finished
	dirty    bool // whether the render has changed since the last squash

	mu sync.Mutex
}

func makeRender(c *camera, w hittable) *render {
	return &render{
		c: c,
		w: w,

		pixels: make([]vec3, c.imageWidth*c.imageHeight),
	}
}

func (r *render) addSample(pixels []vec3) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, pixel := range pixels {
		r.pixels[i].addMut(pixel)
	}

	r.samples++
	r.dirty = true
}

var pixelIntensity = makeInterval(0.000, 0.999)

func (r *render) squash() []vec3 {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]vec3, len(r.pixels))

	for i, pixel := range r.pixels {
		result[i] = pixel.divideScalar(float64(r.samples)).clamp(pixelIntensity).toGammaSpace()
	}

	r.dirty = false

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

		center:      center,
		pixel00Loc:  pixel00Loc,
		pixelDeltaU: pixelDeltaU,
		pixelDeltaV: pixelDeltaV,
	}

}

func (r *render) run(samples int) {
	for sample := 0; sample < samples; sample++ {

		pixels := make([]vec3, r.c.imageWidth*r.c.imageHeight)

		for x := 0; x < imageWidth; x++ {
			for y := 0; y < r.c.imageHeight; y++ {
				ray := r.c.getRay(x, y)

				pixelIndex := getPixelIndex(x, y, imageWidth)
				pixelColor := rayColor(ray, max_depth, r.w)

				pixels[pixelIndex] = pixelColor
			}
		}

		r.addSample(pixels)
	}
}

func (c *camera) render(world hittable) *render {
	return makeRender(c, world)
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

func rayColor(r Ray, depth int, world hittable) vec3 {
	var rec hitRecord

	if depth <= 0 {
		return makeVec3(0, 0, 0)
	}

	if world.hit(r, makeInterval(0.001, infinity), &rec) {
		direction := rec.normal.add(randUnitVector())

		return rayColor(makeRay(rec.p, direction), depth-1, world).multiplyScalar(0.5)
	}

	unitDirection := r.direction.unitVector()
	a := 0.5 * (unitDirection.y() + 1.0)

	return white.multiplyScalar(1.0 - a).add(blue.multiplyScalar(a))
}
