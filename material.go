package main

import "math"

type material interface {
	scatter(r *ray, rec *hitRecord) (absorbed bool, scattered *ray, attenuation vec3)
	emitted(u, v float64, p vec3) vec3
}

type lambertian struct {
	albedo texture
}

func (l lambertian) scatter(r *ray, rec *hitRecord) (absorbed bool, scattered *ray, attenuation vec3) {
	scatterDirection := rec.normal.add(randUnitVector())

	// Catch degenerate scatter direction
	if scatterDirection.nearZero() {
		scatterDirection = rec.normal
	}

	scatteredRay := makeTimedRay(rec.p, scatterDirection, r.time)
	attenuation = l.albedo.value(rec.u, rec.v, rec.p)

	return false, &scatteredRay, attenuation
}

func (l lambertian) emitted(u, v float64, p vec3) vec3 {
	return makeVec3(0, 0, 0)
}

type metal struct {
	albedo vec3
	fuzz   float64
}

func (m metal) scatter(r *ray, rec *hitRecord) (absorbed bool, scattered *ray, attenuation vec3) {
	reflected := r.direction.unitVector().reflect(rec.normal)
	scatteredRay := makeTimedRay(rec.p, reflected.add(randUnitVector().multiplyScalar(m.fuzz)), r.time)
	absorbed = dot(scatteredRay.direction, rec.normal) <= 0

	return absorbed, &scatteredRay, m.albedo
}

func (m metal) emitted(u, v float64, p vec3) vec3 {
	return makeVec3(0, 0, 0)
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
	cosTheta := min(dot(unitDirection.negate(), rec.normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := refractionRatio*sinTheta > 1.0

	var direction vec3

	if cannotRefract || reflectance(cosTheta, refractionRatio) > rand() {
		direction = unitDirection.reflect(rec.normal)
	} else {
		direction = unitDirection.refract(rec.normal, refractionRatio)
	}

	scatteredRay := makeTimedRay(rec.p, direction, r.time)

	return false, &scatteredRay, attenuation
}

func reflectance(cosine, refIdx float64) float64 {
	// Use Schlick's approximation for reflectance
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}

func (d dielectric) emitted(u, v float64, p vec3) vec3 {
	return makeVec3(0, 0, 0)
}

type diffuseLight struct {
	emit texture
}

func (d diffuseLight) scatter(r *ray, rec *hitRecord) (absorbed bool, scattered *ray, attenuation vec3) {
	return true, nil, makeVec3(0, 0, 0)
}

func (d diffuseLight) emitted(u, v float64, p vec3) vec3 {
	return d.emit.value(u, v, p)
}
