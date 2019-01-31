package systems

import (
	"github.com/revfyawo/gogame/components"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"time"
)

const zoomSpeed = 0.125

type Mouse struct {
	ChunkPos components.ChunkPosition
	camera   *Camera
}

func (ms *Mouse) New(world *ecs.World) {
	camera := false
	for _, sys := range world.Systems() {
		switch s := sys.(type) {
		case *Camera:
			camera = true
			ms.camera = s
		}
	}
	if !camera {
		panic("need to add camera system before mouse system")
	}
}

func (ms *Mouse) Update(time.Duration) {
	wheel := engine.Input.Wheel()
	switch wheel {
	case 1:
		ms.camera.Scale *= 1 + zoomSpeed
	case -1:
		ms.camera.Scale *= 1 - zoomSpeed
	}
}

func (*Mouse) RemoveEntity(*ecs.BasicEntity) {}
