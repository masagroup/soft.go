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

type ePackageExtAdapter struct {
	AbstractEAdapter
	pack *EPackageExt
}

func (a *ePackageExtAdapter) NotifyChanged(notification ENotification) {
	eventType := notification.GetEventType()
	if eventType != REMOVING_ADAPTER {
		featureID := notification.GetFeatureID()
		if featureID == EPACKAGE__ECLASSIFIERS {
			a.pack.classifierMutex.Lock()
			a.pack.classifierCache = nil
			a.pack.classifierMutex.Unlock()
		}
	}
}

// EPackageExt is the extension of the model object 'EFactory'
type EPackageExt struct {
	EPackageImpl
	adapter         EAdapter
	classifierCache map[string]EClassifier
	classifierMutex sync.Mutex
}

func newEPackageExt() *EPackageExt {
	pack := new(EPackageExt)
	pack.SetInterfaces(pack)
	pack.Initialize()
	return pack
}

func (pack *EPackageExt) Initialize() {
	pack.EPackageImpl.Initialize()
	pack.SetEFactoryInstance(newEFactoryExt())
	pack.adapter = &ePackageExtAdapter{pack: pack}
	pack.EAdapters().Add(pack.adapter)
}

func (pack *EPackageExt) GetEClassifier(classifier string) EClassifier {
	// retrieve map
	pack.classifierMutex.Lock()
	defer pack.classifierMutex.Unlock()

	// initialize map if needed
	if pack.classifierCache == nil {
		pack.classifierCache = map[string]EClassifier{}
		for itClassifier := pack.GetEClassifiers().Iterator(); itClassifier.HasNext(); {
			classifier := itClassifier.Next().(EClassifier)
			pack.classifierCache[classifier.GetName()] = classifier
		}
	}

	// read map value
	return pack.classifierCache[classifier]
}

func (pack *EPackageExt) CreateResource() EResource {
	resource := pack.EPackageImpl.EResource()
	if resource == nil {
		uri := NewURI(pack.GetNsURI())
		resource = NewEResourceImpl()
		resource.SetURI(uri)
		resource.GetContents().Add(pack)
	}
	return resource
}

func (pack *EPackageExt) InitEClass(eClass EClass, name string, instanceTypeName string, isAbstract bool, isInterface bool) {
	eClass.SetName(name)
	eClass.SetAbstract(isAbstract)
	eClass.SetInterface(isInterface)
	eClass.SetInstanceTypeName(instanceTypeName)
}

func (pack *EPackageExt) initEStructuralFeature(aFeature EStructuralFeature, aClassifier EClassifier, name, defaultValue string, lowerBound, upperBound int, isTransient, isVolatile, isChangeable, isUnSettable, isUnique, isDerived, isOrdered bool) {
	aFeature.SetName(name)
	aFeature.SetEType(aClassifier)
	aFeature.SetDefaultValueLiteral(defaultValue)
	aFeature.SetLowerBound(lowerBound)
	aFeature.SetUpperBound(upperBound)
	aFeature.SetTransient(isTransient)
	aFeature.SetVolatile(isVolatile)
	aFeature.SetChangeable(isChangeable)
	aFeature.SetUnsettable(isUnSettable)
	aFeature.SetUnique(isUnique)
	aFeature.SetDerived(isDerived)
	aFeature.SetOrdered(isOrdered)
}

func (pack *EPackageExt) InitEAttribute(aAttribute EAttribute, aType EClassifier, name, defaultValue string, lowerBound, upperBound int, isTransient, isVolatile, isChangeable, isUnSettable, isUnique, isDerived, isOrdered, isID bool) {
	pack.initEStructuralFeature(aAttribute, aType, name, defaultValue, lowerBound, upperBound, isTransient, isVolatile, isChangeable, isUnSettable, isUnique, isDerived, isOrdered)
	aAttribute.SetID(isID)
}

func (pack *EPackageExt) InitEReference(aReference EReference, aType EClassifier, aOtherEnd EReference, name, defaultValue string, lowerBound, upperBound int, isTransient, isVolatile, isChangeable, isContainment, isResolveProxies, isUnSettable, isUnique, isDerived, isOrdered bool) {
	pack.initEStructuralFeature(aReference, aType, name, defaultValue, lowerBound, upperBound, isTransient, isVolatile, isChangeable, isUnSettable, isUnique, isDerived, isOrdered)
	aReference.SetContainment(isContainment)
	aReference.SetResolveProxies(isResolveProxies)
	aReference.SetEOpposite(aOtherEnd)
}

func (pack *EPackageExt) InitEOperation(aOperation EOperation, aType EClassifier, name string, lowerBound, upperBound int, isUnique, isOrdered bool) {
	aOperation.SetName(name)
	aOperation.SetEType(aType)
	aOperation.SetLowerBound(lowerBound)
	aOperation.SetUpperBound(upperBound)
	aOperation.SetUnique(isUnique)
	aOperation.SetOrdered(isOrdered)
}

func (pack *EPackageExt) AddEParameter(aOperation EOperation, aType EClassifier, name string, lowerBound, upperBound int, isUnique, isOrdered bool) {
	parameter := GetFactory().CreateEParameterFromContainer(aOperation)
	parameter.SetName(name)
	parameter.SetEType(aType)
	parameter.SetLowerBound(lowerBound)
	parameter.SetUpperBound(upperBound)
	parameter.SetUnique(isUnique)
	parameter.SetOrdered(isOrdered)
}

func (pack *EPackageExt) initEClassifier(aClassifier EClassifier, name, instanceTypeName string) {
	aClassifier.SetName(name)
	aClassifier.SetInstanceTypeName(instanceTypeName)
}

func (pack *EPackageExt) InitEDataType(aDataType EDataType, name, instanceTypeName, defaultValue string, isSerializable bool) {
	pack.initEClassifier(aDataType, name, instanceTypeName)
	aDataType.SetSerializable(isSerializable)
	if len(defaultValue) > 0 {
		aDataType.(EDataTypeInternal).SetDefaultValue(pack.GetEFactoryInstance().CreateFromString(aDataType, defaultValue))
	}
}

func (pack *EPackageExt) InitEEnum(aEnum EEnum, name, instanceTypeName string) {
	pack.initEClassifier(aEnum, name, instanceTypeName)
}

func (pack *EPackageExt) AddEEnumLiteral(aEnum EEnum, name, literal string, value int, instance any) {
	enumLiteral := GetFactory().CreateEEnumLiteralFromContainer(aEnum)
	enumLiteral.SetName(name)
	enumLiteral.SetLiteral(literal)
	enumLiteral.SetValue(value)
	enumLiteral.SetInstance(instance)
}
