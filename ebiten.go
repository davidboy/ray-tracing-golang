package main

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	imageWidth, imageHeight int
	imagePixels             []Vec3
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.imageWidth, g.imageHeight
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(image *ebiten.Image) {
	image.Clear()

	for x := 0; x < g.imageWidth; x++ {
		for y := 0; y < g.imageHeight; y++ {
			index := getPixelIndex(x, y, g.imageWidth)
			image.Set(x, y, g.imagePixels[index])
		}
	}
}

func (v Vec3) RGBA() (r, g, b, a uint32) {
	ir := uint32(65534.999 * v.e[0])
	ig := uint32(65534.999 * v.e[1])
	ib := uint32(65534.999 * v.e[2])
	return ir, ig, ib, 1
}
