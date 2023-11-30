package packet

type PlayerIdentification struct {
	ProtocolVersion byte
	Username        string
	VerificationKey string
}

func (PlayerIdentification) ID() byte {
	return 0x00
}

func (p *PlayerIdentification) Decode(r Reader) {
	r.Byte(&p.ProtocolVersion)
	r.String(&p.Username)
	r.String(&p.VerificationKey)
	r.readBytes(1)
}

func (p PlayerIdentification) Encode(w Writer) {
	w.Byte(p.ProtocolVersion)
	w.String(p.Username)
	w.String(p.VerificationKey)
	w.Byte(0)
}
