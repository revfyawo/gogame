package entities

import (
	"github.com/ojrac/opensimplex-go"
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
)

const (
	noiseStep = 0.5
	tileSize  = 32

	waterLevel = 0
	sandDiff   = 0.1
	snowLevel  = 0.8
	rockDiff   = 0.2

	waterColor = 0xff0000ff
	sandColor  = 0xffffff00
	grassColor = 0xff00ff00
	rockColor  = 0xff0f0f0f
	snowColor  = 0xffffffff
)

var (
	noise                                         = opensimplex.New(rand.Int63())
	waterTex, sandTex, grassTex, rockTex, snowTex *sdl.Texture
)

type Chunk struct {
	*ecs.BasicEntity
	*components.Space
	*components.ChunkRender
}

// Generates a chunk and his textures
func (c *Chunk) Generate() {
	tileX := c.Rect.W / tileSize
	tileY := c.Rect.H / tileSize

	// Initialize ChunkRender component
	if c.Textures == nil {
		c.Textures = make([][]*sdl.Texture, tileX)
		for i := range c.Textures {
			c.Textures[i] = make([]*sdl.Texture, tileY)
		}
	}

	// Initialize and compute heightmap
	heightMap := make([][]float64, tileX)
	for i := range heightMap {
		heightMap[i] = make([]float64, tileY)
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

func init() {
	var err error
	tileRect := &sdl.Rect{0, 0, tileSize, tileSize}
	surface, err := sdl.CreateRGBSurface(0, tileSize, tileSize, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
	if err != nil {
		panic(err)
	}
	renderer, err := sdl.CreateSoftwareRenderer(surface)
	if err != nil {
		panic(err)
	}

	err = surface.FillRect(tileRect, waterColor)
	if err != nil {
		panic(err)
	}
	waterTex, err = renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(tileRect, sandColor)
	if err != nil {
		panic(err)
	}
	sandTex, err = renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(tileRect, snowColor)
	if err != nil {
		panic(err)
	}
	snowTex, err = renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(tileRect, rockColor)
	if err != nil {
		panic(err)
	}
	rockTex, err = renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(tileRect, grassColor)
	if err != nil {
		panic(err)
	}
	grassTex, err = renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
}
