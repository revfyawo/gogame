package components

const (
	HeightNoiseStep = 0.01
	TempNoiseStep   = 0.003
	RainNoiseStep   = 0.003
)

type Chunk struct {
	Height [][]float64
	Rain   [][]Rain
	Temp   [][]Temperature
	Biomes [][]Biome
}

func (c *Chunk) Generate(heightNoise, tempNoise, rainNoise Noise, x, y int32) {
	c.Height = make([][]float64, ChunkTile)
	c.Rain = make([][]Rain, ChunkTile)
	c.Temp = make([][]Temperature, ChunkTile)
	c.Biomes = make([][]Biome, ChunkTile)

	var noise float64
	var rain Rain
	var temp Temperature
	for i := 0; i < ChunkTile; i++ {
		c.Height[i] = make([]float64, ChunkTile)
		c.Rain[i] = make([]Rain, ChunkTile)
		c.Temp[i] = make([]Temperature, ChunkTile)
		c.Biomes[i] = make([]Biome, ChunkTile)
		for j := 0; j < ChunkTile; j++ {
			xn := float64(x*ChunkTile + int32(i))
			yn := float64(y*ChunkTile + int32(j))

			noise = heightNoise.EvalOctaves(xn, yn, HeightNoiseStep, 5)
			c.Height[i][j] = noise

			noise = rainNoise.EvalOctaves(xn, yn, RainNoiseStep, 5)
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

			noise = tempNoise.EvalOctaves(xn, yn, TempNoiseStep, 5)
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
		}
	}
	c.generateBiomes()
}

func (c *Chunk) generateBiomes() {
	for x := 0; x < ChunkTile; x++ {
		for y := 0; y < ChunkTile; y++ {
			rain := c.Rain[x][y]
			temp := c.Temp[x][y]
			c.Biomes[x][y] = Biomes[temp][rain]
		}
	}
}
