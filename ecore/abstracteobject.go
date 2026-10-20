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
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// AbstractEObject is the base implementation of an EObject: the reflective API (EGet, ESet,
// EIsSet, EUnset, EInvoke) through the feature and operation IDs, the container and the resource
// of the object, the inverse references and the dynamic features. It accesses the object through
// its interfaces (AsEObject, AsEObjectInternal), so that the methods of the embedding struct
// (the generated class) are the ones called: the struct embedding AbstractEObject provides the
// storage of the container, the resource and the proxy URI (EInternalContainer,
// ESetInternalContainer...) and the access to its static features by ID (EGetFromID...).
type AbstractEObject struct {
	AbstractENotifier
}

// EDynamicProperties stores the values of the dynamic features of an object: the features of its
// class beyond its static features (EStaticFeatureCount), indexed by dynamic feature ID (the
// feature ID minus the static feature count).
type EDynamicProperties interface {
	// EDynamicGet returns the value of the dynamic feature, nil when it is not set.
	EDynamicGet(dynamicFeatureID int) any
	// EDynamicSet sets the value of the dynamic feature.
	EDynamicSet(dynamicFeatureID int, newValue any)
	// EDynamicUnset unsets the value of the dynamic feature.
	EDynamicUnset(dynamicFeatureID int)
	// EDynamicIsSet returns true when the value of the dynamic feature is set.
	EDynamicIsSet(dynamicFeatureID int) bool
}

// EObjectInternal is the implementation protocol of an EObject: its container and resource, the
// inverse references, the access to its features and operations by ID, and its proxy URI.
//
// It is used by the generated code, the lists and the resources of the objects, and the
// implementations of EResource or EList outside this package. It is not meant for application
// code, which uses EObject. It is exported because the generated classes, in other packages,
// implement its methods.
type EObjectInternal interface {
	EObject

	// EDynamicProperties returns the storage of the dynamic features of the object, nil when the
	// object has no dynamic features.
	EDynamicProperties() EDynamicProperties

	// EStaticClass returns the class the object is generated for. EClass returns it unless the
	// object has a dynamic class.
	EStaticClass() EClass
	// EStaticFeatureCount returns the number of features of the static class: the features with a
	// greater ID are dynamic (EDynamicProperties).
	EStaticFeatureCount() int

	// EInternalContainer returns the container of the object, without resolving it when it is a
	// proxy (EContainer resolves it).
	EInternalContainer() EObject
	// EInternalContainerFeatureID returns the feature ID of the container of the object: the ID of
	// the container reference of the object (>= 0) when the containment has an opposite, else
	// EOPPOSITE_FEATURE_BASE minus the ID of the containment reference of the container.
	EInternalContainerFeatureID() int
	// EInternalResource returns the resource the object is directly in (a root of its contents),
	// nil when it is in a resource through its container only (EResource returns both).
	EInternalResource() EResource
	// ESetInternalContainer stores the container and its feature ID, without notification nor
	// inverse update (see EBasicSetContainer).
	ESetInternalContainer(container EObject, containerFeatureID int)
	// ESetInternalResource stores the resource the object is directly in, without notification.
	ESetInternalResource(resource EResource)

	// ESetResource makes resource the resource the object is directly in, called by the contents
	// list of the resource: it removes the object from its old resource, and from its container
	// when the containment does not resolve proxies. It returns the notifications to dispatch.
	ESetResource(resource EResource, notifications ENotificationChain) ENotificationChain

	// EInverseAdd updates the feature of the object whose opposite now references otherEnd: a
	// featureID >= 0 is a feature of the object (EBasicInverseAdd), a negative one makes otherEnd
	// the container of the object. It returns the notifications to dispatch.
	EInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain
	// EInverseRemove updates the feature of the object whose opposite no longer references
	// otherEnd: a featureID >= 0 is a feature of the object (EBasicInverseRemove), a negative one
	// removes the container of the object. It returns the notifications to dispatch.
	EInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain

	// EFeatureID returns the ID of feature in the class of the object. It panics when the class
	// has no such feature.
	EFeatureID(feature EStructuralFeature) int
	// EDerivedFeatureID returns the ID, in the class of the object, of the feature of ID featureID
	// in container (the class or the package declaring the feature).
	EDerivedFeatureID(container EObject, featureID int) int
	// EOperationID returns the ID of operation in the class of the object. It panics when the
	// class has no such operation.
	EOperationID(operation EOperation) int
	// EDerivedOperationID returns the ID, in the class of the object, of the operation of ID
	// operationID in container (the class declaring the operation).
	EDerivedOperationID(container EObject, operationID int) int
	// EGetFromID returns the value of the feature of ID featureID (EGet), resolving the proxies
	// when resolve is true.
	EGetFromID(featureID int, resolve bool) any
	// ESetFromID sets the value of the feature of ID featureID (ESet).
	ESetFromID(featureID int, newValue any)
	// EUnsetFromID unsets the feature of ID featureID (EUnset).
	EUnsetFromID(featureID int)
	// EIsSetFromID returns true when the feature of ID featureID is set (EIsSet).
	EIsSetFromID(featureID int) bool
	// EInvokeFromID invokes the operation of ID operationID with arguments (EInvoke).
	EInvokeFromID(operationID int, arguments EList) any

	// EBasicInverseAdd adds otherEnd to the feature of ID featureID (>= 0) of the object, without
	// updating otherEnd: the inverse side of EInverseAdd. It returns the notifications to dispatch.
	EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain
	// EBasicInverseRemove removes otherEnd from the feature of ID featureID (>= 0) of the object,
	// without updating otherEnd: the inverse side of EInverseRemove. It returns the notifications
	// to dispatch.
	EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain

	// EObjectForFragmentSegment returns the object a segment of a URI fragment designates among
	// the contents of the object: "@feature" or "@feature.index" (see EURIFragmentSegment).
	EObjectForFragmentSegment(string) EObject
	// EURIFragmentSegment returns the segment of URI fragment designating object in the feature of
	// the object: "@feature", or "@feature.index" for a many feature.
	EURIFragmentSegment(EStructuralFeature, EObject) string

	// EProxyURI returns the URI of the object it is a proxy of, nil when it is not a proxy.
	EProxyURI() *URI
	// ESetProxyURI sets the URI of the object it is a proxy of: a nil uri makes it a non-proxy.
	ESetProxyURI(uri *URI)
	// EResolveProxy returns the object proxy designates, resolved in the resource set of the
	// object, or proxy when it cannot be resolved.
	EResolveProxy(proxy EObject) EObject
}

