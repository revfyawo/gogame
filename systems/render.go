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
	chunks      map[sdl.Point]*entities.Chunk
	visible     sdl.Rect
	visibleInfo []ChunkInfo
	lastVisible sdl.Rect
	messages    chan ecs.Message
	camera      *Camera
}

func (c *ChunkRender) New(world *ecs.World) {
	c.chunks = make(map[sdl.Point]*entities.Chunk)
	c.messages = make(chan ecs.Message, 10)
	engine.Message.Listen(NewChunkMessageType, c.messages)

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
	pending := true
	for pending {
		select {
		case message := <-c.messages:
			switch m := message.(type) {
			case *NewChunkMessage:
				c.chunks[sdl.Point{m.Chunk.Rect.X, m.Chunk.Rect.Y}] = m.Chunk
				if m.Chunk.TilesTex != nil {
					err := m.Chunk.TilesTex.Destroy()
					if err != nil {
						panic(err)
					}
					m.Chunk.TilesTex = nil
				}
			}
		default:
			pending = false
		}
	}

	c.getVisibleChunks()
	if c.lastVisible != c.visible {
		c.freeHiddenChunks()
	}
	c.lastVisible = c.visible

	for _, info := range c.visibleInfo {
		chunk := info.Chunk
		if chunk == nil {
			continue
		}
		rect := info.ScreenPos
		scaleCS := int32(components.ChunkSize * c.camera.Scale)
		dst := &sdl.Rect{rect.X, rect.Y, scaleCS, scaleCS}
		if chunk.TilesTex == nil {
			chunk.GenerateTexture()
		}
		err := engine.Renderer.Copy(chunk.TilesTex, nil, dst)
		if err != nil {
			panic(err)
		}
	}
}

func (*ChunkRender) RemoveEntity(e *ecs.BasicEntity) {}

func (c *ChunkRender) getVisibleChunks() {
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
	c.visible = sdl.Rect{camPos.Chunk.X - left, camPos.Chunk.Y - up, left + right + 1, up + down + 1}

	// Fill visible chunk info
	c.visibleInfo = c.visibleInfo[:0]
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
			c.visibleInfo = append(c.visibleInfo, chunkInfo)
		}
	}

}

func (c *ChunkRender) freeHiddenChunks() {
	diff := sdl.Rect{
		c.visible.X - c.lastVisible.X,
		c.visible.Y - c.lastVisible.Y,
		c.lastVisible.X + c.lastVisible.W - c.visible.X - c.visible.W,
		c.lastVisible.Y + c.lastVisible.H - c.visible.Y - c.visible.H,
	}
	if diff.X > 0 {
		for x := c.lastVisible.X; x != c.visible.X; x++ {
			for y := c.lastVisible.Y; y < c.lastVisible.Y+c.lastVisible.H; y++ {
				chunk := c.chunks[sdl.Point{x, y}]
				if chunk == nil {
					continue
				}
				if chunk.TilesTex != nil {
					err := chunk.TilesTex.Destroy()
					if err != nil {
						panic(err)
					}
					chunk.TilesTex = nil
				}
			}
		}
	}
	if diff.Y > 0 {
		for y := c.lastVisible.Y; y != c.visible.Y; y++ {
			for x := c.lastVisible.X; x < c.lastVisible.X+c.lastVisible.W; x++ {
				chunk := c.chunks[sdl.Point{x, y}]
				if chunk == nil {
					continue
				}
				if chunk.TilesTex != nil {
					err := chunk.TilesTex.Destroy()
					if err != nil {
						panic(err)
					}
					chunk.TilesTex = nil
				}
			}
		}
	}
	if diff.W > 0 {
		for x := c.lastVisible.X + c.lastVisible.W; x != c.visible.X+c.visible.W; x-- {
			for y := c.lastVisible.Y; y < c.lastVisible.Y+c.lastVisible.H; y++ {
				chunk := c.chunks[sdl.Point{x, y}]
				if chunk == nil {
					continue
				}
				if chunk.TilesTex != nil {
					err := chunk.TilesTex.Destroy()
					if err != nil {
						panic(err)
					}
					chunk.TilesTex = nil
				}
			}
		}
	}
	if diff.H > 0 {
		for y := c.lastVisible.Y + c.lastVisible.H; y != c.visible.Y+c.visible.H; y-- {
			for x := c.lastVisible.X; x < c.lastVisible.X+c.lastVisible.W; x++ {
				chunk := c.chunks[sdl.Point{x, y}]
				if chunk == nil {
					continue
				}
				if chunk.TilesTex != nil {
					err := chunk.TilesTex.Destroy()
					if err != nil {
						panic(err)
					}
					chunk.TilesTex = nil
				}
			}
		}
	}
}
