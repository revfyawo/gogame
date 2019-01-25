package main

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/systems"
)

type defaultScene struct{}

func (s *defaultScene) Setup(w *ecs.World) {
	w.AddSystem(&systems.ChunkRender{})
	w.AddSystem(&systems.Chunks{})
}

func main() {
	engine.Run(&defaultScene{})
}
