package entities

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/veandco/go-sdl2/sdl"
)

type Landscape struct {
	ecs.BasicEntity
	components.Biome
	components.Landscape
}

func NewLandscape(biome components.Biome) *Landscape {
	return &Landscape{BasicEntity: ecs.NewBasic(), Biome: biome, Landscape: components.Landscape{Tiles: map[sdl.Point]map[sdl.Point]bool{}}}
}

func (l *Landscape) Merge(other *Landscape) {
	if other == nil {
		return
	}
	if l.Biome != other.Biome {
		panic("can't merge two landscapes with different biome")
	}

	for chunkPoint, chunk := range other.Tiles {
		for tilePoint, tile := range chunk {
			l.Tiles[chunkPoint][tilePoint] = tile
		}
	}
}
