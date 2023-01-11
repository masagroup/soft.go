// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

// EStructuralFeatureExt is the extension of the model object 'EStructuralFeature'
type EStructuralFeatureExt struct {
	EStructuralFeatureImpl
	defaultValue        any
	defaultValueFactory EFactory
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
		if ePackage := eType.GetEPackage(); ePackage != nil {
			if factory := ePackage.GetEFactoryInstance(); factory != eStructuralFeature.defaultValueFactory {
				if eDataType.IsSerializable() {
					eStructuralFeature.defaultValue = factory.CreateFromString(eDataType, defaultValueLiteral)
				}
				eStructuralFeature.defaultValueFactory = factory
			}
		}
		return eStructuralFeature.defaultValue
	}
	return nil
}

// SetDefaultValue set the value of defaultValue
func (eStructuralFeature *EStructuralFeatureExt) SetDefaultValue(newDefaultValue any) {
	eType := eStructuralFeature.GetEType()
	if eDataType, _ := eType.(EDataType); eDataType != nil {
		factory := eDataType.GetEPackage().GetEFactoryInstance()
		literal := factory.ConvertToString(eDataType, newDefaultValue)
		eStructuralFeature.EStructuralFeatureImpl.SetDefaultValueLiteral(literal)
		eStructuralFeature.defaultValueFactory = nil // reset default value
	} else {
		panic("Cannot serialize value to object without an EDataType eType")
	}
}

// SetDefaultValueLiteral set the value of defaultValueLiteral
func (eStructuralFeature *EStructuralFeatureExt) SetDefaultValueLiteral(newDefaultValueLiteral string) {
	eStructuralFeature.defaultValueFactory = nil // reset default value
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
