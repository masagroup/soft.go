// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package json

import (
	"github.com/masagroup/soft.go/ecore"
	"io"
)

const (
	OPTION_ID_ATTRIBUTE_NAME = "ID_ATTRIBUTE" // if true, save id attribute of the object
)

type Codec struct {
}

func (jc *Codec) NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) ecore.EEncoder {
	return NewEncoder(resource, w, options)
}
func (jc *Codec) NewDecoder(resource ecore.EResource, r io.Reader, options map[string]any) ecore.EDecoder {
	return NewDecoder(resource, r, options)
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

func getJSONCodecFeatureKind(eFeature ecore.EStructuralFeature) jsonFeatureKind {
	if eFeature.IsTransient() {
		return jfkTransient
	} else if eReference, _ := eFeature.(ecore.EReference); eReference != nil {
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
		// a cross reference is written as a reference, whether it resolves proxies or not: a
		// copy of the referenced object could not be decoded back to the same object
		if eReference.IsMany() {
			return jfkObjectReferenceList
		} else {
			return jfkObjectReference
		}
	} else if eAttribute, _ := eFeature.(ecore.EAttribute); eAttribute != nil {
		if eAttribute.IsMany() {
			return jfkDataList
		} else {
			return jfkData
		}
	}
	return -1
}
