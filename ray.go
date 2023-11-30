package main

type ray struct {
	origin    vec3
	direction vec3
}

func makeRay(origin, direction vec3) ray {
	return ray{origin, direction}
}

func (r ray) at(t float64) vec3 {
	return r.origin.add(r.direction.multiplyScalar(t))
}
