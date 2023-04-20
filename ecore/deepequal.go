// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "reflect"

type deepEqual struct {
	objects map[EObject]EObject
}

func newDeepEqual() *deepEqual {
	return &deepEqual{
		objects: make(map[EObject]EObject),
	}
}

func (dE *deepEqual) equals(eObj1 EObject, eObj2 EObject) bool {
	// If the first object is null, the second object must be null.
	if eObj1 == nil {
		return eObj2 == nil
	}

	// We know the first object isn't null, so if the second one is, it can't be equal.
	if eObj2 == nil {
		return false
	}

	// Both eObject1 and eObject2 are not null.
	// If eObject1 has been compared already...
	eObj1Mapped := dE.objects[eObj1]
	if eObj1Mapped != nil {
		// Then eObject2 must be that previous match.
		return eObj1Mapped == eObj2
	}

	// If eObject2 has been compared already...
	eObj2Mapped := dE.objects[eObj2]
	if eObj2Mapped != nil {
		// Then eObject1 must be that match.
		return eObj2Mapped == eObj1
	}

	// Neither eObject1 nor eObject2 have been compared yet.

	// If eObject1 and eObject2 are the same instance...
	if eObj1 == eObj2 {
		// Match them and return true.
		//
		dE.objects[eObj1] = eObj2
		dE.objects[eObj2] = eObj1
		return true
	}

	// If eObject1 is a proxy...
	if eObj1.EIsProxy() {
		eURI1 := eObj1.(EObjectInternal).EProxyURI()
		eURI2 := eObj2.(EObjectInternal).EProxyURI()
		if (eURI1 != nil && eURI1.Equals(eURI2)) || (eURI1 == nil && eURI2 == nil) {
			dE.objects[eObj1] = eObj2
			dE.objects[eObj2] = eObj1
			return true
		} else {
			return false
		}
	} else if eObj2.EIsProxy() {
		// If eObject1 isn't a proxy but eObject2 is, they can't be equal.
		return false
	}

	// If they don't have the same class, they can't be equal.
	eClass := eObj1.EClass()
	if eClass != eObj2.EClass() {
		return false
	}

	// Assume from now on that they match.
	dE.objects[eObj1] = eObj2
	dE.objects[eObj2] = eObj1

	for it := eClass.GetEStructuralFeatures().Iterator(); it.HasNext(); {
		if eFeature := it.Next().(EStructuralFeature); !eFeature.IsDerived() && !dE.equalsFeature(eObj1, eObj2, eFeature) {
			delete(dE.objects, eObj1)
			delete(dE.objects, eObj2)
			return false
		}
	}

	// There's no reason they aren't equal, so they are.
	return true

}

func (dE *deepEqual) equalsObjectList(l1 EList, l2 EList) bool {
	size := l1.Size()
	if size != l2.Size() {
		return false
	}
	for i := 0; i < size; i++ {
		eObj1 := l1.Get(i).(EObject)
		eObj2 := l2.Get(i).(EObject)
		if !dE.equals(eObj1, eObj2) {
			return false
		}
	}
	return true
}

func (dE *deepEqual) equalsPrimitiveList(l1 EList, l2 EList) bool {
	size := l1.Size()
	if size != l2.Size() {
		return false
	}
	for i := 0; i < size; i++ {
		p1 := l1.Get(i)
		p2 := l2.Get(i)
		if !reflect.DeepEqual(p1, p2) {
			return false
		}
	}
	return true
}

func (dE *deepEqual) equalsFeature(eObj1 EObject, eObj2 EObject, eFeature EStructuralFeature) bool {
	isSet1 := eObj1.EIsSet(eFeature)
	isSet2 := eObj2.EIsSet(eFeature)
	if isSet1 && isSet2 {
		if eAttribute, isAttribute := eFeature.(EAttribute); isAttribute {
			return dE.equalsAttribute(eObj1, eObj2, eAttribute)
		} else if eReference, isReference := eFeature.(EReference); isReference {
			return dE.equalsReference(eObj1, eObj2, eReference)
		}
		panic("invalid feature type")
	}
	return isSet1 == isSet2
}

func (dE *deepEqual) equalsAttribute(eObj1 EObject, eObj2 EObject, eAttribute EAttribute) bool {
	value1 := eObj1.EGet(eAttribute)
	value2 := eObj2.EGet(eAttribute)
	if value1 == nil {
		return value2 == nil
	}
	if value2 == nil {
		return false
	}
	if eAttribute.IsMany() {
		// attribute list
		l1 := value1.(EList)
		l2 := value2.(EList)
		return dE.equalsPrimitiveList(l1, l2)
	} else {
		return reflect.DeepEqual(value1, value2)
	}
}

func (dE *deepEqual) equalsReference(eObj1 EObject, eObj2 EObject, eReference EReference) bool {
	value1 := eObj1.EGet(eReference)
	value2 := eObj2.EGet(eReference)
	if eReference.IsMany() {
		return dE.equalsObjectList(value1.(EList), value2.(EList))
	} else {
		return dE.equals(value1.(EObject), value2.(EObject))
	}
}
