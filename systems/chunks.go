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
	chunks     map[sdl.Point]*entities.Chunk
	seedHeight int64
	seedTemp   int64
	seedRain   int64
}

func (c *Chunks) New(world *ecs.World) {
	c.chunks = make(map[sdl.Point]*entities.Chunk)
	rand.Seed(time.Now().UnixNano())
	c.seedHeight = rand.Int63()
	c.seedTemp = rand.Int63()
	c.seedRain = rand.Int63()
	for x := -20; x <= 20; x++ {
		for y := -20; y <= 20; y++ {
			chunk := entities.NewChunk(components.Space{Rect: sdl.Rect{X: int32(x), Y: int32(y)}})
			chunk.Generate(c.seedHeight, c.seedRain, c.seedTemp)
			engine.Message.Dispatch(&NewChunkMessage{chunk})
			c.chunks[sdl.Point{chunk.Rect.X, chunk.Rect.Y}] = chunk
		}
	}
	engine.Input.Register(sdl.SCANCODE_F5)
}

func (c *Chunks) Update(d time.Duration) {
	if engine.Input.Pressed(sdl.SCANCODE_F5) {
		c.seedHeight = rand.Int63()
		c.seedTemp = rand.Int63()
		c.seedRain = rand.Int63()
		for _, chunk := range c.chunks {
			chunk.Generate(c.seedHeight, c.seedRain, c.seedTemp)
		}
	}
}

func (*Chunks) RemoveEntity(e *ecs.BasicEntity) {}

func (c *Chunks) generateInitialChunks() {

}
