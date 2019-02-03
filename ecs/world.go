package ecs

import "time"

type World struct {
	systems       []System
	renderSystems []RenderSystem
}

func (w *World) AddSystem(s System) {
	init, ok := s.(Initializer)
	if ok {
		init.New(w)
	}
	w.systems = append(w.systems, s)
}

func (w *World) AddRenderSystem(s RenderSystem) {
	init, ok := s.(Initializer)
	if ok {
		init.New(w)
	}
	w.renderSystems = append(w.renderSystems, s)
}

func (w *World) Update(d time.Duration) {
	for _, s := range w.systems {
		s.Update(d)
	}
}

func (w *World) UpdateRender(d time.Duration) {
	for _, s := range w.renderSystems {
		s.UpdateFrame(d)
	}
}

func (w *World) Systems() []System {
	return w.systems
}
