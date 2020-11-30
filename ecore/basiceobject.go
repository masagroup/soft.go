// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

import (
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func newContentsList(adapter *contentsListAdapter, resolve bool) *immutableEList {
	data := []interface{}{}
	o := adapter.obj
	features := adapter.getFeaturesFn(o.EClass())
	for it := features.Iterator(); it.HasNext(); {
		feature := it.Next().(EStructuralFeature)
		if o.EIsSet(feature) {
			value := o.EGetResolve(feature, resolve)
			if feature.IsMany() {
				l := value.(EList)
				data = append(data, l.ToArray()...)
			} else if value != nil {
				data = append(data, value)
			}
		}
	}
	return NewImmutableEList(data)
}

// An unresolved content list
type unResolvedContentsList struct {
	*immutableEList
}

func newUnResolvedContentsList(adapter *contentsListAdapter) *unResolvedContentsList {
	l := new(unResolvedContentsList)
	l.immutableEList = newContentsList(adapter, false)
	return l
}

func (l *unResolvedContentsList) GetUnResolvedList() EList {
	return l
}

// An resolved content list
type resolvedContentsList struct {
	*immutableEList
	adapter *contentsListAdapter
}

func newResolvedContentsList(adapter *contentsListAdapter) *resolvedContentsList {
	l := new(resolvedContentsList)
	l.immutableEList = newContentsList(adapter, true)
	l.adapter = adapter
	return l
}

func (l *resolvedContentsList) GetUnResolvedList() EList {
	return newUnResolvedContentsList(l.adapter)
}

// listen to object features modifications & maintain a list of object contents
type contentsListAdapter struct {
	*Adapter
	obj           *BasicEObject
	getFeaturesFn func(EClass) EList
	list          EList
}

func newContentsListAdapter(obj *BasicEObject, getFeaturesFn func(EClass) EList) *contentsListAdapter {
	a := new(contentsListAdapter)
	a.Adapter = NewAdapter()
	a.obj = obj
	a.getFeaturesFn = getFeaturesFn
	obj.EAdapters().Add(a)
	return a
}

func (a *contentsListAdapter) NotifyChanged(notification ENotification) {
	if a.list != nil {
		feature := notification.GetFeature()
		features := a.getFeaturesFn(a.obj.EClass())
		if features.Contains(feature) {
			a.list = nil
		}
	}
}

func (a *contentsListAdapter) GetList() EList {
	if a.list == nil {
		a.list = newResolvedContentsList(a)
	}
	return a.list
}

// BasicEObject is a basic implementation of an EObject
type BasicEObject struct {
	*BasicNotifier
	resource           EResource
	container          EObject
	containerFeatureID int
	proxyURI           *url.URL
	contents           *contentsListAdapter
	crossReferenceS    *contentsListAdapter
}

type EObjectProperties interface {
	EDynamicGet(dynamicFeatureID int) interface{}
	EDynamicSet(dynamicFeatureID int, newValue interface{})
	EDynamicUnset(dynamicFeatureID int)
}

// EObjectInternal ...
type EObjectInternal interface {
	EObject

	EProperties() EObjectProperties

	EStaticClass() EClass
	EStaticFeatureCount() int

	EInternalContainer() EObject
	EInternalResource() EResource

	ESetResource(resource EResource, notifications ENotificationChain) ENotificationChain

	EInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain
	EInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain

	EDerivedFeatureID(container EObject, featureID int) int
	EDerivedOperationID(container EObject, operationID int) int
	EGetFromID(featureID int, resolve bool) interface{}
	ESetFromID(featureID int, newValue interface{})
	EUnsetFromID(featureID int)
	EIsSetFromID(featureID int) bool
	EInvokeFromID(operationID int, arguments EList) interface{}

	EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain
	EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain

	EObjectForFragmentSegment(string) EObject
	EURIFragmentSegment(EStructuralFeature, EObject) string

	EProxyURI() *url.URL
	ESetProxyURI(uri *url.URL)
	EResolveProxy(proxy EObject) EObject
}

// NewBasicEObject is BasicEObject constructor
func NewBasicEObject() *BasicEObject {
	o := new(BasicEObject)
	o.BasicNotifier = NewBasicNotifier()
	o.interfaces = o
	o.containerFeatureID = -1
	return o
}

// AsEObject ...
func (o *BasicEObject) AsEObject() EObject {
	return o.interfaces.(EObject)
}

// AsEObjectInternal ...
func (o *BasicEObject) AsEObjectInternal() EObjectInternal {
	return o.interfaces.(EObjectInternal)
}

// EClass ...
func (o *BasicEObject) EClass() EClass {
	return o.AsEObjectInternal().EStaticClass()
}

// EStaticClass ...
func (o *BasicEObject) EStaticClass() EClass {
	return GetPackage().GetEObject()
}

func (o *BasicEObject) EStaticFeatureCount() int {
	return o.AsEObjectInternal().EStaticClass().GetFeatureCount()
}

func (o *BasicEObject) EProperties() EObjectProperties {
	return nil
}

// EIsProxy ...
func (o *BasicEObject) EIsProxy() bool {
	return o.proxyURI != nil
}

func (o *BasicEObject) EInternalContainer() EObject {
	return o.container
}

// EContainer ...
func (o *BasicEObject) EContainer() EObject {
	eContainer := o.container
	if eContainer != nil && eContainer.EIsProxy() {
		resolved := o.EResolveProxy(eContainer)
		if resolved != eContainer {
			notifications := o.EBasicRemoveFromContainer(nil)
			o.container = resolved
			if notifications != nil {
				notifications.Dispatch()
			}

			if o.ENotificationRequired() && o.containerFeatureID >= EOPPOSITE_FEATURE_BASE {
				o.ENotify(NewNotificationByFeatureID(o.AsEObject(), RESOLVE, o.containerFeatureID, eContainer, resolved, -1))
			}
		}
		return resolved
	}
	return eContainer
}

// EContainerFeatureID ...
func (o *BasicEObject) EContainerFeatureID() int {
	return o.containerFeatureID
}

// EResource ...
func (o *BasicEObject) EResource() EResource {
	resource := o.resource
	if resource == nil {
		if o.container != nil {
			resource = o.container.EResource()
		}
	}
	return resource
}

// EInternalResource ...
func (o *BasicEObject) EInternalResource() EResource {
	return o.resource
}

// ESetInternalResource ...
func (o *BasicEObject) ESetInternalResource(resource EResource) {
	o.resource = resource
}

// ESetResource ...
func (o *BasicEObject) ESetResource(newResource EResource, n ENotificationChain) ENotificationChain {
	notifications := n
	oldResource := o.EInternalResource()
	// When setting the resource to nil we assume that detach has already been called in the resource implementation
	if oldResource != nil && newResource != nil {
		list := oldResource.GetContents().(ENotifyingList)
		notifications = list.RemoveWithNotification(o.AsEObject(), notifications)
		oldResource.Detached(o.AsEObject())
	}

	eContainer := o.container
	if eContainer != nil {
		if o.EContainmentFeature().IsResolveProxies() {
			if eContainerInternal, _ := eContainer.(EObjectInternal); eContainerInternal != nil {
				oldContainerResource := eContainerInternal.EInternalResource()
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
	o.ESetInternalResource(newResource)
	return notifications
}

// EContainingFeature ...
func (o *BasicEObject) EContainingFeature() EStructuralFeature {
	if o.container != nil {
		if o.containerFeatureID <= EOPPOSITE_FEATURE_BASE {
			return o.container.EClass().GetEStructuralFeature(EOPPOSITE_FEATURE_BASE - o.containerFeatureID)
		} else {
			return o.AsEObject().EClass().GetEStructuralFeature(o.containerFeatureID).(EReference).GetEOpposite()
		}
	}
	return nil
}

// EContainmentFeature ...
func (o *BasicEObject) EContainmentFeature() EReference {
	return eContainmentFeature(o.AsEObject(), o.container, o.containerFeatureID)
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

// EContents ...
func (o *BasicEObject) EContents() EList {
	if o.contents == nil {
		o.contents = newContentsListAdapter(o, func(eClass EClass) EList { return eClass.GetEContainmentFeatures() })
	}
	return o.contents.GetList()
}

// ECrossReferences ...
func (o *BasicEObject) ECrossReferences() EList {
	if o.crossReferenceS == nil {
		o.crossReferenceS = newContentsListAdapter(o, func(eClass EClass) EList { return eClass.GetECrossReferenceFeatures() })
	}
	return o.crossReferenceS.GetList()
}

// EAllContents ...
func (o *BasicEObject) EAllContents() EIterator {
	return newEAllContentsIterator(o)
}

func (o *BasicEObject) eFeatureID(feature EStructuralFeature) int {
	if !o.AsEObject().EClass().GetEAllStructuralFeatures().Contains(feature) {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
	return o.AsEObjectInternal().EDerivedFeatureID(feature.EContainer(), feature.GetFeatureID())
}

func (o *BasicEObject) EDerivedFeatureID(container EObject, featureID int) int {
	return featureID
}

func (o *BasicEObject) eOperationID(operation EOperation) int {
	if !o.AsEObject().EClass().GetEAllOperations().Contains(operation) {
		panic("The operation '" + operation.GetName() + "' is not a valid operation")
	}
	return o.AsEObjectInternal().EDerivedOperationID(operation.EContainer(), operation.GetOperationID())
}

func (o *BasicEObject) EDerivedOperationID(container EObject, operationID int) int {
	return operationID
}

// EGet ...
func (o *BasicEObject) EGet(feature EStructuralFeature) interface{} {
	return o.eGetFromFeature(feature, true)
}

// EGetResolve ...
func (o *BasicEObject) EGetResolve(feature EStructuralFeature, resolve bool) interface{} {
	return o.eGetFromFeature(feature, resolve)
}

func (o *BasicEObject) eGetFromFeature(feature EStructuralFeature, resolve bool) interface{} {
	featureID := o.eFeatureID(feature)
	if featureID >= 0 {
		return o.AsEObjectInternal().EGetFromID(featureID, resolve)
	}
	panic("The feature '" + feature.GetName() + "' is not a valid feature")
}

// EGetFromID ...
func (o *BasicEObject) EGetFromID(featureID int, resolve bool) interface{} {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	if feature == nil {
		panic("Invalid featureID: " + strconv.Itoa(featureID))
	}
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID < 0 {
		return o.AsEObjectInternal().EGetResolve(feature, resolve)
	} else {
		properties := o.AsEObjectInternal().EProperties()
		if properties != nil {
			return o.eDynamicPropertiesGet(properties, feature, dynamicFeatureID, resolve)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
}

func (o *BasicEObject) eDynamicPropertiesGet(properties EObjectProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int, resolve bool) interface{} {
	if isContainer(dynamicFeature) {
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		if o.EContainerFeatureID() == featureID {
			if resolve {
				return o.EContainer()
			} else {
				return o.EInternalContainer()
			}
		}
	} else {
		// retrieve value or compute it if empty
		result := properties.EDynamicGet(dynamicFeatureID)
		if result == nil {
			if dynamicFeature.IsMany() {
				result = o.eDynamicPropertiesCreateList(dynamicFeature)
				properties.EDynamicSet(dynamicFeatureID, result)
			}
		} else if resolve && isProxy(dynamicFeature) {
			oldValue, _ := result.(EObject)
			newValue := o.EResolveProxy(oldValue)
			result = newValue
			if oldValue != newValue {
				properties.EDynamicSet(dynamicFeatureID, newValue)
				if isContains(dynamicFeature) {
					var notifications ENotificationChain
					if !isBidirectional(dynamicFeature) {
						featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
						if oldValue != nil {
							oldObject := oldValue.(EObjectInternal)
							notifications = oldObject.EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
						}
						if newValue != nil {
							newObject := newValue.(EObjectInternal)
							notifications = newObject.EInverseAdd(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
						}
					} else {
						dynamicReference := dynamicFeature.(EReference)
						reverseFeature := dynamicReference.GetEOpposite()
						if oldValue != nil {
							oldObject := oldValue.(EObjectInternal)
							featureID := oldObject.EClass().GetFeatureID(reverseFeature)
							notifications = oldObject.EInverseRemove(o.AsEObject(), featureID, notifications)
						}
						if newValue != nil {
							newObject := newValue.(EObjectInternal)
							featureID := newObject.EClass().GetFeatureID(reverseFeature)
							notifications = newValue.(EObjectInternal).EInverseAdd(o.AsEObject(), featureID, notifications)
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
		return result
	}
	return nil
}

func (o *BasicEObject) eDynamicPropertiesCreateList(feature EStructuralFeature) EList {
	if attribute, isAttribute := feature.(EAttribute); isAttribute {
		if attribute.IsUnique() {
			return NewUniqueBasicEList(nil)
		} else {
			return NewBasicEList(nil)
		}
	} else if ref, isRef := feature.(EReference); isRef {
		inverse := false
		opposite := false
		containment := ref.IsContainment()
		reverseID := -1
		reverseFeature := ref.GetEOpposite()
		if containment {
			if reverseFeature != nil {
				reverseID = reverseFeature.GetFeatureID()
				inverse = true
				opposite = true
			} else {
				inverse = true
				opposite = false
			}
		} else {
			if reverseFeature != nil {
				reverseID = reverseFeature.GetFeatureID()
				inverse = true
				opposite = true
			} else {
				inverse = false
				opposite = false
			}
		}
		return NewBasicEObjectList(o.AsEObjectInternal(), ref.GetFeatureID(), reverseID, containment, inverse, opposite, ref.EIsProxy(), ref.IsUnsettable())
	}
	return nil
}

// ESet ...
func (o *BasicEObject) ESet(feature EStructuralFeature, newValue interface{}) {
	featureID := o.eFeatureID(feature)
	if featureID >= 0 {
		o.AsEObjectInternal().ESetFromID(featureID, newValue)
	} else {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
}

// ESetFromID ...
func (o *BasicEObject) ESetFromID(featureID int, newValue interface{}) {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	if feature == nil {
		panic("Invalid featureID: " + strconv.Itoa(featureID))
	}
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID < 0 {
		o.ESet(feature, newValue)
	} else {
		properties := o.AsEObjectInternal().EProperties()
		if properties != nil {
			o.eDynamicPropertiesSet(properties, feature, dynamicFeatureID, newValue)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
}

func (o *BasicEObject) eDynamicPropertiesSet(properties EObjectProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int, newValue interface{}) {
	if isContainer(dynamicFeature) {
		// container
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		newContainer, _ := newValue.(EObject)
		if newContainer != o.EInternalContainer() || (newContainer != nil && o.EContainerFeatureID() != featureID) {
			var notifications ENotificationChain
			if o.EInternalContainer() != nil {
				notifications = o.EBasicRemoveFromContainer(notifications)
			}
			if newContainer != nil {
				reverseFeature := dynamicFeature.(EReference).GetEOpposite()
				featureID := newContainer.EClass().GetFeatureID(reverseFeature)
				notifications = newContainer.(EObjectInternal).EInverseAdd(o.AsEObject(), featureID, notifications)
			}
			notifications = o.EBasicSetContainer(newContainer, featureID, notifications)
			if notifications != nil {
				notifications.Dispatch()
			}
		} else if o.ENotificationRequired() {
			o.ENotify(NewNotificationByFeature(o.AsEObject(), SET, dynamicFeature, newValue, newValue, NO_INDEX))
		}
	} else if isBidirectional(dynamicFeature) || isContains(dynamicFeature) {
		// inverse - opposite
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		if oldValue != newValue {
			var notifications ENotificationChain
			oldObject, _ := oldValue.(EObject)
			newObject, _ := newValue.(EObject)

			if !isBidirectional(dynamicFeature) {
				featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
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
					featureID := oldObject.EClass().GetFeatureID(reverseFeature)
					notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), featureID, notifications)
				}
				if newObject != nil {
					featureID := newObject.EClass().GetFeatureID(reverseFeature)
					notifications = newObject.(EObjectInternal).EInverseAdd(o.AsEObject(), featureID, notifications)
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
func (o *BasicEObject) EIsSet(feature EStructuralFeature) bool {
	featureID := o.eFeatureID(feature)
	if featureID >= 0 {
		return o.AsEObjectInternal().EIsSetFromID(featureID)
	}
	panic("The feature '" + feature.GetName() + "' is not a valid feature")
}

// EIsSetFromID ...
func (o *BasicEObject) EIsSetFromID(featureID int) bool {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	if feature == nil {
		panic("Invalid featureID: " + strconv.Itoa(featureID))
	}
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID < 0 {
		return o.EIsSet(feature)
	} else {
		properties := o.AsEObjectInternal().EProperties()
		if properties != nil {
			return o.eDynamicPropertiesIsSet(properties, feature, dynamicFeatureID)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
}

func (o *BasicEObject) eDynamicPropertiesIsSet(properties EObjectProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int) bool {
	if isContainer(dynamicFeature) {
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		return o.EContainerFeatureID() == featureID && o.EInternalContainer() != nil
	} else {
		return properties.EDynamicGet(dynamicFeatureID) != nil
	}
}

// EUnset ...
func (o *BasicEObject) EUnset(feature EStructuralFeature) {
	featureID := o.eFeatureID(feature)
	if featureID >= 0 {
		o.AsEObjectInternal().EUnsetFromID(featureID)
	} else {
		panic("The feature '" + feature.GetName() + "' is not a valid feature")
	}
}

// EUnsetFromID ...
func (o *BasicEObject) EUnsetFromID(featureID int) {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	if feature == nil {
		panic("Invalid featureID: " + strconv.Itoa(featureID))
	}
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID < 0 {
		o.EUnset(feature)
	} else {
		properties := o.AsEObjectInternal().EProperties()
		if properties != nil {
			o.eDynamicPropertiesUnset(properties, feature, dynamicFeatureID)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
}

func (o *BasicEObject) eDynamicPropertiesUnset(properties EObjectProperties, dynamicFeature EStructuralFeature, dynamicFeatureID int) {
	if isContainer(dynamicFeature) {
		if o.EInternalContainer() != nil {
			featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
			notifications := o.EBasicRemoveFromContainer(nil)
			notifications = o.EBasicSetContainer(nil, featureID, notifications)
			if notifications != nil {
				notifications.Dispatch()
			}
		} else if o.ENotificationRequired() {
			o.ENotify(NewNotificationByFeature(o.AsEObject(), SET, dynamicFeature, nil, nil, NO_INDEX))
		}
	} else if isBidirectional(dynamicFeature) || isContains(dynamicFeature) {
		// inverse - opposite
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		if oldValue != nil {
			var notifications ENotificationChain
			oldObject, _ := oldValue.(EObject)

			if !isBidirectional(dynamicFeature) {
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
func (o *BasicEObject) EInvoke(operation EOperation, arguments EList) interface{} {
	operationID := o.eOperationID(operation)
	if operationID >= 0 {
		return o.AsEObjectInternal().EInvokeFromID(operationID, arguments)
	}
	panic("The operation '" + operation.GetName() + "' is not a valid operation")
}

// EInvokeFromID ...
func (o *BasicEObject) EInvokeFromID(operationID int, arguments EList) interface{} {
	operation := o.AsEObject().EClass().GetEOperation(operationID)
	if operation == nil {
		panic("Invalid operationID: " + strconv.Itoa(operationID))
	}
	return nil
}

// EInverseAdd ...
func (o *BasicEObject) EInverseAdd(otherEnd EObject, featureID int, n ENotificationChain) ENotificationChain {
	notifications := n
	if featureID >= 0 {
		return o.AsEObjectInternal().EBasicInverseAdd(otherEnd, featureID, notifications)
	} else {
		notifications = o.EBasicRemoveFromContainer(notifications)
		return o.EBasicSetContainer(otherEnd, featureID, notifications)
	}
}

// EInverseRemove ...
func (o *BasicEObject) EInverseRemove(otherEnd EObject, featureID int, n ENotificationChain) ENotificationChain {
	if featureID >= 0 {
		return o.AsEObjectInternal().EBasicInverseRemove(otherEnd, featureID, n)
	} else {
		return o.EBasicSetContainer(nil, featureID, n)
	}
}

// EProxyURI ...
func (o *BasicEObject) EProxyURI() *url.URL {
	return o.proxyURI
}

// ESetProxyURI ...
func (o *BasicEObject) ESetProxyURI(uri *url.URL) {
	o.proxyURI = uri
}

// EResolveProxy ...
func (o *BasicEObject) EResolveProxy(proxy EObject) EObject {
	return ResolveInObject(proxy, o.GetInterfaces().(EObject))
}

// EBasicInverseAdd ...
func (o *BasicEObject) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		properties := o.AsEObjectInternal().EProperties()
		if properties != nil {
			return o.eDynamicPropertiesInverseAdd(properties, otherEnd, feature, dynamicFeatureID, notifications)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
	return notifications
}

func (o *BasicEObject) eDynamicPropertiesInverseAdd(properties EObjectProperties, otherEnd EObject, dynamicFeature EStructuralFeature, dynamicFeatureID int, notifications ENotificationChain) ENotificationChain {
	if dynamicFeature.IsMany() {
		value := properties.EDynamicGet(dynamicFeatureID)
		if value == nil {
			value = o.eDynamicPropertiesCreateList(dynamicFeature)
			properties.EDynamicSet(dynamicFeatureID, value)
		}
		list := value.(ENotifyingList)
		return list.AddWithNotification(otherEnd, notifications)
	} else if isContainer(dynamicFeature) {
		msgs := notifications
		if o.EContainer() != nil {
			msgs = o.EBasicRemoveFromContainer(msgs)
		}
		featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
		return o.EBasicSetContainer(otherEnd, featureID, msgs)
	} else {
		// inverse - opposite
		oldValue := properties.EDynamicGet(dynamicFeatureID)
		oldObject, _ := oldValue.(EObject)
		if oldObject != nil {
			if isContains(dynamicFeature) {
				featureID := o.AsEObject().EClass().GetFeatureID(dynamicFeature)
				notifications = oldObject.(EObjectInternal).EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-featureID, notifications)
			} else if isBidirectional(dynamicFeature) {
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
func (o *BasicEObject) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	feature := o.AsEObject().EClass().GetEStructuralFeature(featureID)
	dynamicFeatureID := featureID - o.AsEObjectInternal().EStaticFeatureCount()
	if dynamicFeatureID >= 0 {
		properties := o.AsEObjectInternal().EProperties()
		if properties != nil {
			return o.eDynamicPropertiesInverseRemove(properties, otherEnd, feature, dynamicFeatureID, notifications)
		} else {
			panic("EObject doesn't define any dynamic properties")
		}
	}
	return notifications
}

func (o *BasicEObject) eDynamicPropertiesInverseRemove(properties EObjectProperties, otherEnd EObject, dynamicFeature EStructuralFeature, dynamicFeatureID int, notifications ENotificationChain) ENotificationChain {
	if dynamicFeature.IsMany() {
		value := properties.EDynamicGet(dynamicFeatureID)
		if value != nil {
			list := value.(ENotifyingList)
			return list.RemoveWithNotification(otherEnd, notifications)
		}
	} else if isContainer(dynamicFeature) {
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

// ESetContainer ...
func (o *BasicEObject) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.container = newContainer
	o.containerFeatureID = newContainerFeatureID
}

// EBasicSetContainer ...
func (o *BasicEObject) EBasicSetContainer(newContainer EObject, newContainerFeatureID int, n ENotificationChain) ENotificationChain {
	notifications := n
	oldResource := o.EInternalResource()
	oldContainer := o.EInternalContainer()
	oldContainerFeatureID := o.containerFeatureID
	oldContainerInternal, _ := oldContainer.(EObjectInternal)
	newContainerInternal, _ := newContainer.(EObjectInternal)

	// resource
	var newResource EResource
	if oldResource != nil {
		if newContainerInternal != nil && !eContainmentFeature(o.AsEObject(), newContainerInternal, newContainerFeatureID).IsResolveProxies() {
			list := oldResource.GetContents().(ENotifyingList)
			notifications = list.RemoveWithNotification(o.AsEObject(), notifications)
			o.ESetInternalResource(nil)
			newResource = newContainerInternal.EInternalResource()
		} else {
			oldResource = nil
		}
	} else {
		if oldContainerInternal != nil {
			oldResource = oldContainerInternal.EInternalResource()
		}

		if newContainerInternal != nil {
			newResource = newContainerInternal.EInternalResource()
		}
	}

	if oldResource != nil && oldResource != newResource {
		oldResource.Detached(o.AsEObject())
	}

	if newResource != nil && newResource != oldResource {
		newResource.Attached(o.AsEObject())
	}

	// internal set
	o.ESetInternalContainer(newContainer, newContainerFeatureID)

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
func (o *BasicEObject) EBasicRemoveFromContainer(notifications ENotificationChain) ENotificationChain {
	if o.containerFeatureID >= 0 {
		return o.EBasicRemoveFromContainerFeature(notifications)
	} else {
		if o.container != nil {
			return o.AsEObjectInternal().EInverseRemove(o.AsEObject(), EOPPOSITE_FEATURE_BASE-o.containerFeatureID, notifications)
		}
	}
	return notifications
}

// EBasicRemoveFromContainerFeature ...
func (o *BasicEObject) EBasicRemoveFromContainerFeature(notifications ENotificationChain) ENotificationChain {
	reference, isReference := o.AsEObject().EClass().GetEStructuralFeature(o.containerFeatureID).(EReference)
	if isReference {
		inverseFeature := reference.GetEOpposite()
		if containerInternal, _ := o.container.(EObjectInternal); containerInternal != nil && inverseFeature != nil {
			return containerInternal.EInverseRemove(o.AsEObject(), inverseFeature.GetFeatureID(), notifications)
		}
	}
	return notifications
}

func (o *BasicEObject) eStructuralFeature(featureName string) EStructuralFeature {
	eFeature := o.AsEObject().EClass().GetEStructuralFeatureFromName(featureName)
	if eFeature == nil {
		panic("The feature " + featureName + " is not a valid feature")
	}
	return eFeature
}

func (o *BasicEObject) EObjectForFragmentSegment(uriSegment string) EObject {

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
			list := o.AsEObject().EGetResolve(eFeature, false).(EList)
			if pos < list.Size() {
				return list.Get(pos).(EObject)
			}
		}
	}
	if index == -1 {
		eFeature := o.eStructuralFeature(uriSegment)
		return o.AsEObject().EGetResolve(eFeature, false).(EObject)
	}
	return nil
}

func (o *BasicEObject) EURIFragmentSegment(feature EStructuralFeature, object EObject) string {
	s := "@"
	s += feature.GetName()
	if feature.IsMany() {
		v := o.AsEObject().EGetResolve(feature, false)
		i := v.(EList).IndexOf(object)
		s += "." + strconv.Itoa(i)
	}
	return s
}
