package ecore

import "io"

type BinaryDecoder struct {
	r io.Reader
}

func NewBinaryDecoder(r io.Reader, options map[string]interface{}) *BinaryDecoder {
	return &BinaryDecoder{}
}

func (bd *BinaryDecoder) DecodeResource(resource EResource) {

}

func (bd *BinaryDecoder) DecodeObject(object *EObject, resource EResource) error {
	return nil
}
