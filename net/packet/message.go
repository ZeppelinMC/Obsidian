package packet

type Message struct {
	PlayerID int8
	Message  string
}

func (Message) ID() byte {
	return 0x0D
}

func (m *Message) Decode(r Reader) {
	r.SByte(&m.PlayerID)
	r.String(&m.Message)
}

func (m Message) Encode(w Writer) {
	w.SByte(m.PlayerID)
	w.String(m.Message)
}
