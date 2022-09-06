// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// ReflectiveEObjectImpl ...
type ReflectiveEObjectImpl struct {
	EObjectImpl
	class      EClass
	properties []any
}

// NewReflectiveEObjectImpl is the constructor of a ReflectiveEObjectImpl
func NewReflectiveEObjectImpl() *ReflectiveEObjectImpl {
	o := new(ReflectiveEObjectImpl)
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *ReflectiveEObjectImpl) EClass() EClass {
	if o.class == nil {
		return o.AsEObjectInternal().EStaticClass()
	}
	return o.class
}

// SetEClass ...
func (o *ReflectiveEObjectImpl) SetEClass(class EClass) {
	o.class = class
}

func (o *ReflectiveEObjectImpl) EStaticFeatureCount() int {
	return 0
}

func (o *ReflectiveEObjectImpl) EDynamicProperties() EDynamicProperties {
	return o.GetInterfaces().(EDynamicProperties)
}

func (o *ReflectiveEObjectImpl) EDynamicGet(dynamicFeatureID int) any {
	return o.getProperties()[dynamicFeatureID]
}

func (o *ReflectiveEObjectImpl) EDynamicSet(dynamicFeatureID int, newValue any) {
	o.getProperties()[dynamicFeatureID] = newValue
}

func (o *ReflectiveEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.getProperties()[dynamicFeatureID] = nil
}

func (o *ReflectiveEObjectImpl) getProperties() []any {
	if o.properties == nil {
		o.properties = make([]any, o.EClass().GetFeatureCount())
	}
	return o.properties
}
