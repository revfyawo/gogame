package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
)

const zoomSpeed = 0.125

type CameraZoom struct {
	ChunkPos components.ChunkPosition
	camera   *Camera
}

func (cz *CameraZoom) New(world *ecs.World) {
	camera := false
	for _, sys := range world.UpdateSystems() {
		switch s := sys.(type) {
		case *Camera:
			camera = true
			cz.camera = s
		}
	}
	if !camera {
		panic("need to add camera system before mouse system")
	}
}

func (cz *CameraZoom) Update() {
	wheel := engine.Input.Wheel()
	switch wheel {
	case 1:
		engine.Message.Dispatch(ChangeScaleMessage{cz.camera.Scale() * (1 + zoomSpeed)})
	case -1:
		engine.Message.Dispatch(ChangeScaleMessage{cz.camera.Scale() * (1 - zoomSpeed)})
	}
}

func (*CameraZoom) RemoveEntity(*ecs.BasicEntity) {}
