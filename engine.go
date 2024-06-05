package main

import "github.com/gdamore/tcell/v2"

const UPDATE_TICK_RATE_MS float64 = 1000.0 / 240.0

var gridChars = []rune{
	'_', 'X', 'O',
}

const (
	InProgress int = iota
	P1Win
	P2Win
	Tie
)

var XMark = GridFromStrings(
	"\\ /",
	" X ",
	"/ \\",
)

var OMark = GridFromStrings(
	".-.",
	"| |",
	"._.",
)

type WinningLine []Position

type EngineState struct {
	LastRenderDuration float64
	LastUpdateDuration float64

	grid Grid[int]

	focusX       int
	focusY       int
	player       int
	toWin        int
	outcome      int
	winningLines []WinningLine
}

func IsRune(ev *tcell.EventKey, r rune) bool {
	return (ev.Key() == tcell.KeyRune && ev.Rune() == r)
}

func InitEngineState() *EngineState {
	return &EngineState{
		LastUpdateDuration: UPDATE_TICK_RATE_MS,
		grid:               MakeGrid(3, 3, 0),
		player:             1,
		toWin:              3,
		outcome:            InProgress,
		winningLines:       make([]WinningLine, 0),
	}
}

func (es *EngineState) ResetGame() {
	es.grid = MakeGrid(3, 3, 0)
	es.winningLines = make([]WinningLine, 0)
	es.outcome = InProgress
	es.player = 1
	es.focusX = 0
	es.focusY = 0
}

func (es *EngineState) HandleInput(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyUp || IsRune(ev, 'w') || IsRune(ev, 'W') {
			es.HandleMove(0, -1)
		} else if ev.Key() == tcell.KeyDown || IsRune(ev, 's') || IsRune(ev, 'S') {
			es.HandleMove(0, 1)
		} else if ev.Key() == tcell.KeyLeft || IsRune(ev, 'a') || IsRune(ev, 'A') {
			es.HandleMove(-1, 0)
		} else if ev.Key() == tcell.KeyRight || IsRune(ev, 'd') || IsRune(ev, 'D') {
			es.HandleMove(1, 0)
		} else if IsRune(ev, ' ') {
			es.HandlePlace()
		} else if IsRune(ev, 'r') || IsRune(ev, 'R') {
			es.HandleReset()
		}
	}
}

func (es *EngineState) Update() {
	// Handle input
}

func (es *EngineState) Draw(lag float64) {
	Screen.Clear()
	defer Screen.Show()
	sw, sh := Screen.Size()
	if sw < MIN_WIDTH || sh < MIN_HEIGHT {
		ShowResizeScreen(sw, sh, defStyle)
		return
	}

	rr := Area{
		X:      (sw - MIN_WIDTH) / 2,
		Y:      (sh - MIN_HEIGHT) / 2,
		Width:  MIN_WIDTH,
		Height: MIN_HEIGHT,
	}

	BorderBox(Area{
		X:      rr.X - 1,
		Y:      rr.Y - 1,
		Width:  rr.Width + 2,
		Height: rr.Height + 2,
	}, defStyle)

	es.DrawGridLines(rr)
	es.DrawGrid(rr)
}

func (es *EngineState) DrawGridLines(rr Area) {
	// horizontal
	for y := 1; y < es.grid.Height; y++ {
		for x := 0; x < es.grid.Width*4-1; x++ {
			Screen.SetContent(
				rr.X+x,
				rr.Y+y*4-1,
				'#',
				nil, defStyle)
		}
	}
	// vertical
	for x := 1; x < es.grid.Width; x++ {
		for y := 0; y < es.grid.Height*4-1; y++ {
			Screen.SetContent(
				rr.X+x*4-1,
				rr.Y+y,
				'#',
				nil, defStyle)
		}
	}
}

