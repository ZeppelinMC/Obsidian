package EnvMapAspect

import "obsidian/net/packet"

type SetMapEnvUrl struct {
	TexturePackURL string
}

func (SetMapEnvUrl) ID() byte {
	return 0x28
}

func (e *SetMapEnvUrl) Decode(r packet.Reader) {
	r.String(&e.TexturePackURL)
}

func (e SetMapEnvUrl) Encode(w packet.Writer) {
	w.String(e.TexturePackURL)
}
