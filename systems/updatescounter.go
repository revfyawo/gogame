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
	font        *ttf.Font
	texture     *sdl.Texture
	size        sdl.Rect
	lock        sync.Mutex
	updateRate  int
	frameRate   int
	rateChanged bool
	updateCount int
	frameCount  int
	lastSecond  time.Time
	enabled     bool
	disable     bool
}

func (uc *UpdatesCounter) New(*ecs.World) {
	font, err := ttf.OpenFont("assets/fonts/Go-Mono.ttf", 32)
	if err != nil {
		panic(err)
	}
	uc.font = font
	engine.Input.Register(sdl.SCANCODE_F2)
}

func (uc *UpdatesCounter) UpdateFrame() {
	uc.lock.Lock()
	if !uc.enabled {
		// Destroy texture if just disabled
		if uc.disable && uc.texture != nil {
			err := uc.texture.Destroy()
			if err != nil {
				panic(err)
			}
			uc.texture = nil
			uc.disable = false
		}
		uc.lock.Unlock()
		return
	}
	uc.frameCount++
	update := uc.updateRate
	frame := uc.frameRate
	changed := uc.rateChanged
	uc.lock.Unlock()

	// Regenerate texture if changed, or just enabled
	if changed || uc.texture == nil {
		if changed {
			uc.lock.Lock()
			uc.rateChanged = false
			uc.lock.Unlock()
		}
		text := fmt.Sprintf("%v FPS / %v UPS", frame, update)
		fontSurface, err := uc.font.RenderUTF8Blended(text, sdl.Color{0xff, 0xff, 0xff, 0xff})
		if err != nil {
			panic(err)
		}
		defer fontSurface.Free()

		surface, err := sdl.CreateRGBSurface(0, fontSurface.W, fontSurface.H, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
		if err != nil {
			panic(err)
		}
		defer surface.Free()
		err = surface.FillRect(nil, 0x80000000)
		if err != nil {
			panic(err)
		}
		err = fontSurface.Blit(nil, surface, &sdl.Rect{0, 0, surface.W, surface.H})
		if err != nil {
			panic(err)
		}

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
		uc.size.W, uc.size.H = surface.W, surface.H
	}

	err := engine.Renderer.Copy(uc.texture, nil, &sdl.Rect{10, 10, uc.size.W, uc.size.H})
	if err != nil {
		panic(err)
	}
}

func (uc *UpdatesCounter) Update() {
	uc.lock.Lock()
	defer uc.lock.Unlock()
	if engine.Input.JustPressed(sdl.SCANCODE_F2) {
		uc.enabled = !uc.enabled
		if !uc.enabled {
			uc.disable = true
		}
	}
	if !uc.enabled {
		return
	}
	now := time.Now()
	if now.Sub(uc.lastSecond) > time.Second {
		uc.lastSecond = now
		uc.updateRate = uc.updateCount
		uc.frameRate = uc.frameCount
		uc.updateCount = 0
		uc.frameCount = 0
		uc.rateChanged = true
	}
	uc.updateCount++
}

func (*UpdatesCounter) RemoveEntity(*ecs.BasicEntity) {}
