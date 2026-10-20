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
	"slices"
	"sync"
	"sync/atomic"
)

// EClassExt is the extension of the model object 'EClass'
//
// The lists and maps derived from the class (all its features, attributes, operations, super
// types, ...) are built lazily and each is published in its own atomic pointer: they are read
// without lock once built. setModified invalidates them, and a build that raced an invalidation
// is discarded instead of published stale (see derivedValue).
type EClassExt struct {
	EClassImpl
	adapter *eSuperAdapter

	allStructuralFeatures derivedValue[EList]
	allAttributes         derivedValue[*eClassAttributes]
	allReferences         derivedValue[*eClassReferences]
	allContainments       derivedValue[EList]
	allCrossReferences    derivedValue[EList]
	allOperations         derivedValue[EList]
	allSuperTypes         derivedValue[EList]
	featuresSubSet        derivedValue[*eClassFeaturesSubSet]
	nameToFeature         derivedValue[map[string]EStructuralFeature]
	operationToOverride   derivedValue[map[EOperation]EOperation]

	// buildMutex serializes the builds that read the lists owned by the class (features,
	// operations, super types: created lazily) and set the feature and operation IDs. A build
	// may take the build mutex of a super class, never its own twice. Reads of built values
	// never take it.
	buildMutex sync.Mutex

	subClassesMutex sync.Mutex
	subClasses      []EClassInternal
}

type eClassAttributes struct {
	eAttributes    EList
	eAllAttributes EList
	eIDAttribute   EAttribute
}

type eClassReferences struct {
	eReferences    EList
	eAllReferences EList
}

type eClassFeaturesSubSet struct {
	eContainmentFeatures    EList
	eCrossReferenceFeatures EList
}

// derivedValue is a value built lazily and published for lock-free reads, until invalidated.
//
// Its state is a derivedBuild: either the built value, or none for a generation of the value.
// Each invalidation publishes a new generation without value. A build starts from the state it
// loads, before reading anything, and is published with a compare-and-swap against it: a build
// that raced an invalidation fails to publish, and is built again.
type derivedValue[T any] struct {
	state atomic.Pointer[derivedBuild[T]]
}

// derivedBuild is a state of a derivedValue. It is never changed once published.
type derivedBuild[T any] struct {
	generation uint64 // the number of invalidations of the derivedValue before
	built      bool
	value      T
}

// load returns the published value, if built.
func (d *derivedValue[T]) load() (value T, ok bool) {
	if s := d.state.Load(); s != nil && s.built {
		return s.value, true
	}
	return value, false
}

// get returns the published value, or builds and publishes it.
func (d *derivedValue[T]) get(build func() T) T {
	for {
		start := d.state.Load()
		if start != nil && start.built {
			return start.value
		}
		value := build()
		if d.publish(start, value) {
			return value
		}
		// invalidated, or published by another build meanwhile: load it again
	}
}

// publish publishes value, built from the state start, unless the state changed meanwhile.
func (d *derivedValue[T]) publish(start *derivedBuild[T], value T) bool {
	built := &derivedBuild[T]{built: true, value: value}
	if start != nil {
		built.generation = start.generation
	}
	return d.state.CompareAndSwap(start, built)
}

// invalidate discards the published value: it is built again at its next read. A build in
// progress fails to publish.
func (d *derivedValue[T]) invalidate() {
	for {
		current := d.state.Load()
		next := &derivedBuild[T]{generation: 1}
		if current != nil {
			next.generation = current.generation + 1
		}
		if d.state.CompareAndSwap(current, next) {
			return
		}
	}
}

