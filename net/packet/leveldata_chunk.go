package packet

type LevelDataChunk struct {
	ChunkLength     int16
	ChunkData       []byte
	PercentComplete byte
}

func (LevelDataChunk) ID() byte {
	return 0x03
}

func (s *LevelDataChunk) Decode(r Reader) {
	r.Short(&s.ChunkLength)
	r.ByteArray(&s.ChunkData)
	r.Byte(&s.PercentComplete)
}

func (s LevelDataChunk) Encode(w Writer) {
	w.Short(s.ChunkLength)
	w.ByteArray(s.ChunkData)
	w.Byte(s.PercentComplete)
}
