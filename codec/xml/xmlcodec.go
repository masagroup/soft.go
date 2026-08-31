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
	OPTION_ID_REFERENCE_AS_PROXY         = "ID_REFERENCE_AS_PROXY"         // if true , an id that resolves to nothing in the document becomes a proxy carrying the raw id, for an external resolver
	OPTION_SUPPRESS_DEFAULT_NS           = "SUPPRESS_DEFAULT_NS"           // encode: if true , omit the default xmlns declaration on the root (no-namespace output)
	OPTION_NO_NAMESPACE_SCHEMA_LOCATION  = "NO_NAMESPACE_SCHEMA_LOCATION"  // encode: string, write xsi:noNamespaceSchemaLocation on the root
	OPTION_FEATURE_VALUE_PROVIDER        = "FEATURE_VALUE_PROVIDER"        // encode: func(ecore.EObject, ecore.EStructuralFeature) string, a non-empty value is written verbatim in place of a reference
	OPTION_KEEP_REQUIRED_DEFAULTS        = "KEEP_REQUIRED_DEFAULTS"        // encode: if true , save a required attribute even at its type default (dense form)
	OPTION_DENSE_DEFAULT_FIELDS          = "DENSE_DEFAULT_FIELDS"          // encode: func(ecore.EObject, ecore.EStructuralFeature) bool, an optional attribute it flags is saved at its default, with KEEP_REQUIRED_DEFAULTS
	OPTION_KEEP_COMMENTS                 = "KEEP_COMMENTS"                 // encode: if true , write back the XML comments the document carried, each among the siblings it stood between
)

type Codec struct {
}

func (d *Codec) NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) ecore.EEncoder {
	return NewEncoder(resource, w, options)
}
func (d *Codec) NewDecoder(resource ecore.EResource, r io.Reader, options map[string]any) ecore.EDecoder {
	return NewDecoder(resource, r, options)
}
