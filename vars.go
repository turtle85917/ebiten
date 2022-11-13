package main

import "image/color"

var colors = [...]color.RGBA{red, orange, yellow, green, blue, purple, brown}
var (
	white  = color.RGBA{255, 255, 255, 255}
	tile   = color.RGBA{77, 77, 77, 255}
	ngoal  = color.RGBA{255, 150, 150, 255}
	ygoal  = color.RGBA{65, 130, 35, 255}
	player = color.RGBA{225, 227, 167, 255}

	red    = color.RGBA{255, 0, 0, 255}
	orange = color.RGBA{255, 125, 0, 255}
	yellow = color.RGBA{255, 225, 0, 255}
	green  = color.RGBA{55, 175, 20, 255}
	blue   = color.RGBA{20, 95, 175, 255}
	purple = color.RGBA{85, 20, 175, 255}
	brown  = color.RGBA{110, 60, 20, 255}

	playerPoint = [2]int{2, 3}
	board       = [height][width]int{}
	box         = []Box{}
	goal        = []Goal{}
	steps       = [][]Step{}
	step        = []Step{}
	clr         = 0
)
