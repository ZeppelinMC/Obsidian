package world

import (
	"io"
	"obsidian/server/world/format"
	"os"
	"unsafe"
)

type World struct {
	path   string
	reader string
	Data   format.WorldData
	ogdata any
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

func LoadWorld(path, typ string) *World {
	file, err := os.Open(path)
	if err != nil {
		return &World{path: path, Data: format.DefaultWorldData, reader: typ}
	}

	var data format.WorldData
	var ogdata any

	switch typ {
	case "level":
		l, err := format.ReadLevel(file)
		if err != nil {
			data = format.DefaultWorldData
			break
		}
		ogdata = l
		data = l.ToWorldData()
	default:
		d, err := format.ReadClassicWorld(file)
		if err != nil && err != io.EOF {
			data = format.DefaultWorldData
			break
		}
		data = d
	}

	file.Close()
	return &World{path: path, Data: data, reader: typ, ogdata: ogdata}
}

func (w *World) Save() {
	file, _ := os.Create(w.path)
	switch w.reader {
	case "level":
		if w.ogdata == nil {
			w.ogdata = format.Level{
				Identifier: 252,
				Width:      w.Data.X,
				Height:     w.Data.Y,
				Length:     w.Data.Z,
				SpawnX:     w.Data.Spawn.X,
				SpawnY:     w.Data.Spawn.Y,
				SpawnZ:     w.Data.Spawn.Z,
				SpawnYaw:   byte(w.Data.Spawn.H),
				SpawnPitch: byte(w.Data.Spawn.P),

				Blocks: *(*[]byte)(unsafe.Pointer(&w.Data.BlockArray)),
			}
		}
		format.WriteLevel(file, w.ogdata.(format.Level))
	default:
		format.WriteClassicWorld(file, w.Data)
	}
	file.Close()
}
