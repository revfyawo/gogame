package main

import (
	"github.com/ojrac/opensimplex-go"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	chunksX    = 5
	chunksY    = 5
	waterColor = 0xff0000ff
	waterLevel = -0.5
	sandColor  = 0xffffff00
	sandDiff   = 0.1
	grassColor = 0xff00ff00
	rockColor  = 0xff808080
	rockDiff   = 0.2
	snowColor  = 0xffffffff
	snowLevel  = 0.7
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		chunksX*8*16, chunksY*16*8, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	world := make([][]uint32, chunksX*16)
	for i := range world {
		world[i] = make([]uint32, chunksY*16)
	}

	noise := opensimplex.New(1234567890)
	step := 0.05
	for i := range world {
		for j := range world[i] {
			val := noise.Eval2(float64(i)*step, float64(j)*step)
			var color uint32
			switch {
			case val <= waterLevel:
				color = waterColor
			case val > waterLevel && val <= waterLevel+sandDiff:
				color = sandColor
			case val > snowLevel-rockDiff && val <= snowLevel:
				color = rockColor
			case val > snowLevel:
				color = snowColor
			default:
				color = grassColor
			}
			world[i][j] = color
		}
	}

	for i := range world {
		for j := range world[i] {
			surface.FillRect(&sdl.Rect{X: int32(i) * 8, Y: int32(j) * 8, W: 8, H: 8}, world[i][j])
		}
	}
	window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
