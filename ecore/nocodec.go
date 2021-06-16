package ecore

import "io"

type NoCodec struct {
}

func (nc *NoCodec) NewEncoder(resource EResource, w io.Writer, options map[string]interface{}) EResourceEncoder {
	return &NoEncoder{}
}

func (nc *NoCodec) NewDecoder(resource EResource, r io.Reader, options map[string]interface{}) EResourceDecoder {
	return &NoDecoder{}
}
