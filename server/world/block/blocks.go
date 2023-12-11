package block

const (
	Air = iota
	Stone
	Grass
	Dirt
	Cobblestone
	Planks
	Sapling
	Bedrock
	FlowingWater
	Water
	FlowingLava
	Lava
	Sand
	Gravel
	GoldOre
	IronOre
	CoalOre
	Wood
	Leaves
	Sponge
	Glass

	Red
	Orange
	Yellow
	Lime
	Green
	Teal
	Aqua
	Cyan
	Blue
	Indigo
	Violet
	Magenta
	Pink
	Black
	Gray
	White

	Dandelion
	Rose
	BrownMushroom
	RedMushroom

	Gold
	Iron
	DoubleSlab
	Slab
	Bricks
	TNT
	Bookshelf
	Moss
	Obsidian
)

// Custom Blocks
const (
	CobblestoneSlab = iota + 50
	Rope
	Sandstone
	Snow
	Fire
	LightPinkWool
	ForestGreenWool
	BrownWool
	DeepBlue
	Turquoise
	Ice
	CeramicTile
	Magma
	Pillar
	Crate
	StoneBrick
)

var CustomBlockFallBack = map[byte]byte{
	CobblestoneSlab: Slab,
	Rope:            BrownMushroom,
	Sandstone:       Sand,
	Snow:            Air,
	Fire:            Lava,
	LightPinkWool:   Pink,
	ForestGreenWool: Green,
	BrownWool:       Dirt,
	DeepBlue:        Blue,
	Turquoise:       Cyan,
	Ice:             Glass,
	CeramicTile:     Iron,
	Magma:           Obsidian,
	Pillar:          White,
	Crate:           Planks,
	StoneBrick:      Stone,
}
