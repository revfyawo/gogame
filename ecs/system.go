package ecs

import "time"

type System interface {
	Update(d time.Duration)
	RemoveEntity(e *BasicEntity)
}

type Initializer interface {
	New()
}
