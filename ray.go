package main

type Ray struct {
	Origin    vec3
	Direction vec3
}

func MakeRay(origin, direction vec3) Ray {
	return Ray{origin, direction}
}

func (r Ray) At(t float64) vec3 {
	return r.Origin.add(r.Direction.multiplyScalar(t))
}
