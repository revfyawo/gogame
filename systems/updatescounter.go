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
	lock        sync.Mutex
	updateRate  int
	frameRate   int
	updateCount int
	frameCount  int
	lastSecond  time.Time
	enabled     bool
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
	if !uc.enabled {
		return
	}
	now := time.Now()
	uc.lock.Lock()
	uc.frameCount++
	if now.Sub(uc.lastSecond) > time.Second {
		uc.lastSecond = now
		uc.updateRate = uc.updateCount
		uc.frameRate = uc.frameCount
		uc.updateCount = 0
		uc.frameCount = 0
	}
	update := uc.updateRate
	frame := uc.frameRate
	uc.lock.Unlock()

	text := fmt.Sprintf("%v FPS / %v UPS", frame, update)
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
	defer uc.lock.Unlock()
	if engine.Input.JustPressed(sdl.SCANCODE_F2) {
		uc.enabled = !uc.enabled
	}
	if !uc.enabled {
		return
	}
	uc.updateCount++
}

func (*UpdatesCounter) RemoveEntity(*ecs.BasicEntity) {}
