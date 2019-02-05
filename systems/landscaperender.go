package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type LandscapeRender struct {
	lock          sync.Mutex
	camera        *Camera
	mouseChunk    *MouseChunk
	mousePos      components.ChunkPosition
	chunks        map[sdl.Point]*entities.Chunk
	landscapes    *Landscapes
	messages      chan ecs.Message
	lastLandscape *entities.Landscape
}

func (lr *LandscapeRender) New(world *ecs.World) {
	camera := false
	mouseChunk := false
	for _, sys := range world.UpdateSystems() {
		switch s := sys.(type) {
		case *Camera:
			camera = true
			lr.camera = s
		case *MouseChunk:
			mouseChunk = true
			lr.mouseChunk = s
		}
	}
	if !camera {
		panic("need to add camera system before landscape render system")
	}
	if !mouseChunk {
		panic("need to add mouse chunk system before landscape render system")
	}

	lr.messages = make(chan ecs.Message, 10)
	engine.Message.Listen(GenerateWorldMessageType, lr.messages)
	engine.Message.Listen(NewChunkMessageType, lr.messages)
	engine.Message.Listen(NewLandscapesMessageType, lr.messages)
}

func (lr *LandscapeRender) UpdateFrame() {
	pending := true
	for pending {
		select {
		case mess := <-lr.messages:
			switch m := mess.(type) {
			case GenerateWorldMessage:
				lr.chunks = make(map[sdl.Point]*entities.Chunk)
			case NewChunkMessage:
				chunk := m.Chunk
				lr.chunks[sdl.Point{chunk.Rect.X, chunk.Rect.Y}] = chunk
			case NewLandscapesMessage:
				lr.landscapes = &m.Landscapes
			}
		default:
			pending = false
		}
	}

	lr.lock.Lock()
	mousePos := lr.mousePos
	lr.lock.Unlock()

	chunk := lr.chunks[mousePos.Chunk]
	if chunk == nil {
		return
	}
	chunkPoint := sdl.Point{chunk.Rect.X, chunk.Rect.Y}

	tile := sdl.Point{mousePos.Position.X / components.TileSize, mousePos.Position.Y / components.TileSize}
	landscape := lr.landscapes.Find(chunkPoint, tile, lr.chunks[chunkPoint].Biomes[tile.X][tile.Y])
	if landscape == nil {
		return
	} else if lr.lastLandscape == nil {
		lr.lastLandscape = landscape
	}

	border, changed := landscape.Border()
	if border == nil {
		return
	}

	lr.camera.RLock()
	_, screenPos := lr.camera.GetVisibleChunks()
	scaledCS := int32(components.ChunkSize * lr.camera.Scale())
	lr.camera.RUnlock()

	chunkPos := screenPos[mousePos.Chunk]
	if landscape != lr.lastLandscape || changed || landscape.BorderTex == nil {
		// Generating border texture
		surface, err := sdl.CreateRGBSurface(0, components.ChunkSize, components.ChunkSize, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
		if err != nil {
			panic(err)
		}
		defer surface.Free()

		// Make surface transparent
		err = surface.FillRect(nil, 0)
		if err != nil {
			panic(err)
		}

		// Color border tiles white
		for _, tile := range border {
			if tile.Chunk != chunkPoint {
				continue
			}
			err = surface.FillRect(&sdl.Rect{tile.Tile.X * components.TileSize, tile.Tile.Y * components.TileSize, components.TileSize, components.TileSize}, 0xffffffff)
			if err != nil {
				panic(err)
			}
		}

		texture, err := engine.Renderer.CreateTextureFromSurface(surface)
		if err != nil {
			panic(err)
		}

		// Free last used texture, and set current one
		if landscape.BorderTex != nil {
			err := landscape.BorderTex.Destroy()
			if err != nil {
				panic(err)
			}
		}
		landscape.BorderTex = texture

		// Free last selected landscape texture
		if landscape != lr.lastLandscape && lr.lastLandscape != nil {
			err := lr.lastLandscape.BorderTex.Destroy()
			if err != nil {
				panic(err)
			}
			lr.lastLandscape.BorderTex = nil
		}
		lr.lastLandscape = landscape
	}

	// Rendering landscape border
	err := engine.Renderer.Copy(landscape.BorderTex, nil, &sdl.Rect{chunkPos.X, chunkPos.Y, scaledCS, scaledCS})
	if err != nil {
		panic(err)
	}
}

func (lr *LandscapeRender) Update() {
	lr.lock.Lock()
	defer lr.lock.Unlock()
	lr.mousePos = lr.mouseChunk.chunkPos
}

func (*LandscapeRender) RemoveEntity(*ecs.BasicEntity) {}
