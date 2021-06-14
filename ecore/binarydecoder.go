package ecore

import "io"

type BinaryDecoder struct {
	r io.Reader
}

func NewBinaryDecoder(r io.Reader, options map[string]interface{}) *BinaryDecoder {
	return &BinaryDecoder{}
}

func (bd *BinaryDecoder) Decode(resource EResource) {

}
