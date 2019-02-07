package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type MouseChunk struct {
	camera       *Camera
	lastCamPos   components.ChunkPosition
	lastCamScale float64
	chunkPos     components.ChunkPosition
	lastPos      sdl.Point
	position     sdl.Point
}

func (mc *MouseChunk) New(world *ecs.World) {
	camera := false
	for _, sys := range world.UpdateSystems() {
		switch s := sys.(type) {
		case *Camera:
			mc.camera = s
			camera = true
		}
	}
	if !camera {
		panic("need to add camera system before mouse chunk system")
	}
}

func (mc *MouseChunk) Update() {
	update := false
	mc.position = engine.Input.MousePosition()
	if mc.lastPos != mc.position {
		mc.lastPos = mc.position
		update = true
	}
	if mc.lastCamPos != mc.camera.Position() {
		mc.lastCamPos = mc.camera.Position()
		update = true
	}
	if mc.lastCamScale != mc.camera.Scale() {
		mc.lastCamScale = mc.camera.Scale()
		update = true
	}

	if update {
		camChunkPos := mc.camera.Position()
		scale := mc.camera.Scale()
		w, h, err := engine.Renderer.GetOutputSize()
		if err != nil {
			panic(err)
		}
		diff := sdl.Point{mc.position.X - w/2, mc.position.Y - h/2}
		camChunkPos.MoveX(int32(float64(diff.X) / scale))
		camChunkPos.MoveY(int32(float64(diff.Y) / scale))
		mc.chunkPos = camChunkPos
	}
}

func (*MouseChunk) RemoveEntity(*ecs.BasicEntity) {}
