package components

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Landscape struct {
	biome Biome
	tiles map[sdl.Point]map[sdl.Point]bool
}

func NewLandscape(biome Biome) *Landscape {
	return &Landscape{biome: biome, tiles: make(map[sdl.Point]map[sdl.Point]bool)}
}

func (l *Landscape) AddTile(chunk, tile sdl.Point) {
	if _, ok := l.tiles[chunk]; !ok {
		l.tiles[chunk] = make(map[sdl.Point]bool)
	}
	l.tiles[chunk][tile] = true
}

func (l *Landscape) Contains(chunk, tile sdl.Point) bool {
	c, ok := l.tiles[chunk]
	if !ok {
		return false
	} else {
		t, ok := c[tile]
		if !ok {
			return false
		}
		return t
	}
}

func (l *Landscape) Merge(other *Landscape) {
	if other == nil {
		return
	}
	if l.biome != other.biome {
		panic("can't merge two landscapes with different biome")
	}

	for chunkPoint, chunk := range other.tiles {
		for tilePoint, tile := range chunk {
			l.tiles[chunkPoint][tilePoint] = tile
		}
	}
}

func (l *Landscape) Size() int {
	var size int
	for _, chunk := range l.tiles {
		size += len(chunk)
	}
	return size
}

func (l *Landscape) RemoveChunk(chunk sdl.Point) {
	delete(l.tiles, chunk)
}

type Landscapes struct {
	landscapes map[Biome][]*Landscape
}

func NewLandscapes() *Landscapes {
	return &Landscapes{landscapes: make(map[Biome][]*Landscape)}
}

func (ls *Landscapes) Add(l *Landscape) {
	ls.landscapes[l.biome] = append(ls.landscapes[l.biome], l)
}

func (ls *Landscapes) Find(chunk, tile sdl.Point, biome Biome) *Landscape {
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
