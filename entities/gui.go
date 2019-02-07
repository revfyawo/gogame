package entities

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type GUIPosition int

const (
	GUIPosTop GUIPosition = iota
)

type GUI struct {
	Name     string
	Position GUIPosition
	components.GUI
	components.Rect
	components.Render
}

func (gui *GUI) DestroyTexture() {
	if gui.Texture == nil {
		return
	}
	err := gui.Texture.Destroy()
	if err != nil {
		panic(err)
	}
	gui.Texture = nil
}

func (gui *GUI) GenerateTexture() {
	switch gui.Type {
	case components.GUILabel:
		fontSurface, err := gui.Font.RenderUTF8Blended(gui.Label, sdl.Color{0xff, 0xff, 0xff, 0xff})
		if err != nil {
			panic(err)
		}
		defer fontSurface.Free()

		var surface *sdl.Surface
		if !gui.Background {
			surface = fontSurface
		} else {
			surface, err = sdl.CreateRGBSurface(0, fontSurface.W+int32(gui.Font.Height()), fontSurface.H, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			if err != nil {
				panic(err)
			}
			defer surface.Free()

			// Fill background
			err = surface.FillRect(nil, 0x80000000)
			if err != nil {
				panic(err)
			}

			// Blit text
			err = fontSurface.Blit(nil, surface, &sdl.Rect{int32(gui.Font.Height()) / 2, 0, fontSurface.W, fontSurface.H})
			if err != nil {
				panic(err)
			}
		}

		// Destroy old texture, and create new one
		gui.DestroyTexture()
		texture, err := engine.Renderer.CreateTextureFromSurface(surface)
		if err != nil {
			panic(err)
		}

		// Assign gui texture, width and height
		gui.Texture, gui.W, gui.H = texture, surface.W, surface.H
	}
}
