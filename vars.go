package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenSizeX = tileSize * (width + 2)
	screenSizeY = tileSize * (height + 2)
	width       = 15
	height      = 12
	tileSize    = 50
)

const (
	playerMove = iota
	boxMove
	boxGoal
)

const (
	none = iota
	tplayer
	ttile
	tngoal
	tygoal
)

var colors = [...]color.RGBA{red, orange, yellow, green, blue, purple, brown}
var textures = [...]*ebiten.Image{nil, nil, nil, nil, nil}
var (
	background = color.RGBA{120, 120, 120, 255}
	red        = color.RGBA{255, 0, 0, 255}
	orange     = color.RGBA{255, 125, 0, 255}
	yellow     = color.RGBA{255, 225, 0, 255}
	green      = color.RGBA{55, 175, 20, 255}
	blue       = color.RGBA{20, 95, 175, 255}
	purple     = color.RGBA{85, 20, 175, 255}
	brown      = color.RGBA{110, 60, 20, 255}

	playerPoint = [2]int{2, 3}
	board       = [height][width]int{}
	box         = []Box{}
	goal        = []Goal{}
	steps       = [][]Step{}
	step        = []Step{}
	clr         = 0

	player *ebiten.Image
	tile   *ebiten.Image
	ngoal  *ebiten.Image
	ygoal  *ebiten.Image
)
