package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Landscapes struct {
	chunks     map[sdl.Point]*entities.Chunk
	messages   chan ecs.Message
	landscapes *components.Landscapes
	toGenerate []*entities.Chunk
	workChan   chan *entities.Chunk
	doneChan   chan *components.Landscapes
}

func (ls *Landscapes) New(*ecs.World) {
	ls.messages = make(chan ecs.Message, 10)
	ls.workChan = make(chan *entities.Chunk, parallelGen)
	ls.doneChan = make(chan *components.Landscapes, parallelGen)
	engine.Message.Listen(GenerateWorldMessageType, ls.messages)
	engine.Message.Listen(NewChunkMessageType, ls.messages)
}

func (ls *Landscapes) Update() {
	pending := true
	for pending {
		select {
		case message := <-ls.messages:
			switch m := message.(type) {
			case NewChunkMessage:
				chunk := m.Chunk
				ls.chunks[sdl.Point{chunk.Rect.X, chunk.Rect.Y}] = chunk
				ls.toGenerate = append(ls.toGenerate, chunk)
			case GenerateWorldMessage:
				ls.chunks = make(map[sdl.Point]*entities.Chunk)
				ls.landscapes = components.NewLandscapes()
				ls.toGenerate = []*entities.Chunk{}
			}
		default:
			pending = false
		}
	}

	if ls.toGenerate != nil && len(ls.toGenerate) > 0 {
		max := int(math.Min(float64(len(ls.toGenerate)), parallelGen))
		for i := 0; i < max; i++ {
			ls.workChan <- ls.toGenerate[i]
			go ls.generateLandscape()
		}
		ls.toGenerate = ls.toGenerate[max:]
		for i := 0; i < max; i++ {
			landscapes := <-ls.doneChan
			ls.landscapes.Merge(landscapes)
		}
	}
}

func (*Landscapes) RemoveEntity(*ecs.BasicEntity) {}

func (ls *Landscapes) generateLandscape() {
	chunk := <-ls.workChan

	chunkPoint := sdl.Point{chunk.Rect.X, chunk.Rect.Y}
	var tile, tileLeft, tileUp, tileRight, tileDown sdl.Point
	var biome, biomeLeft, biomeUp, biomeRight, biomeDown components.Biome
	var landscape, landscapeLeft, landscapeUp, landscapeRight, landscapeDown *components.Landscape
	landscape = components.NewLandscape(chunk.Biomes[0][0])
	landscape.AddTile(chunkPoint, sdl.Point{0, 0})
	landscapes := components.NewLandscapes()
	landscapes.Add(landscape)

	var i, j int32
	for i = 0; i < components.ChunkTile; i++ {
		for j = 0; j < components.ChunkTile; j++ {
			biome = chunk.Biomes[i][j]
			tile = sdl.Point{i, j}
			landscape = landscapes.Find(chunkPoint, tile, biome)
			if i != 0 {
				biomeLeft = chunk.Biomes[i-1][j]
				tileLeft = sdl.Point{i - 1, j}
				landscapeLeft = landscapes.Find(chunkPoint, tileLeft, biomeLeft)
			}
			if j != 0 {
				biomeUp = chunk.Biomes[i][j-1]
				tileUp = sdl.Point{i, j - 1}
				landscapeUp = landscapes.Find(chunkPoint, tileUp, biomeUp)
			}
			if i != components.ChunkTile-1 {
				biomeRight = chunk.Biomes[i+1][j]
				tileRight = sdl.Point{i + 1, j}
			}
			if j != components.ChunkTile-1 {
				biomeDown = chunk.Biomes[i][j+1]
				tileDown = sdl.Point{i, j + 1}
			}

			if i != 0 && biome == biomeLeft && landscape != landscapeLeft {
				if landscape.Size() > landscapeLeft.Size() {
					landscape.Merge(landscapeLeft)
				} else {
					landscapeLeft.Merge(landscape)
					landscape = landscapeLeft
				}
			}
			if j != 0 && biome == biomeUp && landscape != landscapeUp {
				if landscape.Size() > landscapeUp.Size() {
					landscape.Merge(landscapeUp)
				} else {
					landscapeUp.Merge(landscape)
					landscape = landscapeUp
				}
			}

			if i != components.ChunkTile-1 && biome == biomeRight {
				landscape.AddTile(chunkPoint, tileRight)
			} else if i != components.ChunkTile-1 {
				landscapeRight = components.NewLandscape(biomeRight)
				landscapeRight.AddTile(chunkPoint, tileRight)
				landscapes.Add(landscapeRight)
			}
			if j != components.ChunkTile-1 && biome == biomeDown {
				landscape.AddTile(chunkPoint, tileDown)
			} else if j != components.ChunkTile-1 {
				landscapeDown = components.NewLandscape(biomeDown)
				landscapeDown.AddTile(chunkPoint, tileDown)
				landscapes.Add(landscapeDown)
			}

		}
	}
	ls.doneChan <- landscapes
}
