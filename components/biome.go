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
	TropicalDry
	TropicalWet
	Grassland
	TemperateForest
	TemperateWet
	TemperateDesert
	RockDesert
	BorealDry
	BorealWet
	ColdDesert
	FrostDesert
	Tundra
	ShallowWater
	DeepWater
)

var BiomeColors = []uint32{
	0xffe34234, // Hot desert: Vermillon
	0xffe4d96f, // Savanna: Straw
	0xff0bda51, // Tropical dry: Malachite
	0xff008000, // Tropical wet: Office green
	0xff7cfc00, // Grassland: Lawn green
	0xff228b22, // Temperate forest: Forest green
	0xff00ff7f, // Temperate wet: Spring green
	0xfffdee00, // Temperate desert: Cobalt yellow
	0xffc0c0c0, // Rock desert: Silver
	0xff01796f, // Boreal dry: Pine green
	0xff29ab87, // Boreal wet: Jungle green
	0xfff0ffff, // Cold desert: Azure white
	0xff99ffff, // Frost desert: Ice blue
	0xfffffafa, // Tundra: Snow white
	0xff00ffff, // Shallow water: Cyan
	0xff009dc4, // Deep water: Pacific blue
}

var Biomes = [][]Biome{
	{FrostDesert, FrostDesert, Tundra, Tundra, Tundra},
	{ColdDesert, BorealDry, BorealDry, BorealWet, BorealWet},
	{RockDesert, Grassland, TemperateForest, TemperateForest, TemperateWet},
	{TemperateDesert, Grassland, TemperateForest, TemperateForest, TemperateWet},
	{HotDesert, HotDesert, Savanna, TropicalDry, TropicalWet},
}
