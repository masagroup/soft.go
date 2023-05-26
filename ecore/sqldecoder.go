package ecore

import "io"

type SQLDecoder struct {
}

func NewSQLDecoder(resource EResource, r io.Reader, options map[string]any) *SQLDecoder {
	return &SQLDecoder{}
}

func (d *SQLDecoder) DecodeResource() {

}

func (d *SQLDecoder) DecodeObject() (EObject, error) {
	return nil, nil
}
