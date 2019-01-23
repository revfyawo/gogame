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

func drawWorld(renderer *sdl.Renderer, world [][]uint32) {
	var tileSize int32 = 32
	for i := range world {
		for j := range world[i] {
			renderer.SetDrawColor(uint8(world[i][j]>>16&0xFF), uint8(world[i][j]>>8&0xFF), uint8(world[i][j]&0xFF), sdl.ALPHA_OPAQUE)
			renderer.FillRect(&sdl.Rect{X: int32(i) * tileSize, Y: int32(j) * tileSize, W: tileSize, H: tileSize})
		}
	}
	renderer.SetDrawColor(0, 0, 0, sdl.ALPHA_OPAQUE)
}

func mainLoop(window *sdl.Window, world [][]uint32) {
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second / FPSLimit)

	counter := 0
	now := time.Now()
	last := now

	running := true
	for running {
		renderer.Clear()
		drawWorld(renderer, world)
		for polled := sdl.PollEvent(); polled != nil; polled = sdl.PollEvent() {
			switch polled.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		renderer.Present()

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
