package entities

import (
	"github.com/ojrac/opensimplex-go"
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	noiseStep = 0.05

	waterLevel = -0.5
	sandDiff   = 0.2
	snowLevel  = 0.5
	rockDiff   = 0.2

	waterColor = 0xff0000ff
	sandColor  = 0xffffff00
	grassColor = 0xff00ff00
	rockColor  = 0xff808080
	snowColor  = 0xffffffff
)

var (
	noise = opensimplex.New(1234567890)
)

type Chunk struct {
	ecs.BasicEntity
	components.Space
	components.ChunkRender
}

func NewChunk(space components.Space) *Chunk {
	chunk := Chunk{BasicEntity: ecs.NewBasic(), Space: space}
	chunk.Rect.W = components.ChunkSize
	chunk.Rect.H = components.ChunkSize
	chunk.generate()
	return &chunk
}

// Generates a chunk and his textures
func (c *Chunk) generate() {
	// Initialize and compute heightmap
	heightMap := make([][]float64, components.ChunkTile)
	for i := range heightMap {
		heightMap[i] = make([]float64, components.ChunkTile)
	}
	for i := range heightMap {
		for j := range heightMap[i] {
			heightMap[i][j] = noise.Eval2(float64(c.Rect.X*components.ChunkTile+int32(i))*noiseStep, float64(c.Rect.Y*components.ChunkTile+int32(j))*noiseStep)
		}
	}

	// Initialize chunk tile surface
	chunkSurface, err := sdl.CreateRGBSurface(0, components.ChunkSize, components.ChunkSize, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
	if err != nil {
		panic(err)
	}

	// Assign textures and create chunk texture
	for i := range heightMap {
		for j := range heightMap[i] {
			height := heightMap[i][j]
			var color uint32
			switch {
			case height <= waterLevel:
				color = waterColor
			case height > waterLevel && height <= waterLevel+sandDiff:
				color = sandColor
			case height >= snowLevel:
				color = snowColor
			case height < snowLevel && height >= snowLevel-rockDiff:
				color = rockColor
			default:
				color = grassColor
			}
			err = chunkSurface.FillRect(&sdl.Rect{X: components.TileSize * int32(i), Y: components.TileSize * int32(j), W: components.TileSize, H: components.TileSize}, color)
			if err != nil {
				panic(err)
			}
		}
	}
	c.TilesTex, err = engine.Renderer.CreateTextureFromSurface(chunkSurface)
	if err != nil {
		panic(err)
	}
}
