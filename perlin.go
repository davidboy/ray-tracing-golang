package main

import "math"

const POINT_COUNT = 256

type perlin struct {
	ranfloat            []float64
	permX, permY, permZ []int
}

func makePerlin() perlin {
	ranfloat := make([]float64, POINT_COUNT)
	for i := 0; i < POINT_COUNT; i++ {
		ranfloat[i] = rand()
	}

	permX := perlinGeneratePerm()
	permY := perlinGeneratePerm()
	permZ := perlinGeneratePerm()

	return perlin{ranfloat, permX, permY, permZ}
}

func (per *perlin) noise(p vec3) float64 {
	u := p.x() - math.Floor(p.x())
	v := p.y() - math.Floor(p.y())
	w := p.z() - math.Floor(p.z())

	i := int(math.Floor(p.x()))
	j := int(math.Floor(p.y()))
	k := int(math.Floor(p.z()))

	c := make([][][]float64, 2)
	for i := range c {
		c[i] = make([][]float64, 2)
		for j := range c[i] {
			c[i][j] = make([]float64, 2)
		}
	}

	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = per.ranfloat[per.permX[(i+di)&255]^per.permY[(j+dj)&255]^per.permZ[(k+dk)&255]]
			}
		}
	}

	return trilinearInterpolation(c, u, v, w)
}

func perlinGeneratePerm() []int {
	p := make([]int, POINT_COUNT)

	for i := 0; i < POINT_COUNT; i++ {
		p[i] = i
	}

	permute(&p, POINT_COUNT) // TODO: reference needed?

	return p
}

func permute(p *[]int, n int) {
	for i := n - 1; i > 0; i-- {
		target := randIntB(0, i+1) // TODO: check upper bound
		(*p)[i], (*p)[target] = (*p)[target], (*p)[i]
	}
}

func trilinearInterpolation(c [][][]float64, u, v, w float64) float64 {
	accum := 0.0

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				accum += (float64(i)*u + (1-float64(i))*(1-u)) *
					(float64(j)*v + (1-float64(j))*(1-v)) *
					(float64(k)*w + (1-float64(k))*(1-w)) * c[i][j][k]
			}
		}
	}

	return accum
}
