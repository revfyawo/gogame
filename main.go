package main

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/systems"
	"github.com/veandco/go-sdl2/sdl"
)

type defaultScene struct{}

func (s *defaultScene) Setup(w *ecs.World) {
	engine.Input.Register(sdl.SCANCODE_W)
	engine.Input.Register(sdl.SCANCODE_A)
	engine.Input.Register(sdl.SCANCODE_S)
	engine.Input.Register(sdl.SCANCODE_D)

	w.AddSystem(&systems.ChunkRender{})
	w.AddSystem(&systems.Chunks{})
}

func main() {
	engine.Run(&defaultScene{})
}
