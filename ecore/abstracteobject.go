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
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// AbstractEObject is a basic implementation of an EObject
type AbstractEObject struct {
	AbstractENotifier
}

type EDynamicProperties interface {
	EDynamicGet(dynamicFeatureID int) any
	EDynamicSet(dynamicFeatureID int, newValue any)
	EDynamicUnset(dynamicFeatureID int)
	EDynamicIsSet(dynamicFeatureID int) bool
}

// EObjectInternal ...
type EObjectInternal interface {
	EObject

	EDynamicProperties() EDynamicProperties

	EStaticClass() EClass
	EStaticFeatureCount() int

	EInternalContainer() EObject
	EInternalContainerFeatureID() int
	EInternalResource() EResource
	ESetInternalContainer(container EObject, containerFeatureID int)
	ESetInternalResource(resource EResource)

	ESetResource(resource EResource, notifications ENotificationChain) ENotificationChain

	EInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain
	EInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain

	EFeatureID(feature EStructuralFeature) int
	EDerivedFeatureID(container EObject, featureID int) int
	EOperationID(operation EOperation) int
	EDerivedOperationID(container EObject, operationID int) int
	EGetFromID(featureID int, resolve bool) any
	ESetFromID(featureID int, newValue any)
	EUnsetFromID(featureID int)
	EIsSetFromID(featureID int) bool
	EInvokeFromID(operationID int, arguments EList) any

	EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain
	EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain

	EObjectForFragmentSegment(string) EObject
	EURIFragmentSegment(EStructuralFeature, EObject) string

	EProxyURI() *URI
	ESetProxyURI(uri *URI)
	EResolveProxy(proxy EObject) EObject
}

// AsEObject ...
func (o *AbstractEObject) AsEObject() EObject {
	return o.interfaces.(EObject)
}

// AsEObjectInternal ...
func (o *AbstractEObject) AsEObjectInternal() EObjectInternal {
	return o.interfaces.(EObjectInternal)
}

// EClass ...
func (o *AbstractEObject) EClass() EClass {
	return o.AsEObjectInternal().EStaticClass()
}

// EStaticClass ...
func (o *AbstractEObject) EStaticClass() EClass {
	return GetPackage().GetEObject()
}

func (o *AbstractEObject) EStaticFeatureCount() int {
	return o.AsEObjectInternal().EStaticClass().GetFeatureCount()
}

func (o *AbstractEObject) EDynamicProperties() EDynamicProperties {
	return nil
}

// EContainer ...
func (o *AbstractEObject) EContainer() EObject {
	eContainer := o.AsEObjectInternal().EInternalContainer()
	if eContainer != nil && eContainer.EIsProxy() {
		resolved := o.EResolveProxy(eContainer)
		if resolved != eContainer {
			notifications := o.EBasicRemoveFromContainer(nil)
			objectInternal := o.AsEObjectInternal()
			containerFeatureID := objectInternal.EInternalContainerFeatureID()
			objectInternal.ESetInternalContainer(resolved, containerFeatureID)
			if notifications != nil {
				notifications.Dispatch()
			}

			if o.ENotificationRequired() && containerFeatureID >= EOPPOSITE_FEATURE_BASE {
				o.ENotify(NewNotificationByFeatureID(o.AsEObject(), RESOLVE, containerFeatureID, eContainer, resolved, -1))
			}
		}
		return resolved
	}
	return eContainer
}

func (o *AbstractEObject) EContainerFeatureID() int {
	return o.AsEObjectInternal().EInternalContainerFeatureID()
}

// EResource ...
func (o *AbstractEObject) EResource() EResource {
	resource := o.AsEObjectInternal().EInternalResource()
	if resource == nil {
		if container := o.AsEObjectInternal().EInternalContainer(); container != nil && container != o.AsEObject() {
			resource = container.EResource()
		}
	}
	return resource
}

