package entities

import (
	"github.com/ojrac/opensimplex-go"
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	TileSize  = 32
	ChunkTile = 32
	ChunkSize = ChunkTile * ChunkTile
	noiseStep = 0.1

	waterLevel = 0
	sandDiff   = 0.2
	snowLevel  = 0.8
	rockDiff   = 0.2

	waterColor = 0xff0000ff
	sandColor  = 0xffffff00
	grassColor = 0xff00ff00
	rockColor  = 0xff808080
	snowColor  = 0xffffffff
)

var (
	noise                                         = opensimplex.New(1234567890)
	waterTex, sandTex, grassTex, rockTex, snowTex *sdl.Texture
)

type Chunk struct {
	ecs.BasicEntity
	components.Space
	components.ChunkRender
}

func NewChunk(space components.Space) *Chunk {
	chunk := Chunk{BasicEntity: ecs.NewBasic(), Space: space}
	chunk.Rect.W = ChunkSize
	chunk.Rect.H = ChunkSize
	chunk.generate()
	return &chunk
}

// Generates a chunk and his textures
func (c *Chunk) generate() {
	// Initialize textures if needed
	if waterTex == nil {
		initTextures()
	}

	// Initialize ChunkRender component
	if c.Textures == nil {
		c.Textures = make([][]*sdl.Texture, ChunkTile)
		for i := range c.Textures {
			c.Textures[i] = make([]*sdl.Texture, ChunkTile)
		}
	}

	// Initialize and compute heightmap
	heightMap := make([][]float64, ChunkTile)
	for i := range heightMap {
		heightMap[i] = make([]float64, ChunkTile)
	}
	for i := range heightMap {
		for j := range heightMap[i] {
			heightMap[i][j] = noise.Eval2(float64(i)*noiseStep, float64(j)*noiseStep)
		}
	}

	// Assign textures
	for i := range heightMap {
		for j := range heightMap[i] {
			height := heightMap[i][j]
			switch {
			case height <= waterLevel:
				c.Textures[i][j] = waterTex
			case height > waterLevel && height <= waterLevel+sandDiff:
				c.Textures[i][j] = sandTex
			case height >= snowLevel:
				c.Textures[i][j] = snowTex
			case height < snowLevel && height >= snowLevel-rockDiff:
				c.Textures[i][j] = rockTex
			default:
				c.Textures[i][j] = grassTex
			}
		}
	}
}

func initTextures() {
	var err error
	tileRect := &sdl.Rect{0, 0, TileSize, TileSize}
	surface, err := sdl.CreateRGBSurface(0, TileSize, TileSize, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
	if err != nil {
		panic(err)
	}

	err = surface.FillRect(tileRect, waterColor)
	if err != nil {
		panic(err)
	}
	waterTex, err = engine.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(tileRect, sandColor)
	if err != nil {
		panic(err)
	}
	sandTex, err = engine.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(tileRect, snowColor)
	if err != nil {
		panic(err)
	}
	snowTex, err = engine.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(tileRect, rockColor)
	if err != nil {
		panic(err)
	}
	rockTex, err = engine.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(tileRect, grassColor)
	if err != nil {
		panic(err)
	}
	grassTex, err = engine.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
}
