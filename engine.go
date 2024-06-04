package main

import "github.com/gdamore/tcell/v2"

const UPDATE_TICK_RATE_MS float64 = 1000.0 / 240.0

var gridChars = []rune {
	'_', 'X', 'O',
}

type EngineState struct {
	LastRenderDuration float64
	LastUpdateDuration float64

	grid Grid[int]

	focusX int
	focusY int
}

func IsRune(ev *tcell.EventKey, r rune) bool {
	return (ev.Key() == tcell.KeyRune && ev.Rune() == r)
}

func InitEngineState() *EngineState {
	return &EngineState{
		LastUpdateDuration: UPDATE_TICK_RATE_MS,
		grid: MakeGrid(3, 3, 0),
	}
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

	BorderBox(Area {
		X: rr.X - 1,
		Y: rr.Y - 1,
		Width: rr.Width + 2,
		Height: rr.Height + 2,
	}, defStyle)

	es.DrawGrid(rr)
}

func (es *EngineState) DrawGrid(rr Area) {
	for y := 0; y < es.grid.Height; y++ {
		for x := 0; x < es.grid.Width; x++ {
			style := defStyle
			if x == es.focusX && y == es.focusY {
				style = style.Reverse(true)
			}
			Screen.SetContent(
				rr.X + 2 * x + 1,
				rr.Y + 2 * y + 1,
				gridChars[es.grid.MustGet(x, y)],
				nil, style)	
		}
	}
}

func (es *EngineState) HandleMove(dx, dy int) {
	es.focusX = max(0, min(es.grid.Width-1, es.focusX + dx))
	es.focusY = max(0, min(es.grid.Height-1, es.focusY + dy))
}

func (es *EngineState) HandleReset() {
}

func (es *EngineState) HandlePlace() {
}
