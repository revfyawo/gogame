package systems

import (
	"fmt"
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const speed = 5

type Camera struct {
	ChunkPos *components.ChunkPosition
	scale    float64
}

func (c *Camera) New(world *ecs.World) {
	engine.Input.Register(sdl.SCANCODE_W)
	engine.Input.Register(sdl.SCANCODE_A)
	engine.Input.Register(sdl.SCANCODE_S)
	engine.Input.Register(sdl.SCANCODE_D)
	c.ChunkPos = new(components.ChunkPosition)
	c.scale = 1
}

func (c *Camera) Update(d time.Duration) {
	if engine.Input.Pressed(sdl.SCANCODE_W) {
		c.ChunkPos.MoveY(-speed)
	}
	if engine.Input.Pressed(sdl.SCANCODE_A) {
		c.ChunkPos.MoveX(-speed)
	}
	if engine.Input.Pressed(sdl.SCANCODE_S) {
		c.ChunkPos.MoveY(speed)
	}
	if engine.Input.Pressed(sdl.SCANCODE_D) {
		c.ChunkPos.MoveX(speed)
	}
}

func (*Camera) RemoveEntity(e *ecs.BasicEntity) {}

func (c *Camera) Scale() float64 {
	return c.scale
}

func (c *Camera) setScale() {
	scale := float32(c.scale)
	err := engine.Renderer.SetScale(scale, scale)
	if err != nil {
		panic(err)
	}
	fmt.Println("Setting scale to", scale)
}

func (c *Camera) IncScale(inc float64) {
	fmt.Println(c.scale)
	c.scale += inc
	c.setScale()
}

func (c *Camera) DecScale(dec float64) {
	fmt.Println(c.scale)
	c.scale -= dec
	c.setScale()
}
