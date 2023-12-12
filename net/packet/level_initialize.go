package packet

type LevelInitialize struct {
	FastMap bool
	MapSize int32
}

func (LevelInitialize) ID() byte {
	return 0x02
}

func (l *LevelInitialize) Decode(r Reader) {
	if l.FastMap {
		r.Int(&l.MapSize)
	}
}

func (l LevelInitialize) Encode(w Writer) {
	if l.FastMap {
		w.Int(l.MapSize)
	}
}
