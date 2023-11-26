package main

type hitRecord struct {
	p         Vec3
	normal    Vec3
	t         float64
	frontFace bool
}

func (rec *hitRecord) setFaceNormal(r Ray, outwardNormal Vec3) {
	rec.frontFace = r.Direction.Dot(outwardNormal) < 0

	if rec.frontFace {
		rec.normal = outwardNormal
	} else {
		rec.normal = outwardNormal.Negate()
	}
}
