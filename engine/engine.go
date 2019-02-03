package engine

import (
	"flag"
	"github.com/revfyawo/gogame/ecs"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	FPSLimit = 60
	UPSLimit = 60
)

var (
	Input    *InputManager
	Message  *MessageManager
	Renderer *sdl.Renderer

	currentScene ecs.Scene
	currentWorld *ecs.World
	quit         = make(chan bool, 1)
	updateDone   = make(chan bool, 2*UPSLimit)
)

func Run(scene ecs.Scene) {
	flag.Parse()
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("gogame", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		0, 0, sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = window.Destroy()
		if err != nil {
			panic(err)
		}
	}()

	Renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	// Initialize Input, Message, World and Scene
	Input = NewInputManager()
	Message = NewMessageManager()
	currentWorld = &ecs.World{}
	currentScene = scene
	currentScene.Setup(currentWorld)

	go runUpdateLoop()
	runFrameLoop()
}

func runUpdateLoop() {
	var counter int
	var now, start, lastUpdate, lastSecond time.Time
	var delta time.Duration
	start = time.Now()
	lastUpdate = start
	lastSecond = start
	var ticker = time.NewTicker(time.Second / UPSLimit)
	defer ticker.Stop()
	for {
		<-ticker.C
		counter++
		// Notify render loop that update succeeded
		updateDone <- true
		now = time.Now()
		if lastSecond.Add(time.Second).Before(time.Now()) {
			lastSecond = now
			counter = 0
		}
		delta = now.Sub(lastUpdate)
		lastUpdate = now

		// SDL uses same address for each event: need to copy value before passing it to input manager
		// can't group cases, because copy wouldn't work because Event is an interface
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				quit <- true
				return
			case *sdl.KeyboardEvent:
				newEvent := *e
				Input.PushEvent(&newEvent)
			case *sdl.MouseButtonEvent:
				newEvent := *e
				Input.PushEvent(&newEvent)
			case *sdl.MouseMotionEvent:
				newEvent := *e
				Input.PushEvent(&newEvent)
			case *sdl.MouseWheelEvent:
				newEvent := *e
				Input.PushEvent(&newEvent)
			}
		}

		Input.Update()
		currentWorld.Update(delta)
	}
}

func runFrameLoop() {
	var err error
	var counter int
	var now, start, lastFrame, lastSecond time.Time
	var delta time.Duration
	start = time.Now()
	lastFrame = start
	lastSecond = start
	var ticker = time.NewTicker(time.Second / FPSLimit)
	defer ticker.Stop()
	for {
		<-ticker.C
		counter++
		// Wait for update, so frame rate can't be higher than update rate
		<-updateDone

		// Empty channel, if frame rate lower than update rate
		pending := true
		for pending {
			select {
			case <-updateDone:
				continue
			default:
				pending = false
			}
		}

		now = time.Now()
		if lastSecond.Add(time.Second).Before(time.Now()) {
			lastSecond = now
			counter = 0
		}
		delta = now.Sub(lastFrame)
		lastFrame = now

		select {
		case <-quit:
			return
		default:
		}

		err = Renderer.Clear()
		if err != nil {
			panic(err)
		}

		currentWorld.UpdateRender(delta)
		Renderer.Present()
	}
}
