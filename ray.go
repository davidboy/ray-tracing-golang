package main

type Ray struct {
	origin    vec3
	direction vec3
}

func makeRay(origin, direction vec3) Ray {
	return Ray{origin, direction}
}

func (r Ray) at(t float64) vec3 {
	return r.origin.add(r.direction.multiplyScalar(t))
}
