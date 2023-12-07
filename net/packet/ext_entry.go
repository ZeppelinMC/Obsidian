package packet

type ExtEntry struct {
	ExtName string
	Version int32
}

func (ExtEntry) ID() byte {
	return 0x11
}

func (s *ExtEntry) Decode(r Reader) {
	r.String(&s.ExtName)
	r.Int(&s.Version)
}

func (s ExtEntry) Encode(w Writer) {
	w.String(s.ExtName)
	w.Int(s.Version)
}
