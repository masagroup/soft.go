// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xml

import (
	"github.com/masagroup/soft.go/ecore"
	"io"
)

const (
	OPTION_EXTENDED_META_DATA            = "EXTENDED_META_DATA"            // ExtendedMetaData pointer
	OPTION_SUPPRESS_DOCUMENT_ROOT        = "SUPPRESS_DOCUMENT_ROOT"        // if true , suppress document root if found
	OPTION_DEFERRED_REFERENCE_RESOLUTION = "DEFERRED_REFERENCE_RESOLUTION" // if true , defer id ref resolution
	OPTION_DEFERRED_ROOT_ATTACHMENT      = "DEFERRED_ROOT_ATTACHMENT"      // if true , defer id ref resolution
	OPTION_ID_ATTRIBUTE_NAME             = "ID_ATTRIBUTE_NAME"             // value of the id attribute
	OPTION_ROOT_OBJECTS                  = "ROOT_OBJECTS"                  // list of root objects to save
	OPTION_USE_ENCODED_ATTRIBUTE_STYLE   = "USE_ENCODED_ATTRIBUTE_STYLE"   // if true , save cross resource references as attributes
	OPTION_LINE_WIDTH                    = "LINE_WIDTH"                    // int, start a new line before an attribute once a line is longer
	OPTION_KEEP_DEFAULTS                 = "KEEP_DEFAULTS"                 // if true , save unset features with a default value
)

type Codec struct {
}

func (d *Codec) NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) ecore.EEncoder {
	return NewEncoder(resource, w, options)
}
func (d *Codec) NewDecoder(resource ecore.EResource, r io.Reader, options map[string]any) ecore.EDecoder {
	return NewDecoder(resource, r, options)
}
