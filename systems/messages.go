package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/entities"
)

const (
	NewChunkMessageType = iota
	ChangeScaleMessageType
	SetCameraPositionMessageType
)

type NewChunkMessage struct {
	Chunk *entities.Chunk
}

func (*NewChunkMessage) Type() ecs.MessageType {
	return NewChunkMessageType
}

type ChangeScaleMessage struct {
	Scale float64
}

func (*ChangeScaleMessage) Type() ecs.MessageType {
	return ChangeScaleMessageType
}

type SetCameraPositionMessage struct {
	Position components.ChunkPosition
}

func (*SetCameraPositionMessage) Type() ecs.MessageType {
	return SetCameraPositionMessageType
}
