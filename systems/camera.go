package systems

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"time"
)

const speed = 5

type Camera struct {
	Chunk    sdl.Point
	Position sdl.Point
}

func (c *Camera) New(world *ecs.World) {
	engine.Input.Register(sdl.SCANCODE_W)
	engine.Input.Register(sdl.SCANCODE_A)
	engine.Input.Register(sdl.SCANCODE_S)
	engine.Input.Register(sdl.SCANCODE_D)
}

func (c *Camera) Update(d time.Duration) {
	if engine.Input.Pressed(sdl.SCANCODE_W) {
		c.moveY(-speed)
	}
	if engine.Input.Pressed(sdl.SCANCODE_A) {
		c.moveX(-speed)
	}
	if engine.Input.Pressed(sdl.SCANCODE_S) {
		c.moveY(speed)
	}
	if engine.Input.Pressed(sdl.SCANCODE_D) {
		c.moveX(speed)
	}
}

func (*Camera) RemoveEntity(e *ecs.BasicEntity) {}

func (c *Camera) moveX(speed int32) {
	newX := c.Position.X + speed
	if newX >= entities.ChunkSize {
		c.Position.X = newX % entities.ChunkSize
		c.Chunk.X += int32(math.Floor(float64(newX) / entities.ChunkSize))
	} else if newX < 0 {
		div := int32(math.Floor(float64(newX) / entities.ChunkSize))
		newX -= div * entities.ChunkSize
		c.Chunk.X += div
		c.Position.X = newX
	} else {
		c.Position.X = newX
	}
}

func (c *Camera) moveY(speed int32) {
	newY := c.Position.Y + speed
	if newY >= entities.ChunkSize {
		c.Position.Y = newY % entities.ChunkSize
		c.Chunk.Y += int32(math.Floor(float64(newY) / entities.ChunkSize))
	} else if newY < 0 {
		div := int32(math.Floor(float64(newY) / entities.ChunkSize))
		newY -= div * entities.ChunkSize
		c.Chunk.Y += div
		c.Position.Y = newY
	} else {
		c.Position.Y = newY
	}
}
