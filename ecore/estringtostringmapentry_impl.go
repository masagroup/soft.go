// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

// *****************************************************************************
//
// Warning: This file was generated by soft.generator.go Generator
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
	return eStringToStringMapEntry.GetStringKey()
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) GetValue() interface{} {
	return eStringToStringMapEntry.GetStringValue()
}

// GetStringKey get the value of key
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) GetStringKey() string {
	return eStringToStringMapEntry.key
}

// SetStringKey set the value of key
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) SetStringKey(newKey string) {
	oldKey := eStringToStringMapEntry.key
	eStringToStringMapEntry.key = newKey
	if eStringToStringMapEntry.ENotificationRequired() {
		eStringToStringMapEntry.ENotify(NewNotificationByFeatureID(eStringToStringMapEntry.AsEObject(), SET, ESTRING_TO_STRING_MAP_ENTRY__KEY, oldKey, newKey, NO_INDEX))
	}
}

// GetStringValue get the value of value
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) GetStringValue() string {
	return eStringToStringMapEntry.value
}

// SetStringValue set the value of value
func (eStringToStringMapEntry *eStringToStringMapEntryImpl) SetStringValue(newValue string) {
	oldValue := eStringToStringMapEntry.value
	eStringToStringMapEntry.value = newValue
	if eStringToStringMapEntry.ENotificationRequired() {
		eStringToStringMapEntry.ENotify(NewNotificationByFeatureID(eStringToStringMapEntry.AsEObject(), SET, ESTRING_TO_STRING_MAP_ENTRY__VALUE, oldValue, newValue, NO_INDEX))
	}
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) EGetFromID(featureID int, resolve bool) interface{} {
	switch featureID {
	case ESTRING_TO_STRING_MAP_ENTRY__KEY:
		return eStringToStringMapEntry.asEStringToStringMapEntry().GetStringKey()
	case ESTRING_TO_STRING_MAP_ENTRY__VALUE:
		return eStringToStringMapEntry.asEStringToStringMapEntry().GetStringValue()
	default:
		return eStringToStringMapEntry.CompactEObjectContainer.EGetFromID(featureID, resolve)
	}
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) ESetFromID(featureID int, newValue interface{}) {
	switch featureID {
	case ESTRING_TO_STRING_MAP_ENTRY__KEY:
		eStringToStringMapEntry.asEStringToStringMapEntry().SetStringKey(newValue.(string))
	case ESTRING_TO_STRING_MAP_ENTRY__VALUE:
		eStringToStringMapEntry.asEStringToStringMapEntry().SetStringValue(newValue.(string))
	default:
		eStringToStringMapEntry.CompactEObjectContainer.ESetFromID(featureID, newValue)
	}
}

func (eStringToStringMapEntry *eStringToStringMapEntryImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case ESTRING_TO_STRING_MAP_ENTRY__KEY:
		eStringToStringMapEntry.asEStringToStringMapEntry().SetStringKey("")
	case ESTRING_TO_STRING_MAP_ENTRY__VALUE:
		eStringToStringMapEntry.asEStringToStringMapEntry().SetStringValue("")
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
