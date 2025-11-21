package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

// Point represents a single point on the grid
type Point struct {
	X, Y int
}

const (
	width     = 40
	height    = 20
	gameSpeed = 120 * time.Millisecond
)

var (
	snake     []Point
	food      Point
	direction Point
	score     int
	gameOver  bool
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Use a goroutine to listen for keyboard events
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	rand.Seed(time.Now().UnixNano())
	resetGame()

gameLoop:
	for {
		// Non-blocking input check
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowUp:
					if direction.Y == 0 {
						direction = Point{X: 0, Y: -1}
					}
				case termbox.KeyArrowDown:
					if direction.Y == 0 {
						direction = Point{X: 0, Y: 1}
					}
				case termbox.KeyArrowLeft:
					if direction.X == 0 {
						direction = Point{X: -1, Y: 0}
					}
				case termbox.KeyArrowRight:
					if direction.X == 0 {
						direction = Point{X: 1, Y: 0}
					}
				case termbox.KeyCtrlQ:
					break gameLoop
				case termbox.KeyEnter:
					if gameOver {
						resetGame()
					}
				}
				if ev.Ch == 'q' {
					break gameLoop
				}
			}
		default:
			// No event, continue game logic
			if !gameOver {
				update()
			}
			draw()
			time.Sleep(gameSpeed)
		}
	}
}

func resetGame() {
	gameOver = false
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	snake = []Point{
		{X: width / 2, Y: height / 2},
		{X: width/2 - 1, Y: height / 2},
		{X: width/2 - 2, Y: height / 2},
	}
	direction = Point{X: 1, Y: 0}
	score = 0
	placeFood()
}

func update() {
	head := snake[0]
	newHead := Point{X: head.X + direction.X, Y: head.Y + direction.Y}

	// Check for wall collision
	if newHead.X <= 0 || newHead.X >= width-1 || newHead.Y <= 0 || newHead.Y >= height-1 {
		gameOver = true
		return
	}

	// Check for self collision
	for _, p := range snake {
		if p == newHead {
			gameOver = true
			return
		}
	}

	snake = append([]Point{newHead}, snake...)

	// Check for food
	if newHead == food {
		score++
		placeFood() // Let the snake grow by not removing the tail
	} else {
		snake = snake[:len(snake)-1] // Remove tail
	}
}

func placeFood() {
	for {
		food = Point{
			X: rand.Intn(width-2) + 1,
			Y: rand.Intn(height-2) + 1,
		}
		// Ensure food doesn't spawn on the snake
		inSnake := false
		for _, p := range snake {
			if p == food {
				inSnake = true
				break
			}
		}
		if !inSnake {
			break
		}
	}
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawWalls()
	drawSnake()
	termbox.SetCell(food.X, food.Y, '●', termbox.ColorRed, termbox.ColorDefault)
	drawInfo()
	if gameOver {
		showGameOver()
	}
	termbox.Flush()
}

func drawWalls() {
	for i := 0; i < width; i++ {
		termbox.SetCell(i, 0, '═', termbox.ColorBlue, termbox.ColorDefault)
		termbox.SetCell(i, height-1, '═', termbox.ColorBlue, termbox.ColorDefault)
	}
	for i := 1; i < height-1; i++ {
		termbox.SetCell(0, i, '║', termbox.ColorBlue, termbox.ColorDefault)
		termbox.SetCell(width-1, i, '║', termbox.ColorBlue, termbox.ColorDefault)
	}
	termbox.SetCell(0, 0, '╔', termbox.ColorBlue, termbox.ColorDefault)
	termbox.SetCell(width-1, 0, '╗', termbox.ColorBlue, termbox.ColorDefault)
	termbox.SetCell(0, height-1, '╚', termbox.ColorBlue, termbox.ColorDefault)
	termbox.SetCell(width-1, height-1, '╝', termbox.ColorBlue, termbox.ColorDefault)
}

func drawSnake() {
	for i, p := range snake {
		char := '■'
		if i == 0 {
			char = '☻'
		}
		termbox.SetCell(p.X, p.Y, char, termbox.ColorGreen, termbox.ColorDefault)
	}
}

func drawInfo() {
	scoreText := fmt.Sprintf("Score: %d", score)
	drawText(width+2, 2, scoreText, termbox.ColorWhite, termbox.ColorDefault)
	drawText(width+2, 4, "Controls:", termbox.ColorWhite, termbox.ColorDefault)
	drawText(width+2, 5, "Arrow keys", termbox.ColorWhite, termbox.ColorDefault)
	drawText(width+2, 6, "'q' to quit", termbox.ColorWhite, termbox.ColorDefault)
}

func showGameOver() {
	gameOverText := "GAME OVER"
	finalScoreText := fmt.Sprintf("Final Score: %d", score)
	restartText := "Press ENTER to restart"
	drawText((width-len(gameOverText))/2, height/2-1, gameOverText, termbox.ColorRed, termbox.ColorDefault)
	drawText((width-len(finalScoreText))/2, height/2, finalScoreText, termbox.ColorWhite, termbox.ColorDefault)
	drawText((width-len(restartText))/2, height/2+2, restartText, termbox.ColorWhite, termbox.ColorDefault)
}

func drawText(x, y int, msg string, fg, bg termbox.Attribute) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
