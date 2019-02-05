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

func (ls *Landscapes) Index(l *entities.Landscape) int {
	var index int
	var landscape *entities.Landscape
	for index, landscape = range ls.landscapes[l.Biome] {
		if l == landscape {
			return index
		}
	}
	return index
}

func (ls *Landscapes) Merge(other *Landscapes) {
	for biome, list := range other.landscapes {
		ls.landscapes[biome] = append(ls.landscapes[biome], list...)
	}
}

func (ls *Landscapes) MergeLandscape(l1, l2 *entities.Landscape) *entities.Landscape {
	if l1 == nil {
		return l2
	} else if l2 == nil {
		return l1
	}
	if l1.Biome != l2.Biome {
		panic("can't merge two landscapes with different biome")
	}
	biome := l1.Biome
	l1Index := ls.Index(l1)
	if l1Index == len(ls.landscapes[biome]) {
		panic("can't find l1 in landscapes")
	}
	l2Index := ls.Index(l2)
	if l2Index == len(ls.landscapes[biome]) {
		panic("can't find l2 in landscapes")
	}

	var small, big *entities.Landscape
	var removeIndex int
	if l1.Size() < l2.Size() {
		small = l1
		removeIndex = l1Index
		big = l2
	} else {
		small = l2
		removeIndex = l2Index
		big = l1
	}
	big.Merge(small)
	ls.landscapes[biome] = append(ls.landscapes[biome][:removeIndex], ls.landscapes[biome][removeIndex+1:]...)
	return big
}
