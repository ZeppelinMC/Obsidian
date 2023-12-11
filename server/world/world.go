package world

import (
	"compress/gzip"
	"fmt"
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

type World struct {
	Data WorldData
}

func (w *World) SetBlock(x, y, z int16, blockType byte) {
	w.Data.BlockArray[w.GetIndex(x, y, z)] = int8(blockType)
}

func (w *World) GetIndex(x, y, z int16) int {
	return int(x) + int(w.Data.X)*(int(z)+int(w.Data.Z)*int(y))
}

func (w *World) XYZ(index int) (int16, int16, int16) {
	x := index % int(w.Data.X)
	y := index / (int(w.Data.X) * int(w.Data.Z))
	z := (index / int(w.Data.X)) % int(w.Data.Y)
	return int16(x), int16(y), int16(z)
}

func LoadWorld() *World {
	d1, err := os.Open("world/main.cw")
	if err != nil {
		return &World{DefaultWorldData}
	}

	dat, err := gzip.NewReader(d1)
	if err != nil {
		return &World{DefaultWorldData}
	}

	var d WorldData

	if err := nbt.NewDecoder(dat).Decode(&d); err != nil {
		fmt.Println(err)
		d = DefaultWorldData
	}

	return &World{d}
}

func (w *World) Save() {
	os.Mkdir("world", 0755)
	file, _ := os.Create("world/main.cw")
	g := gzip.NewWriter(file)

	nbt.NewEncoder(g).Encode(w.Data)

	g.Close()
	file.Close()
}
