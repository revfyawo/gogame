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

type ChunkInfo struct {
	Chunk     *entities.Chunk
	ScreenPos *sdl.Rect
}

type ChunkRender struct {
	chunks   map[sdl.Point]*entities.Chunk
	messages []*NewChunkMessage
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

	camera := false
	for _, sys := range world.Systems() {
		switch s := sys.(type) {
		case *Camera:
			camera = true
			c.camera = s
		}
	}
	if !camera {
		panic("need to add camera system before render system")
	}
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
		scaleCS := int32(components.ChunkSize * c.camera.Scale)
		dst := &sdl.Rect{rect.X, rect.Y, scaleCS, scaleCS}
		err := engine.Renderer.Copy(chunk.TilesTex, nil, dst)
		if err != nil {
			panic(err)
		}
	}
}

func (*ChunkRender) RemoveEntity(e *ecs.BasicEntity) {}

func (c *ChunkRender) getVisibleChunks() []ChunkInfo {
	w, h, err := engine.Renderer.GetOutputSize()
	if err != nil {
		panic(err)
	}

	camPos := c.camera.ChunkPos
	scale := c.camera.Scale
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

	// Fill and return visible chunk info
	var visible []ChunkInfo
	for x := camPos.Chunk.X - left; x <= camPos.Chunk.X+right; x++ {
		for y := camPos.Chunk.Y - up; y <= camPos.Chunk.Y+down; y++ {
			screenPos := &sdl.Rect{
				camChunkScreen.X + scaledCS*(x-camPos.Chunk.X),
				camChunkScreen.Y + scaledCS*(y-camPos.Chunk.Y),
				scaledCS,
				scaledCS,
			}
			chunkInfo := ChunkInfo{
				Chunk:     c.chunks[sdl.Point{x, y}],
				ScreenPos: screenPos,
			}
			visible = append(visible, chunkInfo)
		}
	}
	return visible
}
