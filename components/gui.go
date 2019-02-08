package components

import (
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type GUIComponentType int

const (
	GUILabel GUIComponentType = iota
	GUIVLayout
	GUIButton
)

type GUI struct {
	Type       GUIComponentType
	Text       string
	Font       *ttf.Font
	Background bool
	Children   []*GUI
}

func (gui *GUI) GenerateTexture() (*sdl.Texture, int32, int32) {
	var surface *sdl.Surface
	var err error
	switch gui.Type {
	case GUILabel, GUIButton:
		surface = gui.generateSurface(0, 0)
	case GUIVLayout:
		w, h := gui.getSize()
		var surfaces []*sdl.Surface
		for _, child := range gui.Children {
			surface := child.generateSurface(w, 0)
			surfaces = append(surfaces, surface)
		}

		surface, err = sdl.CreateRGBSurface(0, w, h, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
		if err != nil {
			panic(err)
		}

		var y int32
		for _, s := range surfaces {
			err = s.Blit(nil, surface, &sdl.Rect{0, y, s.W, s.H})
			if err != nil {
				panic(err)
			}
			y += s.H
			s.Free()
		}
	}

	if gui.Background {
		bgSurface, err := sdl.CreateRGBSurface(0, surface.W, surface.H, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
		if err != nil {
			panic(err)
		}

		// Fill background
		err = bgSurface.FillRect(nil, 0x80000000)
		if err != nil {
			panic(err)
		}

		// Blit surface
		err = surface.Blit(nil, bgSurface, &sdl.Rect{0, 0, surface.W, surface.H})
		if err != nil {
			panic(err)
		}
		surface.Free()
		surface = bgSurface
	}

	if surface != nil {
		// Create texture
		texture, err := engine.Renderer.CreateTextureFromSurface(surface)
		if err != nil {
			panic(err)
		}
		surface.Free()
		// Assign gui texture, width and height
		return texture, surface.W, surface.H
	}
	return nil, 0, 0
}

func (gui *GUI) generateSurface(w, h int32) *sdl.Surface {
	var surface *sdl.Surface
	switch gui.Type {
	case GUILabel:
		fontSurface, err := gui.Font.RenderUTF8Blended(gui.Text, sdl.Color{0xff, 0xff, 0xff, 0xff})
		if err != nil {
			panic(err)
		}
		defer fontSurface.Free()

		if w > 0 && h > 0 {
			surface, err = sdl.CreateRGBSurface(0, w, h, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			err = fontSurface.Blit(nil, surface, &sdl.Rect{(w - fontSurface.W) / 2, (h - fontSurface.H) / 2, fontSurface.W, fontSurface.H})
		} else if w > 0 {
			surface, err = sdl.CreateRGBSurface(0, w, fontSurface.H, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			err = fontSurface.Blit(nil, surface, &sdl.Rect{(w - fontSurface.W) / 2, 0, fontSurface.W, fontSurface.H})
		} else if h > 0 {
			surface, err = sdl.CreateRGBSurface(0, fontSurface.W, h, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			err = fontSurface.Blit(nil, surface, &sdl.Rect{0, (h - fontSurface.H) / 2, fontSurface.W, fontSurface.H})
		} else {
			surface, err = sdl.CreateRGBSurface(0, fontSurface.W, fontSurface.H, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			err = fontSurface.Blit(nil, surface, &sdl.Rect{0, 0, fontSurface.W, fontSurface.H})
		}
		if err != nil {
			panic(err)
		}

	case GUIButton:
		fontSurface, err := gui.Font.RenderUTF8Blended(gui.Text, sdl.Color{0xff, 0xff, 0xff, 0xff})
		if err != nil {
			panic(err)
		}
		defer fontSurface.Free()

		if w > 0 && h > 0 {
			surface, err = sdl.CreateRGBSurface(0, w, h, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			err = fontSurface.Blit(nil, surface, &sdl.Rect{(w - fontSurface.W) / 2, (h - fontSurface.H) / 2, fontSurface.W, fontSurface.H})
		} else if w > 0 {
			surface, err = sdl.CreateRGBSurface(0, w, fontSurface.H, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			err = fontSurface.Blit(nil, surface, &sdl.Rect{(w - fontSurface.W) / 2, 0, fontSurface.W, fontSurface.H})
		} else if h > 0 {
			surface, err = sdl.CreateRGBSurface(0, fontSurface.W, h, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			err = fontSurface.Blit(nil, surface, &sdl.Rect{0, (h - fontSurface.H) / 2, fontSurface.W, fontSurface.H})
		} else {
			surface, err = sdl.CreateRGBSurface(0, fontSurface.W, fontSurface.H, 32, 0xff0000, 0xff00, 0xff, 0xff000000)
			err = fontSurface.Blit(nil, surface, &sdl.Rect{0, 0, fontSurface.W, fontSurface.H})
		}
		if err != nil {
			panic(err)
		}

		err = surface.FillRect(nil, 0x80808080)
		if err != nil {
			panic(err)
		}

		if w > 0 && h > 0 {
			err = fontSurface.Blit(nil, surface, &sdl.Rect{(w - fontSurface.W) / 2, (h - fontSurface.H) / 2, fontSurface.W, fontSurface.H})
		} else if w > 0 {
			err = fontSurface.Blit(nil, surface, &sdl.Rect{(w - fontSurface.W) / 2, 0, fontSurface.W, fontSurface.H})
		} else if h > 0 {
			err = fontSurface.Blit(nil, surface, &sdl.Rect{0, (h - fontSurface.H) / 2, fontSurface.W, fontSurface.H})
		} else {
			err = fontSurface.Blit(nil, surface, &sdl.Rect{0, 0, fontSurface.W, fontSurface.H})
		}
		if err != nil {
			panic(err)
		}
	}
	return surface
}

func (gui *GUI) getSize() (int32, int32) {
	switch gui.Type {
	case GUILabel, GUIButton:
		w, h, err := gui.Font.SizeUTF8(gui.Text)
		if err != nil {
			panic(err)
		}
		return int32(w), int32(h)
	case GUIVLayout:
		var w, h int32
		for _, child := range gui.Children {
			wc, hc := child.getSize()
			if wc > w {
				w = wc
			}
			h += hc
		}
		return w, h
	}
	return 0, 0
}
