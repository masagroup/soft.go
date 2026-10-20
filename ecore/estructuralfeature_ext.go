// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

import "sync/atomic"

// EStructuralFeatureExt is the extension of the model object 'EStructuralFeature'
type EStructuralFeatureExt struct {
	EStructuralFeatureImpl
	// defaultValue is the default value created from the literal by a factory: published in one
	// atomic pointer, as GetDefaultValue is read by concurrent goroutines
	defaultValue atomic.Pointer[structuralFeatureDefaultValue]
}

// structuralFeatureDefaultValue is a default value and the factory that created it. It is never
// changed once published.
type structuralFeatureDefaultValue struct {
	factory EFactory
	value   any
}

func newEStructuralFeatureExt() *EStructuralFeatureExt {
	eStructuralFeature := new(EStructuralFeatureExt)
	eStructuralFeature.SetInterfaces(eStructuralFeature)
	eStructuralFeature.Initialize()
	return eStructuralFeature
}

// GetDefaultValue get the value of defaultValue
func (eStructuralFeature *EStructuralFeatureExt) GetDefaultValue() any {
	eType := eStructuralFeature.GetEType()
	defaultValueLiteral := eStructuralFeature.GetDefaultValueLiteral()
	if eType != nil && len(defaultValueLiteral) == 0 {
		if eStructuralFeature.IsMany() {
			return nil
		} else {
			return eType.GetDefaultValue()
		}
	} else if eDataType, _ := eType.(EDataType); eDataType != nil {
		current := eStructuralFeature.defaultValue.Load()
		if ePackage := eType.GetEPackage(); ePackage != nil {
			if factory := ePackage.GetEFactoryInstance(); current == nil || factory != current.factory {
				next := &structuralFeatureDefaultValue{factory: factory}
				if eDataType.IsSerializable() {
					next.value = createDefaultValue(factory, eDataType, defaultValueLiteral)
				}
				eStructuralFeature.defaultValue.Store(next)
				return next.value
			}
		}
		if current != nil {
			return current.value
		}
	}
	return nil
}

// createDefaultValue creates the default value of eDataType from its literal: an invalid literal,
// on which the factory panics, gives the default value of the data type, as the default value is
// read at each access to an unset feature.
func createDefaultValue(factory EFactory, eDataType EDataType, literal string) (value any) {
	defer func() {
		if r := recover(); r != nil {
			value = eDataType.GetDefaultValue()
		}
	}()
	return factory.CreateFromString(eDataType, literal)
}

// SetDefaultValue set the value of defaultValue
func (eStructuralFeature *EStructuralFeatureExt) SetDefaultValue(newDefaultValue any) {
	eType := eStructuralFeature.GetEType()
	if eDataType, _ := eType.(EDataType); eDataType != nil {
		factory := eDataType.GetEPackage().GetEFactoryInstance()
		literal := factory.ConvertToString(eDataType, newDefaultValue)
		eStructuralFeature.EStructuralFeatureImpl.SetDefaultValueLiteral(literal)
		eStructuralFeature.defaultValue.Store(nil) // reset default value
	} else {
		panic("Cannot serialize value to object without an EDataType eType")
	}
}

// SetDefaultValueLiteral set the value of defaultValueLiteral
func (eStructuralFeature *EStructuralFeatureExt) SetDefaultValueLiteral(newDefaultValueLiteral string) {
	eStructuralFeature.defaultValue.Store(nil) // reset default value
	eStructuralFeature.EStructuralFeatureImpl.SetDefaultValueLiteral(newDefaultValueLiteral)
}

// GetFeatureID get the value of featureID
func (eStructuralFeature *EStructuralFeatureExt) GetFeatureID() int {
	return eStructuralFeature.featureID
}

// SetFeatureID set the value of featureID
func (eStructuralFeature *EStructuralFeatureExt) SetFeatureID(newFeatureID int) {
	eStructuralFeature.featureID = newFeatureID
}

func IsBidirectional(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		return ref.GetEOpposite() != nil
	}
	return false
}

func IsContainer(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		opposite := ref.GetEOpposite()
		if opposite != nil {
			return opposite.IsContainment()
		}
	}
	return false
}

func IsContains(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		return ref.IsContainment()
	}
	return false
}

func IsProxy(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		return ref.IsResolveProxies()
	}
	return false
}

func IsMapType(feature EStructuralFeature) bool {
	if eClass, _ := feature.GetEType().(EClass); eClass != nil {
		return IsMapEntry(eClass)
	}
	return false
}
