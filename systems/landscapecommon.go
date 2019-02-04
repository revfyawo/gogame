package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/entities"
	"github.com/veandco/go-sdl2/sdl"
)

type Landscapes struct {
	landscapes map[components.Biome][]*entities.Landscape
}

func NewLandscapes() *Landscapes {
	return &Landscapes{landscapes: make(map[components.Biome][]*entities.Landscape)}
}

func (ls *Landscapes) Add(l *entities.Landscape) {
	ls.landscapes[l.Biome] = append(ls.landscapes[l.Biome], l)
}

func (ls *Landscapes) Find(chunk, tile sdl.Point, biome components.Biome) *entities.Landscape {
	biomeLs, ok := ls.landscapes[biome]
	if !ok || biomeLs == nil {
		return nil
	}
	for _, l := range biomeLs {
		if l.Contains(chunk, tile) {
			return l
		}
	}
	return nil
}

func (ls *Landscapes) Merge(other *Landscapes) {
	for biome, list := range other.landscapes {
		ls.landscapes[biome] = append(ls.landscapes[biome], list...)
	}
}
