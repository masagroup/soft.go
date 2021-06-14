package ecore

import "io"

type BinaryCodec struct {
}

func (bc *BinaryCodec) NewEncoder(w io.Writer, options map[string]interface{}) EResourceEncoder {
	return NewBinaryEncoder(w, options)
}
func (bc *BinaryCodec) NewDecoder(r io.Reader, options map[string]interface{}) EResourceDecoder {
	return NewBinaryDecoder(r, options)
}
