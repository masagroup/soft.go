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
		return o.EStaticClass()
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

// EGetFromID ...
func (o *DynamicEObjectImpl) EGetFromID(featureID int, resolve bool) interface{} {
	dynamicFeatureID := featureID - o.eStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		dynamicFeature := o.eDynamicFeature(featureID)
		if o.isContainer(dynamicFeature) {
			if o.EContainerFeatureID() == dynamicFeature.GetFeatureID() {
				if resolve {
					return o.EContainer()
				} else {
					return o.EInternalContainer()
				}
			}
		} else {
			// retrieve value or compute it if empty
			result := o.properties[dynamicFeatureID]
			if result == nil {
				if dynamicFeature.IsMany() {
					result = o.createList(dynamicFeature)
					o.properties[dynamicFeatureID] = result
				}

			} else if resolve && o.isProxy(dynamicFeature) {
				oldValue, _ := result.(EObject)
				newValue := o.EResolveProxy(oldValue)
				result = newValue

				if oldValue != newValue {
					o.properties[dynamicFeatureID] = newValue
					if o.isContains(dynamicFeature) {
						var notifications ENotificationChain
						if !o.isBidirectional(dynamicFeature) {
							if oldValue != nil {
								notifications = oldValue.(EObjectInternal).EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
							}
							if newValue != nil {
								notifications = newValue.(EObjectInternal).EInverseAdd(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
							}
						} else {
							dynamicReference := dynamicFeature.(EReference)
							reverseFeature := dynamicReference.GetEOpposite()
							if oldValue != nil {
								notifications = oldValue.(EObjectInternal).EInverseRemove(o.AsEObject(), reverseFeature.GetFeatureID(), notifications)
							}
							if newValue != nil {
								notifications = newValue.(EObjectInternal).EInverseAdd(o.AsEObject(), reverseFeature.GetFeatureID(), notifications)
							}
						}
						if notifications != nil {
							notifications.Dispatch()
						}
					}
					if o.ENotificationRequired() {
						o.ENotify(NewNotificationByFeatureID(o.AsEObject(), RESOLVE, featureID, oldValue, newValue, NO_INDEX))
					}
				}
			}
			return result
		}
	}
	return o.EObjectImpl.EGetFromID(featureID, resolve)
}

// ESetFromID ...
func (o *DynamicEObjectImpl) ESetFromID(featureID int, newValue interface{}) {
	dynamicFeatureID := featureID - o.eStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		dynamicFeature := o.eDynamicFeature(featureID)
		if o.isContainer(dynamicFeature) {
			// container
			newContainer := newValue.(EObject)
			if newContainer != o.EContainer() || (newContainer != nil && o.EContainerFeatureID() != featureID) {
				var notifications ENotificationChain
				if o.EContainer() != nil {
					notifications = o.EBasicRemoveFromContainer(notifications)
				}
				if newContainer != nil {
					notifications = newContainer.(EObjectInternal).EInverseAdd(o.AsEObject(), featureID, notifications)
				}
				notifications = o.EBasicSetContainer(newContainer, featureID, notifications)
				if notifications != nil {
					notifications.Dispatch()
				}
			} else if o.ENotificationRequired() {
				o.ENotify(NewNotificationByFeatureID(o.AsEObject(), SET, featureID, newValue, newValue, NO_INDEX))
			}
		} else if o.isBidirectional(dynamicFeature) || o.isContains(dynamicFeature) {
			// inverse - opposite
			oldValue := o.properties[dynamicFeatureID]
			if oldValue != newValue {
				var notifications ENotificationChain
				oldObject, _ := oldValue.(EObject)
				newObject, _ := newValue.(EObject)

				if !o.isBidirectional(dynamicFeature) {
					if oldObject != nil {
						notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
					}
					if newObject != nil {
						notifications = newObject.(EObjectInternal).EInverseAdd(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
					}
				} else {
					dynamicReference := dynamicFeature.(EReference)
					reverseFeature := dynamicReference.GetEOpposite()
					if oldObject != nil {
						notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), reverseFeature.GetFeatureID(), notifications)
					}
					if newObject != nil {
						notifications = newObject.(EObjectInternal).EInverseAdd(o.AsEObject(), reverseFeature.GetFeatureID(), notifications)
					}
				}
				// basic set
				o.properties[dynamicFeatureID] = newValue

				// create notification
				if o.ENotificationRequired() {
					notification := NewNotificationByFeatureID(o.AsEObject(), SET, featureID, oldValue, newValue, NO_INDEX)
					if notifications != nil {
						notifications.Add(notification)
					} else {
						notifications = notification
					}
				}

				// notify
				if notifications != nil {
					notifications.Dispatch()
				}
			}
		} else {
			// basic set
			oldValue := o.properties[dynamicFeatureID]
			o.properties[dynamicFeatureID] = newValue

			// notify
			if o.ENotificationRequired() {
				o.ENotify(NewNotificationByFeatureID(o.AsEObject(), SET, featureID, oldValue, newValue, NO_INDEX))
			}
		}
	} else {
		o.EObjectImpl.ESetFromID(featureID, newValue)
	}
}

