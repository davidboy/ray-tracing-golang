package main

import "math"

type interval struct {
	min, max float64
}

func makeInterval(min, max float64) interval {
	return interval{min, max}
}

func (i interval) size() float64 {
	return i.max - i.min
}

func (i interval) expand(delta float64) interval {
	padding := delta / 2
	return makeInterval(i.min-padding, i.max+padding)
}

func (i interval) contains(x float64) bool {
	return i.min <= x && x <= i.max
}

func (i interval) surrounds(x float64) bool {
	return i.min < x && x < i.max
}

func (i interval) clamp(x float64) float64 {
	if x < i.min {
		return i.min
	}
	if x > i.max {
		return i.max
	}
	return x
}

func (v vec3) clamp(i interval) vec3 {
	return makeVec3(
		i.clamp(v.e[0]),
		i.clamp(v.e[1]),
		i.clamp(v.e[2]),
	)
}

func (v *vec3) clampMut(i interval) *vec3 {
	v.e[0] = i.clamp(v.e[0])
	v.e[1] = i.clamp(v.e[1])
	v.e[2] = i.clamp(v.e[2])
	return v
}

func (i interval) addScalar(v float64) interval {
	return makeInterval(i.min+v, i.max+v)
}

var empty = makeInterval(math.Inf(1), math.Inf(-1))
var universe = makeInterval(math.Inf(-1), math.Inf(1))
