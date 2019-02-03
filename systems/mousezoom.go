package systems

import (
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type MouseZoom struct {
	camera     *Camera
	mouseChunk *MouseChunk
}

func (mz *MouseZoom) New(world *ecs.World) {
	camera := false
	mouseChunk := false
	for _, sys := range world.UpdateSystems() {
		switch s := sys.(type) {
		case *Camera:
			mz.camera = s
			camera = true
		case *MouseChunk:
			mz.mouseChunk = s
			mouseChunk = true
		}
	}
	if !camera {
		panic("need to add camera system before mouse zoom system")
	}
	if !mouseChunk {
		panic("need to add mouse chunk system before mouse zoom system")
	}
}

func (mz *MouseZoom) Update() {
	wheel := engine.Input.Wheel()
	var newScale float64
	switch wheel {
	case 1:
		newScale = mz.camera.Scale() * (1 + zoomSpeed)
	case -1:
		newScale = mz.camera.Scale() * (1 - zoomSpeed)
	default:
		return
	}
	engine.Message.Dispatch(ChangeScaleMessage{newScale})

	mousePos := engine.Input.MousePosition()
	w, h, err := engine.Renderer.GetOutputSize()
	if err != nil {
		panic(err)
	}
	diff := sdl.Point{w/2 - mousePos.X, h/2 - mousePos.Y}
	newCamPos := mz.mouseChunk.chunkPos
	newCamPos.MoveX(int32(float64(diff.X) / newScale))
	newCamPos.MoveY(int32(float64(diff.Y) / newScale))
	engine.Message.Dispatch(SetCameraPositionMessage{newCamPos})
}

func (*MouseZoom) RemoveEntity(*ecs.BasicEntity) {}
