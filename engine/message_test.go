package engine

import (
	"github.com/revfyawo/gogame/ecs"
	"testing"
)

type MessageTest struct{}

func (MessageTest) Type() ecs.MessageType {
	return 0
}

type MessageSystemTest struct {
	received ecs.Message
}

func (sys *MessageSystemTest) PushMessage(m ecs.Message) {
	sys.received = m
}

func TestMessageManager(t *testing.T) {
	mm := NewMessageManager()
	mess := &MessageTest{}
	sys := &MessageSystemTest{}

	mm.Dispatch(mess)
	if sys.received != nil {
		t.Error("sys received a message without listening")
	}

	mm.Listen(mess.Type(), sys)
	if mm.listeners[mess.Type()] == nil {
		t.Fatal("mm did not register sys")
	}

	mm.Dispatch(mess)
	if sys.received == nil {
		t.Error("sys did not receive the message")
	}
}
