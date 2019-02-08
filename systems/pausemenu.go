package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const PauseMenuName = "pause-menu"

type PauseMenu struct {
	lock  sync.Mutex
	pause bool
	gui   *entities.GUI
}

func (pm *PauseMenu) New(*ecs.World) {
	engine.Input.Register(sdl.SCANCODE_ESCAPE)

	font := components.GetFont(components.MonoRegular, 30)
	pause := &components.GUI{Type: components.GUILabel, Font: font, Text: "Pause menu"}
	resume := &components.GUI{Type: components.GUIButton, Font: font, Text: "Resume"}
	exit := &components.GUI{Type: components.GUIButton, Font: font, Text: "Exit"}

	menu := components.GUI{Type: components.GUIVLayout, Background: true, Children: []*components.GUI{pause, resume, exit}}
	pm.gui = &entities.GUI{Name: PauseMenuName, Position: entities.GUIPosCenter, GUI: menu}
}

func (pm *PauseMenu) Update() {
	if engine.Input.JustPressed(sdl.SCANCODE_ESCAPE) {
		pm.pause = !pm.pause
		if pm.pause {
			engine.Message.Dispatch(GUIAddMessage{pm.gui})
		} else {
			engine.Message.Dispatch(GUIRemoveMessage{PauseMenuName})
		}
	}
}

func (pm *PauseMenu) RemoveEntity(*ecs.BasicEntity) {}