func (es *EngineState) DrawGrid(rr Area) {
	for y := 0; y < es.grid.Height; y++ {
		for x := 0; x < es.grid.Width; x++ {
			style := defStyle
			if x == es.focusX && y == es.focusY {
				style = style.Reverse(true)
			}
			// Screen.SetContent(
			// 	rr.X+2*x+1,
			// 	rr.Y+2*y+1,
			// 	gridChars[es.grid.MustGet(x, y)],
			// 	nil, style)
			es.DrawCell(
				rr,
				x, y,
				es.grid.MustGet(x, y),
				style)
		}
	}

	// Draw game status
	statusY := rr.Y + es.grid.Height*4
	switch es.outcome {
	case InProgress:
		if es.player == 1 {
			SetString(rr.X, statusY, "Player X to move", defStyle)
		} else {
			SetString(rr.X, statusY, "Player O to move", defStyle)
		}
	case P1Win:
		SetString(rr.X, statusY, "Player X wins", defStyle)
	case P2Win:
		SetString(rr.X, statusY, "Player O wins", defStyle)
	case Tie:
		SetString(rr.X, statusY, "Tie", defStyle)
	}
}

func (es *EngineState) DrawCell(
	rr Area, x, y int, player int, style tcell.Style) {
	var grid Grid[rune]

	if player == 1 {
		grid = XMark
	} else if player == 2 {
		grid = OMark
	} else {
		FillRegion(
			rr.X + x*4,
			rr.Y + y*4,
			3,
			3,
			' ',
			style)
		return
	}

	SetGrid(
		rr.X+x*4,
		rr.Y+y*4,
		grid,
		style)

}

func (es *EngineState) HandleMove(dx, dy int) {
	es.focusX = max(0, min(es.grid.Width-1, es.focusX+dx))
	es.focusY = max(0, min(es.grid.Height-1, es.focusY+dy))
}

func (es *EngineState) HandleReset() {
	es.ResetGame()
}

func (es *EngineState) HandlePlace() {
	if es.outcome != InProgress {
		es.ResetGame()
		return
	}
	curTile := es.grid.MustGet(es.focusX, es.focusY)
	if curTile != 0 {
		return
	}

	es.grid.Set(es.focusX, es.focusY, es.player)
	oldPlayer := es.player
	es.player = 3 - es.player

	// Check for wins
	newLines := es.AllKsInARow(oldPlayer, es.toWin)
	if len(newLines) > 0 {
		es.outcome = oldPlayer
		es.winningLines = newLines
	} else if es.BoardFull() {
		es.outcome = Tie
	}
}

func (es *EngineState) AllKsInARow(val int, length int) []WinningLine {
	lines := make([]WinningLine, 0)
	for y := 0; y < es.grid.Height; y++ {
		for x := 0; x < es.grid.Width; x++ {
			if line, ok := es.FindKInARow(x, y, 1, 0, val, length); ok {
				lines = append(lines, line)
			}
			if line, ok := es.FindKInARow(x, y, 0, 1, val, length); ok {
				lines = append(lines, line)
			}
			if line, ok := es.FindKInARow(x, y, 1, -1, val, length); ok {
				lines = append(lines, line)
			}
			if line, ok := es.FindKInARow(x, y, 1, 1, val, length); ok {
				lines = append(lines, line)
			}
		}
	}
	return lines
}

func (es *EngineState) FindKInARow(x, y int, dx, dy int, val int, length int) (WinningLine, bool) {
	positions := make(WinningLine, 0)
	xx := x
	yy := y
	for {
		if !es.grid.InBounds(xx, yy) || es.grid.MustGet(xx, yy) != val {
			return nil, false
		}
		positions = append(positions, Position{X: xx, Y: yy})
		if len(positions) == length {
			return positions, true
		}

		xx += dx
		yy += dy
	}
}

func (es *EngineState) BoardFull() bool {
	for y := 0; y < es.grid.Height; y++ {
		for x := 0; x < es.grid.Width; x++ {
			if es.grid.MustGet(x, y) == 0 {
				return false
			}
		}
	}

	return true
}
