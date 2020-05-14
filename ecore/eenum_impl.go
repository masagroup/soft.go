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

// eEnumImpl is the implementation of the model object 'EEnum'
type eEnumImpl struct {
	*eDataTypeImpl
	eLiterals EList
}

// newEEnumImpl is the constructor of a eEnumImpl
func newEEnumImpl() *eEnumImpl {
	eEnum := new(eEnumImpl)
	eEnum.eDataTypeImpl = newEDataTypeImpl()
	eEnum.SetInterfaces(eEnum)

	return eEnum
}

type eEnumImplInitializers interface {
	initELiterals() EList
}

func (eEnum *eEnumImpl) getInitializers() eEnumImplInitializers {
	return eEnum.AsEObject().(eEnumImplInitializers)
}

func (eEnum *eEnumImpl) asEEnum() EEnum {
	return eEnum.GetInterfaces().(EEnum)
}

func (eEnum *eEnumImpl) EStaticClass() EClass {
	return GetPackage().GetEEnum()
}

// GetEEnumLiteralByLiteral default implementation
func (eEnum *eEnumImpl) GetEEnumLiteralByLiteral(string) EEnumLiteral {
	panic("GetEEnumLiteralByLiteral not implemented")
}

// GetEEnumLiteralByName default implementation
func (eEnum *eEnumImpl) GetEEnumLiteralByName(string) EEnumLiteral {
	panic("GetEEnumLiteralByName not implemented")
}

// GetEEnumLiteralByValue default implementation
func (eEnum *eEnumImpl) GetEEnumLiteralByValue(int) EEnumLiteral {
	panic("GetEEnumLiteralByValue not implemented")
}

// GetELiterals get the value of eLiterals
func (eEnum *eEnumImpl) GetELiterals() EList {
	if eEnum.eLiterals == nil {
		eEnum.eLiterals = eEnum.getInitializers().initELiterals()
	}
	return eEnum.eLiterals
}

func (eEnum *eEnumImpl) initELiterals() EList {
	return NewBasicEObjectList(eEnum.AsEObjectInternal(), EENUM__ELITERALS, EENUM_LITERAL__EENUM, true, true, true, false, false)
}

func (eEnum *eEnumImpl) EGetFromID(featureID int, resolve bool) interface{} {
	switch featureID {
	case EENUM__ELITERALS:
		return eEnum.asEEnum().GetELiterals()
	default:
		return eEnum.eDataTypeImpl.EGetFromID(featureID, resolve)
	}
}

func (eEnum *eEnumImpl) ESetFromID(featureID int, newValue interface{}) {
	switch featureID {
	case EENUM__ELITERALS:
		list := eEnum.asEEnum().GetELiterals()
		list.Clear()
		list.AddAll(newValue.(EList))
	default:
		eEnum.eDataTypeImpl.ESetFromID(featureID, newValue)
	}
}

func (eEnum *eEnumImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case EENUM__ELITERALS:
		eEnum.asEEnum().GetELiterals().Clear()
	default:
		eEnum.eDataTypeImpl.EUnsetFromID(featureID)
	}
}

func (eEnum *eEnumImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case EENUM__ELITERALS:
		return eEnum.eLiterals != nil && eEnum.eLiterals.Size() != 0
	default:
		return eEnum.eDataTypeImpl.EIsSetFromID(featureID)
	}
}

func (eEnum *eEnumImpl) EInvokeFromID(operationID int, arguments EList) interface{} {
	switch operationID {
	case EENUM__GET_EENUM_LITERAL_BY_LITERAL_ESTRING:
		return eEnum.asEEnum().GetEEnumLiteralByLiteral(arguments.Get(0).(string))
	case EENUM__GET_EENUM_LITERAL_ESTRING:
		return eEnum.asEEnum().GetEEnumLiteralByName(arguments.Get(0).(string))
	case EENUM__GET_EENUM_LITERAL_EINT:
		return eEnum.asEEnum().GetEEnumLiteralByValue(arguments.Get(0).(int))
	default:
		return eEnum.eDataTypeImpl.EInvokeFromID(operationID, arguments)
	}
}

func (eEnum *eEnumImpl) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EENUM__ELITERALS:
		list := eEnum.GetELiterals().(ENotifyingList)
		return list.AddWithNotification(otherEnd, notifications)
	default:
		return eEnum.eDataTypeImpl.EBasicInverseAdd(otherEnd, featureID, notifications)
	}
}

func (eEnum *eEnumImpl) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EENUM__ELITERALS:
		list := eEnum.GetELiterals().(ENotifyingList)
		return list.RemoveWithNotification(otherEnd, notifications)
	default:
		return eEnum.eDataTypeImpl.EBasicInverseRemove(otherEnd, featureID, notifications)
	}
}
