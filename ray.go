package main

type ray struct {
	origin    vec3
	direction vec3
	time      float64
}

func makeRay(origin, direction vec3) ray {
	return ray{origin, direction, 0.0}
}

func makeTimedRay(origin, direction vec3, time float64) ray {
	return ray{origin, direction, time}
}

func (r ray) at(t float64) vec3 {
	return r.origin.add(r.direction.multiplyScalar(t))
}
