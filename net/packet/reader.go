package packet

import (
	"io"
	"strings"
)

var serverBoundPool = map[byte]func() Packet{
	0x00: func() Packet { return &PlayerIdentification{} },
	0x05: func() Packet { return &SetBlockServer{} },
	0x08: func() Packet { return &PlayerPositionOrientation{} },
	0x0d: func() Packet { return &Message{} },
	0x10: func() Packet { return &ExtInfo{} },
	0x11: func() Packet { return &ExtEntry{} },
}

func ReadPacket(r io.Reader) Packet {
	var id [1]byte

	_, err := r.Read(id[:])

	if err != nil {
		return nil
	}

	p, ok := serverBoundPool[id[0]]
	if !ok {
		io.ReadAll(r)
		return nil
	}

	pk := p()

	pk.Decode(Reader{r})

	return pk
}

type Reader struct {
	buf io.Reader
}

func (r Reader) readBytes(l int) []byte {
	var i = make([]byte, l)

	r.buf.Read(i)

	return i
}

func (r Reader) Byte(i *uint8) {
	*i = r.readBytes(1)[0]
}

func (r Reader) SByte(i *int8) {
	*i = int8(r.readBytes(1)[0])
}

func (r Reader) Short(i *int16) {
	d := r.readBytes(2)

	*i = int16(d[0])<<8 | int16(d[1])
}

func (r Reader) FShort(i *float32) {
	d := r.readBytes(2)

	x := int16(d[0])<<8 | int16(d[1])

	*i = float32(x / 32)
}

func (r Reader) Int(i *int32) {
	d := r.readBytes(4)
	*i = int32(d[0])<<24 | int32(d[0])<<16 | int32(d[0])<<8 | int32(d[0])
}

func (r Reader) String(s *string) {
	*s = strings.TrimSpace(string(r.readBytes(64)))
}

func (r Reader) ByteArray(b *[]byte) {
	*b = r.readBytes(1024)
}
