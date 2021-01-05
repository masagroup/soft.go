// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

// ReflectiveEObjectImpl ...
type ReflectiveEObjectImpl struct {
	EObjectImpl
	class      EClass
	properties []interface{}
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

func (o *ReflectiveEObjectImpl) EDynamicGet(dynamicFeatureID int) interface{} {
	return o.getProperties()[dynamicFeatureID]
}

func (o *ReflectiveEObjectImpl) EDynamicSet(dynamicFeatureID int, newValue interface{}) {
	o.getProperties()[dynamicFeatureID] = newValue
}

func (o *ReflectiveEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.getProperties()[dynamicFeatureID] = nil
}

func (o *ReflectiveEObjectImpl) getProperties() []interface{} {
	if o.properties == nil {
		o.properties = make([]interface{}, o.EClass().GetFeatureCount())
	}
	return o.properties
}
