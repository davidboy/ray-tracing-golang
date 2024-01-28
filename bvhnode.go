package main

import (
	"sort"
)

type bvhNode struct {
	left  hittable
	right hittable
	box   aabb
}

func makeBvhNode(srcObjects []hittable, start, end int) bvhNode {

	objects := make([]hittable, len(srcObjects))
	copy(objects, srcObjects)

	axis := randIntB(0, 3)

	var comparator func(a, b hittable) bool

	if axis == 0 {
		comparator = boxXCompare
	} else if axis == 1 {
		comparator = boxYCompare
	} else {
		comparator = boxZCompare
	}

	objectSpan := end - start

	var left, right hittable

	if objectSpan == 1 {
		left = objects[start]
		right = objects[start]

		// fmt.Println("Making single object bvh node")
	} else if objectSpan == 2 {
		if comparator(objects[start], objects[start+1]) {
			left = objects[start]
			right = objects[start+1]
		} else {
			left = objects[start+1]
			right = objects[start]
		}

		// fmt.Println("Making two object bvh node")
	} else {

		sort.Slice(objects[start:end], func(i, j int) bool {
			return comparator(objects[i], objects[j])
		})

		mid := start + objectSpan/2
		left = makeBvhNode(objects, start, mid)
		right = makeBvhNode(objects, mid, end)

		// fmt.Printf("Making multi object bvh node; left = %d; right = %d\n", mid-start, end-mid)
	}

	box := makeAABBFromBoxes(left.boundingBox(), right.boundingBox())

	return bvhNode{left, right, box}
}

func makeBvhNodeFromList(list *hittableList) bvhNode {
	return makeBvhNode(list.hittables, 0, len(list.hittables))
}

func (n bvhNode) hit(r ray, ray_t interval, rec *hitRecord) bool {
	if !n.box.hit(r, ray_t) {
		return false
	}

	hitLeft := n.left.hit(r, ray_t, rec)

	var rightRayTMax float64
	if hitLeft {
		rightRayTMax = rec.t
	} else {
		rightRayTMax = ray_t.max
	}

	rightRayT := makeInterval(ray_t.min, rightRayTMax)
	hitRight := n.right.hit(r, rightRayT, rec)

	return hitLeft || hitRight
}

func (n bvhNode) boundingBox() aabb {
	return n.box
}

func boxCompare(a, b hittable, axisIndex int) bool {
	return a.boundingBox().axis(axisIndex).min < b.boundingBox().axis(axisIndex).min
}

func boxXCompare(a, b hittable) bool {
	return boxCompare(a, b, 0)
}

func boxYCompare(a, b hittable) bool {
	return boxCompare(a, b, 1)
}

func boxZCompare(a, b hittable) bool {
	return boxCompare(a, b, 2)
}
