package main

type material interface {
	scatter(r *ray, rec *hitRecord) (absorbed bool, scattered *ray, attenuation vec3)
}

type lambertian struct {
	albedo vec3
}

func (l lambertian) scatter(r *ray, rec *hitRecord) (absorbed bool, scattered *ray, attenuation vec3) {
	scatterDirection := rec.normal.add(randUnitVector())

	// Catch degenerate scatter direction
	if scatterDirection.nearZero() {
		scatterDirection = rec.normal
	}

	scatteredRay := makeRay(rec.p, scatterDirection)

	return false, &scatteredRay, l.albedo
}

type metal struct {
	albedo vec3
	fuzz   float64
}

func (m metal) scatter(r *ray, rec *hitRecord) (absorbed bool, scattered *ray, attenuation vec3) {
	reflected := r.direction.unitVector().reflect(rec.normal)
	scatteredRay := makeRay(rec.p, reflected.add(randUnitVector().multiplyScalar(m.fuzz)))
	absorbed = dot(scatteredRay.direction, rec.normal) <= 0

	return absorbed, &scatteredRay, m.albedo
}
