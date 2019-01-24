package ecs

import "sync"

var (
	id   uint64
	lock sync.Mutex
)

type BasicEntity struct {
	id uint64
}

type Identifier interface {
	ID() uint64
}

func (e *BasicEntity) ID() uint64 {
	return e.id
}

func NewBasic() BasicEntity {
	lock.Lock()
	defer lock.Unlock()
	e := BasicEntity{id}
	id++
	return e
}
