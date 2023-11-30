package packet

type SetBlockServer struct {
	X, Y, Z         int16
	Mode, BlockType byte
}

func (SetBlockServer) ID() byte {
	return 0x05
}

func (m *SetBlockServer) Decode(r Reader) {
	r.Short(&m.X)
	r.Short(&m.Y)
	r.Short(&m.Z)
	r.Byte(&m.Mode)
	r.Byte(&m.BlockType)
}

func (m SetBlockServer) Encode(w Writer) {
	w.Short(m.X)
	w.Short(m.Y)
	w.Short(m.Z)
	w.Byte(m.Mode)
	w.Byte(m.BlockType)
}
