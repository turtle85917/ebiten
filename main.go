package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenSizeX = TileSize * (Width + 1)
	screenSizeY = TileSize * (Height + 1)
	Width       = 15
	Height      = 12
	TileSize    = 50
)

var (
	White  = color.RGBA{255, 255, 255, 255}
	Tile   = color.RGBA{77, 77, 77, 255}
	NGoal  = color.RGBA{255, 0, 190, 255}
	YGoal  = color.RGBA{150, 255, 0, 255}
	Player = color.RGBA{225, 227, 167, 255}

	Red    = color.RGBA{255, 0, 0, 255}
	Orange = color.RGBA{255, 125, 0, 255}
	Yellow = color.RGBA{255, 225, 0, 255}
	Green  = color.RGBA{55, 175, 20, 255}
	Blue   = color.RGBA{20, 95, 175, 255}
	Purple = color.RGBA{85, 20, 175, 255}
	Brown  = color.RGBA{110, 60, 20, 255}

	PlayerPosition = map[string]int{
		"x": 2, "y": 3,
	}
	Board = [Height][Width]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}
	_Box  = []Box{}
	_Goal = []Goal{}
	Color = 0
)

type Game struct{}

type Box struct {
	goal bool
	x    int
	y    int
}

type Goal struct {
	x int
	y int
}

// type Step struct {
// 	_type string
// 	idx   int
// 	x     int
// 	y     int
// 	goal  bool
// }

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(White)

	for _, ga := range _Goal {
		Board[ga.y][ga.x] = 2
	}

	for _, bx := range _Box {
		if bx.goal {
			Board[bx.y][bx.x] = 3
		} else {
			Board[bx.y][bx.x] = 1
		}
	}

	for y := -1; y < Height+1; y++ {
		for x := -1; x < Width+1; x++ {
			if x == -1 || y == -1 || x == Width || y == Height {
				// 테두리
				ebitenutil.DrawRect(screen, float64(50*(x+1)), float64(50*(y+1)), 50, 50, getColor(Color))
			} else if x == PlayerPosition["x"] && y == PlayerPosition["y"] {
				// 플레이어
				ebitenutil.DrawRect(screen, float64(50*(x+1)), float64(50*(y+1)), 50, 50, Player)
			} else {
				ebitenutil.DrawRect(screen, float64(50*(x+1)), float64(50*(y+1)), 50, 50, getBlock(Board[y][x], Color))
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}

// func (bx *Box) move(x, y int) {
// 	bx.x += x
// 	bx.y += y
// }

// func (bx *Box) setGoal(goal bool) {
// 	bx.goal = goal
// }

func getColor(color int) color.RGBA {
	switch color {
	case 0:
		return Red
	case 1:
		return Orange
	case 2:
		return Yellow
	case 3:
		return Green
	case 4:
		return Blue
	case 5:
		return Purple
	case 6:
		return Brown
	}

	return Tile
}

func getBlock(tile int, color int) color.RGBA {
	switch tile {
	case 0:
		return Tile
	case 1:
		return getColor(color)
	case 2:
		return NGoal
	case 3:
		return YGoal
	}

	return Tile
}

func reset(box *[]Box, goal *[]Goal) {
	*box = []Box{}
	*goal = []Goal{}

	*box = append(*box, Box{goal: false, x: 1, y: 3})
	*box = append(*box, Box{goal: false, x: 3, y: 2})
	*box = append(*box, Box{goal: false, x: 7, y: 4})
	*box = append(*box, Box{goal: false, x: 4, y: 3})
	*box = append(*box, Box{goal: false, x: 3, y: 1})
	*box = append(*box, Box{goal: false, x: 9, y: 5})
	*box = append(*box, Box{goal: false, x: 12, y: 10})
	*box = append(*box, Box{goal: false, x: 3, y: 11})
	*box = append(*box, Box{goal: false, x: 14, y: 4})
	*box = append(*box, Box{goal: false, x: 3, y: 3})
	*box = append(*box, Box{goal: false, x: 1, y: 10})
	*box = append(*box, Box{goal: false, x: 10, y: 2})
	*box = append(*box, Box{goal: false, x: 10, y: 6})
	*box = append(*box, Box{goal: false, x: 7, y: 10})

	*goal = append(*goal, Goal{x: 2, y: 4})
	*goal = append(*goal, Goal{x: 1, y: 2})
	*goal = append(*goal, Goal{x: 9, y: 4})
	*goal = append(*goal, Goal{x: 7, y: 3})
	*goal = append(*goal, Goal{x: 2, y: 7})
	*goal = append(*goal, Goal{x: 1, y: 5})
	*goal = append(*goal, Goal{x: 14, y: 2})
	*goal = append(*goal, Goal{x: 4, y: 11})
	*goal = append(*goal, Goal{x: 2, y: 9})
	*goal = append(*goal, Goal{x: 5, y: 3})
	*goal = append(*goal, Goal{x: 14, y: 1})
	*goal = append(*goal, Goal{x: 10, y: 3})
	*goal = append(*goal, Goal{x: 10, y: 11})
	*goal = append(*goal, Goal{x: 14, y: 11})
}

func init() {
	rand.Seed(time.Now().UnixNano())
	Color = rand.Intn(7)

	reset(&_Box, &_Goal)
}

func main() {
	ebiten.SetWindowSize(screenSizeX, screenSizeY)
	ebiten.SetWindowTitle("soGOban")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
