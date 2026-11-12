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

	allStructuralFeatures derivedValue[*eClassFeatures]
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
	subClasses      []eClassSubClasses
}

// eClassFeatures is the list of all the features of a class and the same features as a slice: the
// slice serves the lookups by ID without interface calls nor type assertions.
type eClassFeatures struct {
	list  EList
	slice []EStructuralFeature
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

// eClassModifiable is a class whose derived values a change of its features, operations or super
// types invalidates, in the class and in its sub classes: an attribute that becomes, or stops
// being, the ID of its class invalidates the attributes of the class. setModified starts the
// invalidation; setModifiedFrom carries it on in each sub class, visited holding the classes
// invalidated already.
type eClassModifiable interface {
	setModified(featureID int)
	setModifiedFrom(featureID int, visited map[*EClassExt]struct{})
}

// eClassSubClasses is a class that knows its sub classes, to invalidate their derived values
// with its own: a sub class is an eClassModifiable. The super types adapter registers a class with
// its super types (addSubClass, removeSubClass). The values are built the other way, from the
// super types: see eClassHierarchy.
type eClassSubClasses interface {
	eClassModifiable
	getSubClasses() []eClassSubClasses
	addSubClass(eClassSubClasses)
	removeSubClass(eClassSubClasses)
}

// The internal interfaces of a class, checked at compile time: they are asserted at run time.
var (
	_ eClassModifiable = (*EClassExt)(nil)
	_ eClassSubClasses = (*EClassExt)(nil)
	_ eClassHierarchy  = (*EClassExt)(nil)
)

// eClassHierarchy builds the values derived from the super types of a class, for the build of the
// same value in its sub classes: path holds the classes whose build is in progress in the calling
// goroutine. A class of path is in a cycle of super types, an invalid model: its value is nil,
// neither built again nor waited for (its build mutex is held), and the recursion over the super
// types ends, as EMF ends it with the classes whose computation is in progress.
type eClassHierarchy interface {
	getEAllSuperTypes(path []*EClassExt) EList
	getEAllAttributes(path []*EClassExt) EList
	getEAllReferences(path []*EClassExt) EList
	getEAllOperations(path []*EClassExt) EList
	getEAllStructuralFeatures(path []*EClassExt) EList
}

// superValue returns the value derived from superClass for the build of the classes of path: by
// eClassHierarchy when superClass implements it, nil for a class of path
func superValue(superClass EClass, path []*EClassExt, internal func(eClassHierarchy, []*EClassExt) EList, public func(EClass) EList) EList {
	if hierarchy, ok := superClass.(eClassHierarchy); ok {
		return internal(hierarchy, path)
	}
	return public(superClass)
}

type eSuperAdapter struct {
	AbstractEAdapter
	class *EClassExt
}

func (adapter *eSuperAdapter) NotifyChanged(notification ENotification) {
	eventType := notification.GetEventType()
	eNotifier := notification.GetNotifier().(eClassSubClasses)
	if eventType != REMOVING_ADAPTER {
		featureID := notification.GetFeatureID()
		if featureID == ECLASS__ESUPER_TYPES {
			switch eventType {
			case SET:
				fallthrough
			case RESOLVE:
				if oldValue := notification.GetOldValue(); oldValue != nil {
					oldValue.(eClassSubClasses).removeSubClass(eNotifier)
				}
				if newValue := notification.GetNewValue(); newValue != nil {
					newValue.(eClassSubClasses).addSubClass(eNotifier)
				}
			case ADD:
				if newValue := notification.GetNewValue(); newValue != nil {
					newValue.(eClassSubClasses).addSubClass(eNotifier)
				}
			case ADD_MANY:
				if newValue := notification.GetNewValue(); newValue != nil {
					for _, s := range newValue.([]any) {
						s.(eClassSubClasses).addSubClass(eNotifier)
					}
				}
			case REMOVE:
				if oldValue := notification.GetOldValue(); oldValue != nil {
					oldValue.(eClassSubClasses).removeSubClass(eNotifier)
				}
			case REMOVE_MANY:
				if oldValue := notification.GetOldValue(); oldValue != nil {
					for _, s := range oldValue.([]any) {
						s.(eClassSubClasses).removeSubClass(eNotifier)
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
	if s := eClass.allStructuralFeatures.state.Load(); s != nil && s.built {
		return len(s.value.slice)
	}
	return len(eClass.features(nil).slice)
}

func (eClass *EClassExt) GetEStructuralFeature(featureID int) EStructuralFeature {
	// read at each access to a feature of an object: the lock-free path is inlined
	var features []EStructuralFeature
	if s := eClass.allStructuralFeatures.state.Load(); s != nil && s.built {
		features = s.value.slice
	} else {
		features = eClass.features(nil).slice
	}
	if featureID >= 0 && featureID < len(features) {
		return features[featureID]
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
	var features []EStructuralFeature
	if s := eClass.allStructuralFeatures.state.Load(); s != nil && s.built {
		features = s.value.slice
	} else {
		features = eClass.features(nil).slice
	}
	featureID := feature.GetFeatureID()
	if featureID != -1 {
		for ; featureID < len(features); featureID++ {
			if features[featureID] == feature {
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
		return value.list
	}
	return eClass.getEAllStructuralFeatures(nil)
}

func (eClass *EClassExt) getEAllStructuralFeatures(path []*EClassExt) EList {
	if features := eClass.features(path); features != nil {
		return features.list
	}
	return nil
}

// features returns the features of the class, nil for a class of path (see eClassHierarchy)
func (eClass *EClassExt) features(path []*EClassExt) *eClassFeatures {
	if value, ok := eClass.allStructuralFeatures.load(); ok {
		return value
	}
	if slices.Contains(path, eClass) {
		return nil
	}
	path = append(path, eClass)
	return eClass.allStructuralFeatures.get(func() *eClassFeatures { return eClass.buildEAllStructuralFeatures(path) })
}

func (eClass *EClassExt) GetEAllAttributes() EList {
	if value, ok := eClass.allAttributes.load(); ok {
		return value.eAllAttributes
	}
	return eClass.attributes(nil).eAllAttributes
}

func (eClass *EClassExt) getEAllAttributes(path []*EClassExt) EList {
	if attributes := eClass.attributes(path); attributes != nil {
		return attributes.eAllAttributes
	}
	return nil
}

// attributes returns the attributes of the class, nil for a class of path (see eClassHierarchy)
func (eClass *EClassExt) attributes(path []*EClassExt) *eClassAttributes {
	if value, ok := eClass.allAttributes.load(); ok {
		return value
	}
	if slices.Contains(path, eClass) {
		return nil
	}
	path = append(path, eClass)
	return eClass.allAttributes.get(func() *eClassAttributes { return eClass.buildEAllAttributes(path) })
}

func (eClass *EClassExt) GetEAttributes() EList {
	if value, ok := eClass.allAttributes.load(); ok {
		return value.eAttributes
	}
	return eClass.attributes(nil).eAttributes
}

func (eClass *EClassExt) GetEIDAttribute() EAttribute {
	if value, ok := eClass.allAttributes.load(); ok {
		return value.eIDAttribute
	}
	return eClass.attributes(nil).eIDAttribute
}

func (eClass *EClassExt) GetEAllReferences() EList {
	if value, ok := eClass.allReferences.load(); ok {
		return value.eAllReferences
	}
	return eClass.references(nil).eAllReferences
}

func (eClass *EClassExt) getEAllReferences(path []*EClassExt) EList {
	if references := eClass.references(path); references != nil {
		return references.eAllReferences
	}
	return nil
}

// references returns the references of the class, nil for a class of path (see eClassHierarchy)
func (eClass *EClassExt) references(path []*EClassExt) *eClassReferences {
	if value, ok := eClass.allReferences.load(); ok {
		return value
	}
	if slices.Contains(path, eClass) {
		return nil
	}
	path = append(path, eClass)
	return eClass.allReferences.get(func() *eClassReferences { return eClass.buildEAllReferences(path) })
}

func (eClass *EClassExt) GetEReferences() EList {
	if value, ok := eClass.allReferences.load(); ok {
		return value.eReferences
	}
	return eClass.references(nil).eReferences
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
	return eClass.getEAllOperations(nil)
}

func (eClass *EClassExt) getEAllOperations(path []*EClassExt) EList {
	if value, ok := eClass.allOperations.load(); ok {
		return value
	}
	if slices.Contains(path, eClass) {
		return nil
	}
	path = append(path, eClass)
	return eClass.allOperations.get(func() EList { return eClass.buildEAllOperations(path) })
}

func (eClass *EClassExt) GetEAllSuperTypes() EList {
	if value, ok := eClass.allSuperTypes.load(); ok {
		return value
	}
	return eClass.getEAllSuperTypes(nil)
}

func (eClass *EClassExt) getEAllSuperTypes(path []*EClassExt) EList {
	if value, ok := eClass.allSuperTypes.load(); ok {
		return value
	}
	if slices.Contains(path, eClass) {
		return nil
	}
	path = append(path, eClass)
	return eClass.allSuperTypes.get(func() EList { return eClass.buildEAllSuperTypes(path) })
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
		return ok && isListSet(l.list)
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
	// so that EContents & ECrossReferences of an object visit the features of its super classes.
	// Like EMF, the derived references are included: EMF excludes them only when their values also
	// come from a feature map (a mixed class or a group annotation), which soft does not have.
	containments := []any{}
	crossReferences := []any{}
	for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
		ref, isRef := itFeature.Next().(EReference)
		if isRef {
			if ref.IsContainment() {
				containments = append(containments, ref)
			} else if !ref.IsContainer() {
				crossReferences = append(crossReferences, ref)
			}
		}
	}
	return &eClassFeaturesSubSet{
		eContainmentFeatures:    NewImmutableEList(containments),
		eCrossReferenceFeatures: NewImmutableEList(crossReferences),
	}
}

// buildEAllAttributes builds the attributes of the class, the last class of path
func (eClass *EClassExt) buildEAllAttributes(path []*EClassExt) *eClassAttributes {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
	attributes := []any{}
	allAttributes := []any{}
	var eIDAttribute EAttribute = nil
	seen := map[any]struct{}{}
	for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
		superAttributes := superValue(itClass.Next().(EClass), path, eClassHierarchy.getEAllAttributes, EClass.GetEAllAttributes)
		if superAttributes == nil {
			continue
		}
		for itAttribute := superAttributes.Iterator(); itAttribute.HasNext(); {
			attribute := itAttribute.Next().(EAttribute)
			if _, isSeen := seen[attribute]; isSeen {
				// inherited through another super type
				continue
			}
			seen[attribute] = struct{}{}
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

// buildEAllReferences builds the references of the class, the last class of path
func (eClass *EClassExt) buildEAllReferences(path []*EClassExt) *eClassReferences {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
	allReferences := []any{}
	references := []any{}
	seen := map[any]struct{}{}
	for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
		if superReferences := superValue(itClass.Next().(EClass), path, eClassHierarchy.getEAllReferences, EClass.GetEAllReferences); superReferences != nil {
			allReferences = appendUnique(allReferences, seen, superReferences.ToArray())
		}
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

// buildEAllOperations builds the operations of the class, the last class of path
func (eClass *EClassExt) buildEAllOperations(path []*EClassExt) EList {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
	allOperations := []any{}
	seen := map[any]struct{}{}
	for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
		if superOperations := superValue(itClass.Next().(EClass), path, eClassHierarchy.getEAllOperations, EClass.GetEAllOperations); superOperations != nil {
			allOperations = appendUnique(allOperations, seen, superOperations.ToArray())
		}
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

// buildEAllStructuralFeatures builds the structural features of the class, the last class of path
func (eClass *EClassExt) buildEAllStructuralFeatures(path []*EClassExt) *eClassFeatures {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
	allFeatures := []any{}
	seen := map[any]struct{}{}
	for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
		if superFeatures := superValue(itClass.Next().(EClass), path, eClassHierarchy.getEAllStructuralFeatures, EClass.GetEAllStructuralFeatures); superFeatures != nil {
			allFeatures = appendUnique(allFeatures, seen, superFeatures.ToArray())
		}
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
	slice := make([]EStructuralFeature, len(allFeatures))
	for i, f := range allFeatures {
		slice[i] = f.(EStructuralFeature)
	}
	return &eClassFeatures{list: NewImmutableEList(allFeatures), slice: slice}
}

// buildEAllSuperTypes builds the super types of the class, the last class of path
func (eClass *EClassExt) buildEAllSuperTypes(path []*EClassExt) EList {
	eClass.buildMutex.Lock()
	defer eClass.buildMutex.Unlock()
	allSuperTypes := []any{}
	seen := map[any]struct{}{}
	for itClass := eClass.GetESuperTypes().Iterator(); itClass.HasNext(); {
		superClass := itClass.Next().(EClass)
		if superTypes := superValue(superClass, path, eClassHierarchy.getEAllSuperTypes, EClass.GetEAllSuperTypes); superTypes != nil {
			allSuperTypes = appendUnique(allSuperTypes, seen, superTypes.ToArray())
		}
		allSuperTypes = appendUnique(allSuperTypes, seen, []any{superClass})
	}
	return NewImmutableEList(allSuperTypes)
}

// appendUnique appends the elements not seen yet: a class inheriting a super type through several
// of its super types (a diamond) has its features, operations and super types once, as in EMF.
func appendUnique(dst []any, seen map[any]struct{}, elements []any) []any {
	for _, element := range elements {
		if _, isSeen := seen[element]; !isSeen {
			seen[element] = struct{}{}
			dst = append(dst, element)
		}
	}
	return dst
}

func (eClass *EClassExt) getSubClasses() []eClassSubClasses {
	eClass.subClassesMutex.Lock()
	defer eClass.subClassesMutex.Unlock()
	return eClass.subClasses
}

// addSubClass and removeSubClass update the sub classes under the mutex, in a new slice: the
// slice returned by getSubClasses is never changed, setModified iterates it without lock.
func (eClass *EClassExt) addSubClass(subClass eClassSubClasses) {
	eClass.subClassesMutex.Lock()
	eClass.subClasses = append(slices.Clip(eClass.subClasses), subClass)
	eClass.subClassesMutex.Unlock()
}

func (eClass *EClassExt) removeSubClass(subClass eClassSubClasses) {
	eClass.subClassesMutex.Lock()
	if i := slices.Index(eClass.subClasses, subClass); i != -1 {
		eClass.subClasses = slices.Delete(slices.Clone(eClass.subClasses), i, i+1)
	}
	eClass.subClassesMutex.Unlock()
}

// setModified invalidates the values derived from featureID, in this class and its sub classes.
// A value derived from another one is invalidated with it.
func (eClass *EClassExt) setModified(featureID int) {
	eClass.setModifiedFrom(featureID, nil)
}

// setModifiedFrom is setModified in a sub class, visited holding the classes invalidated already
// by this modification: a class inheriting from the modified one through several super types (a
// diamond) or through a cycle of super types is invalidated once, and the recursion ends.
func (eClass *EClassExt) setModifiedFrom(featureID int, visited map[*EClassExt]struct{}) {
	if visited != nil {
		if _, isVisited := visited[eClass]; isVisited {
			return
		}
		visited[eClass] = struct{}{}
	}
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
	if subClasses := eClass.getSubClasses(); len(subClasses) != 0 {
		if visited == nil {
			visited = map[*EClassExt]struct{}{eClass: {}}
		}
		for _, s := range subClasses {
			s.setModifiedFrom(featureID, visited)
		}
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
