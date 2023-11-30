package packet

type DisconnectPlayer struct {
	Reason string
}

func (DisconnectPlayer) ID() byte {
	return 0x0E
}

func (m *DisconnectPlayer) Decode(r Reader) {
	r.String(&m.Reason)
}

func (m DisconnectPlayer) Encode(w Writer) {
	w.String(m.Reason)
}
