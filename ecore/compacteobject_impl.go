package ecore

import (
	"math/bits"
	"net/url"
)

type CompactEObjectImpl struct {
	AbstractEObject
	flags   uint
	storage interface{}
}

const (
	deliver_flag    uint = 1 << 0
	container_flag  uint = 1 << 1
	resource_flag   uint = 1 << 2
	adapters_flag   uint = 1 << 3
	proxy_flag      uint = 1 << 4
	class_flag      uint = 1 << 5
	contents_flag   uint = 1 << 6
	cross_flag      uint = 1 << 7
	properties_flag uint = 1 << 8
	fields_mask     uint = deliver_flag | container_flag | resource_flag | adapters_flag | proxy_flag | class_flag | contents_flag | cross_flag | properties_flag
	first_flag      uint = container_flag
	last_flag       uint = properties_flag
)

func (o *CompactEObjectImpl) Initialize() {
	o.flags = deliver_flag
}

func (o *CompactEObjectImpl) hasField(field uint) bool {
	return (o.flags & field) != 0
}

func (o *CompactEObjectImpl) getField(field uint) interface{} {
	if o.hasField(field) {
		if fieldIndex := o.fieldIndex(field); fieldIndex == -1 {
			return o.storage
		} else {
			return o.storage.([]interface{})[fieldIndex]
		}
	} else {
		return nil
	}
}

func (o *CompactEObjectImpl) setField(field uint, value interface{}) {
	if o.hasField(field) {
		if value == nil {
			o.removeField(field)
		} else {
			if fieldIndex := o.fieldIndex(field); fieldIndex == -1 {
				o.storage = value
			} else {
				o.storage.([]interface{})[fieldIndex] = value
			}
		}
	} else if value != nil {
		o.addField(field, value)
	}
}

func (o *CompactEObjectImpl) fieldIndex(field uint) int {
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

func (o *CompactEObjectImpl) addField(field uint, value interface{}) {
	if fieldCount := bits.OnesCount(o.flags & fields_mask); fieldCount == 0 {
		o.storage = value
	} else if fieldCount == 1 {
		if fieldIndex := o.fieldIndex(field); fieldIndex == 0 {
			o.storage = []interface{}{value, o.storage}
		} else {
			o.storage = []interface{}{o.storage, value}
		}
	} else {
		result := make([]interface{}, fieldCount+1)
		storage := o.storage.([]interface{})
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

func (o *CompactEObjectImpl) removeField(field uint) {
	if fieldCount := bits.OnesCount(o.flags & fields_mask); fieldCount == 1 {
		o.storage = nil
	} else if fieldCount == 2 {
		storage := o.storage.([]interface{})
		if fieldIndex := o.fieldIndex(field); fieldIndex == 0 {
			o.storage = storage[1]
		} else {
			o.storage = storage[0]
		}
	} else {
		result := make([]interface{}, fieldCount-1)
		storage := o.storage.([]interface{})
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
	return (o.flags & deliver_flag) != 0
}

func (o *CompactEObjectImpl) ESetDeliver(deliver bool) {
	if deliver {
		o.flags |= deliver_flag
	} else {
		o.flags &= ^deliver_flag
	}
}

func (o *CompactEObjectImpl) EAdapters() EList {
	adapters := o.getField(adapters_flag)
	if adapters == nil {
		adapters = newNotifierAdapterList(&o.AbstractENotifier)
		o.setField(adapters_flag, adapters)
	}
	return adapters.(EList)
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
func (o *CompactEObjectImpl) EProxyURI() *url.URL {
	if proxyURI := o.getField(proxy_flag); proxyURI != nil {
		return proxyURI.(*url.URL)
	}
	return nil
}

// ESetProxyURI ...
func (o *CompactEObjectImpl) ESetProxyURI(uri *url.URL) {
	o.setField(proxy_flag, uri)
}

// EContents ...
func (o *CompactEObjectImpl) EContents() EList {
	contents := o.getField(contents_flag)
	if contents == nil {
		contents = newContentsListAdapter(&o.AbstractEObject, func(eClass EClass) EList { return eClass.GetEContainmentFeatures() })
		o.setField(contents_flag, contents)
	}
	return contents.(EList)
}

// ECrossReferences ...
func (o *CompactEObjectImpl) ECrossReferences() EList {
	crossReferenceS := o.getField(cross_flag)
	if crossReferenceS == nil {
		crossReferenceS = newContentsListAdapter(&o.AbstractEObject, func(eClass EClass) EList { return eClass.GetECrossReferenceFeatures() })
		o.setField(cross_flag, crossReferenceS)
	}
	return crossReferenceS.(EList)
}

// ESetContainer ...
func (o *CompactEObjectImpl) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.setField(container_flag, newContainer)
	o.flags = uint(newContainerFeatureID)<<16 | (o.flags & 0x00FF)
}

func (o *CompactEObjectImpl) EInternalContainer() EObject {
	if container := o.getField(container_flag); container != nil {
		return container.(EObject)
	}
	return nil
}

func (o *CompactEObjectImpl) EInternalContainerFeatureID() int {
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
