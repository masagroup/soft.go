package ecore

import "io"

type BinaryDecoder struct {
	r io.Reader
}

func NewBinaryDecoder(resource EResource, r io.Reader, options map[string]interface{}) *BinaryDecoder {
	return &BinaryDecoder{}
}

func (bd *BinaryDecoder) Decode() {

}

func (bd *BinaryDecoder) DecodeObject() (EObject, error) {
	return nil, nil
}
