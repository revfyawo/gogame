package engine

import "github.com/revfyawo/gogame/ecs"

type MessageManager struct {
	listeners map[ecs.MessageType][]chan<- ecs.Message
}

func NewMessageManager() *MessageManager {
	mm := MessageManager{}
	mm.listeners = make(map[ecs.MessageType][]chan<- ecs.Message)
	return &mm
}

func (mm *MessageManager) Listen(messType ecs.MessageType, ch chan<- ecs.Message) int {
	mm.listeners[messType] = append(mm.listeners[messType], ch)
	return len(mm.listeners[messType]) - 1
}

// Need fix, stop listen multiple times on message will cause wrong one to stop listening
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
	for _, ch := range listeners {
		go mm.send(message, ch)
	}
}

func (mm *MessageManager) send(message ecs.Message, ch chan<- ecs.Message) {
	ch <- message
}
