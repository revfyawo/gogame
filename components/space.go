package components

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Position struct {
	sdl.Point
}

type Rect struct {
	sdl.Rect
}

const (
	TileSize  = 1
	ChunkTile = 32
	ChunkSize = ChunkTile * TileSize
)

type ChunkPosition struct {
	Chunk    sdl.Point
	Position sdl.Point // Position in that chunk, in pixels (not scaled)
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

type ChunkTilePosition struct {
	Chunk sdl.Point
	Tile  sdl.Point
}

func (ctp *ChunkTilePosition) MoveX(tile int32) {
	newTile := ctp.Tile.X + tile
	if newTile >= ChunkTile {
		ctp.Tile.X = newTile % ChunkTile
		ctp.Chunk.X += int32(math.Floor(float64(newTile) / ChunkTile))
	} else if newTile < 0 {
		div := int32(math.Floor(float64(newTile) / ChunkTile))
		newTile -= div * ChunkTile
		ctp.Chunk.X += div
		ctp.Tile.X = newTile
	} else {
		ctp.Tile.X = newTile
	}
}

func (ctp *ChunkTilePosition) MoveY(tile int32) {
	newTile := ctp.Tile.Y + tile
	if newTile >= ChunkTile {
		ctp.Tile.Y = newTile % ChunkTile
		ctp.Chunk.Y += int32(math.Floor(float64(newTile) / ChunkTile))
	} else if newTile < 0 {
		div := int32(math.Floor(float64(newTile) / ChunkTile))
		newTile -= div * ChunkTile
		ctp.Chunk.Y += div
		ctp.Tile.Y = newTile
	} else {
		ctp.Tile.Y = newTile
	}
}

func (ctp *ChunkTilePosition) Left() ChunkTilePosition {
	left := *ctp
	left.MoveX(-1)
	return left
}

func (ctp *ChunkTilePosition) Up() ChunkTilePosition {
	up := *ctp
	up.MoveY(-1)
	return up
}

func (ctp *ChunkTilePosition) Right() ChunkTilePosition {
	right := *ctp
	right.MoveX(1)
	return right
}

func (ctp *ChunkTilePosition) Down() ChunkTilePosition {
	down := *ctp
	down.MoveY(1)
	return down
}
