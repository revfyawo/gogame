package components

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Landscape struct {
	Tiles            map[sdl.Point]map[sdl.Point]bool
	chunkRect        sdl.Rect
	border           map[sdl.Point]map[sdl.Point]bool
	regenerateBorder bool
}

func (l *Landscape) AddTile(chunk, tile sdl.Point) {
	if _, ok := l.Tiles[chunk]; !ok {
		l.Tiles[chunk] = make(map[sdl.Point]bool)
		if l.chunkRect.W == 0 && l.chunkRect.H == 0 {
			l.chunkRect = sdl.Rect{chunk.X, chunk.Y, 1, 1}
		} else if diff := l.chunkRect.X - chunk.X; diff > 0 {
			l.chunkRect.X -= diff
			l.chunkRect.W += diff
		} else if diff := l.chunkRect.Y - chunk.Y; diff > 0 {
			l.chunkRect.Y -= diff
			l.chunkRect.H += diff
		} else if diff := chunk.X - l.chunkRect.X - l.chunkRect.W + 1; diff > 0 {
			l.chunkRect.W += diff
		} else if diff := chunk.Y - l.chunkRect.Y - l.chunkRect.H + 1; diff > 0 {
			l.chunkRect.H += diff
		}
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

func (l *Landscape) Border() (border map[sdl.Point]map[sdl.Point]bool, changed bool) {
	if !l.regenerateBorder && l.border != nil && len(l.border) > 0 {
		return l.border, false
	} else if l.regenerateBorder {
		l.regenerateBorder = false
	}

	border = make(map[sdl.Point]map[sdl.Point]bool)
	inside := false
	pos := ChunkTilePosition{sdl.Point{l.chunkRect.X, l.chunkRect.Y}, sdl.Point{0, 0}}
	previous := pos.Up()
	for chunkX := l.chunkRect.X; chunkX < l.chunkRect.X+l.chunkRect.W; chunkX++ {
		for chunkY := l.chunkRect.Y; chunkY < l.chunkRect.Y+l.chunkRect.H; chunkY++ {
			for tileX := int32(0); tileX < ChunkTile; tileX++ {
				for tileY := int32(0); tileY < ChunkTile; tileY++ {
					chunkPoint := sdl.Point{chunkX, chunkY}
					tilePoint := sdl.Point{tileX, tileY}
					pos = ChunkTilePosition{chunkPoint, tilePoint}
					previous = pos.Up()

					// Scan for biome
					if border[chunkPoint] != nil && border[chunkPoint][tilePoint] && !inside {
						// Crossing known border
						inside = true
					} else if l.Tiles[chunkPoint] != nil && !l.Tiles[chunkPoint][tilePoint] && inside {
						// Leaving border
						inside = false
					} else if l.Tiles[chunkPoint] != nil && l.Tiles[chunkPoint][tilePoint] && !inside {
						addBorder := l.contour(pos, previous)
						inside = true
						for chunk, maps := range addBorder {
							for tile, val := range maps {
								if border[chunk] == nil {
									border[chunk] = make(map[sdl.Point]bool)
								}
								border[chunk][tile] = val
							}
						}
					}
				}
			}
		}
	}

	l.border = border
	return border, true
}

func (l *Landscape) borderContains(border []ChunkTilePosition, pos ChunkTilePosition) bool {
	for _, borderPos := range border {
		if pos == borderPos {
			return true
		}
	}
	return false
}

func (l *Landscape) contour(first, firstPrevious ChunkTilePosition) map[sdl.Point]map[sdl.Point]bool {
	var pos, previous = first, firstPrevious
	var left, up, right, down, next ChunkTilePosition
	var border = make(map[sdl.Point]map[sdl.Point]bool)
	for {
		left = pos.Left()
		up = pos.Up()
		right = pos.Right()
		down = pos.Down()

		// Find next tile in landscape from previous clockwise
		if l.Tiles[pos.Chunk][pos.Tile] {
			if border[pos.Chunk] == nil {
				border[pos.Chunk] = make(map[sdl.Point]bool)
			}
			border[pos.Chunk][pos.Tile] = true
			// Turn left
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
			previous = pos
			pos = next
		} else {
			// Turn right
			switch previous {
			case left:
				next = down
			case up:
				next = left
			case right:
				next = up
			case down:
				next = right
			}
			previous = pos
			pos = next
		}

		if pos == first && previous == firstPrevious {
			break
		}
	}
	return border
}
