package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type Chunks struct {
	chunks map[sdl.Point]*entities.Chunk
}

func (c *Chunks) New(world *ecs.World) {
	c.chunks = make(map[sdl.Point]*entities.Chunk)
	for x := -5; x <= 5; x++ {
		for y := -5; y <= 5; y++ {
			chunk := entities.NewChunk(components.Space{Rect: sdl.Rect{X: int32(x), Y: int32(y)}})
			engine.Message.Dispatch(&NewChunkMessage{chunk})
			c.chunks[sdl.Point{chunk.Rect.X, chunk.Rect.Y}] = chunk
		}
	}
}

func (*Chunks) Update(d time.Duration) {}

func (*Chunks) RemoveEntity(e *ecs.BasicEntity) {}
