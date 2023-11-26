package main

import "math"

type sphere struct {
	center Vec3
	radius float64
}

func makeSphere(center Vec3, radius float64) sphere {
	return sphere{center, radius}
}

func (s sphere) hit(r Ray, tMin float64, tMax float64, rec *hitRecord) bool {
	oc := r.Origin.Subtract(s.center)

	a := r.Direction.LengthSquared()
	halfB := oc.Dot(r.Direction)
	c := oc.LengthSquared() - s.radius*s.radius

	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return false
	}
	sqrtd := math.Sqrt(discriminant)

	root := (-halfB - sqrtd) / a
	if root <= tMin || tMax <= root {
		root = (-halfB + sqrtd) / a
		if root <= tMin || tMax <= root {
			return false
		}
	}

	rec.t = root
	rec.p = r.At(rec.t)

	outwardNormal := rec.p.Subtract(s.center).DivideScalar(s.radius)
	rec.setFaceNormal(r, outwardNormal)

	return true
}
