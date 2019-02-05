package components

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Landscape struct {
	Tiles            map[sdl.Point]map[sdl.Point]bool
	border           []ChunkTilePosition
	regenerateBorder bool
}

func (l *Landscape) AddTile(chunk, tile sdl.Point) {
	if _, ok := l.Tiles[chunk]; !ok {
		l.Tiles[chunk] = make(map[sdl.Point]bool)
	}
	l.Tiles[chunk][tile] = true
	l.regenerateBorder = true
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
	l.regenerateBorder = true
}

func (l *Landscape) Size() int {
	var size int
	for _, chunk := range l.Tiles {
		size += len(chunk)
	}
	return size
}

func (l *Landscape) Border() (border []ChunkTilePosition, changed bool) {
	if !l.regenerateBorder && l.border != nil && len(l.border) > 0 {
		return l.border, false
	} else if l.regenerateBorder {
		l.regenerateBorder = false
	}

	// Find random point in landscape (first one)
	var randPos ChunkTilePosition
	for chunk, chunks := range l.Tiles {
		for tile := range chunks {
			randPos = ChunkTilePosition{chunk, tile}
			break
		}
		break
	}
	pos := randPos

	// Go left until we leave landscape
	for l.Tiles[pos.Chunk][pos.Tile] != false {
		pos.MoveX(-1)
	}
	pos.MoveX(1)
	// Set first chunk & tile to leftmost point from random first one
	firstPos := pos
	border = append(border, firstPos)

	// Find next one in border
	up := firstPos.Up()
	right := firstPos.Right()
	down := firstPos.Down()
	if l.Tiles[up.Chunk][up.Tile] {
		border = append(border, up)
		pos = up
	} else if l.Tiles[right.Chunk][right.Tile] {
		border = append(border, right)
		pos = right
	} else if l.Tiles[down.Chunk][down.Tile] {
		border = append(border, down)
		pos = down
	}
	previous := firstPos

	// Find previous one in order
	firstPrevious := firstPos
	if l.Tiles[down.Chunk][down.Tile] {
		firstPrevious = down
	} else if l.Tiles[right.Chunk][right.Tile] {
		firstPrevious = right
	} else if l.Tiles[up.Chunk][up.Tile] {
		firstPrevious = up
	}

	var left ChunkTilePosition
	for pos != firstPos || previous != firstPrevious {
		left = pos.Left()
		up = pos.Up()
		right = pos.Right()
		down = pos.Down()

		// Find next tile in landscape from previous clockwise
		var next ChunkTilePosition
		for i := 0; i < 4; i++ {
			switch previous {
			case left:
				next = up
			case up:
				next = right
			case right:
				next = down
			case down:
				next = left
			}

			if l.Tiles[next.Chunk][next.Tile] {
				border = append(border, next)
				previous = pos
				pos = next
				break
			} else {
				previous = next
			}
		}
	}
	l.border = border
	return border, true
}
