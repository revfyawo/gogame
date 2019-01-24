package engine

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type KeyMap map[sdl.Scancode]*KeyState

type KeyState struct {
	State     uint8
	LastState uint8
}

type InputSystem struct {
	queue []sdl.KeyboardEvent // queue of events to be processed next tick
	keys  KeyMap              // current state of registered keys
}

func (is *InputSystem) New() {
	is.keys = make(KeyMap)
}

func NewInputSystem() *InputSystem {
	is := InputSystem{}
	is.New()
	return &is
}

func (is *InputSystem) Update(d time.Duration) {
	// Clear just released
	for sc := range is.keys {
		if is.JustReleased(sc) {
			is.keys[sc].LastState = sdl.RELEASED
		}
	}

	// Process tick event queue
	for _, event := range is.queue {
		sc := event.Keysym.Scancode
		state, ok := is.keys[sc]
		if !ok {
			break
		}
		state.LastState = is.keys[sc].State
		state.State = event.State
	}
	is.queue = is.queue[:0]
}

func (*InputSystem) RemoveEntity(e *ecs.BasicEntity) {}

func (is *InputSystem) PushEvent(e *sdl.KeyboardEvent) {
	is.queue = append(is.queue, *e)
}

func (is *InputSystem) Pressed(sc sdl.Scancode) bool {
	state, ok := is.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.PRESSED
}

func (is *InputSystem) JustPressed(sc sdl.Scancode) bool {
	state, ok := is.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.PRESSED && state.LastState == sdl.RELEASED
}

func (is *InputSystem) JustReleased(sc sdl.Scancode) bool {
	state, ok := is.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.RELEASED && state.LastState == sdl.PRESSED
}

func (is *InputSystem) Register(sc sdl.Scancode) {
	if _, ok := is.keys[sc]; !ok {
		is.keys[sc] = &KeyState{}
	}
}
