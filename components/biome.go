package components

type Temperature int
type Rain int
type Biome int

const (
	Coldest Temperature = iota
	Cold
	Temperate
	Hot
	Hottest
)

const (
	Aridest Rain = iota
	Arid
	Moderate
	Wet
	Wettest
)

const (
	HotDesert Biome = iota
	Savanna
	TropicalForest
	Grassland
	TemperateForest
	TemperateDesert
	RockDesert
	BorealForest
	ColdDesert
	FrostDesert
	Tundra
	ShallowWater
	DeepWater
)

var BiomeColors = []uint32{
	0xffff6347, // Hot desert: Tomato
	0xffe4d96f, // Savanna: Straw
	0xff0bda51, // Tropical Forest: Malachite
	0xff7cfc00, // Grassland: Lawn green
	0xff228b22, // Temperate forest: Forest green
	0xfffdee00, // Temperate desert: Cobalt yellow
	0xffc0c0c0, // Rock desert: Silver
	0xff29ab87, // Boreal forest: Jungle green
	0xffe1a95f, // Cold desert: Earth yellow
	0xff99ffff, // Frost desert: Ice blue
	0xfffffafa, // Tundra: Snow white
	0xff87ceeb, // Shallow water: Sky blue
	0xff4000ff, // Deep water: Ultramarine
}

var Biomes = [][]Biome{
	{FrostDesert, FrostDesert, Tundra, Tundra, Tundra},
	{ColdDesert, ColdDesert, BorealForest, BorealForest, BorealForest},
	{RockDesert, Grassland, Grassland, TemperateForest, TemperateForest},
	{TemperateDesert, TemperateDesert, Grassland, TemperateForest, TropicalForest},
	{HotDesert, HotDesert, Savanna, TropicalForest, TropicalForest},
}
