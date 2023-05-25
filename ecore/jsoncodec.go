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
	JSON_OPTION_ID_ATTRIBUTE_NAME = "ID_ATTRIBUTE" // if true, save id attribute of the object
)

type JSONCodec struct {
}

func (jc *JSONCodec) NewEncoder(resource EResource, w io.Writer, options map[string]interface{}) EEncoder {
	return NewJSONEncoder(resource, w, options)
}
func (jc *JSONCodec) NewDecoder(resource EResource, r io.Reader, options map[string]interface{}) EResourceDecoder {
	panic("not implemented")
}

type jsonFeatureKind int

const (
	jfkTransient jsonFeatureKind = iota
	jfkData
	jfkDataList
	jfkObject
	jfkObjectList
	jfkObjectReference
	jfkObjectReferenceList
)

func getJSONCodecFeatureKind(eFeature EStructuralFeature) jsonFeatureKind {
	if eFeature.IsTransient() {
		return jfkTransient
	} else if eReference, _ := eFeature.(EReference); eReference != nil {
		if eReference.IsContainment() {
			if eReference.IsMany() {
				return jfkObjectList
			} else {
				return jfkObject
			}
		}
		opposite := eReference.GetEOpposite()
		if opposite != nil && opposite.IsContainment() {
			return jfkTransient
		}
		if eReference.IsResolveProxies() {
			if eReference.IsMany() {
				return jfkObjectReferenceList
			} else {
				return jfkObjectReference
			}
		}
		if eReference.IsMany() {
			return jfkObjectList
		} else {
			return jfkObject
		}
	} else if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
		if eAttribute.IsMany() {
			return jfkDataList
		} else {
			return jfkData
		}
	}
	return -1
}