// EIsSetFromID ...
func (o *DynamicEObjectImpl) EIsSetFromID(featureID int) bool {
	dynamicFeatureID := featureID - o.eStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		dynamicFeature := o.eDynamicFeature(featureID)
		if o.isContainer(dynamicFeature) {
			return o.EContainerFeatureID() == featureID && o.EInternalContainer() != nil
		} else {
			return o.properties[dynamicFeatureID] != nil
		}
	}
	return o.EObjectImpl.EIsSetFromID(featureID)
}

// EUnsetFromID ...
func (o *DynamicEObjectImpl) EUnsetFromID(featureID int) {
	dynamicFeatureID := featureID - o.eStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		dynamicFeature := o.eDynamicFeature(featureID)
		if o.isContainer(dynamicFeature) {
			if o.EInternalContainer() != nil {
				notifications := o.EBasicRemoveFromContainer(nil)
				notifications = o.EBasicSetContainer(nil, featureID, notifications)
				if notifications != nil {
					notifications.Dispatch()
				}
			} else if o.ENotificationRequired() {
				o.ENotify(NewNotificationByFeatureID(o.AsEObject(), SET, featureID, nil, nil, NO_INDEX))
			}
		} else if o.isBidirectional(dynamicFeature) || o.isContains(dynamicFeature) {
			// inverse - opposite
			oldValue := o.properties[dynamicFeatureID]
			if oldValue != nil {
				var notifications ENotificationChain
				oldObject, _ := oldValue.(EObject)

				if !o.isBidirectional(dynamicFeature) {
					if oldObject != nil {
						notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
					}
				} else {
					dynamicReference := dynamicFeature.(EReference)
					reverseFeature := dynamicReference.GetEOpposite()
					if oldObject != nil {
						notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), reverseFeature.GetFeatureID(), notifications)
					}
				}
				// basic unset
				o.properties[dynamicFeatureID] = nil

				// create notification
				if o.ENotificationRequired() {
					eventType := SET
					if dynamicFeature.IsUnsettable() {
						eventType = UNSET
					}
					notification := NewNotificationByFeatureID(o.AsEObject(), eventType, featureID, oldValue, nil, NO_INDEX)
					if notifications != nil {
						notifications.Add(notification)
					} else {
						notifications = notification
					}
				}

				// notify
				if notifications != nil {
					notifications.Dispatch()
				}
			}
		} else {
			oldValue := o.properties[dynamicFeatureID]
			o.properties[dynamicFeatureID] = nil
			if o.ENotificationRequired() {
				o.ENotify(NewNotificationByFeatureID(o.AsEObject(), UNSET, featureID, oldValue, nil, NO_INDEX))
			}
		}
	} else {
		o.EObjectImpl.EUnsetFromID(featureID)
	}
}

// EInvokeFromID ...
func (o *DynamicEObjectImpl) EInvokeFromID(operationID int, arguments EList) interface{} {
	return nil
}

