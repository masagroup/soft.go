package ecore

import "math/bits"

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
	contents_flag   uint = 1 << 5
	cross_flag      uint = 1 << 6
	properties_flag uint = 1 << 7
	fields_mask     uint = deliver_flag | proxy_flag | adapters_flag | resource_flag | container_flag | properties_flag
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
	} else {
		if fieldCount == 1 {
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
	}
	o.flags |= field
}

func (o *CompactEObjectImpl) removeField(field uint) {
	if fieldCount := bits.OnesCount(o.flags & fields_mask); fieldCount == 0 {
		o.storage = nil
	} else {
		if fieldCount == 2 {
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
	}
	o.flags &= ^field
}
