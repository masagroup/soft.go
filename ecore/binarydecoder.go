package ecore

import "io"

type BinaryDecoder struct {
}

func NewBinaryDecoder(options map[string]interface{}) *BinaryDecoder {
	return &BinaryDecoder{}
}

func (bd *BinaryDecoder) Decode(resource EResource, r io.Reader) {

}
