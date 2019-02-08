package entities

import (
	"github.com/revfyawo/gogame/components"
)

type GUIPosition int

const (
	GUIPosTop GUIPosition = iota
	GUIPosCenter
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
	texture, w, h := gui.GUI.GenerateTexture()
	gui.Texture = texture
	gui.W = w
	gui.H = h
}
