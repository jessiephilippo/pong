package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const (
	paddleHeight = 4
	paddleSymbol = 0x2588
	ballSymbol   = 0x25CF
)

var (
	screen        tcell.Screen
	player1Paddle *GameObject
	player2Paddle *GameObject
	ball          *GameObject
	gameObject    []*GameObject
)

type GameObject struct {
	row, col, width, height int
	symbol                  rune
}

func main() {
	initScreen()
	initGameState()
	inputChan := initUserInput()

	for {
		drawState()
		time.Sleep(50 * time.Millisecond)

		key := readInput(inputChan)
		handleUserInput(key)
	}
}

func drawState() {
	screen.Clear()
	for _, obj := range gameObject {
		print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	screen.Show()
}

func print(row, col, width, height int, ch rune) {
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

	player1Paddle = &GameObject{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: paddleHeight,
		symbol: paddleSymbol,
	}

	player2Paddle = &GameObject{
		row:    paddleStart,
		col:    screenWidth - 1,
		width:  1,
		height: paddleHeight,
		symbol: paddleSymbol,
	}

	ball = &GameObject{
		row:    screenHeight / 2,
		col:    screenWidth / 2,
		width:  1,
		height: 1,
		symbol: ballSymbol,
	}

	gameObject = []*GameObject{
		player1Paddle, player2Paddle, ball,
	}

}

func initUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func readInput(input chan string) string {
	var key string
	select {
	case key = <-input:
	default:
		key = ""
	}

	return key
}

func handleUserInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(1)
	} else if key == "Rune[w]" && player1Paddle.row > 0 {
		player1Paddle.row--
	} else if key == "Rune[s]" && player1Paddle.row+player1Paddle.height < screenHeight {
		player1Paddle.row++
	} else if key == "Up" && player2Paddle.row > 0 {
		player2Paddle.row--
	} else if key == "Down" && player2Paddle.row+player2Paddle.height < screenHeight {
		player2Paddle.row++
	}
}
