package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
)

type ChunkRender struct {
	chunks      map[sdl.Point]*entities.Chunk
	visible     sdl.Rect
	screenPos   map[sdl.Point]sdl.Point
	lastVisible sdl.Rect
	messages    chan ecs.Message
	camera      *Camera
}

func (c *ChunkRender) New(world *ecs.World) {
	c.chunks = make(map[sdl.Point]*entities.Chunk)
	c.messages = make(chan ecs.Message, 10)
	engine.Message.Listen(GenerateWorldMessageType, c.messages)
	engine.Message.Listen(NewChunkMessageType, c.messages)

	camera := false
	for _, sys := range world.UpdateSystems() {
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

func (c *ChunkRender) UpdateFrame() {
	pending := true
	for pending {
		select {
		case message := <-c.messages:
			switch m := message.(type) {
			case NewChunkMessage:
				c.chunks[sdl.Point{m.Chunk.X, m.Chunk.Y}] = m.Chunk
				if m.Chunk.TilesTex != nil {
					err := m.Chunk.TilesTex.Destroy()
					if err != nil {
						panic(err)
					}
					m.Chunk.TilesTex = nil
				}
			case GenerateWorldMessage:
				for _, chunk := range c.chunks {
					if chunk.TilesTex != nil {
						err := chunk.TilesTex.Destroy()
						if err != nil {
							panic(err)
						}
						chunk.TilesTex = nil
					}
				}
				c.chunks = make(map[sdl.Point]*entities.Chunk)
			}
		default:
			pending = false
		}
	}

	c.camera.RLock()
	c.visible, c.screenPos = c.camera.GetVisibleChunks()
	if c.lastVisible != c.visible {
		c.freeHiddenChunks()
	}
	c.lastVisible = c.visible
	scaledCS := int32(components.ChunkSize * c.camera.Scale())
	c.camera.RUnlock()

	for point, pos := range c.screenPos {
		chunk := c.chunks[point]
		if chunk == nil {
			continue
		}
		dst := &sdl.Rect{pos.X, pos.Y, scaledCS, scaledCS}
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
