package ecore

import (
	"io"

	"github.com/ugorji/go/codec"
)

const (
	checkNothing = iota
	checkDirectResource
	checkResource
	chechContainer
)

type BinaryEncoder struct {
	w        io.Writer
	resource EResource
	baseURI  *URI
	encoder  *codec.Encoder
	version  int
}

func NewBinaryEncoder(w io.Writer, options map[string]interface{}) *BinaryEncoder {
	return &BinaryEncoder{w: w}
}

func (be *BinaryEncoder) Encode(resource EResource) {
	be.resource = resource
	be.encoder = codec.NewEncoder(be.w, &codec.MsgpackHandle{})
	be.encodeSignature()
	be.encodeVersion()
}

func (be *BinaryEncoder) encodeSignature() {
	// Write a signature that will be obviously corrupt
	// if the binary contents end up being UTF-8 encoded
	// or altered by line feed or carriage return changes.
	be.encoder.Encode([]byte{'\211', 'e', 'm', 'f', '\n', '\r', '\032', '\n'})
}

func (be *BinaryEncoder) encodeVersion() {
	be.encoder.Encode(be.version)
}
