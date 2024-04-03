package main

func makeTripleSphereScene() (*hittableList, cameraParameters) {

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

	return &world, parameters
}

func makeBook1CoverScene() (*hittableList, cameraParameters) {

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
		vFov:     20,
		lookFrom: makeVec3(13, 2, 3),
		lookAt:   makeVec3(0, 0, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.6,
		focusDist:    10.0,
	}

	return &world, parameters
}

func makeBook1CoverSceneWithMotion() (*hittableList, cameraParameters) {

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

					albedo := randVec().multiply(randVec())
					center2 := center.add(makeVec3(0, randb(0, 0.5), 0))

					world.add(makeMovingSphere(center, center2, 0.2, lambertian{albedo}))

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

		defocusAngle: 0.02,
		focusDist:    10.0,
	}

	return &world, parameters
}

func makeBook1CoverSceneWithCheckerTexture() (*hittableList, cameraParameters) {

	world := makeHittableList()

	// ground

	checker := makeCheckerTexture(0.32, makeColorTexture(0.2, 0.3, 0.1), makeColorTexture(0.9, 0.9, 0.9))
	world.add(makeSphere(makeVec3(0, -1000, 0), 1000, lambertian{checker}))

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
		vFov:     20,
		lookFrom: makeVec3(13, 2, 3),
		lookAt:   makeVec3(0, 0, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.6,
		focusDist:    10.0,
	}

	return &world, parameters
}

func makeTwoSpheresScene() (*hittableList, cameraParameters) {

	world := makeHittableList()

	checker := makeCheckerTexture(0.8, makeColorTexture(0.2, 0.3, 0.1), makeColorTexture(0.9, 0.9, 0.9))

	world.add(makeSphere(makeVec3(0, -10, 0), 10, lambertian{checker}))
	world.add(makeSphere(makeVec3(0, 10, 0), 10, lambertian{checker}))

	parameters := cameraParameters{
		vFov:     35,
		lookFrom: makeVec3(13, 2, 3),
		lookAt:   makeVec3(0, 0, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.0,
		focusDist:    10.0,
	}

	return &world, parameters
}

func makeTwoPerlinSpheresScene() (*hittableList, cameraParameters) {

	world := makeHittableList()

	pertext := makeNoiseTexture(4)
	world.add(makeSphere(makeVec3(0, -1000, 0), 1000, lambertian{pertext}))
	world.add(makeSphere(makeVec3(0, 2, 0), 2, lambertian{pertext}))

	parameters := cameraParameters{
		vFov:     20,
		lookFrom: makeVec3(13, 2, 3),
		lookAt:   makeVec3(0, 0, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.0,
		focusDist:    10.0,
	}

	return &world, parameters
}

func makeSimpleQuadsScene() (*hittableList, cameraParameters) {
	world := makeHittableList()

	leftRed := lambertian{makeColorTexture(1.0, 0.2, 0.2)}
	backGreen := lambertian{makeColorTexture(0.2, 1.0, 0.2)}
	rightBlue := lambertian{makeColorTexture(0.2, 0.2, 1.0)}
	upperOrange := lambertian{makeColorTexture(1.0, 0.5, 0.0)}
	lowerTeal := lambertian{makeColorTexture(0.2, 0.8, 0.8)}

	world.add(makeQuad(makeVec3(-3, -2, 5), makeVec3(0, 0, -4), makeVec3(0, 4, 0), leftRed))
	world.add(makeQuad(makeVec3(-2, -2, 0), makeVec3(4, 0, 0), makeVec3(0, 4, 0), backGreen))
	world.add(makeQuad(makeVec3(3, -2, 1), makeVec3(0, 0, 4), makeVec3(0, 4, 0), rightBlue))
	world.add(makeQuad(makeVec3(-2, 3, 1), makeVec3(4, 0, 0), makeVec3(0, 0, 4), upperOrange))
	world.add(makeQuad(makeVec3(-2, -3, 5), makeVec3(4, 0, 0), makeVec3(0, 0, -4), lowerTeal))

	parameters := cameraParameters{
		vFov:     80,
		lookFrom: makeVec3(0, 0, 9),
		lookAt:   makeVec3(0, 0, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.0,
		focusDist:    10.0,
	}

	return &world, parameters
}
