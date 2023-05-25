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
	XML_OPTION_EXTENDED_META_DATA            = "EXTENDED_META_DATA"            // ExtendedMetaData pointer
	XML_OPTION_SUPPRESS_DOCUMENT_ROOT        = "SUPPRESS_DOCUMENT_ROOT"        // if true , suppress document root if found
	XML_OPTION_DEFERRED_REFERENCE_RESOLUTION = "DEFERRED_REFERENCE_RESOLUTION" // if true , defer id ref resolution
	XML_OPTION_DEFERRED_ROOT_ATTACHMENT      = "DEFERRED_ROOT_ATTACHMENT"      // if true , defer id ref resolution
	XML_OPTION_ID_ATTRIBUTE_NAME             = "ID_ATTRIBUTE_NAME"             // value of the id attribute
	XML_OPTION_ROOT_OBJECTS                  = "ROOT_OBJECTS"                  // list of root objects to save
)

type XMLCodec struct {
}

func (d *XMLCodec) NewEncoder(resource EResource, w io.Writer, options map[string]any) EEncoder {
	return NewXMLEncoder(resource, w, options)
}
func (d *XMLCodec) NewDecoder(resource EResource, r io.Reader, options map[string]any) EResourceDecoder {
	return NewXMLDecoder(resource, r, options)
}
