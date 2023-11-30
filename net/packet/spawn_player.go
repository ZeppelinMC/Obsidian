package packet

type SpawnPlayer struct {
	PlayerID   int8
	PlayerName string
	X, Y, Z    float32
	Yaw, Pitch byte
}

func (SpawnPlayer) ID() byte {
	return 0x07
}

func (m *SpawnPlayer) Decode(r Reader) {
	r.SByte(&m.PlayerID)
	r.String(&m.PlayerName)
	r.FShort(&m.X)
	r.FShort(&m.Y)
	r.FShort(&m.Z)
	r.Byte(&m.Yaw)
	r.Byte(&m.Pitch)
}

func (m SpawnPlayer) Encode(w Writer) {
	w.SByte(m.PlayerID)
	w.String(m.PlayerName)
	w.Short(int16(m.X))
	w.Short(int16(m.Y))
	w.Short(int16(m.Z))
	w.Byte(m.Yaw)
	w.Byte(m.Pitch)
}
