// Code generated by soft.generator.go. DO NOT EDIT.

// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// eStringToStringMapEntryImpl is the implementation of the model object 'EStringToStringMapEntry'
type eStringToStringMapEntryImpl struct {
	CompactEObjectContainer
	key   string
	value string
}

// newEStringToStringMapEntryImpl is the constructor of a eStringToStringMapEntryImpl
func newEStringToStringMapEntryImpl() *eStringToStringMapEntryImpl {
	eStringToStringMapEntry := new(eStringToStringMapEntryImpl)
	eStringToStringMapEntry.SetInterfaces(eStringToStringMapEntry)
	eStringToStringMapEntry.Initialize()
	return eStringToStringMapEntry
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) Initialize() {
	eStringToStringMapEntry.CompactEObjectContainer.Initialize()
	eStringToStringMapEntry.key = ""
	eStringToStringMapEntry.value = ""

}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) asEStringToStringMapEntry() EStringToStringMapEntry {
	return eStringToStringMapEntry.GetInterfaces().(EStringToStringMapEntry)
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) EStaticClass() EClass {
	return GetPackage().GetEStringToStringMapEntry()
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) EStaticFeatureCount() int {
	return ESTRING_TO_STRING_MAP_ENTRY_FEATURE_COUNT
}
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) GetKey() interface{} {
	return eStringToStringMapEntry.GetTypedKey()
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) SetKey(key interface{}) {
	eStringToStringMapEntry.SetTypedKey(key.(string))
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) GetValue() interface{} {
	return eStringToStringMapEntry.GetTypedValue()
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) SetValue(value interface{}) {
	eStringToStringMapEntry.SetTypedValue(value.(string))
}

// GetTypedKey get the value of key
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) GetTypedKey() string {
	return eStringToStringMapEntry.key
}

// SetTypedKey set the value of key
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) SetTypedKey(newKey string) {
	oldKey := eStringToStringMapEntry.key
	eStringToStringMapEntry.key = newKey
	if eStringToStringMapEntry.ENotificationRequired() {
		eStringToStringMapEntry.ENotify(NewNotificationByFeatureID(eStringToStringMapEntry.AsEObject(), SET, ESTRING_TO_STRING_MAP_ENTRY__KEY, oldKey, newKey, NO_INDEX))
	}
}

// GetTypedValue get the value of value
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) GetTypedValue() string {
	return eStringToStringMapEntry.value
}

// SetTypedValue set the value of value
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) SetTypedValue(newValue string) {
	oldValue := eStringToStringMapEntry.value
	eStringToStringMapEntry.value = newValue
	if eStringToStringMapEntry.ENotificationRequired() {
		eStringToStringMapEntry.ENotify(NewNotificationByFeatureID(eStringToStringMapEntry.AsEObject(), SET, ESTRING_TO_STRING_MAP_ENTRY__VALUE, oldValue, newValue, NO_INDEX))
	}
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) EGetFromID(featureID int, resolve bool) interface{} {
	switch featureID {
	case ESTRING_TO_STRING_MAP_ENTRY__KEY:
		return eStringToStringMapEntry.asEStringToStringMapEntry().GetTypedKey()
	case ESTRING_TO_STRING_MAP_ENTRY__VALUE:
		return eStringToStringMapEntry.asEStringToStringMapEntry().GetTypedValue()
	default:
		return eStringToStringMapEntry.CompactEObjectContainer.EGetFromID(featureID, resolve)
	}
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) ESetFromID(featureID int, newValue interface{}) {
	switch featureID {
	case ESTRING_TO_STRING_MAP_ENTRY__KEY:
		eStringToStringMapEntry.asEStringToStringMapEntry().SetTypedKey(newValue.(string))
	case ESTRING_TO_STRING_MAP_ENTRY__VALUE:
		eStringToStringMapEntry.asEStringToStringMapEntry().SetTypedValue(newValue.(string))
	default:
		eStringToStringMapEntry.CompactEObjectContainer.ESetFromID(featureID, newValue)
	}
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case ESTRING_TO_STRING_MAP_ENTRY__KEY:
		eStringToStringMapEntry.asEStringToStringMapEntry().SetTypedKey("")
	case ESTRING_TO_STRING_MAP_ENTRY__VALUE:
		eStringToStringMapEntry.asEStringToStringMapEntry().SetTypedValue("")
	default:
		eStringToStringMapEntry.CompactEObjectContainer.EUnsetFromID(featureID)
	}
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case ESTRING_TO_STRING_MAP_ENTRY__KEY:
		return eStringToStringMapEntry.key != ""
	case ESTRING_TO_STRING_MAP_ENTRY__VALUE:
		return eStringToStringMapEntry.value != ""
	default:
		return eStringToStringMapEntry.CompactEObjectContainer.EIsSetFromID(featureID)
	}
}
