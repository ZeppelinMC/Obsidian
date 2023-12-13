package format

type worldDataCreatedBy struct {
	Service, Username string
}

type worldDataMapGenerator struct {
	Software, MapGeneratorName string
}

type worldDataSpawn struct {
	X, Y, Z int16
	H, P    int8
}

type WorldData struct {
	FormatVersion                           int8
	Name                                    string
	UUID                                    []float64
	X, Y, Z                                 int16
	CreatedBy                               worldDataCreatedBy
	MapGenerator                            worldDataMapGenerator
	TimeCreated, LastAccessed, LastModified int64
	Spawn                                   worldDataSpawn
	BlockArray                              []int8
}

var DefaultWorldData = WorldData{
	FormatVersion: 1,
	Name:          "ObsidianWorld",
	X:             512,
	Y:             256,
	Z:             512,
	MapGenerator: worldDataMapGenerator{
		Software:         "Obsidian",
		MapGeneratorName: "Default",
	},
	Spawn: worldDataSpawn{
		X: 150,
		Y: 50,
		Z: 150,
	},
}
