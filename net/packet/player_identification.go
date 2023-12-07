package packet

type PlayerIdentification struct {
	ProtocolVersion byte
	Username        string
	VerificationKey string
	//https://wiki.vg/Classic_Protocol_Extension
	CPE bool
}

func (PlayerIdentification) ID() byte {
	return 0x00
}

func (p *PlayerIdentification) Decode(r Reader) {
	r.Byte(&p.ProtocolVersion)
	r.String(&p.Username)
	r.String(&p.VerificationKey)
	p.CPE = r.readBytes(1)[0] == 0x42
}

func (p PlayerIdentification) Encode(w Writer) {
	w.Byte(p.ProtocolVersion)
	w.String(p.Username)
	w.String(p.VerificationKey)
	cpe := byte(0x00)
	if p.CPE {
		cpe = 0x42
	}
	w.Byte(cpe)
}
