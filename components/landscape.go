package components

import (
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type Landscape struct {
	Tiles            map[sdl.Point]map[sdl.Point]bool
	ChunkRect        sdl.Rect
	size             int
	lock             sync.Mutex
	border           map[sdl.Point]map[sdl.Point]bool
	regenerateBorder bool
}

func (l *Landscape) AddTile(chunk, tile sdl.Point) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.addtile(chunk, tile)
	l.regenerateBorder = true
}

func (l *Landscape) addtile(chunk, tile sdl.Point) {
	if _, ok := l.Tiles[chunk]; !ok {
		l.Tiles[chunk] = make(map[sdl.Point]bool)
		if l.ChunkRect.W == 0 && l.ChunkRect.H == 0 {
			l.ChunkRect = sdl.Rect{chunk.X, chunk.Y, 1, 1}
		} else if diff := l.ChunkRect.X - chunk.X; diff > 0 {
			l.ChunkRect.X -= diff
			l.ChunkRect.W += diff
		} else if diff := l.ChunkRect.Y - chunk.Y; diff > 0 {
			l.ChunkRect.Y -= diff
			l.ChunkRect.H += diff
		} else if diff := chunk.X - l.ChunkRect.X - l.ChunkRect.W + 1; diff > 0 {
			l.ChunkRect.W += diff
		} else if diff := chunk.Y - l.ChunkRect.Y - l.ChunkRect.H + 1; diff > 0 {
			l.ChunkRect.H += diff
		}
	}
	l.Tiles[chunk][tile] = true
	l.size++
}

func (l *Landscape) Contains(chunk, tile sdl.Point) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
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
	l.lock.Lock()
	defer l.lock.Unlock()
	other.lock.Lock()
	defer other.lock.Unlock()
	for chunkPoint, chunk := range other.Tiles {
		for tilePoint, tile := range chunk {
			if tile {
				l.addtile(chunkPoint, tilePoint)
			}
		}
	}
	l.regenerateBorder = true
}

func (l *Landscape) Size() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.size
}

func (l *Landscape) Border() (border map[sdl.Point]map[sdl.Point]bool, changed bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.regenerateBorder && l.border != nil && len(l.border) > 0 {
		return l.border, false
	} else if l.regenerateBorder {
		l.regenerateBorder = false
	}

	border = make(map[sdl.Point]map[sdl.Point]bool)
	inside := false
	var pos, previous ChunkTilePosition
	for chunkX := l.ChunkRect.X; chunkX < l.ChunkRect.X+l.ChunkRect.W; chunkX++ {
		for tileX := int32(0); tileX < ChunkTile; tileX++ {
			for chunkY := l.ChunkRect.Y; chunkY < l.ChunkRect.Y+l.ChunkRect.H; chunkY++ {
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
