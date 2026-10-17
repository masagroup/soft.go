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
	"math/bits"
	"sync"
)

// CompactEObjectImpl packs its optional fields into one storage word read with a flags bitmap:
// mutex makes the pair one critical section, lazy list getters included.
type CompactEObjectImpl struct {
	AbstractEObject
	mutex   sync.Mutex
	storage any
	flags   int
}

const (
	deliver_flag    = 1 << 0
	container_flag  = 1 << 1
	resource_flag   = 1 << 2
	adapters_flag   = 1 << 3
	proxy_flag      = 1 << 4
	class_flag      = 1 << 5
	contents_flag   = 1 << 6
	cross_flag      = 1 << 7
	properties_flag = 1 << 8
	fields_mask     = container_flag | resource_flag | adapters_flag | proxy_flag | class_flag | contents_flag | cross_flag | properties_flag
	first_flag      = container_flag
	last_flag       = properties_flag
)

func (o *CompactEObjectImpl) Initialize() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.flags = deliver_flag | (-1)<<16
}

func (o *CompactEObjectImpl) hasField(field int) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.hasFieldLocked(field)
}

func (o *CompactEObjectImpl) getField(field int) any {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return o.getFieldLocked(field)
}

func (o *CompactEObjectImpl) setField(field int, value any) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.setFieldLocked(field, value)
}

// setFieldIfAbsent publishes value unless the slot is set, and returns what the slot holds.
// Lazy getters build outside the lock, so two first readers agree on one published list.
func (o *CompactEObjectImpl) setFieldIfAbsent(field int, value any) any {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if existing := o.getFieldLocked(field); existing != nil {
		return existing
	}
	o.setFieldLocked(field, value)
	return value
}

// The Locked helpers assume o.mutex is held.

func (o *CompactEObjectImpl) hasFieldLocked(field int) bool {
	return (o.flags & field) != 0
}

func (o *CompactEObjectImpl) getFieldLocked(field int) any {
	if o.hasFieldLocked(field) {
		if fieldIndex := o.fieldIndex(field); fieldIndex == -1 {
			return o.storage
		} else {
			return o.storage.([]any)[fieldIndex]
		}
	} else {
		return nil
	}
}

func (o *CompactEObjectImpl) setFieldLocked(field int, value any) {
	if o.hasFieldLocked(field) {
		if value == nil {
			o.removeField(field)
		} else {
			if fieldIndex := o.fieldIndex(field); fieldIndex == -1 {
				o.storage = value
			} else {
				o.storage.([]any)[fieldIndex] = value
			}
		}
	} else if value != nil {
		o.addField(field, value)
	}
}

func (o *CompactEObjectImpl) fieldIndex(field int) int {
	result := 0
	for bit := first_flag; bit < field; bit <<= 1 {
		if (o.flags & bit) != 0 {
			result++
		}
	}
	if result == 0 {
		field <<= 1
		for bit := field; bit <= last_flag; bit <<= 1 {
			if (o.flags & bit) != 0 {
				return 0
			}
		}
		return -1
	} else {
		return result
	}
}

func (o *CompactEObjectImpl) addField(field int, value any) {
	if fieldCount := bits.OnesCount(uint(o.flags & fields_mask)); fieldCount == 0 {
		o.storage = value
	} else if fieldCount == 1 {
		if fieldIndex := o.fieldIndex(field); fieldIndex == 0 {
			o.storage = []any{value, o.storage}
		} else {
			o.storage = []any{o.storage, value}
		}
	} else {
		result := make([]any, fieldCount+1)
		storage := o.storage.([]any)
		for bit, sourceIndex, targetIndex := first_flag, 0, 0; bit <= last_flag; bit <<= 1 {
			if bit == field {
				result[targetIndex] = value
				targetIndex++
			} else if (o.flags & bit) != 0 {
				result[targetIndex] = storage[sourceIndex]
				targetIndex++
				sourceIndex++
			}
		}
		o.storage = result
	}
	o.flags |= field
}

