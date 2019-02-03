package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
	"time"
)

type Grid struct {
	camera *Camera
	show   bool
	// Need to lock for `show` field because it is written un Update & read in UpdateFrame
	lock sync.RWMutex
}

func (g *Grid) New(world *ecs.World) {
	camera := false
	for _, sys := range world.UpdateSystems() {
		switch s := sys.(type) {
		case *Camera:
			camera = true
			g.camera = s
		}
	}
	if !camera {
		panic("need to add camera system before render system")
	}
	engine.Input.Register(sdl.SCANCODE_F1)
}

func (g *Grid) Update(time.Duration) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if engine.Input.JustPressed(sdl.SCANCODE_F1) {
		g.show = !g.show
	}
}

func (g *Grid) UpdateFrame(time.Duration) {
	g.lock.RLock()
	show := g.show
	g.lock.RUnlock()
	if show {
		w, h, err := engine.Renderer.GetOutputSize()
		if err != nil {
			panic(err)
		}

		g.camera.RLock()
		camPos := g.camera.Position()
		scale := g.camera.Scale()
		g.camera.RUnlock()
		scaledCS := int32(components.ChunkSize * scale)
		var lineWidth int32
		if scaledCS <= 64 {
			lineWidth = 2
		} else {
			lineWidth = 4
		}
		chunkScreenPos := sdl.Point{w/2 - int32(float64(camPos.Position.X)*scale), h/2 - int32(float64(camPos.Position.Y)*scale)}
		initScreenPos := chunkScreenPos
		err = engine.Renderer.SetDrawColor(0, 0, 0, 0xff)
		if err != nil {
			panic(err)
		}

		chunkScreenPos.X -= lineWidth / 2
		chunkScreenPos.Y -= lineWidth / 2
		for chunkScreenPos.X >= 0 || chunkScreenPos.Y >= 0 {
			err = engine.Renderer.FillRect(&sdl.Rect{chunkScreenPos.X, 0, lineWidth, h})
			if err != nil {
				panic(err)
			}
			err = engine.Renderer.FillRect(&sdl.Rect{0, chunkScreenPos.Y, w, lineWidth})
			if err != nil {
				panic(err)
			}
			if chunkScreenPos.X >= 0 {
				chunkScreenPos.X -= scaledCS
			}
			if chunkScreenPos.Y >= 0 {
				chunkScreenPos.Y -= scaledCS
			}
		}

		chunkScreenPos = initScreenPos
		chunkScreenPos.X += scaledCS - lineWidth/2
		chunkScreenPos.Y += scaledCS - lineWidth/2
		for chunkScreenPos.X <= w || chunkScreenPos.Y <= h {
			err = engine.Renderer.FillRect(&sdl.Rect{chunkScreenPos.X, 0, lineWidth, h})
			if err != nil {
				panic(err)
			}
			err = engine.Renderer.FillRect(&sdl.Rect{0, chunkScreenPos.Y, w, lineWidth})
			if err != nil {
				panic(err)
			}
			if chunkScreenPos.X <= w {
				chunkScreenPos.X += scaledCS
			}
			if chunkScreenPos.Y <= h {
				chunkScreenPos.Y += scaledCS
			}
		}

	}
}

func (*Grid) RemoveEntity(*ecs.BasicEntity) {}