// EBasicInverseAdd ...
func (o *DynamicEObjectImpl) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	dynamicFeatureID := featureID - o.eStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		dynamicFeature := o.eDynamicFeature(featureID)
		if dynamicFeature.IsMany() {
			value := o.properties[dynamicFeatureID]
			if value == nil {
				value = o.createList(dynamicFeature)
				o.properties[dynamicFeatureID] = value
			}
			list := value.(ENotifyingList)
			return list.AddWithNotification(otherEnd, notifications)
		} else if o.isContainer(dynamicFeature) {
			msgs := notifications
			if o.EContainer() != nil {
				msgs = o.EBasicRemoveFromContainer(msgs)
			}
			return o.EBasicSetContainer(otherEnd, featureID, msgs)
		} else {
			// inverse - opposite
			oldValue := o.properties[dynamicFeatureID]
			oldObject, _ := oldValue.(EObject)
			if oldObject != nil {
				if o.isContains(dynamicFeature) {
					notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
				} else if o.isBidirectional(dynamicFeature) {
					dynamicReference := dynamicFeature.(EReference)
					reverseFeature := dynamicReference.GetEOpposite()
					notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), reverseFeature.GetFeatureID(), notifications)
				}
			}

			// set current value
			o.properties[dynamicFeatureID] = otherEnd

			// create notification
			if o.ENotificationRequired() {
				notification := NewNotificationByFeatureID(o.AsEObject(), SET, featureID, oldValue, otherEnd, NO_INDEX)
				if notifications != nil {
					notifications.Add(notification)
				} else {
					notifications = notification
				}
			}
		}
	}
	return notifications
}

// EBasicInverseRemove ...
func (o *DynamicEObjectImpl) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	dynamicFeatureID := featureID - o.eStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		dynamicFeature := o.eDynamicFeature(featureID)
		if dynamicFeature.IsMany() {
			value := o.properties[dynamicFeatureID]
			if value != nil {
				list := value.(ENotifyingList)
				return list.RemoveWithNotification(otherEnd, notifications)
			}
		} else if o.isContainer(dynamicFeature) {
			return o.EBasicSetContainer(nil, featureID, notifications)
		} else {
			oldValue := o.properties[dynamicFeatureID]
			o.properties[dynamicFeatureID] = nil

			// create notification
			if o.ENotificationRequired() {
				notification := NewNotificationByFeatureID(o.AsEObject(), SET, featureID, oldValue, nil, NO_INDEX)
				if notifications != nil {
					notifications.Add(notification)
				} else {
					notifications = notification
				}
			}
		}
	}
	return notifications
}

func (o *DynamicEObjectImpl) resizeProperties() {
	newSize := o.EClass().GetFeatureCount() - o.eStaticFeatureCount()
	newProperties := make([]interface{}, newSize)
	copy(newProperties, o.properties)
	o.properties = newProperties
}

func (o *DynamicEObjectImpl) eStaticFeatureCount() int {
	return o.EStaticClass().GetFeatureCount()
}

func (o *DynamicEObjectImpl) eStaticOperationCount() int {
	return o.EStaticClass().GetOperationCount()
}

func (o *DynamicEObjectImpl) eDynamicFeatureID(feature EStructuralFeature) int {
	return o.EClass().GetFeatureID(feature) - o.eStaticFeatureCount()
}

func (o *DynamicEObjectImpl) eDynamicFeature(dynamicFeatureID int) EStructuralFeature {
	return o.EClass().GetEStructuralFeature(dynamicFeatureID + o.eStaticFeatureCount())
}

func (o *DynamicEObjectImpl) isBidirectional(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		return ref.GetEOpposite() != nil
	}
	return false
}

func (o *DynamicEObjectImpl) isContainer(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		opposite := ref.GetEOpposite()
		if opposite != nil {
			return opposite.IsContainment()
		}
	}
	return false
}

func (o *DynamicEObjectImpl) isContains(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		return ref.IsContainment()
	}
	return false
}

func (o *DynamicEObjectImpl) isProxy(feature EStructuralFeature) bool {
	if o.isContainer(feature) || o.isContains(feature) {
		return false
	}

	ref, isRef := feature.(EReference)
	if isRef {
		return ref.IsResolveProxies()
	}
	return false
}

func (o *DynamicEObjectImpl) createList(feature EStructuralFeature) EList {
	if attribute, isAttribute := feature.(EAttribute); isAttribute {
		if attribute.IsUnique() {
			return NewUniqueBasicEList(nil)
		} else {
			return NewBasicEList(nil)
		}
	} else if ref, isRef := feature.(EReference); isRef {
		inverse := false
		opposite := false
		reverseID := -1
		reverseFeature := ref.GetEOpposite()
		if reverseFeature != nil {
			reverseID = reverseFeature.GetFeatureID()
			inverse = true
			opposite = true
		} else if ref.IsContainment() {
			inverse = true
			opposite = false
		}
		return NewBasicEObjectList(o.AsEObjectInternal(), ref.GetFeatureID(), reverseID, ref.IsContainment(), inverse, opposite, ref.EIsProxy(), ref.IsUnsettable())
	}
	return nil
}
