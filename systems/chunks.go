package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"time"
)

type Chunks struct {
	chunks map[engine.Point]*entities.Chunk
}

func (c *Chunks) New(world *ecs.World) {
	c.chunks = make(map[engine.Point]*entities.Chunk)
	chunk := entities.NewChunk(components.Space{})
	engine.Message.Dispatch(&NewChunkMessage{chunk})
	c.chunks[engine.Point{chunk.Rect.X, chunk.Rect.Y}] = chunk
}

func (*Chunks) Update(d time.Duration) {}

func (*Chunks) RemoveEntity(e *ecs.BasicEntity) {}
