package systems

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/entities"
)

const (
	NewChunkMessageType = iota
	ChangeScaleMessageType
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
