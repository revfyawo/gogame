package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"time"
)

type Chunks struct {
	chunks map[sdl.Point]*entities.Chunk
}

func (c *Chunks) New(world *ecs.World) {
	c.chunks = make(map[sdl.Point]*entities.Chunk)
	for x := -5; x <= 5; x++ {
		for y := -5; y <= 5; y++ {
			chunk := entities.NewChunk(components.Space{Rect: sdl.Rect{X: int32(x), Y: int32(y)}}, 1234567890)
			engine.Message.Dispatch(&NewChunkMessage{chunk})
			c.chunks[sdl.Point{chunk.Rect.X, chunk.Rect.Y}] = chunk
		}
	}
	engine.Input.Register(sdl.SCANCODE_F5)
}

func (c *Chunks) Update(d time.Duration) {
	if engine.Input.Pressed(sdl.SCANCODE_F5) {
		seed := rand.Int63()
		for _, chunk := range c.chunks {
			chunk.Generate(seed)
		}
	}
}

func (*Chunks) RemoveEntity(e *ecs.BasicEntity) {}
