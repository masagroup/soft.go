package ecore

import "io"

type BinaryEncoder struct {
}

func NewBinaryEncoder(options map[string]interface{}) *BinaryEncoder {
	return &BinaryEncoder{}
}

func (be *BinaryEncoder) Encode(resource EResource, w io.Writer) {

}
