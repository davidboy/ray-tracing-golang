package main

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	render *render
	image  *ebiten.Image
}

func makeGame(r *render) *Game {

	imageWidth := r.c.imageWidth
	imageHeight := r.c.imageHeight

	return &Game{
		render: r,
		image:  ebiten.NewImage(imageWidth, imageHeight),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.render.c.imageWidth, g.render.c.imageHeight
}

func (g *Game) Update() error {

	if g.render.dirty {
		pixels := g.render.squash()

		for x := 0; x < g.render.c.imageWidth; x++ {
			for y := 0; y < g.render.c.imageHeight; y++ {
				index := getPixelIndex(x, y, g.render.c.imageWidth)
				g.image.Set(x, g.render.c.imageHeight-y, pixels[index])
			}
		}

	}

	return nil
}

func (g *Game) Draw(image *ebiten.Image) {
	image.DrawImage(g.image, nil)
}

func (v vec3) RGBA() (r, g, b, a uint32) {
	ir := uint32(65534.999 * v.e[0])
	ig := uint32(65534.999 * v.e[1])
	ib := uint32(65534.999 * v.e[2])
	return ir, ig, ib, 0xffff
}
