package main

import "math"

type sphere struct {
	center vec3
	radius float64
	mat    material
}

func makeSphere(center vec3, radius float64, mat material) sphere {
	return sphere{center, radius, mat}
}

func (s sphere) hit(r ray, t interval, rec *hitRecord) bool {
	oc := r.origin.subtract(s.center)

	a := r.direction.lengthSquared()
	halfB := dot(oc, r.direction)
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
	rec.p = r.at(rec.t)

	outwardNormal := rec.p.subtract(s.center).divideScalar(s.radius)
	rec.setFaceNormal(r, outwardNormal)
	rec.mat = s.mat

	return true
}
