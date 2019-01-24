package engine

import (
	"github.com/veandco/go-sdl2/sdl"
	"testing"
)

var (
	qPress   = &sdl.KeyboardEvent{Keysym: sdl.Keysym{Scancode: sdl.SCANCODE_Q}, State: sdl.PRESSED}
	qRelease = &sdl.KeyboardEvent{Keysym: sdl.Keysym{Scancode: sdl.SCANCODE_Q}, State: sdl.RELEASED}
	wPress   = &sdl.KeyboardEvent{Keysym: sdl.Keysym{Scancode: sdl.SCANCODE_W}, State: sdl.PRESSED}
	wRelease = &sdl.KeyboardEvent{Keysym: sdl.Keysym{Scancode: sdl.SCANCODE_W}, State: sdl.RELEASED}
)

func TestInputSystem(t *testing.T) {
	is := NewInputManager()
	is.Register(sdl.SCANCODE_W)
	is.Register(sdl.SCANCODE_Q)

	// Press W
	is.PushEvent(wPress)
	is.Update()
	if !is.JustPressed(sdl.SCANCODE_W) || !is.Pressed(sdl.SCANCODE_W) {
		t.Fail()
	}

	// Press W again
	is.PushEvent(wPress)
	is.Update()
	if is.JustPressed(sdl.SCANCODE_W) || !is.Pressed(sdl.SCANCODE_W) {
		t.Fail()
	}

	// Press Q
	is.PushEvent(qPress)
	is.Update()
	if !is.JustPressed(sdl.SCANCODE_Q) || !is.Pressed(sdl.SCANCODE_Q) {
		t.Fail()
	}

	// Release W
	is.PushEvent(wRelease)
	is.Update()
	if !is.JustReleased(sdl.SCANCODE_W) || is.Pressed(sdl.SCANCODE_W) {
		t.Fail()
	}

	// Release Q
	is.PushEvent(qRelease)
	is.Update()
	if is.JustReleased(sdl.SCANCODE_W) || !is.JustReleased(sdl.SCANCODE_Q) || is.Pressed(sdl.SCANCODE_Q) {
		t.Fail()
	}
}
