package systems

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/entities"
)

const (
	NewChunkMessageType = iota
)

type NewChunkMessage struct {
	Chunk *entities.Chunk
}

func (*NewChunkMessage) Type() ecs.MessageType {
	return NewChunkMessageType
}
