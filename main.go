package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// 1 for x and 2 for o
type Game struct {
	board [3][3]int
	turn  int
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBoard(screen)
	g.drawGame(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 600
}

func main() {

	g := &Game{
		turn: 1,
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
