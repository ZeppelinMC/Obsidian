package packet

type ExtInfo struct {
	AppName        string
	ExtensionCount int16
}

func (ExtInfo) ID() byte {
	return 0x10
}

func (s *ExtInfo) Decode(r Reader) {
	r.String(&s.AppName)
	r.Short(&s.ExtensionCount)
}

func (s ExtInfo) Encode(w Writer) {
	w.String(s.AppName)
	w.Short(s.ExtensionCount)
}
