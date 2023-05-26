package ecore

import "io"

const (
	SQL_OPTION_DRIVER            = "DRIVER_NAME"       // value of the sql driver
	SQL_OPTION_ID_ATTRIBUTE_NAME = "ID_ATTRIBUTE_NAME" // value of the id attribute
)

type SQLCodec struct {
}

func (d *SQLCodec) NewEncoder(resource EResource, w io.Writer, options map[string]any) EEncoder {
	return NewSQLEncoder(resource, w, options)
}
func (d *SQLCodec) NewDecoder(resource EResource, r io.Reader, options map[string]any) EDecoder {
	return NewSQLDecoder(resource, r, options)
}
