package world

import (
	"errors"
	"io"
	"obsidian/server/world/format"
	"os"
	"unsafe"
)

var ErrInvalidIndex = errors.New("invalid index")

type World struct {
	path   string
	reader string
	Data   format.WorldData
	ogdata any
}

func (w *World) SetBlock(x, y, z int16, blockType byte) error {
	i := w.GetIndex(x, y, z)

	if i > len(w.Data.BlockArray) || i < 0 {
		return ErrInvalidIndex
	}
	w.Data.BlockArray[i] = int8(blockType)
	return nil
}

// GetIndex returns the index for the x, y, z in the block array
func (w *World) GetIndex(x, y, z int16) int {
	return int(x) + int(w.Data.X)*(int(z)+int(w.Data.Z)*int(y))
}

// XYZ returns the x, y, z for the index
func (w *World) XYZ(index int) (int16, int16, int16) {
	x := index % int(w.Data.X)
	y := index / (int(w.Data.X) * int(w.Data.Z))
	z := (index / int(w.Data.X)) % int(w.Data.Y)
	return int16(x), int16(y), int16(z)
}

// LoadWorld loads a world in the path using the specified reader. Defaults to classicworld.
// If the world failed to load it will use the default world data and the map will then be regenerated
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

// World.Save saves the world. If it's a generated world using the level reader, the identifier will be 252
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
