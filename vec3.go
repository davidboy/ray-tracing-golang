package main

import (
	"math"
)

type Vec3 struct {
	e [3]float64
}

func MakeVec3(x, y, z float64) Vec3 {
	return Vec3{[3]float64{x, y, z}}
}

func MakeVec3Ref(x, y, z float64) *Vec3 {
	return &Vec3{[3]float64{x, y, z}}
}

func (v Vec3) Negate() Vec3 {
	return Vec3{[3]float64{-v.e[0], -v.e[1], -v.e[2]}}
}

func (v Vec3) X() float64 {
	return v.e[0]
}

func (v Vec3) Y() float64 {
	return v.e[1]
}

func (v Vec3) Z() float64 {
	return v.e[2]
}

func (a Vec3) Add(b Vec3) Vec3 {
	return Vec3{[3]float64{a.e[0] + b.e[0], a.e[1] + b.e[1], a.e[2] + b.e[2]}}
}

func (a Vec3) AddScalar(t float64) Vec3 {
	return Vec3{[3]float64{a.e[0] + t, a.e[1] + t, a.e[2] + t}}
}

func (a *Vec3) AddMut(b Vec3) *Vec3 {
	a.e[0] += b.e[0]
	a.e[1] += b.e[1]
	a.e[2] += b.e[2]
	return a
}

func (a *Vec3) AddScalarMut(t float64) *Vec3 {
	a.e[0] += t
	a.e[1] += t
	a.e[2] += t
	return a
}

func (a Vec3) Subtract(b Vec3) Vec3 {
	return Vec3{[3]float64{a.e[0] - b.e[0], a.e[1] - b.e[1], a.e[2] - b.e[2]}}
}

func (a *Vec3) SubtractMut(b Vec3) *Vec3 {
	a.e[0] -= b.e[0]
	a.e[1] -= b.e[1]
	a.e[2] -= b.e[2]
	return a
}

func (a Vec3) SubtractScalar(t float64) Vec3 {
	return Vec3{[3]float64{a.e[0] - t, a.e[1] - t, a.e[2] - t}}
}

func (a *Vec3) SubtractScalarMut(t float64) *Vec3 {
	a.e[0] -= t
	a.e[1] -= t
	a.e[2] -= t
	return a
}

func (a Vec3) Multiply(b Vec3) Vec3 {
	return Vec3{[3]float64{a.e[0] * b.e[0], a.e[1] * b.e[1], a.e[2] * b.e[2]}}
}

func (a Vec3) MultiplyScalar(t float64) Vec3 {
	return Vec3{[3]float64{a.e[0] * t, a.e[1] * t, a.e[2] * t}}
}

func (a *Vec3) MultiplyMut(b Vec3) *Vec3 {
	a.e[0] *= b.e[0]
	a.e[1] *= b.e[1]
	a.e[2] *= b.e[2]
	return a
}

func (a *Vec3) MultiplyScalarMut(t float64) *Vec3 {
	a.e[0] *= t
	a.e[1] *= t
	a.e[2] *= t
	return a
}

func (a Vec3) Divide(b Vec3) Vec3 {
	return Vec3{[3]float64{a.e[0] / b.e[0], a.e[1] / b.e[1], a.e[2] / b.e[2]}}
}

func (a Vec3) DivideScalar(t float64) Vec3 {
	return Vec3{[3]float64{a.e[0] / t, a.e[1] / t, a.e[2] / t}}
}

func (a *Vec3) DivideMut(b Vec3) *Vec3 {
	a.e[0] /= b.e[0]
	a.e[1] /= b.e[1]
	a.e[2] /= b.e[2]
	return a
}

func (a *Vec3) DivideScalarMut(t float64) *Vec3 {
	a.e[0] /= t
	a.e[1] /= t
	a.e[2] /= t
	return a
}

func (a Vec3) Dot(b Vec3) float64 {
	return a.e[0]*b.e[0] + a.e[1]*b.e[1] + a.e[2]*b.e[2]
}

func Cross(a Vec3, b Vec3) Vec3 {
	return Vec3{[3]float64{
		a.e[1]*b.e[2] - a.e[2]*b.e[1],
		a.e[2]*b.e[0] - a.e[0]*b.e[2],
		a.e[0]*b.e[1] - a.e[1]*b.e[0],
	}}
}

func (a Vec3) Length() float64 {
	return math.Sqrt(a.LengthSquared())
}

func (a Vec3) LengthSquared() float64 {
	return a.e[0]*a.e[0] + a.e[1]*a.e[1] + a.e[2]*a.e[2]
}

func (a Vec3) UnitVector() Vec3 {
	return a.DivideScalar(a.Length())
}
