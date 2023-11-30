package packet

import (
	"io"
)

func WritePacket(w io.Writer, pk Packet) {
	w.Write([]byte{pk.ID()})
	pk.Encode(Writer{w})
}

type Writer struct {
	buf io.Writer
}

func (w Writer) Byte(i uint8) {
	w.buf.Write([]byte{i})
}

func (w Writer) SByte(i int8) {
	w.buf.Write([]byte{byte(i)})
}

func (w Writer) Short(i int16) {
	w.buf.Write([]byte{byte(i >> 8), byte(i)})
}

func (w Writer) String(s string) {
	b := []byte(s)
	for len(b) < 64 {
		b = append(b, 0x20)
	}
	w.buf.Write(b)
}

func (w Writer) ByteArray(b []byte) {
	if len(b) < 1024 {
		b = append(b, make([]byte, 1024-len(b))...)
	}
	w.buf.Write(b)
}
