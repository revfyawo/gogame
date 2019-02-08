package systems

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const guiGap = 10

type guiPos struct {
	pos   entities.GUIPosition
	index int
}

type GUI struct {
	lock      sync.Mutex
	elements  map[entities.GUIPosition][]*entities.GUI
	positions map[string]guiPos
	messages  chan ecs.Message
}

func (gui *GUI) New(world *ecs.World) {
	gui.elements = make(map[entities.GUIPosition][]*entities.GUI)
	gui.positions = make(map[string]guiPos)
	gui.messages = make(chan ecs.Message, 10)
	engine.Message.Listen(GUIAddMessageType, gui.messages)
	engine.Message.Listen(GUIRemoveMessageType, gui.messages)
}

func (gui *GUI) UpdateFrame() {
	gui.lock.Lock()
	defer gui.lock.Unlock()

	pending := true
	for pending {
		select {
		case message := <-gui.messages:
			switch m := message.(type) {
			case GUIAddMessage:
				gui.positions[m.Element.Name] = guiPos{m.Element.Position, len(gui.elements[m.Element.Position])}
				gui.elements[m.Element.Position] = append(gui.elements[m.Element.Position], m.Element)
			case GUIRemoveMessage:
				pos, ok := gui.positions[m.Name]
				if ok {
					elem := gui.elements[pos.pos][pos.index]
					elem.DestroyTexture()
					gui.elements[pos.pos] = append(gui.elements[pos.pos][:pos.index], gui.elements[pos.pos][pos.index+1:]...)
					delete(gui.positions, m.Name)
				}
			}
		default:
			pending = false
		}
	}

	top := sdl.Rect{X: guiGap, Y: guiGap}
	for _, e := range gui.elements[entities.GUIPosTop] {
		if e.Texture == nil {
			e.GenerateTexture()
		}
		if e.H > top.H {
			top.H = e.H
		}
		err := engine.Renderer.Copy(e.Texture, nil, &sdl.Rect{top.X, top.Y, e.W, e.H})
		if err != nil {
			panic(err)
		}
		top.X += e.W + guiGap
	}

	w, h, err := engine.Renderer.GetOutputSize()
	if err != nil {
		panic(err)
	}
	center := sdl.Point{w / 2, h / 2}
	for _, e := range gui.elements[entities.GUIPosCenter] {
		if e.Texture == nil {
			e.GenerateTexture()
		}
		err = engine.Renderer.Copy(e.Texture, nil, &sdl.Rect{center.X - e.W/2, center.Y - e.H/2, e.W, e.H})
		if err != nil {
			panic(err)
		}
	}
}

func (gui *GUI) RemoveEntity(*ecs.BasicEntity) {}
