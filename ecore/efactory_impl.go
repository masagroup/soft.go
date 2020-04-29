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

// eFactoryImpl is the implementation of the model object 'EFactory'
type eFactoryImpl struct {
	*eModelElementExt
	ePackage EPackage
}

// newEFactoryImpl is the constructor of a eFactoryImpl
func newEFactoryImpl() *eFactoryImpl {
	eFactory := new(eFactoryImpl)
	eFactory.eModelElementExt = newEModelElementExt()
	eFactory.SetInterfaces(eFactory)

	return eFactory
}

func (eFactory *eFactoryImpl) EStaticClass() EClass {
	return GetPackage().GetEFactory()
}

// ConvertToString default implementation
func (eFactory *eFactoryImpl) ConvertToString(EDataType, interface{}) string {
	panic("ConvertToString not implemented")
}

// Create default implementation
func (eFactory *eFactoryImpl) Create(EClass) EObject {
	panic("Create not implemented")
}

// CreateFromString default implementation
func (eFactory *eFactoryImpl) CreateFromString(EDataType, string) interface{} {
	panic("CreateFromString not implemented")
}

// GetEPackage get the value of ePackage
func (eFactory *eFactoryImpl) GetEPackage() EPackage {
	if eFactory.EContainerFeatureID() == EFACTORY__EPACKAGE {
		return eFactory.EContainer().(EPackage)
	}
	return nil
}

// SetEPackage set the value of ePackage
func (eFactory *eFactoryImpl) SetEPackage(newEPackage EPackage) {
	if newEPackage != eFactory.EContainer() || (newEPackage != nil && eFactory.EContainerFeatureID() != EFACTORY__EPACKAGE) {
		var notifications ENotificationChain
		if eFactory.EContainer() != nil {
			notifications = eFactory.EBasicRemoveFromContainer(notifications)
		}
		if newEPackageInternal, _ := newEPackage.(EObjectInternal); newEPackageInternal != nil {
			notifications = newEPackageInternal.EInverseAdd(eFactory.AsEObject(), EFACTORY__EPACKAGE, notifications)
		}
		notifications = eFactory.basicSetEPackage(newEPackage, notifications)
		if notifications != nil {
			notifications.Dispatch()
		}
	} else if eFactory.ENotificationRequired() {
		eFactory.ENotify(NewNotificationByFeatureID(eFactory, SET, EFACTORY__EPACKAGE, newEPackage, newEPackage, NO_INDEX))
	}
}

func (eFactory *eFactoryImpl) basicSetEPackage(newEPackage EPackage, msgs ENotificationChain) ENotificationChain {
	return eFactory.EBasicSetContainer(newEPackage, EFACTORY__EPACKAGE, msgs)
}

func (eFactory *eFactoryImpl) EGetFromID(featureID int, resolve, coreType bool) interface{} {
	switch featureID {
	case EFACTORY__EPACKAGE:
		return eFactory.GetEPackage()
	default:
		return eFactory.eModelElementExt.EGetFromID(featureID, resolve, coreType)
	}
}

func (eFactory *eFactoryImpl) ESetFromID(featureID int, newValue interface{}) {
	switch featureID {
	case EFACTORY__EPACKAGE:
		e := newValue.(EPackage)
		eFactory.SetEPackage(e)
	default:
		eFactory.eModelElementExt.ESetFromID(featureID, newValue)
	}
}

func (eFactory *eFactoryImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case EFACTORY__EPACKAGE:
		eFactory.SetEPackage(nil)
	default:
		eFactory.eModelElementExt.EUnsetFromID(featureID)
	}
}

func (eFactory *eFactoryImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case EFACTORY__EPACKAGE:
		return eFactory.GetEPackage() != nil
	default:
		return eFactory.eModelElementExt.EIsSetFromID(featureID)
	}
}

func (eFactory *eFactoryImpl) EInvokeFromID(operationID int, arguments EList) interface{} {
	switch operationID {
	case EFACTORY__CONVERT_TO_STRING_EDATATYPE_EJAVAOBJECT:
		return eFactory.ConvertToString(arguments.Get(0).(EDataType), arguments.Get(1))
	case EFACTORY__CREATE_ECLASS:
		return eFactory.Create(arguments.Get(0).(EClass))
	case EFACTORY__CREATE_FROM_STRING_EDATATYPE_ESTRING:
		return eFactory.CreateFromString(arguments.Get(0).(EDataType), arguments.Get(1).(string))
	default:
		return eFactory.eModelElementExt.EInvokeFromID(operationID, arguments)
	}
}

func (eFactory *eFactoryImpl) EBasicInverseAdd(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EFACTORY__EPACKAGE:
		msgs := notifications
		if eFactory.EContainer() != nil {
			msgs = eFactory.EBasicRemoveFromContainer(msgs)
		}
		return eFactory.basicSetEPackage(otherEnd.(EPackage), msgs)
	default:
		return eFactory.eModelElementExt.EBasicInverseAdd(otherEnd, featureID, notifications)
	}
}

func (eFactory *eFactoryImpl) EBasicInverseRemove(otherEnd EObject, featureID int, notifications ENotificationChain) ENotificationChain {
	switch featureID {
	case EFACTORY__EPACKAGE:
		return eFactory.basicSetEPackage(nil, notifications)
	default:
		return eFactory.eModelElementExt.EBasicInverseRemove(otherEnd, featureID, notifications)
	}
}
