package main

func createRender(world hittable, parameters cameraParameters, threads int) *render {

	camera := makeCamera(parameters, &quality)
	render := camera.render(world)

	tasks := make([](chan bool), threads)

	for i := 0; i < threads; i++ {
		taskFinished := make(chan bool)
		tasks[i] = taskFinished

		go render.run(max(1, quality.samples/threads), i, taskFinished)
	}

	go func() {
		for _, taskFinished := range tasks {
			<-taskFinished
		}

		close(render.finished)
	}()

	return render
}