// ESetResource ...
func (o *AbstractEObject) ESetResource(newResource EResource, n ENotificationChain) ENotificationChain {
	notifications := n
	oldResource := o.AsEObjectInternal().EInternalResource()
	// When setting the resource to nil we assume that detach has already been called in the resource implementation
	if oldResource != nil && newResource != nil {
		list := oldResource.GetContents().(ENotifyingList)
		notifications = list.RemoveWithNotification(o.AsEObject(), notifications)
		oldResource.Detached(o.AsEObject())
	}

	eContainer := o.AsEObjectInternal().EInternalContainer()
	if eContainer != nil {
		if o.EContainmentFeature().IsResolveProxies() {
			if eContainerInternal, _ := eContainer.(EObjectInternal); eContainerInternal != nil {
				oldContainerResource := eContainerInternal.EResource()
				if oldContainerResource != nil {
					if newResource == nil {
						// If we're not setting a new resource, attach it to the old container's resource.
						oldContainerResource.Attached(o.AsEObject())
					} else if oldResource == nil {
						// If we didn't detach it from an old resource already, detach it from the old container's resource.
						oldContainerResource.Detached(o.AsEObject())
					}
				}
			}
		} else {
			notifications = o.EBasicRemoveFromContainer(notifications)
			notifications = o.EBasicSetContainer(nil, -1, notifications)
		}
	}
	o.AsEObjectInternal().ESetInternalResource(newResource)
	return notifications
}

// EContainingFeature ...
func (o *AbstractEObject) EContainingFeature() EStructuralFeature {
	objectInternal := o.AsEObjectInternal()
	if container := objectInternal.EInternalContainer(); container != nil {
		containerFeatureID := objectInternal.EInternalContainerFeatureID()
		if containerFeatureID <= EOPPOSITE_FEATURE_BASE {
			return container.EClass().GetEStructuralFeature(EOPPOSITE_FEATURE_BASE - containerFeatureID)
		} else {
			return o.AsEObject().EClass().GetEStructuralFeature(containerFeatureID).(EReference).GetEOpposite()
		}
	}
	return nil
}

// EContainmentFeature ...
func (o *AbstractEObject) EContainmentFeature() EReference {
	objectInternal := o.AsEObjectInternal()
	return eContainmentFeature(objectInternal, objectInternal.EInternalContainer(), objectInternal.EInternalContainerFeatureID())
}

func eContainmentFeature(o EObject, container EObject, containerFeatureID int) EReference {
	if container != nil {
		if containerFeatureID <= EOPPOSITE_FEATURE_BASE {
			feature := container.EClass().GetEStructuralFeature(EOPPOSITE_FEATURE_BASE - containerFeatureID)
			reference, isReference := feature.(EReference)
			if isReference {
				return reference
			}
		} else {
			feature := o.EClass().GetEStructuralFeature(containerFeatureID)
			reference, isReference := feature.(EReference)
			if isReference {
				return reference
			}
		}
		panic("The containment feature could not be located")
	}
	return nil
}

// EAllContents ...
func (o *AbstractEObject) EAllContents() EIterator {
	return newEAllContentsIterator(o.AsEObject())
}

