package main

import "math"

type quad struct {
	q    vec3
	u, v vec3

	mat  material
	bbox aabb

	normal vec3
	d      float64
	w      vec3
}

func makeQuad(q, u, v vec3, mat material) quad {

	n := cross(u, v)
	normal := n.unitVector()
	d := dot(q, normal)
	w := n.divideScalar(dot(n, n))

	result := quad{q, u, v, mat, aabb{}, normal, d, w}
	result.setBoundingBox()

	return result
}

func (q *quad) setBoundingBox() {
	q.bbox = makeAABBFromExtremaPoints(q.q, q.q.add(q.u).add(q.v)).pad()
}

func (q quad) hit(r ray, rayT interval, rec *hitRecord) bool {
	denom := dot(q.normal, r.direction)

	if math.Abs(denom) < 1e-8 {
		return false
	}

	t := (q.d - dot(q.normal, r.origin)) / denom

	if !rayT.contains(t) {
		return false
	}

	intersection := r.at(t)
	planarHitptVector := intersection.subtract(q.q)
	alpha := dot(q.w, cross(planarHitptVector, q.v))
	beta := dot(q.w, cross(q.u, planarHitptVector))

	if !isInterior(alpha, beta, rec) {
		return false
	}

	rec.t = t
	rec.p = intersection
	rec.mat = q.mat
	rec.setFaceNormal(r, q.normal)

	return true
}

func (q quad) boundingBox() aabb {
	return q.bbox
}

func isInterior(a, b float64, rec *hitRecord) bool {
	if a < 0 || 1 < a || b < 0 || 1 < b {
		return false
	}

	rec.u = a
	rec.v = b

	return true
}
