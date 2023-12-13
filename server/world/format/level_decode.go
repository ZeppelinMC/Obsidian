package format

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"unsafe"
)

// https://github.com/UnknownShadow200/MCGalaxy/wiki/Level-format

type Level struct {
	Identifier int16
	Width, Length, Height,
	SpawnX, SpawnY, SpawnZ int16
	SpawnYaw, SpawnPitch byte
	MinAccess, MinBuild  byte

	Blocks         []byte
	CustomBlockIds map[[3]int][4096]byte

	Physics [][2]int32
}

func (l *Level) FindCustomIds() {
	for i, b := range l.Blocks {
		if b != 163 {
			continue
		}

		fmt.Println("dam!", i)
	}
}

func (l Level) ToWorldData() WorldData {
	return WorldData{
		X:          l.Width,
		Y:          l.Height,
		Z:          l.Length,
		BlockArray: *(*[]int8)(unsafe.Pointer(&l.Blocks)),

		Spawn: worldDataSpawn{
			X: l.SpawnX,
			Y: l.SpawnY,
			Z: l.SpawnZ,
			H: int8(l.SpawnYaw),
			P: int8(l.SpawnPitch),
		},
	}
}

func ReadLevel(r io.Reader) (Level, error) {
	var l = Level{}

	r, err := gzip.NewReader(r)
	d, _ := io.ReadAll(r)
	r = bytes.NewReader(d)
	if err != nil {
		return l, err
	}
	if err = readInt16LE(r, &l.Identifier); err != nil {
		return l, err
	}
	if err = readInt16LE(r, &l.Width); err != nil {
		return l, err
	}
	if err = readInt16LE(r, &l.Length); err != nil {
		return l, err
	}
	if err = readInt16LE(r, &l.Height); err != nil {
		return l, err
	}

	if err = readInt16LE(r, &l.SpawnX); err != nil {
		return l, err
	}
	if err = readInt16LE(r, &l.SpawnZ); err != nil {
		return l, err
	}
	if err = readInt16LE(r, &l.SpawnY); err != nil {
		return l, err
	}
	if err = readByte(r, &l.SpawnYaw); err != nil {
		return l, err
	}
	if err = readByte(r, &l.SpawnPitch); err != nil {
		return l, err
	}

	if l.Identifier == 1874 {
		if err = readByte(r, &l.MinAccess); err != nil {
			return l, err
		}
		if err = readByte(r, &l.MinBuild); err != nil {
			return l, err
		}
	}

	l.Blocks = make([]byte, int(l.Width)*int(l.Height)*int(l.Length))

	if _, err := r.Read(l.Blocks); err != nil {
		return l, err
	}

	var next byte
	readByte(r, &next)

	if next == 0xBD {
		l.CustomBlockIds = make(map[[3]int][4096]byte)
		for y := 0; y < int(math.Ceil(float64(l.Height/16))); y += 1 {
			for z := 0; z < int(math.Ceil(float64(l.Length/16))); z += 1 {
				for x := 0; x < int(math.Ceil(float64(l.Width/16))); x += 1 {
					var hasData byte
					readByte(r, &hasData)

					if hasData != 1 {
						continue
					}

					d := [4096]byte{}

					r.Read(d[:])

					l.CustomBlockIds[[3]int{x, y, z}] = d
				}
			}
		}
	}

	readLevelMetadata(r, &l)

	l.FindCustomIds()

	return l, nil
}

func readLevelMetadata(r io.Reader, l *Level) {
	var id byte

	readByte(r, &id)

	switch id {
	//physics
	case 0xFC:
		var num int32
		readInt32LE(r, &num)

		l.Physics = make([][2]int32, num)

		for i := int32(0); i < num; i++ {
			readInt32LE(r, &l.Physics[i][0])
			readInt32LE(r, &l.Physics[i][1])
		}
	}
}

func readInt16LE(r io.Reader, i *int16) error {
	var d [2]byte

	_, err := r.Read(d[:])

	if err != nil {
		return err
	}

	*i = int16(d[0]) | int16(d[1])<<8
	return nil
}

func readInt32LE(r io.Reader, i *int32) error {
	var d [4]byte

	_, err := r.Read(d[:])

	if err != nil {
		return err
	}

	*i = int32(d[0]) | int32(d[1])<<8 | int32(d[2])<<16 | int32(d[3])<<24
	return nil
}

func readByte(r io.Reader, i *byte) error {
	var d [1]byte
	_, err := r.Read(d[:])

	if err != nil {
		return err
	}

	*i = d[0]
	return nil
}
