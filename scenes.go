package main

func makeTripleSphereScene() (hittable, cameraParameters) {

	material_ground := lambertian{makeVec3(0.8, 0.8, 0.0)}

	material_center := lambertian{makeVec3(0.1, 0.2, 0.5)}
	material_left := dielectric{1.5}
	material_right := metal{makeVec3(0.8, 0.6, 0.2), 0.0}

	world := makeHittableList()
	world.add(makeSphere(makeVec3(0, -100.5, -1), 100, material_ground))
	world.add(makeSphere(makeVec3(0, 0, -1), 0.5, material_center))
	world.add(makeSphere(makeVec3(-1, 0, -1), 0.5, material_left))
	world.add(makeSphere(makeVec3(-1, 0, -1), -0.4, material_left))
	world.add(makeSphere(makeVec3(1, 0, -1), 0.5, material_right))

	parameters := cameraParameters{
		vFov:     30,
		lookFrom: makeVec3(-2, 2, 1),
		lookAt:   makeVec3(0, 0, -1),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 10,
		focusDist:    3.4,
	}

	return world, parameters
}

func makeBook1CoverScene() (hittable, cameraParameters) {

	world := makeHittableList()

	// ground
	world.add(makeSphere(makeVec3(0, -1000, 0), 1000, lambertian{makeVec3(0.5, 0.5, 0.5)}))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand()
			center := makeVec3(float64(a)+0.9*rand(), 0.2, float64(b)+0.9*rand())

			if center.subtract(makeVec3(4, 0.2, 0)).length() > 0.9 {
				if chooseMat < 0.8 {
					// diffuse
					world.add(makeSphere(center, 0.2, lambertian{randVec().multiply(randVec())}))
				} else if chooseMat < 0.95 {
					// metal
					world.add(makeSphere(center, 0.2, metal{randVecB(0.5, 1), randb(0, 0.5)}))
				} else {
					// glass
					world.add(makeSphere(center, 0.2, dielectric{1.5}))
				}
			}
		}
	}

	world.add(makeSphere(makeVec3(0, 1, 0), 1.0, dielectric{1.5}))
	world.add(makeSphere(makeVec3(-4, 1, 0), 1.0, lambertian{makeVec3(0.4, 0.2, 0.1)}))
	world.add(makeSphere(makeVec3(4, 1, 0), 1.0, metal{makeVec3(0.7, 0.6, 0.5), 0.0}))

	parameters := cameraParameters{
		vFov:     30,
		lookFrom: makeVec3(13, 2, 3),
		lookAt:   makeVec3(0, 0, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.6,
		focusDist:    10.0,
	}

	return world, parameters
}
