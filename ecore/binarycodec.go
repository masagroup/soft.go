package ecore

import "io"

type BinaryCodec struct {
}

func (bc *BinaryCodec) NewEncoder(resource EResource, w io.Writer, options map[string]interface{}) EResourceEncoder {
	return NewBinaryEncoder(resource, w, options)
}
func (bc *BinaryCodec) NewDecoder(resource EResource, r io.Reader, options map[string]interface{}) EResourceDecoder {
	return NewBinaryDecoder(resource, r, options)
}
