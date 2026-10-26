// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// DynamicEObjectImpl ...
type DynamicEObjectImpl struct {
	EObjectImpl
	class      EClass
	properties []any
}

// NewDynamicEObjectImpl is the constructor of a DynamicEObjectImpl
func NewDynamicEObjectImpl() *DynamicEObjectImpl {
	o := new(DynamicEObjectImpl)
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *DynamicEObjectImpl) Initialize() {
	o.EObjectImpl.Initialize()
	o.resizeProperties()
}

// EClass ...
func (o *DynamicEObjectImpl) EClass() EClass {
	if o.class == nil {
		return o.AsEObjectInternal().EStaticClass()
	}
	return o.class
}

// SetEClass ...
func (o *DynamicEObjectImpl) SetEClass(class EClass) {
	if class != o.class {
		o.class = class
		o.resizeProperties()
	}
}

func (o *DynamicEObjectImpl) EStaticClass() EClass {
	return GetPackage().GetEObject()
}

func (o *DynamicEObjectImpl) EStaticFeatureCount() int {
	return 0
}

func (o *DynamicEObjectImpl) EDynamicProperties() EDynamicProperties {
	return o.GetInterfaces().(EDynamicProperties)
}

func (o *DynamicEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	return o.properties[dynamicFeatureID]
}

func (o *DynamicEObjectImpl) EDynamicSet(dynamicFeatureID int, newValue any) {
	o.properties[dynamicFeatureID] = newValue
}

func (o *DynamicEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.properties[dynamicFeatureID] = nil
}

// eKeepsNil tells that the properties keep NIL, the value of a feature set to nil.
func (o *DynamicEObjectImpl) eKeepsNil() {}

func (o *DynamicEObjectImpl) EDynamicIsSet(dynamicFeatureID int) bool {
	// as in EMF, an attribute (not unsettable) set to its default value is not set, one set to nil
	// (NIL) is
	value := o.properties[dynamicFeatureID]
	return value != nil && isSetValue(o.EClass().GetEStructuralFeature(dynamicFeatureID), value)
}

func (o *DynamicEObjectImpl) EFeatureID(feature EStructuralFeature) int {
	return o.EClass().GetFeatureID(feature)
}

func (o *DynamicEObjectImpl) EOperationID(operation EOperation) int {
	return o.EClass().GetOperationID(operation)
}

func (o *DynamicEObjectImpl) resizeProperties() {
	newSize := o.EClass().GetFeatureCount()
	newProperties := make([]any, newSize)
	copy(newProperties, o.properties)
	o.properties = newProperties
}
