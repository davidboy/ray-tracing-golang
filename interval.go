package main

import "math"

type interval struct {
	min, max float64
}

func makeInterval(min, max float64) interval {
	return interval{min, max}
}

func (i interval) contains(x float64) bool {
	return i.min <= x && x <= i.max
}

func (i interval) surrounds(x float64) bool {
	return i.min < x && x < i.max
}

var empty = makeInterval(math.Inf(1), math.Inf(-1))
var universe = makeInterval(math.Inf(-1), math.Inf(1))
