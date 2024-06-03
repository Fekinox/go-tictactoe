package main

import (
	"github.com/gdamore/tcell/v2"
)

const MIN_WIDTH = 80
const MIN_HEIGHT = 24

var (
	defStyle tcell.Style
)

func main() {
	a := NewApp()
	defer a.Quit()
	a.Loop()
}
