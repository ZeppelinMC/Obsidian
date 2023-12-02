package packet

type SpawnPlayer struct {
	PlayerID   int8
	PlayerName string
	X, Y, Z    int16
	Yaw, Pitch byte
}

func (SpawnPlayer) ID() byte {
	return 0x07
}

func (m *SpawnPlayer) Decode(r Reader) {
	r.SByte(&m.PlayerID)
	r.String(&m.PlayerName)
	r.Short(&m.X)
	r.Short(&m.Y)
	r.Short(&m.Z)
	r.Byte(&m.Yaw)
	r.Byte(&m.Pitch)
}

func (m SpawnPlayer) Encode(w Writer) {
	w.SByte(m.PlayerID)
	w.String(m.PlayerName)
	w.Short(m.X)
	w.Short(m.Y)
	w.Short(m.Z)
	w.Byte(m.Yaw)
	w.Byte(m.Pitch)
}
