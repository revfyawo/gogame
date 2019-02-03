package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type MouseChunk struct {
	camera   *Camera
	chunkPos components.ChunkPosition
	lastPos  sdl.Point
	position sdl.Point
}

func (mc *MouseChunk) New(world *ecs.World) {
	camera := false
	for _, sys := range world.Systems() {
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

func (mc *MouseChunk) Update(time.Duration) {
	mc.position = engine.Input.MousePosition()
	if mc.position != mc.lastPos {
		mc.lastPos = mc.position
		camChunkPos := mc.camera.ChunkPos
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