func (o *AbstractEObject) EFeatureID(feature EStructuralFeature) int {
	if !o.AsEObject().EClass().GetEAllStructuralFeatures().Contains(feature) {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
	return o.AsEObjectInternal().EDerivedFeatureID(feature.EContainer(), feature.GetFeatureID())
}

func (o *AbstractEObject) EDerivedFeatureID(container EObject, featureID int) int {
	return featureID
}

func (o *AbstractEObject) EOperationID(operation EOperation) int {
	if !o.AsEObject().EClass().GetEAllOperations().Contains(operation) {
		panic("The operation '" + operation.GetName() + "' is not a valid operation")
	}
	return o.AsEObjectInternal().EDerivedOperationID(operation.EContainer(), operation.GetOperationID())
}

func (o *AbstractEObject) EDerivedOperationID(container EObject, operationID int) int {
	return operationID
}

// EGet ...
func (o *AbstractEObject) EGet(feature EStructuralFeature) any {
	return o.eGetFromFeature(feature, true)
}

// EGetResolve ...
func (o *AbstractEObject) EGetResolve(feature EStructuralFeature, resolve bool) any {
	return o.eGetFromFeature(feature, resolve)
}

func (o *AbstractEObject) eGetFromFeature(feature EStructuralFeature, resolve bool) any {
	featureID := o.AsEObjectInternal().EFeatureID(feature)
	if featureID >= 0 {
		return o.AsEObjectInternal().EGetFromID(featureID, resolve)
	}
	panic("The feature '" + feature.GetName() + "' is not a valid feature")
}

// EGetFromID ...
func (o *AbstractEObject) EGetFromID(featureID int, resolve bool) any {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	if feature == nil {
		panic("Invalid featureID: " + strconv.Itoa(featureID))
	}
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID < 0 {
		return o.AsEObjectInternal().EGetResolve(feature, resolve)
	} else {
		properties := o.AsEObjectInternal().EDynamicProperties()
		if properties != nil {
			return o.eDynamicPropertiesGet(properties, feature, dynamicFeatureID, resolve)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
}

func (o *AbstractEObject) eDynamicPropertiesGet(properties EDynamicProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int, resolve bool) any {
	if IsContainer(dynamicFeature) {
		objInternal := o.AsEObjectInternal()
		featureID := objInternal.EClass().GetFeatureID(dynamicFeature)
		if objInternal.EInternalContainerFeatureID() == featureID {
			if resolve {
				return objInternal.EContainer()
			} else {
				return objInternal.EInternalContainer()
			}
		}
	} else {
		// retrieve value or compute it if empty
		result := properties.EDynamicGet(dynamicFeatureID)
		if result == nil {
			if dynamicFeature.IsMany() {
				if IsMapType(dynamicFeature) {
					result = o.eDynamicPropertiesCreateMap(dynamicFeature)
				} else {
					result = o.eDynamicPropertiesCreateList(dynamicFeature)
				}
				properties.EDynamicSet(dynamicFeatureID, result)
			} else if defaultValue := dynamicFeature.GetDefaultValue(); defaultValue != nil {
				result = defaultValue
			}
		} else if resolve && IsProxy(dynamicFeature) {
			oldValue, _ := result.(EObject)
			if oldValue != nil && oldValue.EIsProxy() {
				newValue := o.EResolveProxy(oldValue)
				result = newValue
				if oldValue != newValue {
					properties.EDynamicSet(dynamicFeatureID, newValue)
					if IsContains(dynamicFeature) {
						var notifications ENotificationChain
						if !IsBidirectional(dynamicFeature) {
							featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
							oldObject := oldValue.(EObjectInternal)
							notifications = oldObject.EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
							if newValue != nil {
								newObject := newValue.(EObjectInternal)
								notifications = newObject.EInverseAdd(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
							}
						} else {
							dynamicReference := dynamicFeature.(EReference)
							reverseFeature := dynamicReference.GetEOpposite()
							oldObject := oldValue.(EObjectInternal)
							featureID := oldObject.EClass().GetFeatureID(reverseFeature)
							notifications = oldObject.EInverseRemove(o.AsEObject(), featureID, notifications)
							if newValue != nil {
								newObject := newValue.(EObjectInternal)
								featureID := newObject.EClass().GetFeatureID(reverseFeature)
								notifications = newObject.EInverseAdd(o.AsEObject(), featureID, notifications)
							}
						}
						if notifications != nil {
							notifications.Dispatch()
						}
					}
					if o.ENotificationRequired() {
						o.ENotify(NewNotificationByFeature(o.AsEObject(), RESOLVE, dynamicFeature, oldValue, newValue, NO_INDEX))
					}
				}
			}
		}
		return result
	}
	return nil
}

func (o *AbstractEObject) eDynamicPropertiesCreateMap(feature EStructuralFeature) EMap {
	eClass := feature.GetEType().(EClass)
	reverseFeatureID := -1
	if ref, isRef := feature.(EReference); isRef {
		if reverseFeature := ref.GetEOpposite(); reverseFeature != nil {
			reverseFeatureID = reverseFeature.GetFeatureID()
		}
	}
	return NewBasicEObjectMap(eClass, o.AsEObjectInternal(), feature.GetFeatureID(), reverseFeatureID, feature.IsUnsettable())
}

func (o *AbstractEObject) eDynamicPropertiesCreateList(feature EStructuralFeature) EList {
	if attribute, isAttribute := feature.(EAttribute); isAttribute {
		return NewBasicEDataTypeList(o.AsEObjectInternal(), attribute.GetFeatureID(), attribute.IsUnique())
	} else if ref, isRef := feature.(EReference); isRef {
		inverse := false
		opposite := false
		containment := ref.IsContainment()
		reverseFeature := ref.GetEOpposite()
		reverseFeatureID := -1
		if containment {
			if reverseFeature != nil {
				reverseFeatureID = reverseFeature.GetFeatureID()
				inverse = true
				opposite = true
			} else {
				inverse = true
				opposite = false
			}
		} else {
			if reverseFeature != nil {
				reverseFeatureID = reverseFeature.GetFeatureID()
				inverse = true
				opposite = true
			} else {
				inverse = false
				opposite = false
			}
		}
		return NewBasicEObjectList(o.AsEObjectInternal(), ref.GetFeatureID(), reverseFeatureID, containment, inverse, opposite, ref.EIsProxy(), ref.IsUnsettable())
	}
	return nil
}

// ESet ...
func (o *AbstractEObject) ESet(feature EStructuralFeature, newValue any) {
	featureID := o.AsEObjectInternal().EFeatureID(feature)
	if featureID >= 0 {
		o.AsEObjectInternal().ESetFromID(featureID, newValue)
	} else {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
}

// ESetFromID ...
func (o *AbstractEObject) ESetFromID(featureID int, newValue any) {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	if feature == nil {
		panic("Invalid featureID: " + strconv.Itoa(featureID))
	}
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID < 0 {
		o.ESet(feature, newValue)
	} else {
		properties := o.AsEObjectInternal().EDynamicProperties()
		if properties != nil {
			o.eDynamicPropertiesSet(properties, feature, dynamicFeatureID, newValue)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
}

func (o *AbstractEObject) eDynamicPropertiesSet(properties EDynamicProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int, newValue any) {
	if IsContainer(dynamicFeature) {
		// container
		objInternal := o.AsEObjectInternal()
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		oldContainer := objInternal.EInternalContainer()
		newContainer, _ := newValue.(EObjectInternal)
		if newContainer != oldContainer || (newContainer != nil && objInternal.EInternalContainerFeatureID() != featureID) {
			var notifications ENotificationChain
			if oldContainer != nil {
				notifications = o.EBasicRemoveFromContainer(notifications)
			}
			if newContainer != nil {
				reverseFeature := dynamicFeature.(EReference).GetEOpposite()
				featureID := newContainer.EClass().GetFeatureID(reverseFeature)
				notifications = newContainer.EInverseAdd(o.AsEObject(), featureID, notifications)
			}
			notifications = o.EBasicSetContainer(newContainer, featureID, notifications)
			if notifications != nil {
				notifications.Dispatch()
			}
		} else if o.ENotificationRequired() {
			o.ENotify(NewNotificationByFeature(o.AsEObject(), SET, dynamicFeature, newValue, newValue, NO_INDEX))
		}
	} else if IsBidirectional(dynamicFeature) || IsContains(dynamicFeature) {
		// inverse - opposite
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		if oldValue != newValue {
			var notifications ENotificationChain
			oldObject, _ := oldValue.(EObjectInternal)
			newObject, _ := newValue.(EObjectInternal)

			if !IsBidirectional(dynamicFeature) {
				featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
				if oldObject != nil {
					notifications = oldObject.EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
				}
				if newObject != nil {
					notifications = newObject.EInverseAdd(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
				}
			} else {
				dynamicReference := dynamicFeature.(EReference)
				reverseFeature := dynamicReference.GetEOpposite()
				if oldObject != nil {
					featureID := oldObject.EClass().GetFeatureID(reverseFeature)
					notifications = oldObject.EInverseRemove(o.AsEObject(), featureID, notifications)
				}
				if newObject != nil {
					featureID := newObject.EClass().GetFeatureID(reverseFeature)
					notifications = newObject.EInverseAdd(o.AsEObject(), featureID, notifications)
				}
			}
			// basic set
			properties.EDynamicSet(dynamicFeatureID, newValue)

			// create notification
			if o.ENotificationRequired() {
				notification := NewNotificationByFeature(o.AsEObject(), SET, dynamicFeature, oldValue, newValue, NO_INDEX)
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
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		properties.EDynamicSet(dynamicFeatureID, newValue)

		// notify
		if o.ENotificationRequired() {
			o.ENotify(NewNotificationByFeature(o.AsEObject(), SET, dynamicFeature, oldValue, newValue, NO_INDEX))
		}
	}
}

// EIsSet ...
func (o *AbstractEObject) EIsSet(feature EStructuralFeature) bool {
	featureID := o.AsEObjectInternal().EFeatureID(feature)
	if featureID >= 0 {
		return o.AsEObjectInternal().EIsSetFromID(featureID)
	}
	panic("The feature '" + feature.GetName() + "' is not a valid feature")
}

// EIsSetFromID ...
func (o *AbstractEObject) EIsSetFromID(featureID int) bool {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	if feature == nil {
		panic("Invalid featureID: " + strconv.Itoa(featureID))
	}
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID < 0 {
		return o.EIsSet(feature)
	} else {
		properties := o.AsEObjectInternal().EDynamicProperties()
		if properties != nil {
			return o.eDynamicPropertiesIsSet(properties, feature, dynamicFeatureID)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
}

func (o *AbstractEObject) eDynamicPropertiesIsSet(properties EDynamicProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int) bool {
	if IsContainer(dynamicFeature) {
		objInternal := o.AsEObjectInternal()
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		return objInternal.EInternalContainerFeatureID() == featureID && objInternal.EInternalContainer() != nil
	} else {
		return properties.EDynamicIsSet(dynamicFeatureID)
	}
}

// EUnset ...
func (o *AbstractEObject) EUnset(feature EStructuralFeature) {
	featureID := o.AsEObjectInternal().EFeatureID(feature)
	if featureID >= 0 {
		o.AsEObjectInternal().EUnsetFromID(featureID)
	} else {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
}

// EUnsetFromID ...
func (o *AbstractEObject) EUnsetFromID(featureID int) {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	if feature == nil {
		panic("Invalid featureID: " + strconv.Itoa(featureID))
	}
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID < 0 {
		o.EUnset(feature)
	} else {
		properties := o.AsEObjectInternal().EDynamicProperties()
		if properties != nil {
			o.eDynamicPropertiesUnset(properties, feature, dynamicFeatureID)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
}

func (o *AbstractEObject) eDynamicPropertiesUnset(properties EDynamicProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int) {
	if IsContainer(dynamicFeature) {
		if o.AsEObjectInternal().EInternalContainer() != nil {
			featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
			notifications := o.EBasicRemoveFromContainer(nil)
			notifications = o.EBasicSetContainer(nil, featureID, notifications)
			if notifications != nil {
				notifications.Dispatch()
			}
		} else if o.ENotificationRequired() {
			o.ENotify(NewNotificationByFeature(o.AsEObject(), SET, dynamicFeature, nil, nil, NO_INDEX))
		}
	} else if IsBidirectional(dynamicFeature) || IsContains(dynamicFeature) {
		// inverse - opposite
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		if oldValue != nil {
			var notifications ENotificationChain
			oldObject, _ := oldValue.(EObject)

			if !IsBidirectional(dynamicFeature) {
				if oldObject != nil {
					featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
					notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
				}
			} else {
				dynamicReference := dynamicFeature.(EReference)
				reverseFeature := dynamicReference.GetEOpposite()
				if oldObject != nil {
					featureID := oldObject.EClass().GetFeatureID(reverseFeature)
					notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), featureID, notifications)
				}
			}
			// basic unset
			properties.EDynamicUnset(dynamicFeatureID)

			// create notification
			if o.ENotificationRequired() {
				eventType := SET
				if dynamicFeature.IsUnsettable() {
					eventType = UNSET
				}
				notification := NewNotificationByFeature(o.AsEObject(), eventType, dynamicFeature, oldValue, nil, NO_INDEX)
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
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		properties.EDynamicUnset(dynamicFeatureID)
		if o.ENotificationRequired() {
			o.ENotify(NewNotificationByFeature(o.AsEObject(), UNSET, dynamicFeature, oldValue, nil, NO_INDEX))
		}
	}
}

// EInvoke ...
func (o *AbstractEObject) EInvoke(operation EOperation, arguments EList) any {
	operationID := o.AsEObjectInternal().EOperationID(operation)
	if operationID >= 0 {
		return o.AsEObjectInternal().EInvokeFromID(operationID, arguments)
	}
	panic("The operation '" + operation.GetName() + "' is not a valid operation")
}

// EInvokeFromID ...
func (o *AbstractEObject) EInvokeFromID(operationID int, arguments EList) any {
	operation := o.AsEObject().EClass().GetEOperation(operationID)
	if operation == nil {
		panic("Invalid operationID: " + strconv.Itoa(operationID))
	}
	return nil
}

// EInverseAdd ...
func (o *AbstractEObject) EInverseAdd(otherEnd EObject, featureID int, n ENotificationChain) ENotificationChain {
	notifications := n
	if featureID >= 0 {
		return o.AsEObjectInternal().EBasicInverseAdd(otherEnd, featureID, notifications)
	} else {
		notifications = o.EBasicRemoveFromContainer(notifications)
		return o.EBasicSetContainer(otherEnd, featureID, notifications)
	}
}

// EInverseRemove ...
func (o *AbstractEObject) EInverseRemove(otherEnd EObject, featureID int, n ENotificationChain) ENotificationChain {
	if featureID >= 0 {
		return o.AsEObjectInternal().EBasicInverseRemove(otherEnd, featureID, n)
	} else {
		return o.EBasicSetContainer(nil, featureID, n)
	}
}

// EResolveProxy ...
func (o *AbstractEObject) EResolveProxy(proxy EObject) EObject {
	return ResolveInObject(proxy, o.GetInterfaces().(EObject))
}

// EBasicInverseAdd ...
func (o *AbstractEObject) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		properties := o.AsEObjectInternal().EDynamicProperties()
		if properties != nil {
			return o.eDynamicPropertiesInverseAdd(properties, otherEnd, feature, dynamicFeatureID, notifications)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
	return notifications
}

func (o *AbstractEObject) eDynamicPropertiesInverseAdd(properties EDynamicProperties, otherEnd EObject, dynamicFeature EStructuralFeature, dynamicFeatureID int, notifications ENotificationChain) ENotificationChain {
	if dynamicFeature.IsMany() {
		value := properties.EDynamicGet(dynamicFeatureID)
		if value == nil {
			value = o.eDynamicPropertiesCreateList(dynamicFeature)
			properties.EDynamicSet(dynamicFeatureID, value)
		}
		list := value.(ENotifyingList)
		return list.AddWithNotification(otherEnd, notifications)
	} else if IsContainer(dynamicFeature) {
		msgs := notifications
		if o.AsEObjectInternal().EInternalContainer() != nil {
			msgs = o.EBasicRemoveFromContainer(msgs)
		}
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		return o.EBasicSetContainer(otherEnd, featureID, msgs)
	} else {
		// inverse - opposite
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		oldObject, _ := oldValue.(EObject)
		if oldObject != nil {
			if IsContains(dynamicFeature) {
				featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
				notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
			} else if IsBidirectional(dynamicFeature) {
				dynamicReference := dynamicFeature.(EReference)
				reverseFeature := dynamicReference.GetEOpposite()
				featureID := oldObject.EClass().GetFeatureID(reverseFeature)
				notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), featureID, notifications)
			}
		}

		// set current value
		properties.EDynamicSet(dynamicFeatureID, otherEnd)

		// create notification
		if o.ENotificationRequired() {
			notification := NewNotificationByFeature(o.AsEObject(), SET, dynamicFeature, oldValue, otherEnd, NO_INDEX)
			if notifications != nil {
				notifications.Add(notification)
			} else {
				notifications = notification
			}
		}
	}
	return notifications
}

// EBasicInverseRemove ...
func (o *AbstractEObject) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		properties := o.AsEObjectInternal().EDynamicProperties()
		if properties != nil {
			return o.eDynamicPropertiesInverseRemove(properties, otherEnd, feature, dynamicFeatureID, notifications)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
	return notifications
}

func (o *AbstractEObject) eDynamicPropertiesInverseRemove(properties EDynamicProperties, otherEnd EObject, dynamicFeature EStructuralFeature, dynamicFeatureID int, notifications ENotificationChain) ENotificationChain {
	if dynamicFeature.IsMany() {
		value := properties.EDynamicGet(dynamicFeatureID)
		if value != nil {
			list := value.(ENotifyingList)
			return list.RemoveWithNotification(otherEnd, notifications)
		}
	} else if IsContainer(dynamicFeature) {
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		return o.EBasicSetContainer(nil, featureID, notifications)
	} else {
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		properties.EDynamicUnset(dynamicFeatureID)

		// create notification
		if o.ENotificationRequired() {
			notification := NewNotificationByFeature(o.AsEObject(), SET, dynamicFeature, oldValue, nil, NO_INDEX)
			if notifications != nil {
				notifications.Add(notification)
			} else {
				notifications = notification
			}
		}
	}
	return notifications
}

// EBasicSetContainer ...
func (o *AbstractEObject) EBasicSetContainer(newContainer EObject, newContainerFeatureID int, n ENotificationChain) ENotificationChain {
	notifications := n
	objInternal := o.AsEObjectInternal()
	oldResource := objInternal.EInternalResource()
	oldContainer := objInternal.EInternalContainer()
	oldContainerFeatureID := objInternal.EInternalContainerFeatureID()
	oldContainerInternal, _ := oldContainer.(EObjectInternal)
	newContainerInternal, _ := newContainer.(EObjectInternal)

	// resource
	var newResource EResource
	if oldResource != nil {
		if newContainerInternal != nil && !eContainmentFeature(o.AsEObject(), newContainerInternal, newContainerFeatureID).IsResolveProxies() {
			list := oldResource.GetContents().(ENotifyingList)
			notifications = list.RemoveWithNotification(o.AsEObject(), notifications)
			objInternal.ESetInternalResource(nil)
			newResource = newContainerInternal.EResource()
		} else {
			oldResource = nil
		}
	} else {
		if oldContainerInternal != nil {
			oldResource = oldContainerInternal.EResource()
		}

		if newContainerInternal != nil {
			newResource = newContainerInternal.EResource()
		}
	}

	// detach object from resource
	// container is correctly set to the old value
	if oldResource != nil && oldResource != newResource {
		oldResource.Detached(o.AsEObject())
	}

	// internal set
	objInternal.ESetInternalContainer(newContainer, newContainerFeatureID)

	// attach object to resource
	// container is correctly set to the new value
	if newResource != nil && newResource != oldResource {
		newResource.Attached(o.AsEObject())
	}

	// notification
	if o.ENotificationRequired() {
		if oldContainer != nil && oldContainerFeatureID >= 0 && oldContainerFeatureID != newContainerFeatureID {
			notification := NewNotificationByFeatureID(o.AsEObject(), SET, oldContainerFeatureID, oldContainer, nil, -1)
			if notifications != nil {
				notifications.Add(notification)
			} else {
				notifications = notification
			}
		}
		if newContainerFeatureID >= 0 {
			var c EObject
			if oldContainerFeatureID == newContainerFeatureID {
				c = oldContainer
			}
			notification := NewNotificationByFeatureID(o.AsEObject(), SET, newContainerFeatureID, c, newContainer, -1)
			if notifications != nil {
				notifications.Add(notification)
			} else {
				notifications = notification
			}
		}
	}
	return notifications
}

// EBasicRemoveFromContainer ...
func (o *AbstractEObject) EBasicRemoveFromContainer(notifications ENotificationChain) ENotificationChain {
	objInternal := o.AsEObjectInternal()
	if objInternal.EInternalContainerFeatureID() >= 0 {
		return o.EBasicRemoveFromContainerFeature(notifications)
	} else {
		if containerInternal, _ := objInternal.EInternalContainer().(EObjectInternal); containerInternal != nil {
			return containerInternal.EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-objInternal.EInternalContainerFeatureID(), notifications)
		}
	}
	return notifications
}

// EBasicRemoveFromContainerFeature ...
func (o *AbstractEObject) EBasicRemoveFromContainerFeature(notifications ENotificationChain) ENotificationChain {
	objInternal := o.AsEObjectInternal()
	reference, isReference := o.AsEObject().EClass().GetEStructuralFeature(objInternal.EInternalContainerFeatureID()).(EReference)
	if isReference {
		inverseFeature := reference.GetEOpposite()
		if containerInternal, _ := objInternal.EInternalContainer().(EObjectInternal); containerInternal != nil && inverseFeature != nil {
			return containerInternal.EInverseRemove(o.AsEObject(), inverseFeature.GetFeatureID(), notifications)
		}
	}
	return notifications
}

func (o *AbstractEObject) eStructuralFeature(featureName string) EStructuralFeature {
	eFeature := o.AsEObject().EClass().GetEStructuralFeatureFromName(featureName)
	if eFeature == nil {
		panic("The feature " + featureName + " is not a valid feature")
	}
	return eFeature
}

func (o *AbstractEObject) EObjectForFragmentSegment(uriSegment string) EObject {

	lastIndex := len(uriSegment) - 1
	if lastIndex == -1 || uriSegment[0] != '@' {
		panic("Expecting @ at index 0 of '" + uriSegment + "'")
	}

	index := -1
	r, _ := utf8.DecodeLastRuneInString(uriSegment)
	if unicode.IsDigit(r) {
		if index = strings.LastIndex(uriSegment, "."); index != -1 {
			pos, _ := strconv.Atoi(uriSegment[index+1:])
			eFeatureName := uriSegment[1:index]
			eFeature := o.eStructuralFeature(eFeatureName)
			list, _ := o.AsEObject().EGetResolve(eFeature, false).(EList)
			if list != nil && pos < list.Size() {
				return list.Get(pos).(EObject)
			}
		}
	}
	if index == -1 {
		eFeature := o.eStructuralFeature(uriSegment[1:])
		eObject, _ := o.AsEObject().EGetResolve(eFeature, false).(EObject)
		return eObject
	}
	return nil
}

func (o *AbstractEObject) EURIFragmentSegment(feature EStructuralFeature, object EObject) string {
	s := "@"
	s += feature.GetName()
	if feature.IsMany() {
		v := o.AsEObject().EGetResolve(feature, false)
		i := v.(EList).IndexOf(object)
		s += "." + strconv.Itoa(i)
	}
	return s
}

func (o *AbstractEObject) String() string {
	return fmt.Sprintf("%s(%p)", o.AsEObject().EClass().GetName(), o.AsEObject())
}
