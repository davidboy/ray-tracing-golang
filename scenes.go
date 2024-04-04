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

		backgroundColor: makeVec3(0.7, 0.8, 1.0),
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

		backgroundColor: makeVec3(0.7, 0.8, 1.0),
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

		backgroundColor: makeVec3(0.7, 0.8, 1.0),
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

		backgroundColor: makeVec3(0.7, 0.8, 1.0),
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

		backgroundColor: makeVec3(0.7, 0.8, 1.0),
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

		backgroundColor: makeVec3(0.7, 0.8, 1.0),
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

		backgroundColor: makeVec3(0.7, 0.8, 1.0),
	}

	return &world, parameters
}

func makeSimpleLightScene() (*hittableList, cameraParameters) {
	world := makeHittableList()

	pertext := makeNoiseTexture(4)
	world.add(makeSphere(makeVec3(0, -1000, 0), 1000, lambertian{pertext}))
	world.add(makeSphere(makeVec3(0, 2, 0), 2, lambertian{pertext}))

	difflight := diffuseLight{makeColorTexture(4, 4, 4)}
	world.add(makeSphere(makeVec3(0, 7, 0), 2, difflight))
	world.add(makeQuad(makeVec3(3, 1, -2), makeVec3(2, 0, 0), makeVec3(0, 2, 0), difflight))

	parameters := cameraParameters{
		vFov:     20,
		lookFrom: makeVec3(26, 3, 6),
		lookAt:   makeVec3(0, 2, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.0,
		focusDist:    10.0,

		backgroundColor: makeVec3(0, 0, 0),
	}

	return &world, parameters
}

func makeCornellBoxScene() (*hittableList, cameraParameters) {
	world := makeHittableList()

	red := lambertian{makeColorTexture(0.65, 0.05, 0.05)}
	white := lambertian{makeColorTexture(0.73, 0.73, 0.73)}
	green := lambertian{makeColorTexture(0.12, 0.45, 0.15)}
	light := diffuseLight{makeColorTexture(15, 15, 15)}

	world.add(makeQuad(makeVec3(555, 0, 0), makeVec3(0, 555, 0), makeVec3(0, 0, 555), green))
	world.add(makeQuad(makeVec3(0, 0, 0), makeVec3(0, 555, 0), makeVec3(0, 0, 555), red))
	world.add(makeQuad(makeVec3(343, 554, 332), makeVec3(0-130, 0, 0), makeVec3(0, 0, -105), light))
	world.add(makeQuad(makeVec3(0, 0, 0), makeVec3(555, 0, 0), makeVec3(0, 0, 555), white))
	world.add(makeQuad(makeVec3(555, 555, 555), makeVec3(-555, 0, 0), makeVec3(0, 0, -555), white))
	world.add(makeQuad(makeVec3(0, 0, 555), makeVec3(555, 0, 0), makeVec3(0, 555, 0), white))

	world.add(makeBox(makeVec3(130, 0, 65), makeVec3(295, 165, 230), white))
	world.add(makeBox(makeVec3(265, 0, 295), makeVec3(430, 330, 460), white))

	parameters := cameraParameters{
		vFov:     40,
		lookFrom: makeVec3(278, 278, -800),
		lookAt:   makeVec3(278, 278, 0),
		vUp:      makeVec3(0, 1, 0),

		defocusAngle: 0.0,
		focusDist:    10.0,

		backgroundColor: makeVec3(0, 0, 0),
	}

	return &world, parameters
}

func makeBox(a, b vec3, mat material) hittableList {
	sides := makeHittableList()

	// TODO: access elements instead of using methods

	min := makeVec3(min(a.x(), b.x()), min(a.y(), b.y()), min(a.z(), b.z()))
	max := makeVec3(max(a.x(), b.x()), max(a.y(), b.y()), max(a.z(), b.z()))

	dx := makeVec3(max.x()-min.x(), 0, 0)
	dy := makeVec3(0, max.y()-min.y(), 0)
	dz := makeVec3(0, 0, max.z()-min.z())

	sides.add(makeQuad(makeVec3(min.x(), min.y(), max.z()), dx, dy, mat))          // front
	sides.add(makeQuad(makeVec3(max.x(), min.y(), max.z()), dz.negate(), dy, mat)) // right
	sides.add(makeQuad(makeVec3(max.x(), min.y(), min.z()), dx.negate(), dy, mat)) // back
	sides.add(makeQuad(makeVec3(min.x(), min.y(), min.z()), dz, dy, mat))          // left
	sides.add(makeQuad(makeVec3(min.x(), max.y(), max.z()), dx, dz.negate(), mat)) // top
	sides.add(makeQuad(makeVec3(min.x(), min.y(), min.z()), dx, dz, mat))          // bottom

	return sides
}
