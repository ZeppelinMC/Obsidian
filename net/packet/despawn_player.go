package packet

type DespawnPlayer struct {
	PlayerID int8
}

func (DespawnPlayer) ID() byte {
	return 0x0C
}

func (m *DespawnPlayer) Decode(r Reader) {
	r.SByte(&m.PlayerID)
}

func (m DespawnPlayer) Encode(w Writer) {
	w.SByte(m.PlayerID)
}
