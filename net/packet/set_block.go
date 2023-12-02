package packet

type SetBlock struct {
	X, Y, Z   int16
	BlockType byte
}

func (SetBlock) ID() byte {
	return 0x06
}

func (m *SetBlock) Decode(r Reader) {
	r.Short(&m.X)
	r.Short(&m.Y)
	r.Short(&m.Z)
	r.Byte(&m.BlockType)
}

func (m SetBlock) Encode(w Writer) {
	w.Short(m.X)
	w.Short(m.Y)
	w.Short(m.Z)
	w.Byte(m.BlockType)
}