// AsEObject returns the object as an EObject: its interfaces, the embedding struct.
func (o *AbstractEObject) AsEObject() EObject {
	return o.interfaces.(EObject)
}

// AsEObjectInternal returns the object as an EObjectInternal: its interfaces, the embedding struct.
func (o *AbstractEObject) AsEObjectInternal() EObjectInternal {
	return o.interfaces.(EObjectInternal)
}

// EClass returns the class of the object, its static class (EStaticClass).
func (o *AbstractEObject) EClass() EClass {
	return o.AsEObjectInternal().EStaticClass()
}

// EStaticClass returns the class the object is generated for: EObject, the generated classes
// override it.
func (o *AbstractEObject) EStaticClass() EClass {
	return GetPackage().GetEObject()
}

// EStaticFeatureCount returns the number of features of the static class (EStaticClass).
func (o *AbstractEObject) EStaticFeatureCount() int {
	return o.AsEObjectInternal().EStaticClass().GetFeatureCount()
}

// EDynamicProperties returns nil: the object has no dynamic features (see EDynamicProperties).
func (o *AbstractEObject) EDynamicProperties() EDynamicProperties {
	return nil
}

// EContainer returns the container of the object. When the container is a proxy, it resolves it,
// replaces it by the resolved object and notifies the RESOLVE when the containment has an opposite.
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

// EContainerFeatureID returns the feature ID of the container of the object
// (EInternalContainerFeatureID).
func (o *AbstractEObject) EContainerFeatureID() int {
	return o.AsEObjectInternal().EInternalContainerFeatureID()
}

// EResource returns the resource of the object: the resource it is directly in, else the resource
// of its container.
func (o *AbstractEObject) EResource() EResource {
	resource := o.AsEObjectInternal().EInternalResource()
	if resource == nil {
		if container := o.AsEObjectInternal().EInternalContainer(); container != nil && container != o.AsEObject() {
			resource = container.EResource()
		}
	}
	return resource
}

