package CustomBlocks

import "obsidian/net/packet"

type CustomBlockSupportLevel struct {
	SupportLevel byte
}

func (CustomBlockSupportLevel) ID() byte {
	return 0x13
}

func (c *CustomBlockSupportLevel) Decode(r packet.Reader) {
	r.Byte(&c.SupportLevel)
}

func (c CustomBlockSupportLevel) Encode(w packet.Writer) {
	w.Byte(c.SupportLevel)
}
