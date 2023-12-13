package format

import (
	"compress/gzip"
	"io"

	"github.com/aimjel/minecraft/nbt"
)

func ReadClassicWorld(r io.Reader) (WorldData, error) {
	var d WorldData
	r, err := gzip.NewReader(r)
	if err != nil {
		return d, err
	}

	err = nbt.NewDecoder(r).Decode(&d)

	return d, err
}

func WriteClassicWorld(w io.Writer, d WorldData) error {
	gun := gzip.NewWriter(w)

	err := nbt.NewEncoder(gun).Encode(&d)

	if err != nil {
		return err
	}

	return gun.Close()
}
