package main

type hittableList struct {
	hittables []hittable
}

func makeHittableList() hittableList {
	return hittableList{make([]hittable, 0)}
}

func (l *hittableList) clear() {
	l.hittables = make([]hittable, 0)
}

func (l *hittableList) add(h hittable) {
	l.hittables = append(l.hittables, h)
}

func (l hittableList) hit(r ray, t interval, rec *hitRecord) bool {

	hitAnything := false
	closestSoFar := t.max

	for _, object := range l.hittables {
		if object.hit(r, makeInterval(t.min, closestSoFar), rec) {
			hitAnything = true
			closestSoFar = rec.t
		}
	}

	return hitAnything
}
