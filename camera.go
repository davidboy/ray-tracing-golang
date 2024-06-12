package main

import (
	"math"
)

type cameraParameters struct {
	vFov     float64
	lookFrom vec3
	lookAt   vec3
	vUp      vec3

	defocusAngle float64
	focusDist    float64

	backgroundColor vec3
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

	backgroundColor vec3

	center      vec3 // camera center
	pixel00Loc  vec3 // location of pixel (0, 0)
	pixelDeltaU vec3 // offset to pixel to the right
	pixelDeltaV vec3 // offset to pixel below
	u, v, w     vec3 // camera frame basis vectors

	defocusAngle float64
	defocusDiskU vec3 // defocus disk horizontal radius
	defocusDiskV vec3 // defocus disk vertical radius
}

var pixelIntensity = makeInterval(0.000, 0.999)

func makeCamera(parameters cameraParameters, q *qualityParameters) *camera {
	imageWidth := q.imageWidth
	imageHeight := q.imageHeight

	aspectRatio := float64(imageWidth) / float64(imageHeight)

	if !quality.dof {
		parameters.defocusAngle = 0 // disable depth of field
	}

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

		backgroundColor: parameters.backgroundColor,

		center:      center,
		pixel00Loc:  pixel00Loc,
		pixelDeltaU: pixelDeltaU,
		pixelDeltaV: pixelDeltaV,

		defocusAngle: parameters.defocusAngle,
		defocusDiskU: defocusDiskU,
		defocusDiskV: defocusDiskV,
	}

}

var white = makeVec3(1.0, 1.0, 1.0)
var blue = makeVec3(0.5, 0.7, 1.0)

func (c *camera) getRay(x, y int) ray {
	offset := pixelSampleSquare()

	pixelCenter := c.pixel00Loc.
		add(c.pixelDeltaU.multiplyScalar(float64(x) + offset.e[0])).
		add(c.pixelDeltaV.multiplyScalar(float64(y) + offset.e[1]))

	var rayOrigin vec3
	if c.defocusAngle <= 0 {
		rayOrigin = c.center
	} else {
		p := randVecInUnitDisk()
		rayOrigin = c.center.add(c.defocusDiskU.multiplyScalar(p.x())).add(c.defocusDiskV.multiplyScalar(p.y()))
	}

	rayDirection := pixelCenter.subtract(rayOrigin)
	rayTime := rand()

	return makeTimedRay(rayOrigin, rayDirection, rayTime)
}

func pixelSampleSquare() vec3 {
	px := -0.5 + rand()
	py := -0.5 + rand()
	return makeVec3(px, py, 0)
}

func (c *camera) rayColor(r ray, depth int, world hittable) vec3 {
	var rec hitRecord

	if depth <= 0 {
		return makeVec3(0, 0, 0)
	}

	if !world.hit(r, makeInterval(0.001, infinity), &rec) {
		return c.backgroundColor
	}

	colorFromEmission := rec.mat.emitted(rec.u, rec.v, rec.p)

	absorbed, scattered, attenuation := rec.mat.scatter(&r, &rec)

	// TODO: flip boolean ? idk ??
	if absorbed {
		return colorFromEmission
	}

	colorFromScatter := attenuation.multiply(c.rayColor(*scattered, depth-1, world))

	return colorFromEmission.add(colorFromScatter)

	// absorbed, scattered, attenuation := rec.mat.scatter(&r, &rec)

	// if absorbed {
	// 	return makeVec3(0, 0, 0)
	// }

	// return attenuation.multiply(rayColor(*scattered, depth-1, world))

	// unitDirection := r.direction.unitVector()
	// a := 0.5 * (unitDirection.y() + 1.0)

	// return white.multiplyScalar(1.0 - a).add(blue.multiplyScalar(a))
}
