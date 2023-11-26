package main

import "math"

type sphere struct {
	center vec3
	radius float64
}

func makeSphere(center vec3, radius float64) sphere {
	return sphere{center, radius}
}

func (s sphere) hit(r Ray, t interval, rec *hitRecord) bool {
	oc := r.Origin.subtract(s.center)

	a := r.Direction.lengthSquared()
	halfB := dot(oc, r.Direction)
	c := oc.lengthSquared() - s.radius*s.radius

	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return false
	}
	sqrtd := math.Sqrt(discriminant)

	root := (-halfB - sqrtd) / a
	if !t.surrounds(root) {
		root = (-halfB + sqrtd) / a
		if !t.surrounds(root) {
			return false
		}
	}

	rec.t = root
	rec.p = r.At(rec.t)

	outwardNormal := rec.p.subtract(s.center).divideScalar(s.radius)
	rec.setFaceNormal(r, outwardNormal)

	return true
}
