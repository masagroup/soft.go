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

type NoCodec struct {
}

func (nc *NoCodec) NewEncoder(resource EResource, w io.Writer, options map[string]any) EEncoder {
	return &NoEncoder{}
}

func (nc *NoCodec) NewDecoder(resource EResource, r io.Reader, options map[string]any) EResourceDecoder {
	return &NoDecoder{}
}
