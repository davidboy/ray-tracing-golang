package main

import "sync"

type render struct {
	c *camera
	w hittable

	samples []int  // number of samples that have been taken for rendered pixels
	pixels  []vec3 // rendered pixels

	finished chan bool // whether the render has finished
	dirty    bool      // whether the render has changed since the last squash

	mu sync.Mutex
}

func makeRender(c *camera, w hittable) *render {
	return &render{
		c: c,
		w: w,

		finished: make(chan bool),

		samples: make([]int, c.imageWidth*c.imageHeight),
		pixels:  make([]vec3, c.imageWidth*c.imageHeight),
	}
}

func (r *render) runWholeImage(samples int, finished chan bool) {

	for sample := 0; sample < samples; sample++ {

		pixels := make([]vec3, r.c.imageWidth*r.c.imageHeight)

		for x := 0; x < r.c.imageWidth; x++ {
			for y := 0; y < r.c.imageHeight; y++ {
				ray := r.c.getRay(x, y)

				pixelIndex := getPixelIndex(x, y, r.c.imageWidth)
				pixelColor := r.c.rayColor(ray, r.c.q.samples, r.w)

				pixels[pixelIndex] = pixelColor
			}
		}

		r.addWholeImageSample(pixels)
	}

	close(finished)
}

func (r *render) runSinglePixel(x, y, samples int) {

	pixel := makeVec3(0, 0, 0)
	pixelIndex := getPixelIndex(x, y, r.c.imageWidth)

	for i := 0; i < samples; i++ {
		pixel.addMut(r.c.rayColor(r.c.getRay(x, y), r.c.q.samples, r.w))
	}

	r.addSinglePixelSample(pixel, samples, pixelIndex)
}

func (r *render) addSinglePixelSample(pixel vec3, samples int, index int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.samples[index] = r.samples[index] + samples
	r.pixels[index].addMut(pixel)

	r.dirty = true
}

func (r *render) addWholeImageSample(pixels []vec3) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, pixel := range pixels {
		r.samples[i]++
		r.pixels[i].addMut(pixel)
	}

	r.dirty = true
}

func (r *render) getTotalSamplesTaken() (result int) {
	for _, samples := range r.samples {
		result += samples
	}

	return result
}

func (r *render) squash() []vec3 {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]vec3, len(r.pixels))

	for i, pixel := range r.pixels {
		result[i] = pixel.divideScalar(float64(r.samples[i])).clamp(pixelIntensity).toGammaSpace()
	}

	r.dirty = false

	return result
}
