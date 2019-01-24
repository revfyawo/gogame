package engine

import "github.com/revfyawo/gogame/ecs"

type MessageManager struct {
	listeners map[ecs.MessageType][]ecs.MessageSystem
}

func NewMessageManager() *MessageManager {
	mm := MessageManager{}
	mm.listeners = make(map[ecs.MessageType][]ecs.MessageSystem)
	return &mm
}

func (mm *MessageManager) Listen(messType ecs.MessageType, system ecs.MessageSystem) int {
	mm.listeners[messType] = append(mm.listeners[messType], system)
	return len(mm.listeners[messType]) - 1
}

func (mm *MessageManager) StopListen(message ecs.Message, id int) {
	messType := message.Type()
	if id == len(mm.listeners[messType])-1 {
		mm.listeners[messType] = mm.listeners[messType][:id]
	} else {
		mm.listeners[messType] = append(mm.listeners[messType][:id], mm.listeners[messType][id+1:]...)
	}
}

func (mm *MessageManager) Dispatch(message ecs.Message) {
	messType := message.Type()
	listeners := mm.listeners[messType]
	if listeners == nil {
		return
	}
	for _, sys := range listeners {
		sys.PushMessage(message)
	}
}
