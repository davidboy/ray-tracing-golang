package main

type hittable interface {
	hit(r ray, t interval, rec *hitRecord) bool
	boundingBox() aabb
}
