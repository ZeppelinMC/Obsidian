package packet

type LevelFinalize struct {
	XSize, YSize, ZSize int16
}

func (LevelFinalize) ID() byte {
	return 0x04
}

func (l *LevelFinalize) Decode(r Reader) {
	r.Short(&l.XSize)
	r.Short(&l.YSize)
	r.Short(&l.ZSize)
}

func (l LevelFinalize) Encode(w Writer) {
	w.Short(l.XSize)
	w.Short(l.YSize)
	w.Short(l.ZSize)
}
