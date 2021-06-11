package ecore

import "io"

type EResourceEncoder interface {
	Encode(resource EResource, w io.Writer)
}