type EClassInternal interface {
	getSubClasses() []EClassInternal
	addSubClass(EClassInternal)
	removeSubClass(EClassInternal)
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
				if oldValue := notification.GetOldValue(); oldValue != nil {
					oldValue.(EClassInternal).removeSubClass(eNotifier)
				}
				if newValue := notification.GetNewValue(); newValue != nil {
					newValue.(EClassInternal).addSubClass(eNotifier)
				}
			case ADD:
				if newValue := notification.GetNewValue(); newValue != nil {
					newValue.(EClassInternal).addSubClass(eNotifier)
				}
			case ADD_MANY:
				if newValue := notification.GetNewValue(); newValue != nil {
					for _, s := range newValue.([]any) {
						s.(EClassInternal).addSubClass(eNotifier)
					}
				}
			case REMOVE:
				if oldValue := notification.GetOldValue(); oldValue != nil {
					oldValue.(EClassInternal).removeSubClass(eNotifier)
				}
			case REMOVE_MANY:
				if oldValue := notification.GetOldValue(); oldValue != nil {
					for _, s := range oldValue.([]any) {
						s.(EClassInternal).removeSubClass(eNotifier)
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
	// created now, not at the first access: the lists the derived values are built from are read
	// without lock, and a lazy creation would race between readers. The operations stay lazy: they
	// are read under the build mutex of the class, and most models have none.
	eClass.GetESuperTypes()
	eClass.GetEStructuralFeatures()
}

func (eClass *EClassExt) IsSuperTypeOf(someClass EClass) bool {
	return someClass == eClass || (someClass != nil && someClass.GetEAllSuperTypes().Contains(eClass))
}

func (eClass *EClassExt) GetFeatureCount() int {
	// read at each access to a feature of an object: the lock-free path is inlined
	var features EList
	if s := eClass.allStructuralFeatures.state.Load(); s != nil && s.built {
		features = s.value
	} else {
		features = eClass.GetEAllStructuralFeatures()
	}
	return features.Size()
}

func (eClass *EClassExt) GetEStructuralFeature(featureID int) EStructuralFeature {
	// read at each access to a feature of an object: the lock-free path is inlined
	var features EList
	if s := eClass.allStructuralFeatures.state.Load(); s != nil && s.built {
		features = s.value
	} else {
		features = eClass.GetEAllStructuralFeatures()
	}
	if featureID >= 0 && featureID < features.Size() {
		return features.Get(featureID).(EStructuralFeature)
	}
	return nil
}

func (eClass *EClassExt) GetEStructuralFeatureFromName(featureName string) EStructuralFeature {
	if value, ok := eClass.nameToFeature.load(); ok {
		return value[featureName]
	}
	return eClass.nameToFeature.get(eClass.buildNameToFeatureMap)[featureName]
}

func (eClass *EClassExt) GetFeatureID(feature EStructuralFeature) int {
	// read at each access to a feature of an object: the lock-free path is inlined
	var features EList
	if s := eClass.allStructuralFeatures.state.Load(); s != nil && s.built {
		features = s.value
	} else {
		features = eClass.GetEAllStructuralFeatures()
	}
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
	if value, ok := eClass.operationToOverride.load(); ok {
		return value[operation]
	}
	return eClass.operationToOverride.get(eClass.buildOperationToOverrideMap)[operation]
}

func (eClass *EClassExt) GetEAllStructuralFeatures() EList {
	if value, ok := eClass.allStructuralFeatures.load(); ok {
		return value
	}
	return eClass.allStructuralFeatures.get(eClass.buildEAllStructuralFeatures)
}

func (eClass *EClassExt) GetEAllAttributes() EList {
	if value, ok := eClass.allAttributes.load(); ok {
		return value.eAllAttributes
	}
	return eClass.allAttributes.get(eClass.buildEAllAttributes).eAllAttributes
}

func (eClass *EClassExt) GetEAttributes() EList {
	if value, ok := eClass.allAttributes.load(); ok {
		return value.eAttributes
	}
	return eClass.allAttributes.get(eClass.buildEAllAttributes).eAttributes
}

func (eClass *EClassExt) GetEIDAttribute() EAttribute {
	if value, ok := eClass.allAttributes.load(); ok {
		return value.eIDAttribute
	}
	return eClass.allAttributes.get(eClass.buildEAllAttributes).eIDAttribute
}

func (eClass *EClassExt) GetEAllReferences() EList {
	if value, ok := eClass.allReferences.load(); ok {
		return value.eAllReferences
	}
	return eClass.allReferences.get(eClass.buildEAllReferences).eAllReferences
}

func (eClass *EClassExt) GetEReferences() EList {
	if value, ok := eClass.allReferences.load(); ok {
		return value.eReferences
	}
	return eClass.allReferences.get(eClass.buildEAllReferences).eReferences
}

func (eClass *EClassExt) GetEAllContainments() EList {
	if value, ok := eClass.allContainments.load(); ok {
		return value
	}
	return eClass.allContainments.get(eClass.buildEAllContainments)
}

func (eClass *EClassExt) GetEAllCrossReferences() EList {
	if value, ok := eClass.allCrossReferences.load(); ok {
		return value
	}
	return eClass.allCrossReferences.get(eClass.buildEAllCrossReferences)
}

func (eClass *EClassExt) GetEAllOperations() EList {
	if value, ok := eClass.allOperations.load(); ok {
		return value
	}
	return eClass.allOperations.get(eClass.buildEAllOperations)
}

func (eClass *EClassExt) GetEAllSuperTypes() EList {
	if value, ok := eClass.allSuperTypes.load(); ok {
		return value
	}
	return eClass.allSuperTypes.get(eClass.buildEAllSuperTypes)
}

func (eClass *EClassExt) GetEContainmentFeatures() EList {
	if value, ok := eClass.featuresSubSet.load(); ok {
		return value.eContainmentFeatures
	}
	return eClass.featuresSubSet.get(eClass.buildFeaturesSubSet).eContainmentFeatures
}

func (eClass *EClassExt) GetECrossReferenceFeatures() EList {
	if value, ok := eClass.featuresSubSet.load(); ok {
		return value.eCrossReferenceFeatures
	}
	return eClass.featuresSubSet.get(eClass.buildFeaturesSubSet).eCrossReferenceFeatures
}

// EIsSetFromID reports the derived features as set when they are built and not empty.
func (eClass *EClassExt) EIsSetFromID(featureID int) bool {
	switch featureID {
	case ECLASS__EALL_ATTRIBUTES:
		attributes, ok := eClass.allAttributes.load()
		return ok && isListSet(attributes.eAllAttributes)
	case ECLASS__EATTRIBUTES:
		attributes, ok := eClass.allAttributes.load()
		return ok && isListSet(attributes.eAttributes)
	case ECLASS__EID_ATTRIBUTE:
		attributes, ok := eClass.allAttributes.load()
		return ok && attributes.eIDAttribute != nil
	case ECLASS__EALL_REFERENCES:
		references, ok := eClass.allReferences.load()
		return ok && isListSet(references.eAllReferences)
	case ECLASS__EREFERENCES:
		references, ok := eClass.allReferences.load()
		return ok && isListSet(references.eReferences)
	case ECLASS__EALL_CONTAINMENTS:
		l, ok := eClass.allContainments.load()
		return ok && isListSet(l)
	case ECLASS__EALL_CROSS_REFERENCES:
		l, ok := eClass.allCrossReferences.load()
		return ok && isListSet(l)
	case ECLASS__EALL_OPERATIONS:
		l, ok := eClass.allOperations.load()
		return ok && isListSet(l)
	case ECLASS__EALL_STRUCTURAL_FEATURES:
		l, ok := eClass.allStructuralFeatures.load()
		return ok && isListSet(l)
	case ECLASS__EALL_SUPER_TYPES:
		l, ok := eClass.allSuperTypes.load()
		return ok && isListSet(l)
	case ECLASS__ECONTAINMENT_FEATURES:
		subSet, ok := eClass.featuresSubSet.load()
		return ok && isListSet(subSet.eContainmentFeatures)
	case ECLASS__ECROSS_REFERENCE_FEATURES:
		subSet, ok := eClass.featuresSubSet.load()
		return ok && isListSet(subSet.eCrossReferenceFeatures)
	default:
		return eClass.EClassImpl.EIsSetFromID(featureID)
	}
}

func isListSet(l EList) bool {
	return l != nil && l.Size() != 0
}

func (eClass *EClassExt) buildNameToFeatureMap() map[string]EStructuralFeature {
	nameToFeatureMap := make(map[string]EStructuralFeature)
	for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
		feature := itFeature.Next().(EStructuralFeature)
		nameToFeatureMap[feature.GetName()] = feature
	}
	return nameToFeatureMap
}

func (eClass *EClassExt) buildOperationToOverrideMap() map[EOperation]EOperation {
	operationToOverrideMap := make(map[EOperation]EOperation)
	operations := eClass.GetEAllOperations()
	size := operations.Size()
	for i := range size {
		for j := size - 1; j > i; j-- {
			oi := operations.Get(i).(EOperation)
			oj := operations.Get(j).(EOperation)
			if oj.IsOverrideOf(oi) {
				operationToOverrideMap[oi] = oj
			}
		}
	}
	return operationToOverrideMap
}

func (eClass *EClassExt) buildFeaturesSubSet() *eClassFeaturesSubSet {
	// like EMF , the subsets are computed from all structural features ( inherited ones included )
	// so that EContents & ECrossReferences of an object visit the features of its super classes
	containments := []any{}
	crossReferences := []any{}
	for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
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
	return &eClassFeaturesSubSet{
		eContainmentFeatures:    NewImmutableEList(containments),
		eCrossReferenceFeatures: NewImmutableEList(crossReferences),
	}
}

func (eClass *EClassExt) buildEAllAttributes() *eClassAttributes {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
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
	return &eClassAttributes{
		eAttributes:    NewImmutableEList(attributes),
		eAllAttributes: NewImmutableEList(allAttributes),
		eIDAttribute:   eIDAttribute,
	}
}

func (eClass *EClassExt) buildEAllReferences() *eClassReferences {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
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
	return &eClassReferences{
		eReferences:    NewImmutableEList(references),
		eAllReferences: NewImmutableEList(allReferences),
	}
}

func (eClass *EClassExt) buildEAllContainments() EList {
	allContainments := []any{}
	for itReference := eClass.GetEAllReferences().Iterator(); itReference.HasNext(); {
		reference := itReference.Next().(EReference)
		if reference.IsContainment() {
			allContainments = append(allContainments, reference)
		}
	}
	return NewImmutableEList(allContainments)
}

func (eClass *EClassExt) buildEAllCrossReferences() EList {
	allCrossReferences := []any{}
	for itReference := eClass.GetEAllReferences().Iterator(); itReference.HasNext(); {
		reference := itReference.Next().(EReference)
		if !reference.IsContainment() && !reference.IsContainer() {
			allCrossReferences = append(allCrossReferences, reference)
		}
	}
	return NewImmutableEList(allCrossReferences)
}

func (eClass *EClassExt) buildEAllOperations() EList {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
	allOperations := []any{}
	for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
		superOperations := itClass.Next().(EClass).GetEAllOperations()
		allOperations = append(allOperations, superOperations.ToArray()...)
	}

	operationID := len(allOperations)
	for itFeature := eClass.GetEOperations().Iterator(); itFeature.HasNext(); {
		operation, isOperation := itFeature.Next().(EOperation)
		if isOperation {
			// set only when changed: the ID is read without lock
			if operation.GetOperationID() != operationID {
				operation.SetOperationID(operationID)
			}
			operationID++
			allOperations = append(allOperations, operation)
		}
	}
	return NewImmutableEList(allOperations)
}

