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
)

const (
	JSON_OPTION_ID_ATTRIBUTE = "ID_ATTRIBUTE" // if true, save id attribute of the object
)

type JSONCodec struct {
}

func (jc *JSONCodec) NewEncoder(resource EResource, w io.Writer, options map[string]interface{}) EResourceEncoder {
	return nil
}
func (jc *JSONCodec) NewDecoder(resource EResource, r io.Reader, options map[string]interface{}) EResourceDecoder {
	return nil
}
