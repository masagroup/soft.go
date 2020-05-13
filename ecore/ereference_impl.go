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

// eReferenceImpl is the implementation of the model object 'EReference'
type eReferenceImpl struct {
	*eStructuralFeatureExt
	eKeys            EList
	eOpposite        EReference
	eReferenceType   EClass
	isContainer      bool
	isContainment    bool
	isResolveProxies bool
}

// newEReferenceImpl is the constructor of a eReferenceImpl
func newEReferenceImpl() *eReferenceImpl {
	eReference := new(eReferenceImpl)
	eReference.eStructuralFeatureExt = newEStructuralFeatureExt()
	eReference.SetInterfaces(eReference)
	eReference.isContainment = false
	eReference.isResolveProxies = true

	return eReference
}

type eReferenceImplInitializers interface {
	initEKeys() EList
}

func (eReference *eReferenceImpl) getInitializers() eReferenceImplInitializers {
	return eReference.AsEObject().(eReferenceImplInitializers)
}

func (eReference *eReferenceImpl) EStaticClass() EClass {
	return GetPackage().GetEReference()
}

// GetEKeys get the value of eKeys
func (eReference *eReferenceImpl) GetEKeys() EList {
	if eReference.eKeys == nil {
		eReference.eKeys = eReference.getInitializers().initEKeys()
	}
	return eReference.eKeys
}

// GetEOpposite get the value of eOpposite
func (eReference *eReferenceImpl) GetEOpposite() EReference {
	if eReference.eOpposite != nil && eReference.eOpposite.EIsProxy() {
		oldEOpposite := eReference.eOpposite
		newEOpposite := eReference.EResolveProxy(oldEOpposite).(EReference)
		eReference.eOpposite = newEOpposite
		if newEOpposite != oldEOpposite {
			if eReference.ENotificationRequired() {
				eReference.ENotify(NewNotificationByFeatureID(eReference, RESOLVE, EREFERENCE__EOPPOSITE, oldEOpposite, newEOpposite, NO_INDEX))
			}
		}
	}
	return eReference.eOpposite
}

func (eReference *eReferenceImpl) basicGetEOpposite() EReference {
	return eReference.eOpposite
}

// SetEOpposite set the value of eOpposite
func (eReference *eReferenceImpl) SetEOpposite(newEOpposite EReference) {
	oldEOpposite := eReference.eOpposite
	eReference.eOpposite = newEOpposite
	if eReference.ENotificationRequired() {
		eReference.ENotify(NewNotificationByFeatureID(eReference.AsEObject(), SET, EREFERENCE__EOPPOSITE, oldEOpposite, newEOpposite, NO_INDEX))
	}
}

// GetEReferenceType get the value of eReferenceType
func (eReference *eReferenceImpl) GetEReferenceType() EClass {
	panic("GetEReferenceType not implemented")
}

func (eReference *eReferenceImpl) basicGetEReferenceType() EClass {
	panic("GetEReferenceType not implemented")
}

// IsContainer get the value of isContainer
func (eReference *eReferenceImpl) IsContainer() bool {
	panic("IsContainer not implemented")
}

// IsContainment get the value of isContainment
func (eReference *eReferenceImpl) IsContainment() bool {
	return eReference.isContainment
}

// SetContainment set the value of isContainment
func (eReference *eReferenceImpl) SetContainment(newIsContainment bool) {
	oldIsContainment := eReference.isContainment
	eReference.isContainment = newIsContainment
	if eReference.ENotificationRequired() {
		eReference.ENotify(NewNotificationByFeatureID(eReference.AsEObject(), SET, EREFERENCE__CONTAINMENT, oldIsContainment, newIsContainment, NO_INDEX))
	}
}

// IsResolveProxies get the value of isResolveProxies
func (eReference *eReferenceImpl) IsResolveProxies() bool {
	return eReference.isResolveProxies
}

// SetResolveProxies set the value of isResolveProxies
func (eReference *eReferenceImpl) SetResolveProxies(newIsResolveProxies bool) {
	oldIsResolveProxies := eReference.isResolveProxies
	eReference.isResolveProxies = newIsResolveProxies
	if eReference.ENotificationRequired() {
		eReference.ENotify(NewNotificationByFeatureID(eReference.AsEObject(), SET, EREFERENCE__RESOLVE_PROXIES, oldIsResolveProxies, newIsResolveProxies, NO_INDEX))
	}
}

func (eReference *eReferenceImpl) initEKeys() EList {
	return NewBasicEObjectList(eReference.AsEObjectInternal(), EREFERENCE__EKEYS, -1, false, false, false, true, false)
}

func (eReference *eReferenceImpl) EGetFromID(featureID int, resolve, coreType bool) interface{} {
	switch featureID {
	case EREFERENCE__CONTAINER:
		return eReference.IsContainer()
	case EREFERENCE__CONTAINMENT:
		return eReference.IsContainment()
	case EREFERENCE__EKEYS:
		return eReference.GetEKeys()
	case EREFERENCE__EOPPOSITE:
		return eReference.GetEOpposite()
	case EREFERENCE__EREFERENCE_TYPE:
		return eReference.GetEReferenceType()
	case EREFERENCE__RESOLVE_PROXIES:
		return eReference.IsResolveProxies()
	default:
		return eReference.eStructuralFeatureExt.EGetFromID(featureID, resolve, coreType)
	}
}

func (eReference *eReferenceImpl) ESetFromID(featureID int, newValue interface{}) {
	switch featureID {
	case EREFERENCE__CONTAINMENT:
		c := newValue.(bool)
		eReference.SetContainment(c)
	case EREFERENCE__EKEYS:
		e := newValue.(EList)
		eReference.GetEKeys().Clear()
		eReference.GetEKeys().AddAll(e)
	case EREFERENCE__EOPPOSITE:
		e := newValue.(EReference)
		eReference.SetEOpposite(e)
	case EREFERENCE__RESOLVE_PROXIES:
		r := newValue.(bool)
		eReference.SetResolveProxies(r)
	default:
		eReference.eStructuralFeatureExt.ESetFromID(featureID, newValue)
	}
}

func (eReference *eReferenceImpl) EUnsetFromID(featureID int) {
	switch featureID {
	case EREFERENCE__CONTAINMENT:
		eReference.SetContainment(false)
	case EREFERENCE__EKEYS:
		eReference.GetEKeys().Clear()
	case EREFERENCE__EOPPOSITE:
		eReference.SetEOpposite(nil)
	case EREFERENCE__RESOLVE_PROXIES:
		eReference.SetResolveProxies(true)
	default:
		eReference.eStructuralFeatureExt.EUnsetFromID(featureID)
	}
}

func (eReference *eReferenceImpl) EIsSetFromID(featureID int) bool {
	switch featureID {
	case EREFERENCE__CONTAINER:
		return eReference.IsContainer() != false
	case EREFERENCE__CONTAINMENT:
		return eReference.isContainment != false
	case EREFERENCE__EKEYS:
		return eReference.eKeys != nil && eReference.eKeys.Size() != 0
	case EREFERENCE__EOPPOSITE:
		return eReference.eOpposite != nil
	case EREFERENCE__EREFERENCE_TYPE:
		return eReference.GetEReferenceType() != nil
	case EREFERENCE__RESOLVE_PROXIES:
		return eReference.isResolveProxies != true
	default:
		return eReference.eStructuralFeatureExt.EIsSetFromID(featureID)
	}
}