func (eClass *EClassExt) buildEAllStructuralFeatures() EList {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
	allFeatures := []any{}
	for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
		superFeatures := itClass.Next().(EClass).GetEAllStructuralFeatures()
		allFeatures = append(allFeatures, superFeatures.ToArray()...)
	}

	featureID := len(allFeatures)
	for itFeature := eClass.GetEStructuralFeatures().Iterator(); itFeature.HasNext(); {
		feature := itFeature.Next().(EStructuralFeature)
		// set only when changed: the ID is read without lock
		if feature.GetFeatureID() != featureID {
			feature.SetFeatureID(featureID)
		}
		featureID++
		allFeatures = append(allFeatures, feature)
	}
	return NewImmutableEList(allFeatures)
}

func (eClass *EClassExt) buildEAllSuperTypes() EList {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
	allSuperTypes := []any{}
	for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
		superClass := itClass.Next().(EClass)
		superTypes := superClass.GetEAllSuperTypes()
		allSuperTypes = append(allSuperTypes, superTypes.ToArray()...)
		allSuperTypes = append(allSuperTypes, superClass)
	}
	return NewImmutableEList(allSuperTypes)
}

func (eClass *EClassExt) getSubClasses() []EClassInternal {
	eClass.subClassesMutex.Lock()
	defer eClass.subClassesMutex.Unlock()
	return eClass.subClasses
}

