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
	"bytes"
	"errors"
	"io"

	"github.com/ugorji/go/codec"
)

type BinaryDecoder struct {
	resource EResource
	r        io.Reader
	decoder  *codec.Decoder
	baseURI  *URI
}

func NewBinaryDecoder(resource EResource, r io.Reader, options map[string]interface{}) *BinaryDecoder {
	mh := &codec.MsgpackHandle{}
	d := &BinaryDecoder{
		resource: resource,
		r:        r,
		decoder:  codec.NewDecoder(r, mh),
	}
	if uri := resource.GetURI(); uri != nil && uri.IsAbsolute() {
		d.baseURI = uri
	}
	return d
}

func (d *BinaryDecoder) Decode() {
	defer func() {
		if err, _ := recover().(error); err != nil {
			resourcePath := ""
			if d.resource.GetURI() != nil {
				resourcePath = d.resource.GetURI().String()
			}
			d.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), resourcePath, 0, 0))
		}
	}()
	d.decodeSignature()
	d.decodeVersion()
}

func (d *BinaryDecoder) DecodeObject() (eObject EObject, err error) {
	defer func() {
		if panicErr, _ := recover().(error); panicErr != nil {
			err = panicErr
		}
	}()
	d.decodeSignature()
	d.decodeVersion()
	return
}

func (d *BinaryDecoder) decode(v interface{}) {
	if err := d.decoder.Decode(v); err != nil {
		panic(err)
	}
}

func (d *BinaryDecoder) decodeSignature() {
	signature := make([]byte, 8)
	d.decode(signature)
	if bytes.Compare(signature, binarySignature) != 0 {
		panic(errors.New("Invalid signature for a binary EMF serialization"))
	}
}

func (d *BinaryDecoder) decodeVersion() {
	var version int
	d.decode(&version)
	if version != binaryVersion {
		panic(errors.New("Invalid version for binary EMF serialization"))
	}
}
