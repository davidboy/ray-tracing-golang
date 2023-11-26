package main

import (
	"math"
)

type vec3 struct {
	e [3]float64
}

func makeVec3(x, y, z float64) vec3 {
	return vec3{[3]float64{x, y, z}}
}

func makeVec3Ref(x, y, z float64) *vec3 {
	return &vec3{[3]float64{x, y, z}}
}

func (v vec3) negate() vec3 {
	return vec3{[3]float64{-v.e[0], -v.e[1], -v.e[2]}}
}

func (v vec3) x() float64 {
	return v.e[0]
}

func (v vec3) y() float64 {
	return v.e[1]
}

func (v vec3) z() float64 {
	return v.e[2]
}

func (a vec3) add(b vec3) vec3 {
	return vec3{[3]float64{a.e[0] + b.e[0], a.e[1] + b.e[1], a.e[2] + b.e[2]}}
}

func (a vec3) addScalar(t float64) vec3 {
	return vec3{[3]float64{a.e[0] + t, a.e[1] + t, a.e[2] + t}}
}

func (a *vec3) addMut(b vec3) *vec3 {
	a.e[0] += b.e[0]
	a.e[1] += b.e[1]
	a.e[2] += b.e[2]
	return a
}

func (a *vec3) addScalarMut(t float64) *vec3 {
	a.e[0] += t
	a.e[1] += t
	a.e[2] += t
	return a
}

func (a vec3) subtract(b vec3) vec3 {
	return vec3{[3]float64{a.e[0] - b.e[0], a.e[1] - b.e[1], a.e[2] - b.e[2]}}
}

func (a *vec3) subtractMut(b vec3) *vec3 {
	a.e[0] -= b.e[0]
	a.e[1] -= b.e[1]
	a.e[2] -= b.e[2]
	return a
}

func (a vec3) subtractScalar(t float64) vec3 {
	return vec3{[3]float64{a.e[0] - t, a.e[1] - t, a.e[2] - t}}
}

func (a *vec3) subtractScalarMut(t float64) *vec3 {
	a.e[0] -= t
	a.e[1] -= t
	a.e[2] -= t
	return a
}

func (a vec3) multiply(b vec3) vec3 {
	return vec3{[3]float64{a.e[0] * b.e[0], a.e[1] * b.e[1], a.e[2] * b.e[2]}}
}

func (a vec3) multiplyScalar(t float64) vec3 {
	return vec3{[3]float64{a.e[0] * t, a.e[1] * t, a.e[2] * t}}
}

func (a *vec3) multiplyMut(b vec3) *vec3 {
	a.e[0] *= b.e[0]
	a.e[1] *= b.e[1]
	a.e[2] *= b.e[2]
	return a
}

func (a *vec3) multiplyScalarMut(t float64) *vec3 {
	a.e[0] *= t
	a.e[1] *= t
	a.e[2] *= t
	return a
}

func (a vec3) divide(b vec3) vec3 {
	return vec3{[3]float64{a.e[0] / b.e[0], a.e[1] / b.e[1], a.e[2] / b.e[2]}}
}

func (a vec3) divideScalar(t float64) vec3 {
	return vec3{[3]float64{a.e[0] / t, a.e[1] / t, a.e[2] / t}}
}

func (a *vec3) divideMut(b vec3) *vec3 {
	a.e[0] /= b.e[0]
	a.e[1] /= b.e[1]
	a.e[2] /= b.e[2]
	return a
}

func (a *vec3) divideScalarMut(t float64) *vec3 {
	a.e[0] /= t
	a.e[1] /= t
	a.e[2] /= t
	return a
}

func dot(a vec3, b vec3) float64 {
	return a.e[0]*b.e[0] + a.e[1]*b.e[1] + a.e[2]*b.e[2]
}

func cross(a vec3, b vec3) vec3 {
	return vec3{[3]float64{
		a.e[1]*b.e[2] - a.e[2]*b.e[1],
		a.e[2]*b.e[0] - a.e[0]*b.e[2],
		a.e[0]*b.e[1] - a.e[1]*b.e[0],
	}}
}

func (a vec3) length() float64 {
	return math.Sqrt(a.lengthSquared())
}

func (a vec3) lengthSquared() float64 {
	return a.e[0]*a.e[0] + a.e[1]*a.e[1] + a.e[2]*a.e[2]
}

func (a vec3) unitVector() vec3 {
	return a.divideScalar(a.length())
}
