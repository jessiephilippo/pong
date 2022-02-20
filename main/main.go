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

	initialBallVelocityRow = 1
	initialBallVelocityCol = 2
)

var (
	screen tcell.Screen
	paused bool

	player1Paddle *GameObject
	player2Paddle *GameObject
	ball          *GameObject

	gameObject []*GameObject
)

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

func main() {
	initScreen()
	initGameobjects()

	for !isGameOver() {
		readUserInput()

		updateState()
		drawState()

		time.Sleep(75 * time.Millisecond)
	}

	screenWidth, screenHeight := screen.Size()
	winner := getWinner()
	printStringCentered(screenHeight/2-1, screenWidth/2, "Game Over!")
	printStringCentered(screenHeight/2, screenWidth/2, fmt.Sprintf("%s wins...", winner))

	screen.Show()

	time.Sleep(3 * time.Second)

	screen.Fini()

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

func initGameobjects() {
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
		velRow: initialBallVelocityRow,
		velCol: initialBallVelocityCol,
		symbol: ballSymbol,
	}

	gameObject = []*GameObject{
		player1Paddle, player2Paddle, ball,
	}
}

func drawState() {
	if paused {
		return
	}

	screen.Clear()
	for _, obj := range gameObject {
		printGameobjects(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	screen.Show()
}

func updateState() {
	if paused {
		return
	}

	for i := range gameObject {
		gameObject[i].row += gameObject[i].velRow
		gameObject[i].col += gameObject[i].velCol
	}

	if collideWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if collideWithPaddle(ball, player1Paddle) || collideWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
	}
}

func readUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				_, screenHeight := screen.Size()
				// inputChan <- ev.Name()
				if ev.Rune() == 'q' {
					screen.Fini()
					os.Exit(1)
				} else if ev.Rune() == 'w' && player1Paddle.row > 0 {
					player1Paddle.row--
				} else if ev.Rune() == 's' && player1Paddle.row+player1Paddle.height < screenHeight {
					player1Paddle.row++
				} else if ev.Rune() == 'p' {
					paused = !paused
				} else if ev.Key() == tcell.KeyUp && player2Paddle.row > 0 {
					player2Paddle.row--
				} else if ev.Key() == tcell.KeyDown && player2Paddle.row+player2Paddle.height < screenHeight {
					player2Paddle.row++
				}
			}
		}
	}()

	return inputChan
}

func collideWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}

func collideWithPaddle(ball *GameObject, paddle *GameObject) bool {
	var collidesCol bool
	if ball.col < paddle.col {
		collidesCol = ball.col+ball.velCol >= paddle.col
	} else {
		collidesCol = ball.col+ball.velCol <= paddle.col
	}

	return collidesCol &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height
}

func isGameOver() bool {
	return getWinner() != ""
}

func getWinner() string {
	screenWidth, _ := screen.Size()
	if ball.col < 0 {
		return "Player 1"
	} else if ball.col >= screenWidth {
		return "Player 2"
	} else {
		return ""
	}
}

func printGameobjects(row, col, width, height int, ch rune) {
	// col = x as
	// row = y as
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func printStringCentered(row, col int, str string) {
	col = col - len(str)/2
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}
