package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenSizeX = TileSize * (Width + 2)
	screenSizeY = TileSize * (Height + 2)
	Width       = 15
	Height      = 12
	TileSize    = 50
)

var (
	White  = color.RGBA{255, 255, 255, 255}
	Tile   = color.RGBA{77, 77, 77, 255}
	NGoal  = color.RGBA{255, 150, 150, 255}
	YGoal  = color.RGBA{65, 130, 35, 255}
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
	_Box     = []Box{}
	_Goal    = []Goal{}
	_Steps   = [][]Step{}
	TempStep = []Step{}
	Color    = 0

	Gameover = false
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

type Step struct {
	_type string
	idx   int
	x     int
	y     int
	goal  bool
}

func (g *Game) Update() error {
	if Gameover {
		return nil
	}

	var directionX int
	var directionY int

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		directionX = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		directionX = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		directionY = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		directionY = 1
	}

	if !ebiten.IsKeyPressed(ebiten.KeyU) && !ebiten.IsKeyPressed(ebiten.KeyR) && (directionX != 0 || directionY != 0) {
		TempStep = append(TempStep, Step{_type: "player-move", x: PlayerPosition["x"], y: PlayerPosition["y"]})
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyU) && len(_Steps) > 0 {
		step := _Steps[len(_Steps)-1]

		for _, st := range step {
			switch st._type {
			case "player-move":
				PlayerPosition["x"] = st.x
				PlayerPosition["y"] = st.y
			case "box-move":
				_Box[st.idx].x = st.x
				_Box[st.idx].y = st.y
			case "box-goal":
				_Box[st.idx].setGoal(st.goal)
			}
		}

		_Steps = _Steps[:len(_Steps)-1]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		PlayerPosition = map[string]int{
			"x": 2, "y": 3,
		}
		_Steps = [][]Step{}
		TempStep = []Step{}
		Color = rand.Intn(7)
		initialization()
	}

	PlayerPosition["x"] += directionX
	PlayerPosition["y"] += directionY

	for idx := 0; idx < len(_Box); idx++ {
		newbox := BoxFilter(_Box, func(box Box) bool {
			return box.x == _Box[idx].x+directionX && box.y == _Box[idx].y+directionY
		})

		if _Box[idx].x == PlayerPosition["x"] && _Box[idx].y == PlayerPosition["y"] && len(newbox) == 0 {
			TempStep = append(TempStep, Step{_type: "box-move", idx: idx, x: _Box[idx].x, y: _Box[idx].y})
			_Box[idx].move(directionX, directionY)

			if _Box[idx].x < 0 {
				PlayerPosition["x"] -= directionX
				_Box[idx].x = 0
			}

			if _Box[idx].x > Width-1 {
				PlayerPosition["x"] -= directionX
				_Box[idx].x = Width - 1
			}

			if _Box[idx].y < 0 {
				PlayerPosition["y"] -= directionY
				_Box[idx].y = 0
			}

			if _Box[idx].y > Height-1 {
				PlayerPosition["y"] -= directionY
				_Box[idx].y = Height - 1
			}
		}
	}

	colidbox := BoxFilter(_Box, func(box Box) bool {
		return box.x == PlayerPosition["x"] && box.y == PlayerPosition["y"]
	})
	if len(colidbox) != 0 {
		PlayerPosition["x"] -= directionX
		PlayerPosition["y"] -= directionY
	}

	if PlayerPosition["x"] < 0 {
		PlayerPosition["x"] = 0
	}
	if PlayerPosition["x"] > Width-1 {
		PlayerPosition["x"] = Width - 1
	}
	if PlayerPosition["y"] < 0 {
		PlayerPosition["y"] = 0
	}
	if PlayerPosition["y"] > Height-1 {
		PlayerPosition["y"] = Height - 1
	}

	cancelGoal()
	Gameover = checkWin()

	if len(TempStep) > 0 {
		_Steps = append(_Steps, TempStep)
	}
	TempStep = []Step{}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	BoardClear()
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
			_x := float64(50 * (x + 1))
			_y := float64(50 * (y + 1))

			if x == -1 || y == -1 || x == Width || y == Height {
				// 테두리
				ebitenutil.DrawRect(screen, _x, _y, 50, 50, getColor(Color))
			} else if x == PlayerPosition["x"] && y == PlayerPosition["y"] {
				// 플레이어
				ebitenutil.DrawRect(screen, _x, _y, 50, 50, Player)
			} else {
				ebitenutil.DrawRect(screen, _x, _y, 50, 50, getBlock(Board[y][x], Color))
			}
		}
	}

	image, _, err := ebitenutil.NewImageFromFile("./assets/game-clear.png")
	if err != nil {
		log.Fatal(err)
	} else if Gameover {
		sizeX := float64(image.Bounds().Size().X / 2)
		sizeY := float64(image.Bounds().Size().Y / 2)

		screenCenterX := float64(screenSizeX / 2)
		screenCenterY := float64(screenSizeY / 2)

		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenCenterX-sizeX, screenCenterY-sizeY)

		screen.DrawImage(image, &op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}

func (box *Box) move(x, y int) {
	box.x += x
	box.y += y
}

func (box *Box) setGoal(goal bool) {
	box.goal = goal
}

func NewBox(point [2]int) {
	_Box = append(_Box, Box{goal: false, x: point[0], y: point[1]})
}

func NewGoal(point [2]int) {
	_Goal = append(_Goal, Goal{x: point[0], y: point[1]})
}

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

func checkWin() bool {
	var stack int
	for idx := 0; idx < len(_Box); idx++ {
		for _, goal := range _Goal {
			if goal.x == _Box[idx].x && goal.y == _Box[idx].y {
				stack++
				if len(TempStep) > 0 {
					TempStep = append(TempStep, Step{_type: "box-goal", idx: idx, goal: _Box[idx].goal})
				}
				_Box[idx].setGoal(true)
			}
		}
	}

	return len(_Goal) == stack
}

func cancelGoal() {
	for idx := 0; idx < len(_Box); idx++ {
		for _, goal := range _Goal {
			if _Box[idx].goal && !(goal.x == _Box[idx].x && goal.y == _Box[idx].y) {
				_Box[idx].setGoal(false)
			}
		}
	}
}

func BoxFilter(vs []Box, f func(Box) bool) []Box {
	vsf := make([]Box, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func BoardClear() {
	Board = [Height][Width]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func initialization() {
	_Box = []Box{}
	_Goal = []Goal{}

	points := [...][2]int{{1, 3}, {3, 2}, {7, 4}, {4, 3}, {3, 1}, {9, 5}, {12, 10}, {3, 11}, {14, 4}, {3, 3}, {1, 10}, {10, 2}, {10, 6}, {7, 10}}
	for _, point := range points {
		NewBox(point)
	}

	points = [...][2]int{{2, 4}, {1, 2}, {9, 4}, {7, 3}, {2, 7}, {1, 6}, {14, 2}, {4, 11}, {2, 9}, {5, 3}, {14, 1}, {10, 3}, {10, 11}, {14, 11}}
	for _, point := range points {
		NewGoal(point)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
	Color = rand.Intn(7)

	initialization()
}

func main() {
	ebiten.SetWindowSize(screenSizeX, screenSizeY)
	ebiten.SetWindowTitle("soGOban")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
