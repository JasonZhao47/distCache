package distCache

// ByteView immutable byte slice
type ByteView interface {
	Len() int
	AsSlice() []byte
}

// slice value is not reliable due to easy modifications
type byteView struct {
	b []byte
}

func NewByteView(src []byte) ByteView {
	dst := make([]byte, len(src))
	copy(dst, src)
	return byteView{b: dst}
}

// Len returns a slice of byteView
func (bv byteView) Len() int {
	return len(bv.b)
}

// AsSlice returns a copy of slice for byteView
func (bv byteView) AsSlice() []byte {
	res := make([]byte, len(bv.b))
	copy(res, bv.b)
	return res
}
