package entities

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type Chunk struct {
	ecs.BasicEntity
	components.Chunk
	components.Space
	components.ChunkRender
}

func NewChunk(space components.Space) *Chunk {
	chunk := Chunk{BasicEntity: ecs.NewBasic(), Space: space}
	chunk.Rect.W = components.ChunkSize
	chunk.Rect.H = components.ChunkSize
	return &chunk
}

// Generates a chunk heigh, rain and temp
func (c *Chunk) Generate(seedHeight, seedTemp, seedRain int64) {
	c.Chunk.Generate(seedHeight, seedTemp, seedRain, c.Rect.X, c.Rect.Y)
	if c.TilesTex != nil {
		err := c.TilesTex.Destroy()
		if err != nil {
			panic(err)
		}
		c.TilesTex = nil
	}
}

// Generates the chunk texture
func (c *Chunk) GenerateTexture() {
	// Initialize chunk tile surface
	chunkSurface, err := sdl.CreateRGBSurface(0, components.ChunkSize, components.ChunkSize, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
	if err != nil {
		panic(err)
	}
	defer chunkSurface.Free()

	// Assign textures and create chunk texture
	for i := 0; i < components.ChunkTile; i++ {
		for j := 0; j < components.ChunkTile; j++ {
			color := components.BiomeColors[c.Biomes[i][j]]
			rect := &sdl.Rect{X: components.TileSize * int32(i), Y: components.TileSize * int32(j), W: components.TileSize, H: components.TileSize}
			err = chunkSurface.FillRect(rect, color)
			if err != nil {
				panic(err)
			}
		}
	}
	if c.TilesTex != nil {
		err = c.TilesTex.Destroy()
		if err != nil {
			panic(err)
		}
	}
	c.TilesTex, err = engine.Renderer.CreateTextureFromSurface(chunkSurface)
	if err != nil {
		panic(err)
	}
}
