package components

import (
	"github.com/ojrac/opensimplex-go"
)

const (
	HeightNoiseStep = 0.01
	TempNoiseStep   = 0.003
	RainNoiseStep   = 0.003

	WaterLevel  = -0.9
	ShallowDiff = 0.2
)

type Chunk struct {
	Height [][]float64
	Rain   [][]Rain
	Temp   [][]Temperature
	Biomes [][]Biome
}

func (c *Chunk) Generate(heightNoise, tempNoise, rainNoise opensimplex.Noise, x, y int32) {
	c.Height = make([][]float64, ChunkTile)
	c.Rain = make([][]Rain, ChunkTile)
	c.Temp = make([][]Temperature, ChunkTile)
	c.Biomes = make([][]Biome, ChunkTile)

	var noise, noise0, noise1, noise2, noise3, noise4 float64
	var rain Rain
	var temp Temperature
	for i := 0; i < ChunkTile; i++ {
		c.Height[i] = make([]float64, ChunkTile)
		c.Rain[i] = make([]Rain, ChunkTile)
		c.Temp[i] = make([]Temperature, ChunkTile)
		c.Biomes[i] = make([]Biome, ChunkTile)
		for j := 0; j < ChunkTile; j++ {
			noise0 = heightNoise.Eval2(float64(x*ChunkTile+int32(i))*HeightNoiseStep, float64(y*ChunkTile+int32(j))*HeightNoiseStep) * 16
			noise1 = heightNoise.Eval2(float64(x*ChunkTile+int32(i))*2*HeightNoiseStep, float64(y*ChunkTile+int32(j))*2*HeightNoiseStep) * 8
			noise2 = heightNoise.Eval2(float64(x*ChunkTile+int32(i))*4*HeightNoiseStep, float64(y*ChunkTile+int32(j))*4*HeightNoiseStep) * 4
			noise3 = heightNoise.Eval2(float64(x*ChunkTile+int32(i))*8*HeightNoiseStep, float64(y*ChunkTile+int32(j))*8*HeightNoiseStep) * 2
			noise4 = heightNoise.Eval2(float64(x*ChunkTile+int32(i))*16*HeightNoiseStep, float64(y*ChunkTile+int32(j))*16*HeightNoiseStep)
			noise = (noise0 + noise1 + noise2 + noise3 + noise4) / 16
			c.Height[i][j] = noise
			water := false
			if noise < WaterLevel {
				water = true
				c.Biomes[i][j] = DeepWater
			} else if noise < WaterLevel+ShallowDiff {
				water = true
				c.Biomes[i][j] = ShallowWater
			}

			noise0 = rainNoise.Eval2(float64(x*ChunkTile+int32(i))*RainNoiseStep, float64(y*ChunkTile+int32(j))*RainNoiseStep) * 16
			noise1 = rainNoise.Eval2(float64(x*ChunkTile+int32(i))*2*RainNoiseStep, float64(y*ChunkTile+int32(j))*2*RainNoiseStep) * 8
			noise2 = rainNoise.Eval2(float64(x*ChunkTile+int32(i))*4*RainNoiseStep, float64(y*ChunkTile+int32(j))*4*RainNoiseStep) * 4
			noise3 = rainNoise.Eval2(float64(x*ChunkTile+int32(i))*8*RainNoiseStep, float64(y*ChunkTile+int32(j))*8*RainNoiseStep) * 2
			noise4 = rainNoise.Eval2(float64(x*ChunkTile+int32(i))*16*RainNoiseStep, float64(y*ChunkTile+int32(j))*16*RainNoiseStep)
			noise = (noise0 + noise1 + noise2 + noise3 + noise4) / 16
			if noise < -0.6 {
				rain = Aridest
			} else if noise < -0.2 {
				rain = Arid
			} else if noise < 0.2 {
				rain = Moderate
			} else if noise < 0.6 {
				rain = Wet
			} else {
				rain = Wettest
			}
			c.Rain[i][j] = rain

			noise0 = tempNoise.Eval2(float64(x*ChunkTile+int32(i))*TempNoiseStep, float64(y*ChunkTile+int32(j))*TempNoiseStep) * 16
			noise1 = tempNoise.Eval2(float64(x*ChunkTile+int32(i))*2*TempNoiseStep, float64(y*ChunkTile+int32(j))*2*TempNoiseStep) * 8
			noise2 = tempNoise.Eval2(float64(x*ChunkTile+int32(i))*4*TempNoiseStep, float64(y*ChunkTile+int32(j))*4*TempNoiseStep) * 4
			noise3 = tempNoise.Eval2(float64(x*ChunkTile+int32(i))*8*TempNoiseStep, float64(y*ChunkTile+int32(j))*8*TempNoiseStep) * 2
			noise4 = tempNoise.Eval2(float64(x*ChunkTile+int32(i))*16*TempNoiseStep, float64(y*ChunkTile+int32(j))*16*TempNoiseStep)
			noise = (noise0 + noise1 + noise2 + noise3 + noise4) / 16
			if noise < -0.6 {
				temp = Coldest
			} else if noise < -0.2 {
				temp = Cold
			} else if noise < 0.2 {
				temp = Temperate
			} else if noise < 0.6 {
				temp = Hot
			} else {
				temp = Hottest
			}
			c.Temp[i][j] = temp

			if !water {
				c.Biomes[i][j] = Biomes[temp][rain]
			}
		}
	}
}
