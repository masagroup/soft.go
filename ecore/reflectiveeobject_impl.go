// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// ReflectiveEObjectImpl is a reflective EObject implementation: all its features are dynamic, and
// their values are kept in a slice by feature ID. Like EMF MinimalEObjectImpl.Container used
// with reflective feature delegation, it is compact: its container is a field, and its adapters,
// resource, proxy URI and values share one packed storage.
type ReflectiveEObjectImpl struct {
	CompactEObjectContainer
	// class is a field, not in the packed storage: it is read at each access to a feature.
	class EClass
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
	if class == o.class {
		return
	}
	o.class = class
	// the values are by feature ID: resized to the features of the new class, as
	// DynamicEObjectImpl.SetEClass does
	if properties, _ := o.getField(properties_flag).([]any); properties != nil {
		newProperties := make([]any, o.EClass().GetFeatureCount())
		copy(newProperties, properties)
		o.setField(properties_flag, newProperties)
	}
}

func (o *ReflectiveEObjectImpl) EStaticFeatureCount() int {
	return 0
}

// EInvokeFromID invokes the operations of EObject, as EObjectImpl.
func (o *ReflectiveEObjectImpl) EInvokeFromID(operationID int, arguments EList) any {
	var eObject EObjectImpl
	eObject.SetInterfaces(o.GetInterfaces())
	return eObject.EInvokeFromID(operationID, arguments)
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

func (o *ReflectiveEObjectImpl) EDynamicIsSet(dynamicFeatureID int) bool {
	return o.getProperties()[dynamicFeatureID] != nil
}

// getProperties returns the values of the features, created at the first access.
func (o *ReflectiveEObjectImpl) getProperties() []any {
	if properties, _ := o.getField(properties_flag).([]any); properties != nil {
		return properties
	}
	properties := make([]any, o.EClass().GetFeatureCount())
	o.setField(properties_flag, properties)
	return properties
}
