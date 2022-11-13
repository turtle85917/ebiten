package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenSizeX = 320
	screenSizeY = 320
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(ebiten.NewImage(5, 5), 0, 0, 5, 5, color.RGBA{255, 255, 255, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}

func main() {
	ebiten.SetWindowSize(screenSizeX, screenSizeY)
	ebiten.SetWindowTitle("soGOban")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
