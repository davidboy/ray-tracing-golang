package main

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
	i := int(4*p.x()) & 255
	j := int(4*p.y()) & 255
	k := int(4*p.z()) & 255

	return per.ranfloat[per.permX[i]^per.permY[j]^per.permZ[k]]
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
