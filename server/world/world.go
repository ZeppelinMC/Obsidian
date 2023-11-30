package world

import (
	"compress/gzip"
	"os"

	"github.com/aimjel/minecraft/nbt"
)

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

type World struct {
	Data WorldData
}

func LoadWorld() *World {
	d1, _ := os.Open("world/main.cw")

	dat, _ := gzip.NewReader(d1)

	var d WorldData

	nbt.NewDecoder(dat).Decode(&d)

	return &World{d}
}
