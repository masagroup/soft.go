// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "sync"

// EClassExt is the extension of the model object 'EClass'
type EClassExt struct {
	EClassImpl
	adapter                *eSuperAdapter
	nameToFeatureMap       map[string]EStructuralFeature
	operationToOverrideMap map[EOperation]EOperation
	subClasses             []EClassInternal
	mutex                  sync.Mutex
}

type EClassInternal interface {
	getSubClasses() []EClassInternal
	setSubClasses([]EClassInternal)
	setModified(featureID int)
}

type eSuperAdapter struct {
	AbstractEAdapter
	class *EClassExt
}

func (adapter *eSuperAdapter) NotifyChanged(notification ENotification) {
	eventType := notification.GetEventType()
	eNotifier := notification.GetNotifier().(EClassInternal)
	if eventType != REMOVING_ADAPTER {
		featureID := notification.GetFeatureID()
		if featureID == ECLASS__ESUPER_TYPES {
			switch eventType {
			case SET:
				fallthrough
			case RESOLVE:
				oldValue := notification.GetOldValue()
				if oldValue != nil {
					class := oldValue.(EClassInternal)
					subClasses := class.getSubClasses()
					for i, s := range subClasses {
						if s == eNotifier {
							class.setSubClasses(append(subClasses[:i], subClasses[i+1:]...))
							break
						}
					}
				}
				newValue := notification.GetNewValue()
				if newValue != nil {
					class := newValue.(EClassInternal)
					class.setSubClasses(append(class.getSubClasses(), eNotifier))
				}
			case ADD:
				newValue := notification.GetNewValue()
				if newValue != nil {
					class := newValue.(EClassInternal)
					class.setSubClasses(append(class.getSubClasses(), eNotifier))
				}
			case ADD_MANY:
				newValue := notification.GetNewValue()
				if newValue != nil {
					collection := newValue.([]any)
					for _, s := range collection {
						class := s.(EClassInternal)
						class.setSubClasses(append(class.getSubClasses(), eNotifier))
					}
				}
			case REMOVE:
				oldValue := notification.GetOldValue()
				if oldValue != nil {
					class := oldValue.(EClassInternal)
					subClasses := class.getSubClasses()
					for i, s := range subClasses {
						if s == eNotifier {
							class.setSubClasses(append(subClasses[:i], subClasses[i+1:]...))
							break
						}
					}
				}
			case REMOVE_MANY:
				oldValue := notification.GetOldValue()
				if oldValue != nil {
					collection := oldValue.([]any)
					for _, s := range collection {
						class := s.(EClassInternal)
						subClasses := class.getSubClasses()
						for i, s := range subClasses {
							if s == eNotifier {
								class.setSubClasses(append(subClasses[:i], subClasses[i+1:]...))
								break
							}
						}
					}
				}
			}
		}
		adapter.class.setModified(featureID)
	}
}

func newEClassExt() *EClassExt {
	eClass := new(EClassExt)
	eClass.SetInterfaces(eClass)
	eClass.Initialize()
	return eClass
}

func (eClass *EClassExt) Initialize() {
	eClass.EClassImpl.Initialize()
	eClass.adapter = &eSuperAdapter{class: eClass}
	eClass.EAdapters().Add(eClass.adapter)
}

func (eClass *EClassExt) IsSuperTypeOf(someClass EClass) bool {
	return someClass == eClass || (someClass != nil && someClass.GetEAllSuperTypes().Contains(eClass))
}

func (eClass *EClassExt) GetFeatureCount() int {
	return eClass.GetEAllStructuralFeatures().Size()
}

func (eClass *EClassExt) GetEStructuralFeature(featureID int) EStructuralFeature {
	features := eClass.GetEAllStructuralFeatures()
	if featureID >= 0 && featureID < features.Size() {
		return features.Get(featureID).(EStructuralFeature)
	}
	return nil
}

func (eClass *EClassExt) GetEStructuralFeatureFromName(featureName string) EStructuralFeature {
	eClass.initNameToFeatureMap()
	eClass.mutex.Lock()
	defer eClass.mutex.Unlock()
	return eClass.nameToFeatureMap[featureName]
}

