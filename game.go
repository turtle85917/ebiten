package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	kind int
	idx  int
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
		step = append(step, Step{kind: playerMove, Point: Point{x: playerPoint[0], y: playerPoint[1]}})
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyU) && len(steps) > 0 {
		step := steps[len(steps)-1]

		for _, st := range step {
			switch st.kind {
			case playerMove:
				playerPoint[0] = st.x
				playerPoint[1] = st.y
			case boxMove:
				box[st.idx].x = st.x
				box[st.idx].y = st.y
			case boxGoal:
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
			step = append(step, Step{kind: boxMove, idx: idx, Point: Point{x: box[idx].x, y: box[idx].y}})
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
	screen.Fill(background)

	image, _, err := ebitenutil.NewImageFromFile("./assets/game-clear.png")
	if err != nil {
		log.Print(err)
	} else if g.clear {
		sizeX := float64(image.Bounds().Size().X / 2)
		sizeY := float64(image.Bounds().Size().Y / 2)

		screenCenterX := float64(screenSizeX / 2)
		screenCenterY := float64(screenSizeY / 2)

		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenCenterX-sizeX, screenCenterY-sizeY-150)

		screen.DrawImage(image, &op)

		return
	}

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
			texture := none
			fx := float64(50 * (x + 1))
			fy := float64(50 * (y + 1))

			if x == -1 || y == -1 || x == width || y == height || board[y][x] == 1 {
				// 테두리
				ebitenutil.DrawRect(screen, fx, fy, 50, 50, colors[clr])
			} else if x == playerPoint[0] && y == playerPoint[1] {
				// 플레이어
				texture = tplayer
			} else {
				texture = getTexture(board[y][x])
			}

			// Texture인 경우에만
			if texture != none {
				texture := textures[texture]

				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(fx, fy)

				screen.DrawImage(texture, &op)
			}
		}
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

func getTexture(block int) int {
	switch block {
	case 0:
		return ttile
	case 2:
		return tngoal
	case 3:
		return tygoal
	default:
		return ttile
	}
}

func checkWin() bool {
	var stack int
	for idx := 0; idx < len(box); idx++ {
		for _, goal := range goal {
			if goal.x == box[idx].x && goal.y == box[idx].y {
				stack++
				if len(step) > 0 {
					step = append(step, Step{kind: boxGoal, idx: idx, goal: box[idx].goal})
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

func loadTexture(path string) *ebiten.Image {
	texture, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Print(err)
		return nil
	}

	return texture
}

func init() {
	rand.Seed(time.Now().UnixNano())
	clr = rand.Intn(7)

	player = loadTexture("./assets/texture/smile.png")
	tile = loadTexture("./assets/texture/tile.png")
	ngoal = loadTexture("./assets/texture/goal-not.png")
	ygoal = loadTexture("./assets/texture/goal-success.png")

	textures = [...]*ebiten.Image{nil, player, tile, ngoal, ygoal}

	reset()
}
