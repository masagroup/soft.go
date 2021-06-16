// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "io"

type BinaryDecoder struct {
	r io.Reader
}

func NewBinaryDecoder(resource EResource, r io.Reader, options map[string]interface{}) *BinaryDecoder {
	return &BinaryDecoder{}
}

func (bd *BinaryDecoder) Decode() {

}

func (bd *BinaryDecoder) DecodeObject() (EObject, error) {
	return nil, nil
}