// ESetResource makes newResource the resource the object is directly in (see
// EObjectInternal.ESetResource) and returns the notifications to dispatch.
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

// EContainingFeature returns the feature of the container holding the object, nil when the object
// has no container.
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

// EContainmentFeature returns the containment reference of the container holding the object, nil
// when the object has no container. It panics when the reference cannot be located.
func (o *AbstractEObject) EContainmentFeature() EReference {
	objectInternal := o.AsEObjectInternal()
	return eContainmentFeature(objectInternal, objectInternal.EInternalContainer(), objectInternal.EInternalContainerFeatureID())
}

// eContainmentFeature returns the containment reference holding o in container, from the feature
// ID of the container of o (see EObjectInternal.EInternalContainerFeatureID), nil when container
// is nil. It panics when the reference cannot be located.
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

// EAllContents returns an iterator on the contents of the object, recursively.
func (o *AbstractEObject) EAllContents() EIterator {
	return newEAllContentsIterator(o.AsEObject())
}

// EFeatureID returns the ID of feature in the class of the object (EDerivedFeatureID). It panics
// when the class has no such feature.
func (o *AbstractEObject) EFeatureID(feature EStructuralFeature) int {
	if !o.AsEObject().EClass().GetEAllStructuralFeatures().Contains(feature) {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
	return o.AsEObjectInternal().EDerivedFeatureID(feature.EContainer(), feature.GetFeatureID())
}

// EDerivedFeatureID returns featureID: the generated classes override it for the features of
// their mixins.
func (o *AbstractEObject) EDerivedFeatureID(container EObject, featureID int) int {
	return featureID
}

// EOperationID returns the ID of operation in the class of the object (EDerivedOperationID). It
// panics when the class has no such operation.
func (o *AbstractEObject) EOperationID(operation EOperation) int {
	if !o.AsEObject().EClass().GetEAllOperations().Contains(operation) {
		panic("The operation '" + operation.GetName() + "' is not a valid operation")
	}
	return o.AsEObjectInternal().EDerivedOperationID(operation.EContainer(), operation.GetOperationID())
}

// EDerivedOperationID returns operationID: the generated classes override it for the operations
// of their mixins.
func (o *AbstractEObject) EDerivedOperationID(container EObject, operationID int) int {
	return operationID
}

// EGet returns the value of feature, resolving the proxies (EGetFromID). It panics when the class
// of the object has no such feature.
func (o *AbstractEObject) EGet(feature EStructuralFeature) any {
	return o.eGetFromFeature(feature, true)
}

// EGetResolve returns the value of feature, resolving the proxies when resolve is true
// (EGetFromID). It panics when the class of the object has no such feature.
func (o *AbstractEObject) EGetResolve(feature EStructuralFeature, resolve bool) any {
	return o.eGetFromFeature(feature, resolve)
}

// eGetFromFeature returns the value of feature (EGetFromID). It panics when the class of the object
// has no such feature.
func (o *AbstractEObject) eGetFromFeature(feature EStructuralFeature, resolve bool) any {
	featureID := o.AsEObjectInternal().EFeatureID(feature)
	if featureID >= 0 {
		return o.AsEObjectInternal().EGetFromID(featureID, resolve)
	}
	panic("The feature '" + feature.GetName() + "' is not a valid feature")
}

// EGetFromID returns the value of the feature of ID featureID: the value of a dynamic feature,
// stored in EDynamicProperties, or EGetResolve for a static feature, whose value the generated
// classes return by overriding EGetFromID. It panics for an invalid featureID, or a dynamic
// feature of an object without dynamic properties.
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

// eDynamicPropertiesGet returns the value of a dynamic feature: the container for a container
// reference, else the stored value, creating the list or the map of a many feature, the default
// value of an unset feature, and resolving a proxy when resolve is true.
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

// eDynamicPropertiesCreateMap returns the map of a dynamic map feature.
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

// eDynamicPropertiesCreateList returns the list of a dynamic many feature: a list of data values
// for an attribute, a list of objects for a reference, with its containment and opposite.
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

// ESet sets the value of feature (ESetFromID). It panics when the class of the object has no such
// feature.
func (o *AbstractEObject) ESet(feature EStructuralFeature, newValue any) {
	featureID := o.AsEObjectInternal().EFeatureID(feature)
	if featureID >= 0 {
		o.AsEObjectInternal().ESetFromID(featureID, newValue)
	} else {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
}

// ESetFromID sets the value of the feature of ID featureID: a dynamic feature in
// EDynamicProperties, updating the opposite or the container and notifying, or ESet for a static
// feature, whose value the generated classes set by overriding ESetFromID. It panics for an
// invalid featureID, or a dynamic feature of an object without dynamic properties.
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

// eDynamicPropertiesSet sets the value of a dynamic feature: the container for a container
// reference, updating the opposite or the containment of the old and the new value, and notifies.
func (o *AbstractEObject) eDynamicPropertiesSet(properties EDynamicProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int, newValue any) {
	checkSameImplementation(o.interfaces, newValue)
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

// EIsSet returns true when feature is set (EIsSetFromID). It panics when the class of the object
// has no such feature.
func (o *AbstractEObject) EIsSet(feature EStructuralFeature) bool {
	featureID := o.AsEObjectInternal().EFeatureID(feature)
	if featureID >= 0 {
		return o.AsEObjectInternal().EIsSetFromID(featureID)
	}
	panic("The feature '" + feature.GetName() + "' is not a valid feature")
}

// EIsSetFromID returns true when the feature of ID featureID is set: a dynamic feature in
// EDynamicProperties, or EIsSet for a static feature, overridden by the generated classes. It
// panics for an invalid featureID, or a dynamic feature of an object without dynamic properties.
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

// eDynamicPropertiesIsSet returns true when a dynamic feature is set: for a container reference,
// when the object has a container through it.
func (o *AbstractEObject) eDynamicPropertiesIsSet(properties EDynamicProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int) bool {
	if IsContainer(dynamicFeature) {
		objInternal := o.AsEObjectInternal()
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		return objInternal.EInternalContainerFeatureID() == featureID && objInternal.EInternalContainer() != nil
	} else {
		return properties.EDynamicIsSet(dynamicFeatureID)
	}
}

// EUnset unsets feature (EUnsetFromID). It panics when the class of the object has no such feature.
func (o *AbstractEObject) EUnset(feature EStructuralFeature) {
	featureID := o.AsEObjectInternal().EFeatureID(feature)
	if featureID >= 0 {
		o.AsEObjectInternal().EUnsetFromID(featureID)
	} else {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
}

// EUnsetFromID unsets the feature of ID featureID: a dynamic feature in EDynamicProperties,
// updating the opposite or the container and notifying, or EUnset for a static feature, overridden
// by the generated classes. It panics for an invalid featureID, or a dynamic feature of an object
// without dynamic properties.
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

// eDynamicPropertiesUnset unsets a dynamic feature, updating the opposite or the containment of the
// old value, and notifies.
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

// EInvoke invokes operation with arguments (EInvokeFromID). It panics when the class of the object
// has no such operation.
func (o *AbstractEObject) EInvoke(operation EOperation, arguments EList) any {
	operationID := o.AsEObjectInternal().EOperationID(operation)
	if operationID >= 0 {
		return o.AsEObjectInternal().EInvokeFromID(operationID, arguments)
	}
	panic("The operation '" + operation.GetName() + "' is not a valid operation")
}

// EInvokeFromID returns nil: the generated classes override it to invoke their operations. It
// panics for an invalid operationID.
func (o *AbstractEObject) EInvokeFromID(operationID int, arguments EList) any {
	operation := o.AsEObject().EClass().GetEOperation(operationID)
	if operation == nil {
		panic("Invalid operationID: " + strconv.Itoa(operationID))
	}
	return nil
}

// EInverseAdd updates the object when the opposite of one of its features now references otherEnd
// (see EObjectInternal.EInverseAdd) and returns the notifications to dispatch.
func (o *AbstractEObject) EInverseAdd(otherEnd EObject, featureID int, n ENotificationChain) ENotificationChain {
	notifications := n
	if featureID >= 0 {
		return o.AsEObjectInternal().EBasicInverseAdd(otherEnd, featureID, notifications)
	} else {
		notifications = o.EBasicRemoveFromContainer(notifications)
		return o.EBasicSetContainer(otherEnd, featureID, notifications)
	}
}

// EInverseRemove updates the object when the opposite of one of its features no longer references
// otherEnd (see EObjectInternal.EInverseRemove) and returns the notifications to dispatch.
func (o *AbstractEObject) EInverseRemove(otherEnd EObject, featureID int, n ENotificationChain) ENotificationChain {
	if featureID >= 0 {
		return o.AsEObjectInternal().EBasicInverseRemove(otherEnd, featureID, n)
	} else {
		return o.EBasicSetContainer(nil, featureID, n)
	}
}

// EResolveProxy returns the object proxy designates, resolved in the resource set of the object,
// or proxy when it cannot be resolved.
func (o *AbstractEObject) EResolveProxy(proxy EObject) EObject {
	return ResolveInObject(proxy, o.GetInterfaces().(EObject))
}

// EBasicInverseAdd adds otherEnd to the dynamic feature of ID featureID, without updating
// otherEnd, and returns the notifications to dispatch: the generated classes override it for their
// static features. It panics for a dynamic feature of an object without dynamic properties.
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

// eDynamicPropertiesInverseAdd adds otherEnd to a dynamic feature, without updating otherEnd: to
// the list of a many feature, as the container for a container reference, else as the value,
// removing the old value from its opposite.
func (o *AbstractEObject) eDynamicPropertiesInverseAdd(properties EDynamicProperties, otherEnd EObject, dynamicFeature EStructuralFeature, dynamicFeatureID int, notifications ENotificationChain) ENotificationChain {
	if dynamicFeature.IsMany() {
		value := properties.EDynamicGet(dynamicFeatureID)
		if value == nil {
			value = o.eDynamicPropertiesCreateList(dynamicFeature)
			properties.EDynamicSet(dynamicFeatureID, value)
		}
		if list, _ := value.(ENotifyingList); list != nil {
			return list.AddWithNotification(otherEnd, notifications)
		}
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

// EBasicInverseRemove removes otherEnd from the dynamic feature of ID featureID, without updating
// otherEnd, and returns the notifications to dispatch: the generated classes override it for their
// static features. It panics for a dynamic feature of an object without dynamic properties.
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

// eDynamicPropertiesInverseRemove removes otherEnd from a dynamic feature, without updating
// otherEnd: from the list of a many feature, as the container for a container reference, else
// unsetting the value.
func (o *AbstractEObject) eDynamicPropertiesInverseRemove(properties EDynamicProperties, otherEnd EObject, dynamicFeature EStructuralFeature, dynamicFeatureID int, notifications ENotificationChain) ENotificationChain {
	if dynamicFeature.IsMany() {
		value := properties.EDynamicGet(dynamicFeatureID)
		if value != nil {
			if list, _ := value.(ENotifyingList); list != nil {
				return list.RemoveWithNotification(otherEnd, notifications)
			}
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

// EBasicSetContainer sets the container of the object and its feature ID, without updating the
// old or the new container: it removes the object from the resource it is directly in when the new
// containment does not resolve proxies, notifies the resources the object leaves or joins
// (Detached, Attached) and adds the SET notifications of the container features to the returned
// notifications.
func (o *AbstractEObject) EBasicSetContainer(newContainer EObject, newContainerFeatureID int, n ENotificationChain) ENotificationChain {
	checkSameImplementation(newContainer, o.interfaces)
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

// EBasicRemoveFromContainer removes the object from the feature of its container holding it
// (EInverseRemove of the container) and returns the notifications to dispatch. The container of the
// object is left unchanged (EBasicSetContainer).
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

// EBasicRemoveFromContainerFeature removes the object from the containment reference of its
// container, the opposite of the container reference of the object, and returns the notifications
// to dispatch.
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

// eStructuralFeature returns the feature of the class of the object named featureName. It panics
// when there is no such feature.
func (o *AbstractEObject) eStructuralFeature(featureName string) EStructuralFeature {
	eFeature := o.AsEObject().EClass().GetEStructuralFeatureFromName(featureName)
	if eFeature == nil {
		panic("The feature " + featureName + " is not a valid feature")
	}
	return eFeature
}

// EObjectForFragmentSegment returns the object a segment of a URI fragment designates among the
// contents of the object: "@feature" or "@feature.index", nil when there is no such object. It
// panics when the segment does not start with '@' or names no feature of the object.
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

// EURIFragmentSegment returns the segment of URI fragment designating object in feature of the
// object: "@feature", or "@feature.index" for a many feature.
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
