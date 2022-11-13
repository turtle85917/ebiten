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
	screenSizeX = tileSize * (width + 2)
	screenSizeY = tileSize * (height + 2)
	width       = 15
	height      = 12
	tileSize    = 50
)

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

type Game struct {
	clear bool
}

type Point struct {
	x int
	y int
}

type Box struct {
	goal bool
	Point
}

type Goal struct {
	Point
}

type Step struct {
	_type string
	idx   int
	Point
	goal bool
}

func (g *Game) Update() error {
	if g.clear {
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
		step = append(step, Step{_type: "player-move"})
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyU) && len(steps) > 0 {
		step := steps[len(steps)-1]

		for _, st := range step {
			switch st._type {
			case "player-move":
				playerPoint[0] = st.x
				playerPoint[1] = st.y
			case "box-move":
				box[st.idx].x = st.x
				box[st.idx].y = st.y
			case "box-goal":
				box[st.idx].setGoal(st.goal)
			}
		}

		steps = steps[:len(steps)-1]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		reset()
	}

	playerPoint[0] += directionX
	playerPoint[1] += directionY

	for idx := 0; idx < len(box); idx++ {
		newbox := BoxFilter(box, func(bx Box) bool {
			return bx.x == box[idx].x+directionX && bx.y == box[idx].y+directionY
		})

		if box[idx].x == playerPoint[0] && box[idx].y == playerPoint[1] && len(newbox) == 0 {
			step = append(step, Step{_type: "box-move", idx: idx, Point: Point{x: box[idx].x, y: box[idx].y}})
			box[idx].move(directionX, directionY)

			if box[idx].x < 0 {
				playerPoint[0] -= directionX
				box[idx].x = 0
			}

			if box[idx].x > width-1 {
				playerPoint[0] -= directionX
				box[idx].x = width - 1
			}

			if box[idx].y < 0 {
				playerPoint[1] -= directionY
				box[idx].y = 0
			}

			if box[idx].y > height-1 {
				playerPoint[1] -= directionY
				box[idx].y = height - 1
			}
		}
	}

	colidbox := BoxFilter(box, func(box Box) bool {
		return box.x == playerPoint[0] && box.y == playerPoint[1]
	})
	if len(colidbox) != 0 {
		playerPoint[0] -= directionX
		playerPoint[1] -= directionY
	}

	if playerPoint[0] < 0 {
		playerPoint[0] = 0
	}
	if playerPoint[0] > width-1 {
		playerPoint[0] = width - 1
	}
	if playerPoint[1] < 0 {
		playerPoint[1] = 0
	}
	if playerPoint[1] > height-1 {
		playerPoint[1] = height - 1
	}

	cancelGoal()
	g.clear = checkWin()

	if len(step) > 0 {
		steps = append(steps, step)
	}
	step = []Step{}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	board = [height][width]int{}
	screen.Fill(white)
	for _, ga := range goal {
		board[ga.y][ga.x] = 2
	}

	for _, bx := range box {
		if bx.goal {
			board[bx.y][bx.x] = 3
		} else {
			board[bx.y][bx.x] = 1
		}
	}

	for y := -1; y < height+1; y++ {
		for x := -1; x < width+1; x++ {
			_x := float64(50 * (x + 1))
			_y := float64(50 * (y + 1))

			if x == -1 || y == -1 || x == width || y == height {
				// 테두리
				ebitenutil.DrawRect(screen, _x, _y, 50, 50, colors[clr])
			} else if x == playerPoint[0] && y == playerPoint[1] {
				// 플레이어
				ebitenutil.DrawRect(screen, _x, _y, 50, 50, player)
			} else {
				ebitenutil.DrawRect(screen, _x, _y, 50, 50, getBlock(board[y][x], clr))
			}
		}
	}

	image, _, err := ebitenutil.NewImageFromFile("./assets/game-clear.png")
	if err != nil {
		log.Fatal(err)
	} else if g.clear {
		sizeX := float64(image.Bounds().Size().X / 2)
		sizeY := float64(image.Bounds().Size().Y / 2)

		screenCenterX := float64(screenSizeX / 2)
		screenCenterY := float64(screenSizeY / 2)

		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenCenterX-sizeX, screenCenterY-sizeY)

		screen.DrawImage(image, &op)
	}
}

func (g *Game) Layout(outsidewidth, outsideheight int) (screenwidth, screenheight int) {
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
	box = append(box, Box{goal: false, Point: Point{x: point[0], y: point[1]}})
}

func NewGoal(point [2]int) {
	goal = append(goal, Goal{Point{x: point[0], y: point[1]}})
}

func getBlock(block int, color int) color.RGBA {
	switch block {
	case 0:
		return tile
	case 1:
		return colors[color]
	case 2:
		return ngoal
	case 3:
		return ygoal
	default:
		return tile
	}
}

func checkWin() bool {
	var stack int
	for idx := 0; idx < len(box); idx++ {
		for _, goal := range goal {
			if goal.x == box[idx].x && goal.y == box[idx].y {
				stack++
				if len(step) > 0 {
					step = append(step, Step{_type: "box-goal", idx: idx, goal: box[idx].goal})
				}
				box[idx].setGoal(true)
			}
		}
	}

	return len(goal) == stack
}

func cancelGoal() {
	for idx := 0; idx < len(box); idx++ {
		for _, goal := range goal {
			if box[idx].goal && !(goal.x == box[idx].x && goal.y == box[idx].y) {
				box[idx].setGoal(false)
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

func reset() {
	board = [height][width]int{}
	box = []Box{}
	goal = []Goal{}

	playerPoint = [2]int{2, 3}

	points := [...][2]int{{1, 3}, {3, 2}, {7, 4}, {4, 3}, {3, 1}, {9, 5}, {12, 10}, {3, 11}, {14, 4}, {3, 3}, {1, 10}, {10, 2}, {10, 6}, {7, 10}}
	for _, point := range points {
		NewBox(point)
	}

	points = [...][2]int{{2, 4}, {1, 2}, {9, 4}, {7, 3}, {2, 7}, {1, 6}, {14, 2}, {4, 11}, {2, 9}, {5, 3}, {14, 1}, {10, 3}, {10, 11}, {14, 11}}
	for _, point := range points {
		NewGoal(point)
	}

	steps = [][]Step{}
	step = []Step{}
	clr = rand.Intn(7)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	clr = rand.Intn(7)

	reset()
}

func main() {
	ebiten.SetWindowSize(screenSizeX, screenSizeY)
	ebiten.SetWindowTitle("soGOban")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
