package ecore

import "io"

type SQLEncoder struct {
}

func NewSQLEncoder(resource EResource, w io.Writer, options map[string]any) *SQLEncoder {
	return &SQLEncoder{}
}

func (e *SQLEncoder) EncodeResource() {

}

func (e *SQLEncoder) EncodeObject(object EObject) error {
	return nil
}
