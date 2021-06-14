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

const (
	OPTION_EXTENDED_META_DATA        = "EXTENDED_META_DATA"        // ExtendedMetaData pointer
	OPTION_SUPPRESS_DOCUMENT_ROOT    = "SUPPRESS_DOCUMENT_ROOT"    // if true , suppress document root if found
	OPTION_IDREF_RESOLUTION_DEFERRED = "IDREF_RESOLUTION_DEFERRED" // if true , defer id ref resolution
	OPTION_ID_ATTRIBUTE_NAME         = "ID_ATTRIBUTE_NAME"         // value of the id attribute
	OPTION_ROOT_OBJECTS              = "ROOT_OBJECTS"              // list of root objects to save
)

type EResourceCodec interface {
	NewEncoder(w io.Writer, options map[string]interface{}) EResourceEncoder
	NewDecoder(r io.Reader, options map[string]interface{}) EResourceDecoder
}
