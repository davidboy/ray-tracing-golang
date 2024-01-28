package main

type aabb struct {
	x, y, z interval
}

func makeAABB(x, y, z interval) aabb {
	return aabb{x, y, z}
}

func makeEmptyAABB() aabb {
	return makeAABB(empty, empty, empty)
}

func makeAABBFromExtremaPoints(a, b vec3) aabb {
	return makeAABB(
		makeInterval(min(a.e[0], b.e[0]), max(a.e[0], b.e[0])),
		makeInterval(min(a.e[1], b.e[1]), max(a.e[1], b.e[1])),
		makeInterval(min(a.e[2], b.e[2]), max(a.e[2], b.e[2])),
	)
}

func makeAABBFromBoxes(box1, box2 aabb) aabb {
	return makeAABB(
		// makeInterval(box1.x, box1.x),
		// makeInterval(box1.y, box1.y),
		// makeInterval(box1.z, box1.z),
		makeInterval(min(box1.x.min, box2.x.min), max(box1.x.max, box2.x.max)),
		makeInterval(min(box1.y.min, box2.y.min), max(box1.y.max, box2.y.max)),
		makeInterval(min(box1.z.min, box2.z.min), max(box1.z.max, box2.z.max)),
	)
}

func (a aabb) axis(n int) interval {
	if n == 1 {
		return a.y
	}

	if n == 2 {
		return a.z
	}

	return a.x
}

func (a aabb) hit(r ray, ray_t interval) bool {

	// for i := 0; i < 3; i++ {
	// 	t0 := min((a.axis(i).min-r.origin.e[i])/r.direction.e[i], (a.axis(i).max-r.origin.e[i])/r.direction.e[i])
	// 	t1 := max((a.axis(i).min-r.origin.e[i])/r.direction.e[i], (a.axis(i).max-r.origin.e[i])/r.direction.e[i])

	// 	ray_t.min = max(t0, ray_t.min)
	// 	ray_t.max = min(t1, ray_t.max)
	// 	if ray_t.max <= ray_t.min {
	// 		return false
	// 	}
	// }

	// return true

	//////////////////////////////

	for i := 0; i < 3; i++ {
		invD := 1 / r.direction.e[i]
		orig := r.origin.e[i]

		t0 := (a.axis(i).min - orig) * invD
		t1 := (a.axis(i).max - orig) * invD

		if invD < 0 {
			t0, t1 = t1, t0
		}

		if t0 > ray_t.min {
			ray_t.min = t0
		}

		if t1 < ray_t.max {
			ray_t.max = t1
		}

		if ray_t.max <= ray_t.min {
			return false
		}
	}

	return true
}
