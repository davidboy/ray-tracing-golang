package main

import "math"

const POINT_COUNT = 256

type perlin struct {
	ranvec              []vec3
	permX, permY, permZ []int
}

func makePerlin() perlin {
	ranvec := make([]vec3, POINT_COUNT)
	for i := 0; i < POINT_COUNT; i++ {
		ranvec[i] = randVecB(-1, 1).unitVector()
	}

	permX := perlinGeneratePerm()
	permY := perlinGeneratePerm()
	permZ := perlinGeneratePerm()

	return perlin{ranvec, permX, permY, permZ} // TODO: need pointer? how do slices work lol
}

func (per *perlin) noise(p vec3) float64 {
	u := p.x() - math.Floor(p.x())
	v := p.y() - math.Floor(p.y())
	w := p.z() - math.Floor(p.z())

	u = u * u * (3 - 2*u)
	v = v * v * (3 - 2*v)
	w = w * w * (3 - 2*w)

	i := int(math.Floor(p.x()))
	j := int(math.Floor(p.y()))
	k := int(math.Floor(p.z()))

	c := make([][][]vec3, 2)
	for i := range c {
		c[i] = make([][]vec3, 2)
		for j := range c[i] {
			c[i][j] = make([]vec3, 2)
		}
	}

	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = per.ranvec[per.permX[(i+di)&255]^per.permY[(j+dj)&255]^per.permZ[(k+dk)&255]]
			}
		}
	}

	return perlinInterp(c, u, v, w)
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

func perlinInterp(c [][][]vec3, u, v, w float64) float64 {
	uu := u * u * (3 - 2*u)
	vv := v * v * (3 - 2*v)
	ww := w * w * (3 - 2*w)

	accum := 0.0

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weightV := makeVec3(u-float64(i), v-float64(j), w-float64(k))
				accum += (float64(i)*uu + (1-float64(i))*(1-uu)) *
					(float64(j)*vv + (1-float64(j))*(1-vv)) *
					(float64(k)*ww + (1-float64(k))*(1-ww)) *
					dot(c[i][j][k], weightV)
			}
		}
	}

	return accum
}
