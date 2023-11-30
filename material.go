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

type dielectric struct {
	ir float64 // Index of Refraction
}

func (d dielectric) scatter(r *ray, rec *hitRecord) (absorbed bool, scattered *ray, attenuation vec3) {
	attenuation = makeVec3(1.0, 1.0, 1.0)

	var refractionRatio float64
	if rec.frontFace {
		refractionRatio = 1.0 / d.ir
	} else {
		refractionRatio = d.ir
	}

	unitDirection := r.direction.unitVector()
	refracted := unitDirection.refract(rec.normal, refractionRatio)

	scatteredRay := makeRay(rec.p, refracted)

	return false, &scatteredRay, attenuation
}
