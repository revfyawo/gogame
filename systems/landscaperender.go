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
				lr.chunks[sdl.Point{chunk.X, chunk.Y}] = chunk
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
	chunkPoint := sdl.Point{chunk.X, chunk.Y}

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
	visible, screenPos := lr.camera.GetVisibleChunks()
	scaledCS := int32(components.ChunkSize * lr.camera.Scale())
	lr.camera.RUnlock()
	// Return if chunk not visible
	if mousePos.Chunk.X < visible.X || mousePos.Chunk.X >= visible.X+visible.W || mousePos.Chunk.Y < visible.Y || mousePos.Chunk.Y >= visible.Y+visible.H {
		return
	}

	clipRect := sdl.Rect{0, 0, landscape.ChunkRect.W, landscape.ChunkRect.H}
	if diff := visible.X - landscape.ChunkRect.X; diff > 0 {
		clipRect.X += diff
		clipRect.W -= diff
	}
	if diff := visible.Y - landscape.ChunkRect.Y; diff > 0 {
		clipRect.Y += diff
		clipRect.H -= diff
	}
	if diff := landscape.ChunkRect.X + landscape.ChunkRect.W - 1 - (visible.X + visible.W - 1); diff > 0 {
		clipRect.W -= diff
	}
	if diff := landscape.ChunkRect.Y + landscape.ChunkRect.H - 1 - (visible.Y + visible.H - 1); diff > 0 {
		clipRect.H -= diff
	}
	landWidth := scaledCS * clipRect.W
	landHeight := scaledCS * clipRect.H
	landPos := screenPos[sdl.Point{landscape.ChunkRect.X + clipRect.X, landscape.ChunkRect.Y + clipRect.Y}]

	clipRect.X *= components.ChunkSize
	clipRect.Y *= components.ChunkSize
	clipRect.W *= components.ChunkSize
	clipRect.H *= components.ChunkSize

	// Generating border texture if different from last frame
	if landscape != lr.lastLandscape || changed || landscape.BorderTex == nil {
		surface, err := sdl.CreateRGBSurface(0, landscape.ChunkRect.W*components.ChunkSize, landscape.ChunkRect.H*components.ChunkSize, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
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
		for x := int32(0); x < landscape.ChunkRect.W; x++ {
			for y := int32(0); y < landscape.ChunkRect.H; y++ {
				chunkPoint := sdl.Point{landscape.ChunkRect.X + x, landscape.ChunkRect.Y + y}
				for tile := range border[chunkPoint] {
					err = surface.FillRect(&sdl.Rect{x*components.ChunkSize + tile.X*components.TileSize, y*components.ChunkSize + tile.Y*components.TileSize, components.TileSize, components.TileSize}, 0xffffffff)
					if err != nil {
						panic(err)
					}
				}
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
		if landscape != lr.lastLandscape && lr.lastLandscape != nil && lr.lastLandscape.BorderTex != nil {
			err := lr.lastLandscape.BorderTex.Destroy()
			if err != nil {
				panic(err)
			}
			lr.lastLandscape.BorderTex = nil
		}
		lr.lastLandscape = landscape
	}

	// Rendering landscape border
	err := engine.Renderer.Copy(landscape.BorderTex, &clipRect, &sdl.Rect{landPos.X, landPos.Y, landWidth, landHeight})
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
