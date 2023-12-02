package packet

type PlayerPositionOrientation struct {
	PlayerID   int8
	X, Y, Z    int16 //float32
	Yaw, Pitch byte
}

func (PlayerPositionOrientation) ID() byte {
	return 0x08
}

func (m *PlayerPositionOrientation) Decode(r Reader) {
	r.SByte(&m.PlayerID)
	r.Short(&m.X)
	r.Short(&m.Y)
	r.Short(&m.Z)
	r.Byte(&m.Yaw)
	r.Byte(&m.Pitch)
}

func (m PlayerPositionOrientation) Encode(w Writer) {
	w.SByte(m.PlayerID)
	w.Short(m.X)
	w.Short(m.Y)
	w.Short(m.Z)
	w.Byte(m.Yaw)
	w.Byte(m.Pitch)
}
