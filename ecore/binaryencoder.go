// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

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

func NewBinaryEncoder(resource EResource, w io.Writer, options map[string]interface{}) *BinaryEncoder {
	return &BinaryEncoder{w: w, resource: resource}
}

func (be *BinaryEncoder) Encode() {
	be.encoder = codec.NewEncoder(be.w, &codec.MsgpackHandle{})
	be.encodeSignature()
	be.encodeVersion()
}

func (be *BinaryEncoder) EncodeObject(object EObject) error {
	return nil
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
