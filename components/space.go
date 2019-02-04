package components

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Space struct {
	Rect sdl.Rect
}

const (
	TileSize  = 1
	ChunkTile = 32
	ChunkSize = ChunkTile * TileSize
)

type ChunkPosition struct {
	Chunk    sdl.Point
	Position sdl.Point
}

func (cp *ChunkPosition) MoveX(speed int32) {
	newX := cp.Position.X + speed
	if newX >= ChunkSize {
		cp.Position.X = newX % ChunkSize
		cp.Chunk.X += int32(math.Floor(float64(newX) / ChunkSize))
	} else if newX < 0 {
		div := int32(math.Floor(float64(newX) / ChunkSize))
		newX -= div * ChunkSize
		cp.Chunk.X += div
		cp.Position.X = newX
	} else {
		cp.Position.X = newX
	}
}

func (cp *ChunkPosition) MoveY(speed int32) {
	newY := cp.Position.Y + speed
	if newY >= ChunkSize {
		cp.Position.Y = newY % ChunkSize
		cp.Chunk.Y += int32(math.Floor(float64(newY) / ChunkSize))
	} else if newY < 0 {
		div := int32(math.Floor(float64(newY) / ChunkSize))
		newY -= div * ChunkSize
		cp.Chunk.Y += div
		cp.Position.Y = newY
	} else {
		cp.Position.Y = newY
	}
}
