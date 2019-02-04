package ecs

import (
	"log"
)

type World struct {
	updateSystems []UpdateSystem
	renderSystems []RenderSystem
}

func (w *World) AddSystem(sys System) {
	init, initOK := sys.(Initializer)
	update, updateOK := sys.(UpdateSystem)
	render, renderOK := sys.(RenderSystem)
	if !updateOK && !renderOK {
		log.Panic("system", sys, "is neither an UpdateSystem nor a RenderSystem")
	}
	if initOK {
		init.New(w)
	}
	if updateOK {
		w.updateSystems = append(w.updateSystems, update)
	}
	if renderOK {
		w.renderSystems = append(w.renderSystems, render)
	}
}

func (w *World) Update() {
	for _, s := range w.updateSystems {
		s.Update()
	}
}

func (w *World) UpdateRender() {
	for _, s := range w.renderSystems {
		s.UpdateFrame()
	}
}

func (w *World) UpdateSystems() []UpdateSystem {
	return w.updateSystems
}

func (w *World) RenderSystems() []RenderSystem {
	return w.renderSystems
}
