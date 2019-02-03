package systems

import (
	"fmt"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"sync"
	"time"
)

type UpdatesCounter struct {
	font    *ttf.Font
	texture *sdl.Texture
	lock    sync.Mutex
	update  time.Duration
	frame   time.Duration
}

func (uc *UpdatesCounter) New(*ecs.World) {
	font, err := ttf.OpenFont("assets/fonts/Go-Mono.ttf", 32)
	if err != nil {
		panic(err)
	}
	uc.font = font
}

func (uc *UpdatesCounter) UpdateFrame() {
	uc.lock.Lock()
	update := uc.update
	uc.frame = engine.FrameDelta
	uc.lock.Unlock()

	var updates, frames float64
	switch update {
	case 0:
		updates = 60
	default:
		updates = float64(time.Second) / float64(update)
	}
	switch engine.FrameDelta {
	case 0:
		frames = 60
	default:
		frames = float64(time.Second) / float64(engine.FrameDelta)
	}
	text := fmt.Sprintf("%.1f FPS / %.1f UPS", frames, updates)
	surface, err := uc.font.RenderUTF8Blended(text, sdl.Color{0xff, 0xff, 0xff, 0xff})
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	if uc.texture != nil {
		err = uc.texture.Destroy()
		if err != nil {
			panic(err)
		}
		uc.texture = nil
	}

	texture, err := engine.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	uc.texture = texture

	err = engine.Renderer.Copy(texture, nil, &sdl.Rect{10, 10, surface.W, surface.H})
	if err != nil {
		panic(err)
	}
}

func (uc *UpdatesCounter) Update() {
	uc.lock.Lock()
	uc.update = engine.UpdateDelta
	uc.lock.Unlock()
}

func (*UpdatesCounter) RemoveEntity(*ecs.BasicEntity) {}
