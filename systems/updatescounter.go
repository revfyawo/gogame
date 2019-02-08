package systems

import (
	"fmt"
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
	"time"
)

const UpdateCounterName = "updates-counter"

type UpdatesCounter struct {
	lock        sync.Mutex
	gui         *entities.GUI
	updateRate  int
	frameRate   int
	rateChanged bool
	updateCount int
	frameCount  int
	lastSecond  time.Time
	enabled     bool
}

func (uc *UpdatesCounter) New(*ecs.World) {
	engine.Input.Register(sdl.SCANCODE_F2)

	gui := &entities.GUI{}

	// Component
	gui.Type = components.GUILabel
	gui.Text = uc.text()
	gui.Font = components.GetFont(components.MonoRegular, 32)
	gui.Background = true

	// Entity
	gui.Name = UpdateCounterName
	gui.Position = entities.GUIPosTop

	uc.gui = gui
}

func (uc *UpdatesCounter) UpdateFrame() {
	uc.lock.Lock()
	defer uc.lock.Unlock()
	if !uc.enabled {
		return
	}

	uc.frameCount++

	// Regenerate text if changed, and destroy texture
	if uc.rateChanged {
		uc.rateChanged = false
		uc.gui.Text = uc.text()
		uc.gui.DestroyTexture()
	}
}

func (uc *UpdatesCounter) Update() {
	uc.lock.Lock()
	defer uc.lock.Unlock()
	if engine.Input.JustPressed(sdl.SCANCODE_F2) {
		uc.enabled = !uc.enabled
		if uc.enabled {
			engine.Message.Dispatch(GUIAddMessage{uc.gui})
		} else {
			engine.Message.Dispatch(GUIRemoveMessage{UpdateCounterName})
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

func (uc *UpdatesCounter) text() string {
	return fmt.Sprintf("%v FPS / %v UPS", uc.frameRate, uc.updateRate)
}
