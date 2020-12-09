// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

// eStructuralFeatureExt is the extension of the model object 'EStructuralFeature'
type eStructuralFeatureExt struct {
	*eStructuralFeatureImpl
	defaultValue        interface{}
	defaultValueFactory EFactory
}

func newEStructuralFeatureExt() *eStructuralFeatureExt {
	eStructuralFeature := new(eStructuralFeatureExt)
	eStructuralFeature.eStructuralFeatureImpl = newEStructuralFeatureImpl()
	eStructuralFeature.interfaces = eStructuralFeature
	return eStructuralFeature
}

// GetDefaultValue get the value of defaultValue
func (eStructuralFeature *eStructuralFeatureExt) GetDefaultValue() interface{} {
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
func (eStructuralFeature *eStructuralFeatureExt) SetDefaultValue(newDefaultValue interface{}) {
	eType := eStructuralFeature.GetEType()
	if eDataType, _ := eType.(EDataType); eDataType != nil {
		factory := eDataType.GetEPackage().GetEFactoryInstance()
		literal := factory.ConvertToString(eDataType, newDefaultValue)
		eStructuralFeature.eStructuralFeatureImpl.SetDefaultValueLiteral(literal)
		eStructuralFeature.defaultValueFactory = nil // reset default value
	} else {
		panic("Cannot serialize value to object without an EDataType eType")
	}
}

// SetDefaultValueLiteral set the value of defaultValueLiteral
func (eStructuralFeature *eStructuralFeatureExt) SetDefaultValueLiteral(newDefaultValueLiteral string) {
	eStructuralFeature.defaultValueFactory = nil // reset default value
	eStructuralFeature.eStructuralFeatureImpl.SetDefaultValueLiteral(newDefaultValueLiteral)
}

func isBidirectional(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		return ref.GetEOpposite() != nil
	}
	return false
}

func isContainer(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		opposite := ref.GetEOpposite()
		if opposite != nil {
			return opposite.IsContainment()
		}
	}
	return false
}

func isContains(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		return ref.IsContainment()
	}
	return false
}

func isProxy(feature EStructuralFeature) bool {
	ref, isRef := feature.(EReference)
	if isRef {
		return ref.IsResolveProxies()
	}
	return false
}
