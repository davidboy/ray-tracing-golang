package main

import "sync"

type render struct {
	c *camera
	w hittable

	samples int    // number of samples that have been taken
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

		pixels: make([]vec3, c.imageWidth*c.imageHeight),
	}
}

func (r *render) run(samples int, id int, finished chan bool) {

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

	close(finished)
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
