package packet

type ServerIdentification struct {
	ProtocolVersion byte
	ServerName      string
	ServerMOTD      string
	UserType        byte
}

func (ServerIdentification) ID() byte {
	return 0x00
}

func (s *ServerIdentification) Decode(r Reader) {
	r.Byte(&s.ProtocolVersion)
	r.String(&s.ServerName)
	r.String(&s.ServerMOTD)
	r.Byte(&s.UserType)
}

func (s ServerIdentification) Encode(w Writer) {
	w.Byte(s.ProtocolVersion)
	w.String(s.ServerName)
	w.String(s.ServerMOTD)
	w.Byte(s.UserType)
}
