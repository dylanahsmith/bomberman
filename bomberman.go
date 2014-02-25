package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

var (
	Wall = termbox.Cell{
		Ch: '▓',
		Fg: termbox.ColorGreen,
		Bg: termbox.ColorBlack,
	}

	Rock = termbox.Cell{
		Ch: '▓',
		Fg: termbox.ColorYellow,
		Bg: termbox.ColorBlack,
	}

	Ground = termbox.Cell{
		Ch: ' ',
		Fg: termbox.ColorDefault,
		Bg: termbox.ColorDefault,
	}

	Player = termbox.Cell{
		Ch: '♨',
		Fg: termbox.ColorWhite,
		Bg: termbox.ColorMagenta,
	}

	Bomb = termbox.Cell{
		Ch: 'ß',
		Fg: termbox.ColorRed,
		Bg: termbox.ColorDefault,
	}
)

var (
	turn  = time.Millisecond * 10
	board [81][25]termbox.Cell

	x, y         int = 1, 1
	lastX, lastY int
	h, w         int

	done = false
)

func setupBoard() {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			switch {
			case i == 0 || i == len(board)-1:
				board[i][j] = Wall
			case j == 0 || j == len(board[0])-1:
				board[i][j] = Wall
			case j%2 == 0 && i%2 == 0:
				board[i][j] = Wall
			default:
				board[i][j] = Rock
			}
		}
	}

	around(x, y, 3, func(in termbox.Cell) termbox.Cell {
		if in == Rock {
			return Ground
		}
		return in
	})
}

func main() {
	setupBoard()
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	evChan := make(chan termbox.Event)

	go func() {
		for {
			evChan <- termbox.PollEvent()
		}
	}()

	w, h = termbox.Size()
	draw()
	for !done {
		select {
		case ev := <-evChan:
			switch ev.Type {
			case termbox.EventResize:
				w, h = ev.Width, ev.Height
			case termbox.EventError:
				done = true
			case termbox.EventKey:
				doKey(ev.Key)
				draw()
			}
		}
	}
}

func draw() {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			cell := board[i][j]
			termbox.SetCell(i*2, j, cell.Ch, cell.Fg, cell.Bg)
			termbox.SetCell(i*2+1, j, cell.Ch, cell.Fg, cell.Bg)
		}
	}
	termbox.SetCell(x*2, y, Player.Ch, Player.Fg, Player.Bg)
	termbox.SetCell(x*2+1, y, Player.Ch, Player.Fg, Player.Bg)
	termbox.Flush()
}

func doKey(key termbox.Key) {
	switch key {
	case termbox.KeyCtrlC,
		termbox.KeyArrowUp,
		termbox.KeyArrowDown,
		termbox.KeyArrowLeft,
		termbox.KeyArrowRight:
		move(key)
	case termbox.KeyEnter:
		placeBomb()
	}
}

func move(key termbox.Key) {
	nextX, nextY := x, y
	switch key {
	case termbox.KeyCtrlC:
		done = true
	case termbox.KeyArrowUp:
		nextY--
	case termbox.KeyArrowDown:
		nextY++
	case termbox.KeyArrowLeft:
		nextX--
	case termbox.KeyArrowRight:
		nextX++
	}

	if canMove(nextX, nextY) {
		lastX, lastY = x, y
		x, y = nextX, nextY
	}
}

func canMove(x, y int) bool {
	switch board[x][y] {
	case Ground:
		return true
	}
	return false
}

func placeBomb() {
	board[x][y] = Bomb
	tmpX, tmpY := x, y

	time.AfterFunc(time.Second*3, func() {
		explode(tmpX, tmpY)
		draw()
	})
}

func explode(x int, y int) {
	board[x][y] = Ground
	cross(x, y, 3, func(cell termbox.Cell) termbox.Cell {
		if cell == Rock {
			return Ground
		}
		return cell
	})
}

func min(n, m int) int {
	if n < m {
		return n
	}
	return m
}

func max(n, m int) int {
	if n > m {
		return n
	}
	return m
}

func around(x, y, rad int, apply func(termbox.Cell) termbox.Cell) {
	for i := max(x-rad, 0); i < min(x+rad, len(board)); i++ {
		for j := max(y-rad, 0); j < min(y+rad, len(board[0])); j++ {
			board[i][j] = apply(board[i][j])
		}
	}
}

func cross(x, y, dist int, apply func(termbox.Cell) termbox.Cell) {
	for i := x; i < min(x+dist, len(board)); i++ {
		if board[i][y] == Wall {
			break
		}

		if i != x {
			board[i][y] = apply(board[i][y])
		}
	}

	for i := x; i > max(x-dist, 0); i-- {
		if board[i][y] == Wall {
			break
		}

		if i != x {
			board[i][y] = apply(board[i][y])
		}
	}

	for j := y; j < min(y+dist, len(board)); j++ {
		if board[x][j] == Wall {
			break
		}

		if j != y {
			board[x][j] = apply(board[x][j])
		}
	}

	for j := y; j > max(y+dist, 0); j-- {
		if board[x][j] == Wall {
			break
		}
		board[x][j] = apply(board[x][j])
	}
}
