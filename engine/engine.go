package engine

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/systems"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	FPSLimit = 60
)

var (
	Input *systems.InputSystem

	gameSystems []ecs.System
)

func Run(argSystems []ecs.System) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("gogame", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_MAXIMIZED)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	// Initialize systems
	for _, sys := range argSystems {
		gameSystems = append(gameSystems, sys)
		switch s := sys.(type) {
		case ecs.Initializer:
			s.New()
		}
	}
	if Input == nil {
		Input = systems.NewInputSystem()
		Input.Register(sdl.SCANCODE_W)
	}

	var counter int
	var now, start, lastFrame, lastSecond time.Time
	var delta time.Duration
	start = time.Now()
	lastFrame = start
	lastSecond = start
	var ticker = time.NewTicker(time.Second / FPSLimit)
	var running = true
	for running {
		<-ticker.C
		counter++
		now = time.Now()
		if lastSecond.Add(time.Second).Before(time.Now()) {
			lastSecond = now
			counter = 0
		}
		delta = now.Sub(lastFrame)
		lastFrame = now

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				Input.PushEvent(e)
			}
		}

		Input.Update(delta)

		renderer.Clear()
		for _, sys := range gameSystems {
			sys.Update(delta)
		}
		renderer.Present()
	}
	ticker.Stop()
}
