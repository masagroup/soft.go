// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

// DynamicEObjectImpl ...
type DynamicEObjectImpl struct {
	*EObjectImpl
	class      EClass
	properties []interface{}
	adapter    *dynamicFeaturesAdapter
}

type dynamicFeaturesAdapter struct {
	*Adapter
	object *DynamicEObjectImpl
}

func (adapter *dynamicFeaturesAdapter) NotifyChanged(notification ENotification) {
	eventType := notification.GetEventType()
	if eventType != REMOVING_ADAPTER {
		featureID := notification.GetFeatureID()
		if featureID == ECLASS__ESTRUCTURAL_FEATURES {
			adapter.object.resizeProperties()
		}
	}
}

// NewDynamicEObjectImpl is the constructor of a DynamicEObjectImpl
func NewDynamicEObjectImpl() *DynamicEObjectImpl {
	o := new(DynamicEObjectImpl)
	o.EObjectImpl = NewEObjectImpl()
	o.adapter = &dynamicFeaturesAdapter{Adapter: NewAdapter(), object: o}
	o.SetInterfaces(o)
	o.SetEClass(nil)
	return o
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
	if o.class != nil {
		o.class.EAdapters().Remove(o.adapter)
	}

	o.class = class
	o.resizeProperties()

	if o.class != nil {
		o.class.EAdapters().Add(o.adapter)
	}
}

func (o *DynamicEObjectImpl) EStaticClass() EClass {
	return GetPackage().GetEObject()
}

func (o *DynamicEObjectImpl) EStaticFeatureCount() int {
	return 0
}

func (o *DynamicEObjectImpl) EProperties() EDynamicProperties {
	return o.GetInterfaces().(EDynamicProperties)
}

func (o *DynamicEObjectImpl) EDynamicGet(dynamicFeatureID int) interface{} {
	return o.properties[dynamicFeatureID]
}

func (o *DynamicEObjectImpl) EDynamicSet(dynamicFeatureID int, newValue interface{}) {
	o.properties[dynamicFeatureID] = newValue
}

func (o *DynamicEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.properties[dynamicFeatureID] = nil
}

func (o *DynamicEObjectImpl) resizeProperties() {
	newSize := o.EClass().GetFeatureCount()
	newProperties := make([]interface{}, newSize)
	copy(newProperties, o.properties)
	o.properties = newProperties
}
