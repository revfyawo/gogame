package components

import "github.com/veandco/go-sdl2/sdl"

type Render struct {
	Texture *sdl.Texture
}

type ChunkRender struct {
	Textures [][]*sdl.Texture
}
