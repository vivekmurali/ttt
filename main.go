package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
)

// 1 for x and 2 for o
type Game struct {
	board [3][3]int
	turn  int
	// 0 - menu
	// 1 - playing
	// 2 - game over
	state int
	// 3 is draw
	winner   int
	timeOver time.Time
}

var (
	fontface font.Face
	rect     image.Rectangle
)

func init() {

	tt, err := opentype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 100

	fontface, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})

	// 100, 290
	rect = text.BoundString(fontface, "play")
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if g.state == 1 {

			i, j := getSquare(x, y)

			if g.board[i][j] < 1 {

				g.board[i][j] = g.turn

				if g.turn == 1 {
					g.turn = 2
				} else {
					g.turn = 1
				}

				g.checkWinCondition()
			}
		}

		if g.state == 0 {
			topleftX := rect.Min.X + 100
			topleftY := rect.Min.Y + 290
			bottomrightX := rect.Max.X + 100
			bottomrightY := rect.Max.Y + 290

			if topleftX < x && x < bottomrightX && topleftY < y && y < bottomrightY {
				g.state = 1
			}
		}

		if g.state == 2 {
			if time.Since(g.timeOver).Milliseconds() > 200 {
				g.state = 0
				g.reset()
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.state == 0 {
		text.Draw(screen, "play", fontface, 100, 290, color.RGBA{255, 255, 255, 1})
	}

	if g.state > 0 {
		drawBoard(screen)
		g.drawGame(screen)
	}

	if g.state == 2 {

		if g.winner < 3 {
			text.Draw(screen, fmt.Sprintf("Player %d wins!!!", g.winner), fontface, 100, 290, color.RGBA{255, 255, 255, 1})
		}

		if g.winner == 3 {
			text.Draw(screen, "Draw!!!!", fontface, 100, 290, color.RGBA{255, 255, 255, 1})
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 600
}

func main() {

	g := &Game{
		turn:  1,
		state: 0,
	}

	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func drawBoard(screen *ebiten.Image) {

	// line 1
	vector.StrokeLine(screen, 200, 0, 200, 600, 2, color.RGBA{255, 255, 255, 1}, false)
	// line 2
	vector.StrokeLine(screen, 400, 0, 400, 600, 2, color.RGBA{255, 255, 255, 1}, false)
	// line 3
	vector.StrokeLine(screen, 0, 200, 600, 200, 2, color.RGBA{255, 255, 255, 1}, false)
	// line 4
	vector.StrokeLine(screen, 0, 400, 600, 400, 2, color.RGBA{255, 255, 255, 1}, false)
}

func (g *Game) drawGame(screen *ebiten.Image) {

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			g.drawPlayer(i, j, screen)
		}
	}

}

func (g *Game) drawPlayer(i, j int, screen *ebiten.Image) {
	if g.board[i][j] == 0 {
		return
	}

	x0, x1, y0, y1 := coord(i, j)

	if g.board[i][j] == 1 {
		vector.StrokeLine(screen, x0, y0, x1, y1, 2, color.RGBA{255, 255, 255, 1}, false)
		vector.StrokeLine(screen, x1, y0, x0, y1, 2, color.RGBA{255, 255, 255, 1}, false)
	}

	if g.board[i][j] == 2 {
		vector.StrokeCircle(screen, (x0+x1)/2, (y0+y1)/2, 80, 2, color.RGBA{255, 255, 255, 1}, false)
	}

}

func coord(i, j int) (float32, float32, float32, float32) {

	x0 := i * 200
	x1 := 200 * (i + 1)
	y0 := j * 200
	y1 := 200 * (j + 1)

	return float32(x0), float32(x1), float32(y0), float32(y1)
}

func getSquare(x, y int) (int, int) {

	var i, j int

	if x < 200 {
		i = 0
	} else if x < 400 {
		i = 1
	} else if x < 600 {
		i = 2
	}

	if y < 200 {
		j = 0
	} else if y < 400 {
		j = 1
	} else if y < 600 {
		j = 2
	}

	return i, j
}

func (g *Game) checkWinCondition() {
	// check columns
	for i := 0; i < 3; i++ {
		if g.board[i][0] == g.board[i][1] && g.board[i][1] == g.board[i][2] && g.board[i][0] > 0 {
			g.state = 2
			g.winner = g.board[i][0]
			g.timeOver = time.Now()
			return
		}
	}

	// check rows
	for j := 0; j < 3; j++ {
		if g.board[0][j] == g.board[1][j] && g.board[1][j] == g.board[2][j] && g.board[0][j] > 0 {
			g.state = 2
			g.winner = g.board[0][j]
			g.timeOver = time.Now()
			return
		}
	}

	// check diagonals
	if g.board[0][0] == g.board[1][1] && g.board[0][0] == g.board[2][2] && g.board[0][0] > 0 {
		g.state = 2
		g.winner = g.board[0][0]
		g.timeOver = time.Now()
		return
	}

	if g.board[0][2] == g.board[1][1] && g.board[0][2] == g.board[2][0] && g.board[0][2] > 0 {
		g.state = 2
		g.winner = g.board[1][1]
		g.timeOver = time.Now()
		return
	}

	flag := false
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.board[i][j] == 0 {
				flag = true
			}
		}
	}

	if !flag {
		g.winner = 3
		g.timeOver = time.Now()
		g.state = 2
		return
	}
}

func (g *Game) reset() {

	a := &Game{
		turn:  1,
		state: 0,
	}
	*g = *a
}
