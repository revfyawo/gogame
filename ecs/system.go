package ecs

import "time"

type System interface {
	Update(time.Duration)
	RemoveEntity(*BasicEntity)
}

type Initializer interface {
	New(*World)
}

type MessageType uint64

type Message interface {
	Type() MessageType
}
