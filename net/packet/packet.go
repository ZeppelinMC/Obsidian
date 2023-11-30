package packet

type Packet interface {
	ID() byte
	Decode(Reader)
	Encode(Writer)
}

func Marshal(id byte, fields ...any) Packet {
	return packet{id, fields}
}

type packet struct {
	id     byte
	fields []any
}

func (p packet) ID() byte {
	return p.id
}

func (p packet) Decode(r Reader) {
	for _, f := range p.fields {
		switch field := f.(type) {
		case uint8:
			r.Byte(&field)
		case int8:
			r.SByte(&field)
		case int16:
			r.Short(&field)
		case []byte:
			r.ByteArray(&field)
		}
	}
}

func (p packet) Encode(w Writer) {
	for _, f := range p.fields {
		switch field := f.(type) {
		case uint8:
			w.Byte(field)
		case int8:
			w.SByte(field)
		case int16:
			w.Short(field)
		case []byte:
			w.ByteArray(field)
		}
	}
}
