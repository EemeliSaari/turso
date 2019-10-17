package buffer

// Data ...
type Data interface {
	Hash() [16]byte
}

// Buffer ...
type Buffer struct {
	signatures map[[16]byte]struct{}
}

// NewBuffer ...
func NewBuffer() *Buffer {
	set := make(map[[16]byte]struct{})
	b := Buffer{signatures: set}
	return &b
}

// Add ...
func (b *Buffer) Add(content Data) bool {
	var exists struct{}
	hash := content.Hash()

	new := true
	if _, ok := b.signatures[hash]; ok {
		new = false
	} else {
		b.signatures[hash] = exists
	}
	return new
}

// Length ...
func (b *Buffer) Length() int {
	return len(b.signatures)
}
