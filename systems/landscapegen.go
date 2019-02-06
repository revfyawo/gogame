package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type landscapeWorkDone struct {
	*Landscapes
	sdl.Point
}

type LandscapeGen struct {
	chunks     map[sdl.Point]*entities.Chunk
	messages   chan ecs.Message
	landscapes *Landscapes
	toGenerate []*entities.Chunk
	generated  map[sdl.Point]bool
	workChan   chan *entities.Chunk
	doneChan   chan landscapeWorkDone
}

func (lg *LandscapeGen) New(*ecs.World) {
	lg.messages = make(chan ecs.Message, 10)
	lg.workChan = make(chan *entities.Chunk, parallelGen)
	lg.doneChan = make(chan landscapeWorkDone, parallelGen)
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
				lg.chunks[sdl.Point{chunk.X, chunk.Y}] = chunk
				lg.toGenerate = append(lg.toGenerate, chunk)
			case GenerateWorldMessage:
				lg.chunks = make(map[sdl.Point]*entities.Chunk)
				lg.generated = make(map[sdl.Point]bool)
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
			lg.landscapes.Merge(landscapes.Landscapes)
			lg.generated[landscapes.Point] = true
			lg.mergeNeighbours(landscapes.Point)
		}
	}
	engine.Message.Dispatch(NewLandscapesMessage{*lg.landscapes})
}

func (*LandscapeGen) RemoveEntity(*ecs.BasicEntity) {}

func (lg *LandscapeGen) generateLandscape() {
	chunk := <-lg.workChan

	chunkPoint := sdl.Point{chunk.X, chunk.Y}
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
				landscapeDown = landscapes.Find(chunkPoint, tileDown, biomeDown)
			}

			if i != 0 && biome == biomeLeft && landscape != landscapeLeft {
				landscape = landscapes.MergeLandscape(landscape, landscapeLeft)
			}
			if j != 0 && biome == biomeUp && landscape != landscapeUp {
				landscape = landscapes.MergeLandscape(landscape, landscapeUp)
			}
			if i != 0 && j != components.ChunkTile-1 && biome == biomeDown && landscape != landscapeDown {
				landscape = landscapes.MergeLandscape(landscape, landscapeDown)
			}

			if i != components.ChunkTile-1 && biome == biomeRight {
				landscape.AddTile(chunkPoint, tileRight)
			} else if i != components.ChunkTile-1 {
				landscapeRight = entities.NewLandscape(biomeRight)
				landscapeRight.AddTile(chunkPoint, tileRight)
				landscapes.Add(landscapeRight)
			}
			if i == 0 && j != components.ChunkTile-1 && biome == biomeDown {
				landscape.AddTile(chunkPoint, tileDown)
			} else if i == 0 && j != components.ChunkTile-1 {
				landscapeDown = entities.NewLandscape(biomeDown)
				landscapeDown.AddTile(chunkPoint, tileDown)
				landscapes.Add(landscapeDown)
			}

		}
	}
	lg.doneChan <- landscapeWorkDone{landscapes, chunkPoint}
}

func (lg *LandscapeGen) mergeNeighbours(chunk sdl.Point) {
	left := sdl.Point{chunk.X - 1, chunk.Y}
	up := sdl.Point{chunk.X, chunk.Y - 1}
	right := sdl.Point{chunk.X + 1, chunk.Y}
	down := sdl.Point{chunk.X, chunk.Y + 1}

	generated := lg.chunks[chunk]
	if lg.generated[left] {
		leftChunk := lg.chunks[left]
		for i := int32(0); i < components.ChunkTile; i++ {
			biome := generated.Biomes[0][i]
			landscape := lg.landscapes.Find(chunk, sdl.Point{0, i}, biome)
			biomeLeft := leftChunk.Biomes[components.ChunkTile-1][i]
			landscapeLeft := lg.landscapes.Find(left, sdl.Point{components.ChunkTile - 1, i}, biomeLeft)
			if biome == biomeLeft && landscape != landscapeLeft {
				lg.landscapes.MergeLandscape(landscape, landscapeLeft)
			}
		}
	}
	if lg.generated[up] {
		upChunk := lg.chunks[up]
		for i := int32(0); i < components.ChunkTile; i++ {
			biome := generated.Biomes[i][0]
			landscape := lg.landscapes.Find(chunk, sdl.Point{i, 0}, biome)
			biomeUp := upChunk.Biomes[i][components.ChunkTile-1]
			landscapeUp := lg.landscapes.Find(up, sdl.Point{i, components.ChunkTile - 1}, biomeUp)
			if biome == biomeUp && landscape != landscapeUp {
				lg.landscapes.MergeLandscape(landscape, landscapeUp)
			}
		}
	}
	if lg.generated[right] {
		rightChunk := lg.chunks[right]
		for i := int32(0); i < components.ChunkTile; i++ {
			biome := generated.Biomes[components.ChunkTile-1][i]
			landscape := lg.landscapes.Find(chunk, sdl.Point{components.ChunkTile - 1, i}, biome)
			biomeRight := rightChunk.Biomes[0][i]
			landscapeRight := lg.landscapes.Find(right, sdl.Point{0, i}, biomeRight)
			if biome == biomeRight && landscape != landscapeRight {
				lg.landscapes.MergeLandscape(landscape, landscapeRight)
			}
		}
	}
	if lg.generated[down] {
		downChunk := lg.chunks[down]
		for i := int32(0); i < components.ChunkTile; i++ {
			biome := generated.Biomes[i][components.ChunkTile-1]
			landscape := lg.landscapes.Find(chunk, sdl.Point{i, components.ChunkTile - 1}, biome)
			biomeDown := downChunk.Biomes[i][0]
			landscapeDown := lg.landscapes.Find(down, sdl.Point{i, 0}, biomeDown)
			if biome == biomeDown && landscape != landscapeDown {
				lg.landscapes.MergeLandscape(landscape, landscapeDown)
			}
		}
	}
}
