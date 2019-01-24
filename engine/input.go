package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type KeyMap map[sdl.Scancode]*KeyState

type KeyState struct {
	State     uint8
	LastState uint8
}

type InputManager struct {
	queue []sdl.KeyboardEvent // queue of events to be processed next tick
	keys  KeyMap              // current state of registered keys
}

func NewInputManager() *InputManager {
	is := InputManager{}
	is.keys = make(KeyMap)
	return &is
}

func (is *InputManager) Update() {
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

func (is *InputManager) PushEvent(e *sdl.KeyboardEvent) {
	is.queue = append(is.queue, *e)
}

func (is *InputManager) Pressed(sc sdl.Scancode) bool {
	state, ok := is.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.PRESSED
}

func (is *InputManager) JustPressed(sc sdl.Scancode) bool {
	state, ok := is.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.PRESSED && state.LastState == sdl.RELEASED
}

func (is *InputManager) JustReleased(sc sdl.Scancode) bool {
	state, ok := is.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.RELEASED && state.LastState == sdl.PRESSED
}

func (is *InputManager) Register(sc sdl.Scancode) {
	if _, ok := is.keys[sc]; !ok {
		is.keys[sc] = &KeyState{}
	}
}
