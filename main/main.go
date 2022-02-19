package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

const (
	paddleHeight = 4
	paddleSymbol = 0x2588
)

var (
	player1 *Paddle
	player2 *Paddle
	screen  tcell.Screen
)

type Paddle struct {
	row, col, width, height int
}

func main() {
	initScreen()
	initGameState()

	for {
		drawState()

		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				screen.Fini()
				os.Exit(1)
			} else if ev.Rune() == 'w' {
				player1.row--
			} else if ev.Rune() == 's' {
				player1.row++
			} else if ev.Key() == tcell.KeyUp {
				player2.row--
			} else if ev.Key() == tcell.KeyDown {
				player2.row++
			}
		}
	}

}

func drawState() {
	screen.Clear()
	printPaddle(player1.row, player1.row, player1.width, player1.height, paddleSymbol)
	printPaddle(player2.row, player2.row, player2.width, player2.height, paddleSymbol)
	screen.Show()

}

func printPaddle(row, col, width, height int, ch rune) {
	// col = x as
	// row = y as
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}

}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func initGameState() {
	screenWidth, screenHeight := screen.Size()
	paddleStart := screenHeight/2 - paddleHeight/2

	// this uses xy quadrants
	// first int value is row = y as
	// second int value is col = x as
	// third int value is the width
	// fourth int value is the height

	player1 = &Paddle{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: paddleHeight,
	}

	player2 = &Paddle{
		row:    paddleStart,
		col:    screenWidth - 1,
		width:  1,
		height: paddleHeight,
	}

}
