package main

type hittable interface {
	hit(r Ray, t interval, rec *hitRecord) bool
}
