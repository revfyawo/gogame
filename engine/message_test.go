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
	ch chan ecs.Message
}

func TestMessageManager(t *testing.T) {
	mm := NewMessageManager()
	mess := &MessageTest{}
	sys := &MessageSystemTest{make(chan ecs.Message, 1)}

	mm.Dispatch(mess)
	select {
	case <-sys.ch:
		t.Error("sys received a message without listening")
	default:
	}

	mm.Listen(mess.Type(), sys.ch)
	if mm.listeners[mess.Type()] == nil {
		t.Fatal("mm did not register sys")
	}

	mm.Dispatch(mess)
	select {
	case <-sys.ch:
		break
	default:
		t.Error("sys did not receive the message")
	}
}
