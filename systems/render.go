package systems

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type ChunkRender struct {
	chunks   map[engine.Point]*entities.Chunk
	messages []*NewChunkMessage
}

func (c *ChunkRender) PushMessage(m ecs.Message) {
	mess, ok := m.(*NewChunkMessage)
	if !ok {
		return
	}
	c.messages = append(c.messages, mess)
}

func (c *ChunkRender) New(world *ecs.World) {
	c.chunks = make(map[engine.Point]*entities.Chunk)
	engine.Message.Listen(NewChunkMessageType, c)
}

func (c *ChunkRender) Update(d time.Duration) {
	for _, m := range c.messages {
		c.chunks[engine.Point{m.Chunk.Rect.X, m.Chunk.Rect.Y}] = m.Chunk
	}
	for point, chunk := range c.chunks {
		for i := range chunk.Textures {
			for j := range chunk.Textures[i] {
				dst := sdl.Rect{
					X: entities.ChunkSize*int32(point.X) + int32(i)*entities.TileSize,
					Y: entities.ChunkSize*int32(point.Y) + int32(j)*entities.TileSize,
					W: entities.TileSize,
					H: entities.TileSize,
				}
				err := engine.Renderer.Copy(chunk.Textures[i][j], nil, &dst)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (*ChunkRender) RemoveEntity(e *ecs.BasicEntity) {}
