package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"time"
)

type ChunkRect struct {
	Chunk     *entities.Chunk
	ScreenPos sdl.Rect
}

type ChunkRender struct {
	chunks   map[sdl.Point]*entities.Chunk
	messages []*NewChunkMessage
	world    *ecs.World
	camera   *Camera
}

func (c *ChunkRender) PushMessage(m ecs.Message) {
	mess, ok := m.(*NewChunkMessage)
	if !ok {
		return
	}
	c.messages = append(c.messages, mess)
}

func (c *ChunkRender) New(world *ecs.World) {
	c.chunks = make(map[sdl.Point]*entities.Chunk)
	engine.Message.Listen(NewChunkMessageType, c)
	c.world = world
	c.addCameraOnce()
}

func (c *ChunkRender) Update(d time.Duration) {
	for _, m := range c.messages {
		c.chunks[sdl.Point{m.Chunk.Rect.X, m.Chunk.Rect.Y}] = m.Chunk
	}
	chunkRects := c.getVisibleChunks()
	for _, chunkRect := range chunkRects {
		chunk := chunkRect.Chunk
		if chunk == nil {
			continue
		}
		rect := chunkRect.ScreenPos
		err := engine.Renderer.Copy(chunk.TilesTex, nil, &sdl.Rect{X: rect.X, Y: rect.Y, W: components.ChunkSize, H: components.ChunkSize})
		if err != nil {
			panic(err)
		}
	}
}

func (*ChunkRender) RemoveEntity(e *ecs.BasicEntity) {}

func (c *ChunkRender) addCameraOnce() {
	for _, sys := range c.world.Systems() {
		switch s := sys.(type) {
		case *Camera:
			c.camera = s
			return
		}
	}
	c.camera = &Camera{}
	c.world.AddSystem(c.camera)
}

func (c *ChunkRender) getVisibleChunks() []ChunkRect {
	w, h, err := engine.Renderer.GetOutputSize()
	if err != nil {
		panic(err)
	}

	// xChunkMin: leftmost chunk X coordinate
	// xChunkMax: rightmost chunk X coordinate
	// yChunkMin: topmost chunk Y coordinate
	// yChunkMax: bottommost chunk Y coordinate
	// xMin, yMin: coordinates on screen of leftmost topmost chunk
	var xChunkMin, xChunkMax, yChunkMin, yChunkMax, xMin, yMin int32
	xChunkMin = int32(math.Floor(float64(c.camera.ChunkPos.Chunk.X*components.ChunkSize-w/2) / components.ChunkSize))
	xChunkMax = int32(math.Ceil(float64(c.camera.ChunkPos.Chunk.X*components.ChunkSize+w/2) / components.ChunkSize))
	yChunkMin = int32(math.Floor(float64(c.camera.ChunkPos.Chunk.Y*components.ChunkSize-h/2) / components.ChunkSize))
	yChunkMax = int32(math.Ceil(float64(c.camera.ChunkPos.Chunk.Y*components.ChunkSize+h/2) / components.ChunkSize))
	xMin = w/2 - c.camera.ChunkPos.Position.X - components.ChunkSize*(c.camera.ChunkPos.Chunk.X-xChunkMin)
	yMin = h/2 - c.camera.ChunkPos.Position.Y - components.ChunkSize*(c.camera.ChunkPos.Chunk.Y-yChunkMin)

	var visible []ChunkRect
	for x := xChunkMin; x <= xChunkMax; x++ {
		for y := yChunkMin; y <= yChunkMax; y++ {
			chunk := c.chunks[sdl.Point{x, y}]
			rect := sdl.Rect{xMin + (x-xChunkMin)*components.ChunkSize, yMin + (y-yChunkMin)*components.ChunkSize, components.ChunkSize, components.ChunkSize}
			visible = append(visible, ChunkRect{chunk, rect})
		}
	}
	return visible
}
