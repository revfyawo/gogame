package ecs

import "time"

type World struct {
	systems []System
}

func (w *World) AddSystem(s System) {
	init, ok := s.(Initializer)
	if ok {
		init.New()
	}
	w.systems = append(w.systems, s)
}

func (w *World) Update(d time.Duration) {
	for _, s := range w.systems {
		s.Update(d)
	}
}
