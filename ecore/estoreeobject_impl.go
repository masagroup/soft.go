package ecore

type EStoreEObjectImpl struct {
	*ReflectiveEObjectImpl
	isCaching bool
}

func NewEStoreEObjectImpl(isCaching bool) *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.ReflectiveEObjectImpl = NewReflectiveEObjectImpl()
	o.isCaching = isCaching
	o.SetInterfaces(o)
	return o
}

func (o *EStoreEObjectImpl) AsEStoreEObject() EStoreEObject {
	return o.GetInterfaces().(EStoreEObject)
}

func (o *EStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) interface{} {
	result := o.getProperties()[dynamicFeatureID]
	if result == nil {
		eFeature := o.eDynamicFeature(dynamicFeatureID)
		if !eFeature.IsTransient() {
			if eFeature.IsMany() {
				result = o.createList(eFeature)
				o.getProperties()[dynamicFeatureID] = result
			} else {
				result = o.AsEStoreEObject().EStore().Get(o.AsEObject(), eFeature, NO_INDEX)
				if o.isCaching {
					o.getProperties()[dynamicFeatureID] = result
				}
			}
		}
	}
	return result
}

func (o *EStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value interface{}) {
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	if eFeature.IsTransient() {
		o.getProperties()[dynamicFeatureID] = value
	} else {
		o.AsEStoreEObject().EStore().Set(o.AsEObject(), eFeature, NO_INDEX, value)
		if o.isCaching {
			o.getProperties()[dynamicFeatureID] = value
		}
	}
}

func (o *EStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	if eFeature.IsTransient() {
		o.getProperties()[dynamicFeatureID] = nil
	} else {
		o.AsEStoreEObject().EStore().UnSet(o.AsEObject(), eFeature)
		o.getProperties()[dynamicFeatureID] = nil
	}
}

func (o *EStoreEObjectImpl) eDynamicFeature(dynamicFeatureID int) EStructuralFeature {
	return o.EClass().GetEStructuralFeature(o.EStaticFeatureCount() + dynamicFeatureID)
}

func (o *EStoreEObjectImpl) createList(eFeature EStructuralFeature) EList {
	return NewBasicEStoreList(o.AsEObject(), eFeature, o.AsEStoreEObject().EStore())
}
