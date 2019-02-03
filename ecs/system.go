package ecs

import "time"

type System interface {
	RemoveEntity(*BasicEntity)
}

type UpdateSystem interface {
	Update(time.Duration)
	RemoveEntity(*BasicEntity)
}

type RenderSystem interface {
	UpdateFrame(time.Duration)
	RemoveEntity(*BasicEntity)
}

type Initializer interface {
	New(*World)
}

type MessageType uint64

type Message interface {
	Type() MessageType
}