func (eClass *EClassExt) initNameToFeatureMap() {
	eClass.initEAllStructuralFeatures()
	eClass.mutex.Lock()
	if eClass.nameToFeatureMap == nil {
		eClass.nameToFeatureMap = make(map[string]EStructuralFeature)
		for itFeature := eClass.eAllStructuralFeatures.Iterator(); itFeature.HasNext(); {
			feature := itFeature.Next().(EStructuralFeature)
			eClass.nameToFeatureMap[feature.GetName()] = feature
		}
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) GetFeatureID(feature EStructuralFeature) int {
	features := eClass.GetEAllStructuralFeatures()
	featureID := feature.GetFeatureID()
	if featureID != -1 {
		for ; featureID < features.Size(); featureID++ {
			if features.Get(featureID) == feature {
				return featureID
			}
		}
	}
	return -1
}

func (eClass *EClassExt) GetOperationCount() int {
	return eClass.GetEAllOperations().Size()
}

func (eClass *EClassExt) GetEOperation(operationID int) EOperation {
	operations := eClass.GetEAllOperations()
	if operationID >= 0 && operationID < operations.Size() {
		return operations.Get(operationID).(EOperation)
	}
	return nil
}

func (eClass *EClassExt) GetOperationID(operation EOperation) int {
	operations := eClass.GetEAllOperations()
	operationID := operation.GetOperationID()
	if operationID != -1 {
		for ; operationID < operations.Size(); operationID++ {
			if operations.Get(operationID) == operation {
				return operationID
			}
		}
	}
	return -1
}

func (eClass *EClassExt) GetOverride(operation EOperation) EOperation {
	eClass.initOperationToOverrideMap()
	eClass.mutex.Lock()
	defer eClass.mutex.Unlock()
	return eClass.operationToOverrideMap[operation]
}

func (eClass *EClassExt) initOperationToOverrideMap() {
	eClass.initEAllOperations()
	eClass.mutex.Lock()
	if eClass.operationToOverrideMap == nil {
		eClass.operationToOverrideMap = make(map[EOperation]EOperation)
		size := eClass.eAllOperations.Size()
		for i := 0; i < size; i++ {
			for j := size - 1; j > i; j-- {
				oi := eClass.eAllOperations.Get(i).(EOperation)
				oj := eClass.eAllOperations.Get(j).(EOperation)
				if oj.IsOverrideOf(oi) {
					eClass.operationToOverrideMap[oi] = oj
				}
			}
		}
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) initEAttributes() {
	eClass.initEAllAttributes()
}

func (eClass *EClassExt) initEReferences() {
	eClass.initEAllReferences()
}

func (eClass *EClassExt) initEContainmentFeatures() {
	eClass.initFeaturesSubSet()
}

func (eClass *EClassExt) initECrossReferenceFeatures() {
	eClass.initFeaturesSubSet()
}

func (eClass *EClassExt) initFeaturesSubSet() {
	eClass.initEAllStructuralFeatures()
	eClass.mutex.Lock()
	if eClass.eContainmentFeatures == nil {
		containments := []any{}
		crossReferences := []any{}
		for itFeature := eClass.GetEStructuralFeatures().Iterator(); itFeature.HasNext(); {
			ref, isRef := itFeature.Next().(EReference)
			if isRef {
				if ref.IsContainment() {
					if !ref.IsDerived() {
						containments = append(containments, ref)
					}
				} else if !ref.IsContainer() {
					if !ref.IsDerived() {
						crossReferences = append(crossReferences, ref)
					}
				}
			}

		}
		eClass.eContainmentFeatures = NewImmutableEList(containments)
		eClass.eCrossReferenceFeatures = NewImmutableEList(crossReferences)
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) initEAllAttributes() {
	eClass.mutex.Lock()
	if eClass.eAllAttributes == nil {

		attributes := []any{}
		allAttributes := []any{}
		var eIDAttribute EAttribute = nil
		for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
			superAttributes := itClass.Next().(EClass).GetEAllAttributes()
			for itAttribute := superAttributes.Iterator(); itAttribute.HasNext(); {
				attribute := itAttribute.Next().(EAttribute)
				allAttributes = append(allAttributes, attribute)
				if attribute.IsID() && eIDAttribute == nil {
					eIDAttribute = attribute
				}
			}
		}

		for itFeature := eClass.GetEStructuralFeatures().Iterator(); itFeature.HasNext(); {
			attribute, isAttribute := itFeature.Next().(EAttribute)
			if isAttribute {
				attributes = append(attributes, attribute)
				allAttributes = append(allAttributes, attribute)
				if attribute.IsID() && eIDAttribute == nil {
					eIDAttribute = attribute
				}
			}
		}
		eClass.eIDAttribute = eIDAttribute
		eClass.eAttributes = NewImmutableEList(attributes)
		eClass.eAllAttributes = NewImmutableEList(allAttributes)
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) initEAllReferences() {
	eClass.mutex.Lock()
	if eClass.eAllReferences == nil {
		allReferences := []any{}
		references := []any{}
		for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
			superReferences := itClass.Next().(EClass).GetEAllReferences()
			allReferences = append(allReferences, superReferences.ToArray()...)
		}

		for itFeature := eClass.GetEStructuralFeatures().Iterator(); itFeature.HasNext(); {
			reference, isReference := itFeature.Next().(EReference)
			if isReference {
				references = append(references, reference)
				allReferences = append(allReferences, reference)
			}
		}

		eClass.eReferences = NewImmutableEList(references)
		eClass.eAllReferences = NewImmutableEList(allReferences)
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) initEAllContainments() {
	eClass.mutex.Lock()
	if eClass.eAllContainments == nil {
		allContainments := []any{}
		for itReference := eClass.GetEAllReferences().Iterator(); itReference.HasNext(); {
			reference := itReference.Next().(EReference)
			if reference.IsContainment() {
				allContainments = append(allContainments, reference)
			}
		}
		eClass.eAllContainments = NewImmutableEList(allContainments)
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) initEAllOperations() {
	eClass.mutex.Lock()
	if eClass.eAllOperations == nil {
		eClass.operationToOverrideMap = nil

		allOperations := []any{}
		for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
			superOperations := itClass.Next().(EClass).GetEAllOperations()
			allOperations = append(allOperations, superOperations.ToArray()...)
		}

		operationID := len(allOperations)
		for itFeature := eClass.GetEOperations().Iterator(); itFeature.HasNext(); {
			operation, isOperation := itFeature.Next().(EOperation)
			if isOperation {
				operation.SetOperationID(operationID)
				operationID++
				allOperations = append(allOperations, operation)
			}
		}
		eClass.eAllOperations = NewImmutableEList(allOperations)
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) initEAllStructuralFeatures() {
	eClass.mutex.Lock()
	if eClass.eAllStructuralFeatures == nil {
		eClass.eCrossReferenceFeatures = nil
		eClass.eContainmentFeatures = nil
		eClass.nameToFeatureMap = nil

		allFeatures := []any{}
		for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
			superFeatures := itClass.Next().(EClass).GetEAllStructuralFeatures()
			allFeatures = append(allFeatures, superFeatures.ToArray()...)
		}

		featureID := len(allFeatures)
		for itFeature := eClass.GetEStructuralFeatures().Iterator(); itFeature.HasNext(); {
			feature := itFeature.Next().(EStructuralFeature)
			feature.SetFeatureID(featureID)
			featureID++
			allFeatures = append(allFeatures, feature)
		}
		eClass.eAllStructuralFeatures = NewImmutableEList(allFeatures)
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) initEAllSuperTypes() {
	eClass.mutex.Lock()
	if eClass.eAllSuperTypes == nil {
		allSuperTypes := []any{}
		for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
			superClass := itClass.Next().(EClass)
			superTypes := superClass.GetEAllSuperTypes()
			allSuperTypes = append(allSuperTypes, superTypes.ToArray()...)
			allSuperTypes = append(allSuperTypes, superClass)
		}
		eClass.eAllSuperTypes = NewImmutableEList(allSuperTypes)
	}
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) initEIDAttribute() {
	eClass.initEAllAttributes()
}

func (eClass *EClassExt) getSubClasses() []EClassInternal {
	eClass.mutex.Lock()
	defer eClass.mutex.Unlock()
	return eClass.subClasses
}

func (eClass *EClassExt) setSubClasses(subClasses []EClassInternal) {
	eClass.mutex.Lock()
	eClass.subClasses = subClasses
	eClass.mutex.Unlock()
}

func (eClass *EClassExt) setModified(featureID int) {
	eClass.mutex.Lock()
	switch featureID {
	case ECLASS__ESTRUCTURAL_FEATURES:
		eClass.eAllAttributes = nil
		eClass.eAllStructuralFeatures = nil
		eClass.eAllReferences = nil
		eClass.eAllContainments = nil
	case ECLASS__EATTRIBUTES:
		eClass.eAllAttributes = nil
		eClass.eAllStructuralFeatures = nil
		eClass.eAllContainments = nil
	case ECLASS__EREFERENCES:
		eClass.eAllReferences = nil
		eClass.eAllStructuralFeatures = nil
		eClass.eAllContainments = nil
	case ECLASS__EOPERATIONS:
		eClass.eAllOperations = nil
		eClass.eAllContainments = nil
	case ECLASS__ESUPER_TYPES:
		eClass.eAllSuperTypes = nil
		eClass.eAllAttributes = nil
		eClass.eAllOperations = nil
		eClass.eAllStructuralFeatures = nil
		eClass.eAllReferences = nil
		eClass.eAllContainments = nil
	}
	for _, s := range eClass.subClasses {
		s.setModified(featureID)
	}
	eClass.mutex.Unlock()
}

func IsMapEntry(eClass EClass) bool {
	instanceTypeName := eClass.GetInstanceTypeName()
	return (instanceTypeName == "java.util.Map.Entry" ||
		instanceTypeName == "java.util.Map$Entry" ||
		instanceTypeName == "ecore.EMapEntry") &&
		eClass.GetEStructuralFeatureFromName("key") != nil &&
		eClass.GetEStructuralFeatureFromName("value") != nil
}
