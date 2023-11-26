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

func (l hittableList) hit(r Ray, t interval, rec *hitRecord) bool {
	var tempRec hitRecord

	hitAnything := false
	closestSoFar := t.max

	for _, hittable := range l.hittables {
		if hittable.hit(r, makeInterval(t.min, closestSoFar), &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.t
			*rec = tempRec
		}
	}

	return hitAnything
}
