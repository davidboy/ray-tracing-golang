package main

type hitRecord struct {
	p         vec3
	normal    vec3
	mat       material
	t         float64
	frontFace bool
}

func (rec *hitRecord) setFaceNormal(r ray, outwardNormal vec3) {
	rec.frontFace = dot(r.direction, outwardNormal) < 0

	if rec.frontFace {
		rec.normal = outwardNormal
	} else {
		rec.normal = outwardNormal.negate()
	}
}
