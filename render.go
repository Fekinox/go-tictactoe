package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

var (
	screen	tcell.Screen
)

func SetString(x int, y int, s string, style tcell.Style) {
	col := x
	for _, ch := range s {
		width := runewidth.RuneWidth(ch)
		screen.SetContent(col, y, ch, nil, style)
		col += width
	}
}

func SetCenteredString(x, y int, s string, style tcell.Style) {
	col := x - runewidth.StringWidth(s)/2
	for _, ch := range s {
		width := runewidth.RuneWidth(ch)
		screen.SetContent(col, y, ch, nil, style)
		col += width
	}
}

func ShowResizeScreen(w, h int, style tcell.Style) {
	SetCenteredString(w/2, h/2, "Screen too small!", style)
	SetCenteredString(w/2, h/2 + 1, fmt.Sprintf("Current: %d x %d", w, h), style)
}
