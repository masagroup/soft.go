package ecore

import "io"

type EResourceDecoder interface {
	Decode(resource EResource, r io.Reader)
}