func (o *CompactEObjectImpl) removeField(field int) {
	if fieldCount := bits.OnesCount(uint(o.flags & fields_mask)); fieldCount == 1 {
		o.storage = nil
	} else if fieldCount == 2 {
		storage := o.storage.([]any)
		if fieldIndex := o.fieldIndex(field); fieldIndex == 0 {
			o.storage = storage[1]
		} else {
			o.storage = storage[0]
		}
	} else {
		result := make([]any, fieldCount-1)
		storage := o.storage.([]any)
		for bit, sourceIndex, targetIndex := first_flag, 0, 0; bit <= last_flag; bit <<= 1 {
			if bit == field {
				sourceIndex++
			} else if (o.flags & bit) != 0 {
				result[targetIndex] = storage[sourceIndex]
				targetIndex++
				sourceIndex++
			}
		}
		o.storage = result
	}
	o.flags &= ^field
}

func (o *CompactEObjectImpl) EClass() EClass {
	class := o.getField(class_flag)
	if class != nil {
		return class.(EClass)
	}
	return o.AsEObjectInternal().EStaticClass()
}

func (o *CompactEObjectImpl) SetEClass(class EClass) {
	o.setField(class_flag, class)
}

func (o *CompactEObjectImpl) EDeliver() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return (o.flags & deliver_flag) != 0
}

func (o *CompactEObjectImpl) ESetDeliver(deliver bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if deliver {
		o.flags |= deliver_flag
	} else {
		o.flags &= ^deliver_flag
	}
}

func (o *CompactEObjectImpl) EAdapters() EList {
	if adapters := o.getField(adapters_flag); adapters != nil {
		return adapters.(EList)
	}
	return o.setFieldIfAbsent(adapters_flag, newNotifierAdapterList(&o.AbstractENotifier)).(EList)
}

func (o *CompactEObjectImpl) EBasicHasAdapters() bool {
	return o.hasField(adapters_flag)
}

func (o *CompactEObjectImpl) EBasicAdapters() EList {
	if adapters := o.getField(adapters_flag); adapters != nil {
		return adapters.(EList)
	}
	return nil
}

func (o *CompactEObjectImpl) EIsProxy() bool {
	return o.hasField(proxy_flag)
}

// EProxyURI ...
func (o *CompactEObjectImpl) EProxyURI() *URI {
	if proxyURI := o.getField(proxy_flag); proxyURI != nil {
		return proxyURI.(*URI)
	}
	return nil
}

// ESetProxyURI ...
func (o *CompactEObjectImpl) ESetProxyURI(uri *URI) {
	o.setField(proxy_flag, uri)
}

// EContents ...
func (o *CompactEObjectImpl) EContents() EList {
	if contents := o.getField(contents_flag); contents != nil {
		return contents.(EList)
	}
	// Built outside the lock on purpose: EClass() takes this same mutex.
	eObject := o.AsEObject()
	contents := newEContentsList(eObject, eObject.EClass().GetEContainmentFeatures(), true)
	return o.setFieldIfAbsent(contents_flag, contents).(EList)
}

// ECrossReferences ...
func (o *CompactEObjectImpl) ECrossReferences() EList {
	if crossReferences := o.getField(cross_flag); crossReferences != nil {
		return crossReferences.(EList)
	}
	eObject := o.AsEObject()
	crossReferences := newEContentsList(eObject, eObject.EClass().GetECrossReferenceFeatures(), true)
	return o.setFieldIfAbsent(cross_flag, crossReferences).(EList)
}

// ESetContainer ...
func (o *CompactEObjectImpl) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	// The container and its feature id, packed in flags, are published together.
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.setFieldLocked(container_flag, newContainer)
	o.flags = newContainerFeatureID<<16 | (o.flags & 0x00FF)
}

func (o *CompactEObjectImpl) EInternalContainer() EObject {
	if container := o.getField(container_flag); container != nil {
		return container.(EObject)
	}
	return nil
}

func (o *CompactEObjectImpl) EInternalContainerFeatureID() int {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	return int(o.flags >> 16)
}

// EInternalResource ...
func (o *CompactEObjectImpl) EInternalResource() EResource {
	if resource := o.getField(resource_flag); resource != nil {
		return resource.(EResource)
	}
	return nil
}

// ESetInternalResource ...
func (o *CompactEObjectImpl) ESetInternalResource(resource EResource) {
	o.setField(resource_flag, resource)
}
