package main

import (
	"fmt"
	"math"
	"sync"
)

type cameraParameters struct {
	vFov     float64
	lookFrom vec3
	lookAt   vec3
	vUp      vec3

	defocusAngle float64
	focusDist    float64
}

type qualityParameters struct {
	imageWidth, imageHeight int

	samples  int
	maxDepth int

	dof bool
}

type camera struct {
	q *qualityParameters

	imageWidth  int
	imageHeight int
	aspectRatio float64

	center      vec3 // camera center
	pixel00Loc  vec3 // location of pixel (0, 0)
	pixelDeltaU vec3 // offset to pixel to the right
	pixelDeltaV vec3 // offset to pixel below
	u, v, w     vec3 // camera frame basis vectors

	defocusAngle float64
	defocusDiskU vec3 // defocus disk horizontal radius
	defocusDiskV vec3 // defocus disk vertical radius
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

func makeCamera(parameters cameraParameters, q *qualityParameters) *camera {
	imageWidth := q.imageWidth
	imageHeight := q.imageHeight

	aspectRatio := float64(imageWidth) / float64(imageHeight)

	center := parameters.lookFrom

	theta := degreesToRadians(parameters.vFov)
	h := math.Tan(theta / 2.0)
	viewportHeight := 2.0 * h * parameters.focusDist
	viewportWidth := viewportHeight * (aspectRatio)

	w := parameters.lookFrom.subtract(parameters.lookAt).unitVector()
	u := cross(parameters.vUp, w).unitVector()
	v := cross(w, u)

	viewportU := u.multiplyScalar(viewportWidth)
	viewportV := v.multiplyScalar(viewportHeight) // FIXME: negate v?

	pixelDeltaU := viewportU.divideScalar(float64(imageWidth))
	pixelDeltaV := viewportV.divideScalar(float64(imageHeight))

	viewportUpperLeft := center.
		subtract(w.multiplyScalar(parameters.focusDist)).
		subtract(viewportU.divideScalar(2.0)).
		subtract(viewportV.divideScalar(2.0))

	viewportUpperLeft = center.
		subtract(w.multiplyScalar(parameters.focusDist)).
		subtract(viewportU.divideScalar(2.0)).
		subtract(viewportV.divideScalar(2.0))
	pixel00Loc := viewportUpperLeft.add(pixelDeltaU.add(pixelDeltaV).multiplyScalar(0.5))

	defocusRadius := parameters.focusDist * math.Tan(degreesToRadians(parameters.defocusAngle)/2.0)
	defocusDiskU := u.multiplyScalar(defocusRadius)
	defocusDiskV := v.multiplyScalar(defocusRadius)

	return &camera{
		q: q,

		imageWidth:  imageWidth,
		imageHeight: imageHeight,
		aspectRatio: aspectRatio,

		center:      center,
		pixel00Loc:  pixel00Loc,
		pixelDeltaU: pixelDeltaU,
		pixelDeltaV: pixelDeltaV,

		defocusAngle: parameters.defocusAngle,
		defocusDiskU: defocusDiskU,
		defocusDiskV: defocusDiskV,
	}

}

func (r *render) run(samples int, id int) {
	for sample := 0; sample < samples; sample++ {

		pixels := make([]vec3, r.c.imageWidth*r.c.imageHeight)

		for x := 0; x < r.c.imageWidth; x++ {
			for y := 0; y < r.c.imageHeight; y++ {
				ray := r.c.getRay(x, y)

				pixelIndex := getPixelIndex(x, y, r.c.imageWidth)
				pixelColor := rayColor(ray, r.c.q.samples, r.w)

				pixels[pixelIndex] = pixelColor
			}
		}

		r.addSample(pixels)
	}

	fmt.Printf("Thread %d finished\n", id)
}

func (c *camera) render(world hittable) *render {
	return makeRender(c, world)
}

var white = makeVec3(1.0, 1.0, 1.0)
var blue = makeVec3(0.5, 0.7, 1.0)

func (c *camera) getRay(x, y int) ray {
	pixelCenter := c.pixel00Loc.
		add(c.pixelDeltaU.multiplyScalar(float64(x))).
		add(c.pixelDeltaV.multiplyScalar(float64(y)))

	pixelSample := pixelCenter.add(c.pixelSampleSquare())

	var rayOrigin vec3
	if c.defocusAngle <= 0 {
		rayOrigin = c.center
	} else {
		p := randVecInUnitDisk()
		rayOrigin = c.center.add(c.defocusDiskU.multiplyScalar(p.x())).add(c.defocusDiskV.multiplyScalar(p.y()))
	}

	rayDirection := pixelSample.subtract(rayOrigin)

	return makeRay(rayOrigin, rayDirection)
}

func (c *camera) pixelSampleSquare() vec3 {
	px := -0.5 + rand()
	py := -0.5 + rand()
	return c.pixelDeltaU.multiplyScalar(px).add(c.pixelDeltaV.multiplyScalar(py))
}

func rayColor(r ray, depth int, world hittable) vec3 {
	var rec hitRecord

	if depth <= 0 {
		return makeVec3(0, 0, 0)
	}

	if world.hit(r, makeInterval(0.001, infinity), &rec) {

		absorbed, scattered, attenuation := rec.mat.scatter(&r, &rec)

		if absorbed {
			return makeVec3(0, 0, 0)
		}

		return attenuation.multiply(rayColor(*scattered, depth-1, world))
	}

	unitDirection := r.direction.unitVector()
	a := 0.5 * (unitDirection.y() + 1.0)

	return white.multiplyScalar(1.0 - a).add(blue.multiplyScalar(a))
}
