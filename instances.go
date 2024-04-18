package main

import "math"

type translate struct {
	h      hittable
	offset vec3
	bbox   aabb
}

func makeTranslate(h hittable, offset vec3) *translate {
	return &translate{h, offset, h.boundingBox().add(offset)}
}

func (t *translate) hit(r ray, rayT interval, rec *hitRecord) bool {
	offsetR := makeTimedRay(r.origin.subtract(t.offset), r.direction, r.time)

	if !t.h.hit(offsetR, rayT, rec) {
		return false
	}

	rec.p = rec.p.add(t.offset)
	return true
}

func (t *translate) boundingBox() aabb {
	return t.bbox
}

type rotateY struct {
	h                  hittable
	sinTheta, cosTheta float64
	bbox               aabb
}

func makeRotateY(h hittable, angle float64) *rotateY {
	radians := degreesToRadians(angle)

	sinTheta := math.Sin(radians)
	cosTheta := math.Cos(radians)
	bbox := h.boundingBox()

	min := makeVec3(math.Inf(1), math.Inf(1), math.Inf(1))
	max := makeVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1))

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*bbox.x.max + float64(1-i)*bbox.x.min
				y := float64(j)*bbox.y.max + float64(1-j)*bbox.y.min
				z := float64(k)*bbox.z.max + float64(1-k)*bbox.z.min

				newX := cosTheta*x + sinTheta*z
				newZ := -sinTheta*x + cosTheta*z

				tester := makeVec3(newX, y, newZ)

				for c := 0; c < 3; c++ {
					min.e[c] = math.Min(min.e[c], tester.e[c])
					max.e[c] = math.Max(max.e[c], tester.e[c])
				}
			}
		}
	}

	return &rotateY{h, sinTheta, cosTheta, makeAABBFromExtremaPoints(min, max)}
}

func (self *rotateY) hit(r ray, rayT interval, rec *hitRecord) bool {
	origin := r.origin
	direction := r.direction

	origin.e[0] = self.cosTheta*r.origin.e[0] - self.sinTheta*r.origin.e[2]
	origin.e[2] = self.sinTheta*r.origin.e[0] + self.cosTheta*r.origin.e[2]

	direction.e[0] = self.cosTheta*r.direction.e[0] - self.sinTheta*r.direction.e[2]
	direction.e[2] = self.sinTheta*r.direction.e[0] + self.cosTheta*r.direction.e[2]

	rotatedR := makeTimedRay(origin, direction, r.time)

	if !self.h.hit(rotatedR, rayT, rec) {
		return false
	}

	p := rec.p
	p.e[0] = self.cosTheta*rec.p.e[0] + self.sinTheta*rec.p.e[2]
	p.e[2] = -self.sinTheta*rec.p.e[0] + self.cosTheta*rec.p.e[2]

	normal := rec.normal
	normal.e[0] = self.cosTheta*rec.normal.e[0] + self.sinTheta*rec.normal.e[2]
	normal.e[2] = -self.sinTheta*rec.normal.e[0] + self.cosTheta*rec.normal.e[2]

	rec.p = p
	rec.normal = normal

	return true
}

func (r *rotateY) boundingBox() aabb {
	return r.bbox
}
