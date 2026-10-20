// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmlcodec

import (
	"github.com/masagroup/soft.go/ecore"
	"io"
)

type XMICodec struct {
}

func (d *XMICodec) NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) ecore.EEncoder {
	return NewXMIEncoder(resource, w, options)
}
func (d *XMICodec) NewDecoder(resource ecore.EResource, r io.Reader, options map[string]any) ecore.EDecoder {
	return NewXMIDecoder(resource, r, options)
}
