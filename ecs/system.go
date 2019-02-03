package ecs

type System interface {
	RemoveEntity(*BasicEntity)
}

type UpdateSystem interface {
	Update()
	RemoveEntity(*BasicEntity)
}

type RenderSystem interface {
	UpdateFrame()
	RemoveEntity(*BasicEntity)
}

type Initializer interface {
	New(*World)
}

type MessageType uint64

type Message interface {
	Type() MessageType
}
