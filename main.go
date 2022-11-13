package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(screenSizeX, screenSizeY)
	ebiten.SetWindowTitle("soGOban")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
