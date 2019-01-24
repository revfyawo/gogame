package ecs

import "time"

type System interface {
	Update(d time.Duration)
	RemoveEntity(e *BasicEntity)
}

type Initializer interface {
	New()
}

type MessageType uint64

type Message interface {
	Type() MessageType
}

type MessageSystem interface {
	PushMessage(m Message)
}
