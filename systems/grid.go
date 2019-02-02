package systems

import (
	"fmt"
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type Grid struct {
	camera *Camera
	show   bool
}

func (g *Grid) New(world *ecs.World) {
	camera := false
	for _, sys := range world.Systems() {
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
	if engine.Input.JustPressed(sdl.SCANCODE_F1) {
		fmt.Println("F1")
		g.show = !g.show
	}
	if g.show {
		w, h, err := engine.Renderer.GetOutputSize()
		if err != nil {
			panic(err)
		}

		camPos := g.camera.ChunkPos
		scale := g.camera.Scale
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
