package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	FPSLimit = 60
)

func Run(window *sdl.Window) {
	world := NewWorld()
	mainLoop(window, world)
}

func drawWorld(surface *sdl.Surface, world [][]uint32) {
	var tileSize int32 = 32
	for i := range world {
		for j := range world[i] {
			surface.FillRect(&sdl.Rect{X: int32(i) * tileSize, Y: int32(j) * tileSize, W: tileSize, H: tileSize}, world[i][j])
		}
	}
}

func mainLoop(window *sdl.Window, world [][]uint32) {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second / FPSLimit)

	counter := 0
	now := time.Now()
	last := now

	running := true
	for running {
		drawWorld(surface, world)
		for polled := sdl.PollEvent(); polled != nil; polled = sdl.PollEvent() {
			switch event := polled.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.WindowEvent:
				if event.Event == sdl.WINDOWEVENT_RESIZED {
					surface, err = window.GetSurface()
					if err != nil {
						panic(err)
					}
				}
			}
		}
		window.UpdateSurface()

		<-ticker.C
		counter++
		now = time.Now()
		if last.Add(time.Second).Before(time.Now()) {
			fmt.Printf("%v: %v frames\n", now.Sub(last), counter)
			last = now
			counter = 0
		}
	}
	ticker.Stop()
}
