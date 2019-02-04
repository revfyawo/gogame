package components

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Landscape struct {
	Tiles map[sdl.Point]map[sdl.Point]bool
}

func (l *Landscape) AddTile(chunk, tile sdl.Point) {
	if _, ok := l.Tiles[chunk]; !ok {
		l.Tiles[chunk] = make(map[sdl.Point]bool)
	}
	l.Tiles[chunk][tile] = true
}

func (l *Landscape) Contains(chunk, tile sdl.Point) bool {
	c, ok := l.Tiles[chunk]
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

	for chunkPoint, chunk := range other.Tiles {
		for tilePoint, tile := range chunk {
			l.Tiles[chunkPoint][tilePoint] = tile
		}
	}
}

func (l *Landscape) Size() int {
	var size int
	for _, chunk := range l.Tiles {
		size += len(chunk)
	}
	return size
}

func (l *Landscape) RemoveChunk(chunk sdl.Point) {
	delete(l.Tiles, chunk)
}
