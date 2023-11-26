package main

type camera struct {
	imageWidth  int
	imageHeight int
	aspectRatio float64

	center      Vec3 // camera center
	pixel00Loc  Vec3 // location of pixel (0, 0)
	pixelDeltaU Vec3 // offset to pixel to the right
	pixelDeltaV Vec3 // offset to pixel below

	onUpdateProgress func(progress float64)
}

func makeCamera(imageWidth int, aspectRatio float64) *camera {
	imageHeight := calculateImageHeight(imageWidth, aspectRatio)

	center := MakeVec3(0, 0, 0)

	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))

	viewportU := MakeVec3(viewportWidth, 0, 0)
	viewportV := MakeVec3(0, viewportHeight, 0)

	pixelDeltaU := viewportU.DivideScalar(float64(imageWidth))
	pixelDeltaV := viewportV.DivideScalar(float64(imageHeight))

	viewportUpperLeft := center.
		Subtract(MakeVec3(0, 0, focalLength)).
		Subtract(viewportU.DivideScalar(2)).
		Subtract(viewportV.DivideScalar(2))

	pixel00Loc := viewportUpperLeft.Add(pixelDeltaU.Add(pixelDeltaV).MultiplyScalar(0.5))

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

func (c *camera) render(world hittable) []Vec3 {
	pixels := make([]Vec3, c.imageWidth*c.imageHeight)

	for x := 0; x < imageWidth; x++ {
		c.onUpdateProgress(float64(x) / float64(imageWidth))

		for y := 0; y < c.imageHeight; y++ {
			pixelCenter := c.pixel00Loc.
				Add(c.pixelDeltaU.MultiplyScalar(float64(x))).
				Add(c.pixelDeltaV.MultiplyScalar(float64(y)))

			rayDirection := pixelCenter.Subtract(c.center)
			ray := MakeRay(c.center, rayDirection)

			pixels[getPixelIndex(x, y, imageWidth)] = rayColor(ray, world)
		}
	}

	return pixels
}

var white = MakeVec3(1.0, 1.0, 1.0)
var blue = MakeVec3(0.5, 0.7, 1.0)

func rayColor(r Ray, world hittable) Vec3 {
	var rec hitRecord

	if world.hit(r, makeInterval(0.0, infinity), &rec) {
		return rec.normal.AddScalar(1.0).MultiplyScalar(0.5)
	}

	unitDirection := r.Direction.UnitVector()
	a := 0.5 * (unitDirection.Y() + 1.0)

	return white.MultiplyScalar(1.0 - a).Add(blue.MultiplyScalar(a))
}
