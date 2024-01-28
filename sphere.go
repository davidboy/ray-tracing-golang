package main

import "math"

type sphere struct {
	center    vec3 // the position of this sphere; if moving, this is the initial position
	centerVec vec3 // if this sphere is moving, the vector from `center` to the final position
	radius    float64
	mat       material
	bbox      aabb
}

func makeSphere(center vec3, radius float64, mat material) sphere {
	rvec := makeVec3(radius, radius, radius)
	bbox := makeAABBFromExtremaPoints(center.subtract(rvec), center.add(rvec))

	return sphere{center, makeVec3(0, 0, 0), radius, mat, bbox}
}

func makeMovingSphere(center1, center2 vec3, radius float64, mat material) sphere {
	rvec := makeVec3(radius, radius, radius)
	box1 := makeAABBFromExtremaPoints(center1.subtract(rvec), center1.add(rvec))
	box2 := makeAABBFromExtremaPoints(center2.subtract(rvec), center2.add(rvec))
	bbox := makeAABBFromBoxes(box1, box2)

	return sphere{center1, center2.subtract(center1), radius, mat, bbox}
}

func (s sphere) isMoving() bool {
	return s.centerVec.lengthSquared() > 0 // TODO: ask a math person?
}

func (s sphere) centerAt(time float64) vec3 {
	if s.isMoving() {
		return s.center.add(s.centerVec.multiplyScalar(time))
	} else {
		return s.center
	}
}

func (s sphere) hit(r ray, t interval, rec *hitRecord) bool {
	center := s.centerAt(r.time)

	oc := r.origin.subtract(center)

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

func (s sphere) boundingBox() aabb {
	return s.bbox
}
