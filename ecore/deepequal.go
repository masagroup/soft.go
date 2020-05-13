// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

type deepEqual struct {
	objects map[EObject]EObject
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
		eObj1Internal := eObj1.(EObjectInternal)
		eObj2Internal := eObj2.(EObjectInternal)
		if eObj1Internal.EProxyURI() == eObj2Internal.EProxyURI() {
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

	for it := eClass.GetEAttributes().Iterator(); it.HasNext(); {
		eAttribute := it.Next().(EAttribute)
		if !eAttribute.IsDerived() && !dE.equalsAttribute(eObj1, eObj2, eAttribute) {
			delete(dE.objects, eObj1)
			delete(dE.objects, eObj2)
			return false
		}
	}
	for it := eClass.GetEReferences().Iterator(); it.HasNext(); {
		eReference := it.Next().(EReference)
		if !eReference.IsDerived() && !dE.equalsReference(eObj1, eObj2, eReference) {
			delete(dE.objects, eObj1)
			delete(dE.objects, eObj2)
			return false
		}
	}

	// There's no reason they aren't equal, so they are.
	return true

}

func (dE *deepEqual) equalsAll(l1 EList, l2 EList) bool {
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

func (dE *deepEqual) equalsAttribute(eObj1 EObject, eObj2 EObject, eAttribute EAttribute) bool {
	isSet1 := eObj1.EIsSet(eAttribute)
	isSet2 := eObj2.EIsSet(eAttribute)
	if isSet1 && isSet2 {
		value1 := eObj1.EGet(eAttribute)
		value2 := eObj2.EGet(eAttribute)
		return value1 == value2
	}
	return isSet1 == isSet2
}

func (dE *deepEqual) equalsReference(eObj1 EObject, eObj2 EObject, eReference EReference) bool {
	isSet1 := eObj1.EIsSet(eReference)
	isSet2 := eObj2.EIsSet(eReference)
	if isSet1 && isSet2 {
		value1 := eObj1.EGet(eReference)
		value2 := eObj2.EGet(eReference)
		if eReference.IsMany() {
			return dE.equalsAll(value1.(EList), value2.(EList))
		} else {
			return dE.equals(value1.(EObject), value2.(EObject))
		}
	}
	return isSet1 == isSet2
}
