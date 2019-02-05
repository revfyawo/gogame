package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type LandscapeGen struct {
	chunks     map[sdl.Point]*entities.Chunk
	messages   chan ecs.Message
	landscapes *Landscapes
	toGenerate []*entities.Chunk
	workChan   chan *entities.Chunk
	doneChan   chan *Landscapes
}

func (lg *LandscapeGen) New(*ecs.World) {
	lg.messages = make(chan ecs.Message, 10)
	lg.workChan = make(chan *entities.Chunk, parallelGen)
	lg.doneChan = make(chan *Landscapes, parallelGen)
	engine.Message.Listen(GenerateWorldMessageType, lg.messages)
	engine.Message.Listen(NewChunkMessageType, lg.messages)
}

func (lg *LandscapeGen) Update() {
	pending := true
	for pending {
		select {
		case message := <-lg.messages:
			switch m := message.(type) {
			case NewChunkMessage:
				chunk := m.Chunk
				lg.chunks[sdl.Point{chunk.Rect.X, chunk.Rect.Y}] = chunk
				lg.toGenerate = append(lg.toGenerate, chunk)
			case GenerateWorldMessage:
				lg.chunks = make(map[sdl.Point]*entities.Chunk)
				lg.landscapes = NewLandscapes()
				lg.toGenerate = []*entities.Chunk{}
			}
		default:
			pending = false
		}
	}

	if lg.toGenerate != nil && len(lg.toGenerate) > 0 {
		max := int(math.Min(float64(len(lg.toGenerate)), parallelGen))
		for i := 0; i < max; i++ {
			lg.workChan <- lg.toGenerate[i]
			go lg.generateLandscape()
		}
		lg.toGenerate = lg.toGenerate[max:]
		for i := 0; i < max; i++ {
			landscapes := <-lg.doneChan
			lg.landscapes.Merge(landscapes)
		}
	}
	engine.Message.Dispatch(NewLandscapesMessage{*lg.landscapes})
}

func (*LandscapeGen) RemoveEntity(*ecs.BasicEntity) {}

func (lg *LandscapeGen) generateLandscape() {
	chunk := <-lg.workChan

	chunkPoint := sdl.Point{chunk.Rect.X, chunk.Rect.Y}
	var tile, tileLeft, tileUp, tileRight, tileDown sdl.Point
	var biome, biomeLeft, biomeUp, biomeRight, biomeDown components.Biome
	var landscape, landscapeLeft, landscapeUp, landscapeRight, landscapeDown *entities.Landscape
	landscape = entities.NewLandscape(chunk.Biomes[0][0])
	landscape.AddTile(chunkPoint, sdl.Point{0, 0})
	landscapes := NewLandscapes()
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
				landscapeRight = entities.NewLandscape(biomeRight)
				landscapeRight.AddTile(chunkPoint, tileRight)
				landscapes.Add(landscapeRight)
			}
			if j != components.ChunkTile-1 && biome == biomeDown {
				landscape.AddTile(chunkPoint, tileDown)
			} else if j != components.ChunkTile-1 {
				landscapeDown = entities.NewLandscape(biomeDown)
				landscapeDown.AddTile(chunkPoint, tileDown)
				landscapes.Add(landscapeDown)
			}

		}
	}
	lg.doneChan <- landscapes
}
