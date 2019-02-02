package systems

import (
	"github.com/ojrac/opensimplex-go"
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
	"time"
)

const parallelGen = 8

type Chunks struct {
	mapSeed     int64
	chunks      map[sdl.Point]*entities.Chunk
	noiseHeight opensimplex.Noise
	noiseTemp   opensimplex.Noise
	noiseRain   opensimplex.Noise
	toGenerate  []sdl.Point
	workChan    chan *entities.Chunk
	chunkChan   chan *entities.Chunk
}

func (c *Chunks) New(world *ecs.World) {
	c.chunks = make(map[sdl.Point]*entities.Chunk)
	c.workChan = make(chan *entities.Chunk, parallelGen)
	c.chunkChan = make(chan *entities.Chunk, parallelGen)

	c.mapSeed = time.Now().UnixNano()
	rand.Seed(c.mapSeed)
	c.noiseHeight = opensimplex.New(rand.Int63())
	c.noiseTemp = opensimplex.New(rand.Int63())
	c.noiseRain = opensimplex.New(rand.Int63())

	for x := -20; x <= 20; x++ {
		for y := -20; y <= 20; y++ {
			chunk := entities.NewChunk(components.Space{Rect: sdl.Rect{X: int32(x), Y: int32(y)}})
			c.chunks[sdl.Point{chunk.Rect.X, chunk.Rect.Y}] = chunk
			c.toGenerate = append(c.toGenerate, sdl.Point{chunk.Rect.X, chunk.Rect.Y})
		}
	}
	engine.Input.Register(sdl.SCANCODE_F5)
}

func (c *Chunks) Update(d time.Duration) {
	if engine.Input.JustPressed(sdl.SCANCODE_F5) {
		c.mapSeed = time.Now().UnixNano()
		rand.Seed(c.mapSeed)
		c.noiseHeight = opensimplex.New(rand.Int63())
		c.noiseTemp = opensimplex.New(rand.Int63())
		c.noiseRain = opensimplex.New(rand.Int63())
		c.toGenerate = nil
		for point := range c.chunks {
			c.toGenerate = append(c.toGenerate, point)
		}
	}

	if len(c.toGenerate) > 0 {
		max := int(math.Min(float64(len(c.toGenerate)), parallelGen))
		for i := 0; i < max; i++ {
			c.workChan <- c.chunks[c.toGenerate[i]]
			go c.generateChunk()
		}
		c.toGenerate = c.toGenerate[max:]
		for i := 0; i < max; i++ {
			chunk := <-c.chunkChan
			engine.Message.Dispatch(&NewChunkMessage{chunk})
		}
	}
}

func (*Chunks) RemoveEntity(e *ecs.BasicEntity) {}

func (c *Chunks) generateChunk() {
	chunk := <-c.workChan
	chunk.Generate(c.noiseHeight, c.noiseRain, c.noiseTemp)
	c.chunkChan <- chunk
}
