package main

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	imageWidth, imageHeight int

	image       *ebiten.Image
	imagePixels []vec3
	imageDirty  bool
}

func makeGame() *Game {
	imageHeight := calculateImageHeight(imageWidth, aspectRatio)

	return &Game{
		imageWidth:  imageWidth,
		imageHeight: imageHeight,

		image:       ebiten.NewImage(imageWidth, imageHeight),
		imagePixels: make([]vec3, imageWidth*imageHeight),
		imageDirty:  true,
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.imageWidth, g.imageHeight
}

func (g *Game) Update() error {
	if g.imageDirty {

		for x := 0; x < g.imageWidth; x++ {
			for y := 0; y < g.imageHeight; y++ {
				index := getPixelIndex(x, y, g.imageWidth)
				g.image.Set(x, g.imageHeight-y, g.imagePixels[index])
			}
		}

		g.imageDirty = false
	}

	return nil
}

func (g *Game) SetPixels(imagePixels *[]vec3) {
	g.imagePixels = *imagePixels
	g.imageDirty = true
}

func (g *Game) Draw(image *ebiten.Image) {
	image.DrawImage(g.image, nil)
}

func (v vec3) RGBA() (r, g, b, a uint32) {
	ir := uint32(65534.999 * v.e[0])
	ig := uint32(65534.999 * v.e[1])
	ib := uint32(65534.999 * v.e[2])
	return ir, ig, ib, 1
}
