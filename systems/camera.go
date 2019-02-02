package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"time"
)

const speed = 5

type Camera struct {
	ChunkPos  *components.ChunkPosition
	messages  chan ecs.Message
	scale     float64
	visible   sdl.Rect
	screenPos map[sdl.Point]sdl.Point
}

func (c *Camera) New(world *ecs.World) {
	engine.Input.Register(sdl.SCANCODE_W)
	engine.Input.Register(sdl.SCANCODE_A)
	engine.Input.Register(sdl.SCANCODE_S)
	engine.Input.Register(sdl.SCANCODE_D)
	c.ChunkPos = new(components.ChunkPosition)
	c.scale = 1
	c.messages = make(chan ecs.Message, 10)
	engine.Message.Listen(ChangeScaleMessageType, c.messages)
}

func (c *Camera) Update(d time.Duration) {
	pending := true
	for pending {
		select {
		case message := <-c.messages:
			switch m := message.(type) {
			case *ChangeScaleMessage:
				c.scale = m.Scale
			}
		default:
			pending = false
		}
	}

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

	c.getVisibleChunks()
}

func (*Camera) RemoveEntity(e *ecs.BasicEntity) {}

func (c *Camera) Scale() float64 {
	return c.scale
}

func (c *Camera) GetVisibleChunks() (sdl.Rect, map[sdl.Point]sdl.Point) {
	return c.visible, c.screenPos
}

func (c *Camera) getVisibleChunks() {
	w, h, err := engine.Renderer.GetOutputSize()
	if err != nil {
		panic(err)
	}

	camPos := c.ChunkPos
	scale := c.scale
	scaledCS := int32(components.ChunkSize * scale)
	// Screen position of the chunk the camera is in
	camChunkScreen := sdl.Point{w/2 - int32(float64(camPos.Position.X)*scale), h/2 - int32(float64(camPos.Position.Y)*scale)}

	// Compute how many chunks left, right, up and down the camera chunk
	var left, right, up, down int32
	if camChunkScreen.X >= 0 {
		left = int32(math.Ceil(float64(camChunkScreen.X) / float64(scaledCS)))
	}
	if camChunkScreen.X+int32(scaledCS) <= w {
		right = int32(math.Ceil(float64(w-camChunkScreen.X-scaledCS) / float64(scaledCS)))
	}
	if camChunkScreen.Y >= 0 {
		up = int32(math.Ceil(float64(camChunkScreen.Y) / float64(scaledCS)))
	}
	if camChunkScreen.Y+int32(scaledCS) <= h {
		down = int32(math.Ceil(float64(h-camChunkScreen.Y-scaledCS) / float64(scaledCS)))
	}
	c.visible = sdl.Rect{camPos.Chunk.X - left, camPos.Chunk.Y - up, left + right + 1, up + down + 1}

	// Fill visible chunk info
	c.screenPos = make(map[sdl.Point]sdl.Point)
	for x := camPos.Chunk.X - left; x <= camPos.Chunk.X+right; x++ {
		for y := camPos.Chunk.Y - up; y <= camPos.Chunk.Y+down; y++ {
			c.screenPos[sdl.Point{x, y}] = sdl.Point{
				camChunkScreen.X + scaledCS*(x-camPos.Chunk.X),
				camChunkScreen.Y + scaledCS*(y-camPos.Chunk.Y),
			}
		}
	}
}
