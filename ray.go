package main

type Ray struct {
	Origin    Vec3
	Direction Vec3
}

func MakeRay(origin, direction Vec3) Ray {
	return Ray{origin, direction}
}

func (r Ray) At(t float64) Vec3 {
	return r.Origin.Add(r.Direction.MultiplyScalar(t))
}
