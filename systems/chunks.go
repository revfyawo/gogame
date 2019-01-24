package systems

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/entities"
	"time"
)

type Chunks struct {
	chunks [][]entities.Chunk
}

func (c *Chunks) New() {
	c.chunks = make([][]entities.Chunk, 1)
	c.chunks[0] = make([]entities.Chunk, 1)
	c.chunks[0][0].Generate()
}

func (*Chunks) Update(d time.Duration) {}

func (*Chunks) RemoveEntity(e *ecs.BasicEntity) {}
