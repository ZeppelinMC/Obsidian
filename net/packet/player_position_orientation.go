package packet

type PlayerPositionOrientation struct {
	PlayerID   int8
	X, Y, Z    float32
	Yaw, Pitch byte
}

func (PlayerPositionOrientation) ID() byte {
	return 0x08
}

func (m *PlayerPositionOrientation) Decode(r Reader) {
	r.SByte(&m.PlayerID)
	r.FShort(&m.X)
	r.FShort(&m.Y)
	r.FShort(&m.Z)
	r.Byte(&m.Yaw)
	r.Byte(&m.Pitch)
}

func (m PlayerPositionOrientation) Encode(w Writer) {
	w.SByte(m.PlayerID)
	w.Short(int16(m.X))
	w.Short(int16(m.Y))
	w.Short(int16(m.Z))
	w.Byte(m.Yaw)
	w.Byte(m.Pitch)
}
