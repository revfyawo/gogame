package main

import (
	"github.com/ojrac/opensimplex-go"
)

const (
	waterColor = 0xff0000ff
	waterLevel = -0.5
	sandColor  = 0xffffff00
	sandDiff   = 0.1
	grassColor = 0xff00ff00
	rockColor  = 0xff808080
	rockDiff   = 0.2
	snowColor  = 0xffffffff
	snowLevel  = 0.7

	chunksX = 2
	chunksY = 2
)

type World [][]uint32

func NewWorld() World {
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
	return world
}
