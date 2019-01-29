package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type KeyMap map[sdl.Scancode]*KeyState

type KeyState struct {
	State     uint8
	LastState uint8
}

type ButtonMap [5]uint8

type MouseState struct {
	pressed ButtonMap
	x, y    int32
	wheel   int32 // Only Y wheel coordinates
}

type InputManager struct {
	queue []sdl.Event // queue of events to be processed next tick
	keys  KeyMap      // current state of registered keys
	mouse MouseState  // current state of mouse
}

func NewInputManager() *InputManager {
	im := InputManager{}
	im.keys = make(KeyMap)
	return &im
}

func (im *InputManager) Update() {
	// Clear just released & mouse wheel
	for sc := range im.keys {
		if im.JustReleased(sc) {
			im.keys[sc].LastState = sdl.RELEASED
		}
	}
	im.mouse.wheel = 0

	// Process tick event queue
	for _, event := range im.queue {
		switch e := event.(type) {
		case *sdl.KeyboardEvent:
			sc := e.Keysym.Scancode
			state, ok := im.keys[sc]
			if !ok {
				break
			}
			state.LastState = im.keys[sc].State
			state.State = e.State
		case *sdl.MouseButtonEvent:
			im.mouse.pressed[e.Button-1] = e.State
		case *sdl.MouseMotionEvent:
			im.mouse.x = e.X
			im.mouse.y = e.Y
		case *sdl.MouseWheelEvent:
			im.mouse.wheel = e.Y
		}
	}
	im.queue = im.queue[:0]
}

func (im *InputManager) PushEvent(e sdl.Event) {
	im.queue = append(im.queue, e)
}

func (im *InputManager) Pressed(sc sdl.Scancode) bool {
	state, ok := im.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.PRESSED
}

func (im *InputManager) JustPressed(sc sdl.Scancode) bool {
	state, ok := im.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.PRESSED && state.LastState == sdl.RELEASED
}

func (im *InputManager) JustReleased(sc sdl.Scancode) bool {
	state, ok := im.keys[sc]
	if !ok {
		return false
	}
	return state.State == sdl.RELEASED && state.LastState == sdl.PRESSED
}

func (im *InputManager) Register(sc sdl.Scancode) {
	if _, ok := im.keys[sc]; !ok {
		im.keys[sc] = &KeyState{}
	}
}

func (im *InputManager) Clicked(button uint8) bool {
	return im.mouse.pressed[button-1] == sdl.PRESSED
}

func (im *InputManager) MousePosition(button uint8) sdl.Point {
	return sdl.Point{X: im.mouse.x, Y: im.mouse.y}
}

func (im *InputManager) Wheel() int32 {
	return im.mouse.wheel
}