// addSubClass and removeSubClass update the sub classes under the mutex, in a new slice: the
// slice returned by getSubClasses is never changed, setModified iterates it without lock.
func (eClass *EClassExt) addSubClass(subClass EClassInternal) {
	eClass.subClassesMutex.Lock()
	eClass.subClasses = append(slices.Clip(eClass.subClasses), subClass)
	eClass.subClassesMutex.Unlock()
}

func (eClass *EClassExt) removeSubClass(subClass EClassInternal) {
	eClass.subClassesMutex.Lock()
	if i := slices.Index(eClass.subClasses, subClass); i != -1 {
		eClass.subClasses = slices.Delete(slices.Clone(eClass.subClasses), i, i+1)
	}
	eClass.subClassesMutex.Unlock()
}

// setModified invalidates the values derived from featureID, in this class and its sub classes.
// A value derived from another one is invalidated with it.
func (eClass *EClassExt) setModified(featureID int) {
	allStructuralFeatures := false
	switch featureID {
	case ECLASS__ESTRUCTURAL_FEATURES:
		eClass.allAttributes.invalidate()
		eClass.allReferences.invalidate()
		eClass.allContainments.invalidate()
		eClass.allCrossReferences.invalidate()
		allStructuralFeatures = true
	case ECLASS__EATTRIBUTES:
		eClass.allAttributes.invalidate()
		eClass.allContainments.invalidate()
		eClass.allCrossReferences.invalidate()
		allStructuralFeatures = true
	case ECLASS__EREFERENCES:
		eClass.allReferences.invalidate()
		eClass.allContainments.invalidate()
		eClass.allCrossReferences.invalidate()
		allStructuralFeatures = true
	case ECLASS__EOPERATIONS:
		eClass.allOperations.invalidate()
		eClass.operationToOverride.invalidate()
		eClass.allContainments.invalidate()
		eClass.allCrossReferences.invalidate()
	case ECLASS__ESUPER_TYPES:
		eClass.allSuperTypes.invalidate()
		eClass.allAttributes.invalidate()
		eClass.allOperations.invalidate()
		eClass.operationToOverride.invalidate()
		eClass.allReferences.invalidate()
		eClass.allContainments.invalidate()
		eClass.allCrossReferences.invalidate()
		allStructuralFeatures = true
	}
	if allStructuralFeatures {
		eClass.allStructuralFeatures.invalidate()
		eClass.featuresSubSet.invalidate()
		eClass.nameToFeature.invalidate()
	}
	for _, s := range eClass.getSubClasses() {
		s.setModified(featureID)
	}
}

func IsMapEntry(eClass EClass) bool {
	instanceTypeName := eClass.GetInstanceTypeName()
	return (instanceTypeName == "java.util.Map.Entry" ||
		instanceTypeName == "java.util.Map$Entry" ||
		instanceTypeName == "github.com/masagroup/soft.go/ecore/ecore.EMapEntry") &&
		eClass.GetEStructuralFeatureFromName("key") != nil &&
		eClass.GetEStructuralFeatureFromName("value") != nil
}
