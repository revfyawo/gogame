package components

import "github.com/veandco/go-sdl2/ttf"

type FontType int

const (
	MonoRegular FontType = iota
)

var fonts = map[FontType]map[int]*ttf.Font{}

func GetFont(fontType FontType, size int) *ttf.Font {
	if fonts[fontType] == nil {
		fonts[fontType] = map[int]*ttf.Font{}
	}
	if fonts[fontType][size] == nil {
		var font *ttf.Font
		var err error
		switch fontType {
		case MonoRegular:
			font, err = ttf.OpenFont("assets/fonts/Go-Mono.ttf", size)
		}
		if err != nil {
			panic(err)
		}
		fonts[fontType][size] = font
	}
	return fonts[fontType][size]
}
