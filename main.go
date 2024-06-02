package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

const UPDATE_TICK_RATE_MS float64 = 1000.0 / 60.0
const MIN_WIDTH = 80
const MIN_HEIGHT = 24

var (
	defStyle tcell.Style
)

type EngineState struct {
	lastRenderDuration float64
	lastUpdateDuration float64
}

func initEngineState() *EngineState {
	return &EngineState{
		lastUpdateDuration: UPDATE_TICK_RATE_MS,
	}
}

func update(es *EngineState) {
}

func draw(es *EngineState, lag float64) {
	screen.Clear()
	defer screen.Show()
	sw, sh := screen.Size()
	if sw < MIN_WIDTH || sh < MIN_HEIGHT {
		ShowResizeScreen(sw, sh, defStyle)
		return
	}
	renderFPS := 1000/es.lastRenderDuration
	updateFPS := 1000/es.lastUpdateDuration
	SetString(0, 0, "Hello World!", defStyle)
	SetString(0, 1, fmt.Sprintf("Render: %2f/s", renderFPS), defStyle)
	SetString(0, 2, fmt.Sprintf("Update: %2f/s", updateFPS), defStyle)
	SetString(0, 3, fmt.Sprintf("Lag: %2f ms", lag), defStyle)
	SetString(0, 4, fmt.Sprintf("Resolution: %d %d", sw, sh), defStyle)
}

func main() {
	defStyle :=
		tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	screen = s

	screen.SetStyle(defStyle)
	screen.EnableMouse()
	screen.EnablePaste()
	screen.Clear()

	quit := func() {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// main loop
	lag := 0.0
	prevTime := time.Now()
	lastRender := time.Now()

	es := initEngineState()

	for {
		currTime := time.Now()
		elapsed := float64(currTime.Sub(prevTime).Nanoseconds()) / (1000 * 1000)
		lag += elapsed
		prevTime = currTime

		// Event handling
		for screen.HasPendingEvent() {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventResize:
				screen.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					return
				} else if ev.Key() == tcell.KeyCtrlL {
					screen.Sync()
				}
			}
		}

		dirty := false
		for lag >= UPDATE_TICK_RATE_MS {
			dirty = true
			update(es)
			lag -= UPDATE_TICK_RATE_MS
		}

		if dirty {
			draw(es, lag)
			currRender := time.Now()
			es.lastRenderDuration = float64(currRender.Sub(lastRender).Nanoseconds()) / (1000 * 1000)
			lastRender = currRender
		}
	}
}
