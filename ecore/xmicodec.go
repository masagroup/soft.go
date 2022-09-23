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

type XMICodec struct {
}

func (d *XMICodec) NewEncoder(resource EResource, w io.Writer, options map[string]any) EResourceEncoder {
	return NewXMIEncoder(resource, w, options)
}
func (d *XMICodec) NewDecoder(resource EResource, r io.Reader, options map[string]any) EResourceDecoder {
	return NewXMIDecoder(resource, r, options)
}
