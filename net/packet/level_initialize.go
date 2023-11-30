package packet

type LevelInitialize struct{}

func (LevelInitialize) ID() byte {
	return 0x02
}

func (LevelInitialize) Decode(Reader) {}

func (LevelInitialize) Encode(Writer) {}
