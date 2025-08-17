package trgm

type block struct {
	offset  uint32
	tri     uint32
	entries set
}

func (b *block) encode() []byte {
	return nil
}

func (b *block) decode([]byte) {

}
