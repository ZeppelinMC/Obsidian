package format

import (
	"compress/gzip"
	"errors"
	"io"
	"math"
)

// https://github.com/UnknownShadow200/MCGalaxy/wiki/Level-format

var ErrBlocksCountMismatch = errors.New("block count must be width*length*height")

func WriteLevel(w io.Writer, l Level) error {
	gun := gzip.NewWriter(w)

	writeInt16LE(gun, l.Identifier)

	writeInt16LE(gun, l.Width)
	writeInt16LE(gun, l.Length)
	writeInt16LE(gun, l.Height)

	writeInt16LE(gun, l.SpawnX)
	writeInt16LE(gun, l.SpawnZ)
	writeInt16LE(gun, l.SpawnY)
	writeByte(gun, l.SpawnYaw)
	writeByte(gun, l.SpawnPitch)

	if l.Identifier == 1874 {
		writeByte(gun, l.MinAccess)
		writeByte(gun, l.MinBuild)
	}

	if len(l.Blocks) != int(l.Width)*int(l.Length)*int(l.Height) {
		return ErrBlocksCountMismatch
	}

	if l.CustomBlockIds != nil {
		writeByte(gun, 0xBD)

		for y := 0; y < int(math.Ceil(float64(l.Height/16))); y += 1 {
			for z := 0; z < int(math.Ceil(float64(l.Length/16))); z += 1 {
				for x := 0; x < int(math.Ceil(float64(l.Width/16))); x += 1 {
					data, ok := l.CustomBlockIds[[3]int{x, y, z}]
					if !ok {
						writeByte(gun, 0)
						continue
					}
					writeByte(gun, 1)

					gun.Write(data[:])
				}
			}
		}
	} else {
		writeByte(gun, 0)
	}

	if l.Physics != nil {
		writeByte(gun, 0xFC)
		writeInt32LE(gun, int32(len(l.Physics)))

		for _, e := range l.Physics {
			writeInt32LE(gun, e[0])
			writeInt32LE(gun, e[1])
		}
	}

	gun.Close()

	return nil
}

func writeInt16LE(w io.Writer, i int16) {
	w.Write([]byte{byte(i), byte(i >> 8)})
}

func writeInt32LE(w io.Writer, i int32) {
	w.Write([]byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)})
}

func writeByte(w io.Writer, i byte) {
	w.Write([]byte{i})
}
