package ecore

type EStoreEObjectImpl struct {
	*ReflectiveEObjectImpl
	isCaching bool
}

func NewEStoreEObjectImpl() *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.isCaching = true
	o.SetInterfaces(o)
	return o
}

func (o *EStoreEObjectImpl) AsStoreEObject() EStoreEObject {
	return o.GetInterfaces().(EStoreEObject)
}

func (o *EStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) interface{} {
	result := o.properties[dynamicFeatureID]
	if result == nil {
		eFeature := o.eDynamicFeature(dynamicFeatureID)
		if !eFeature.IsTransient() {
			if eFeature.IsMany() {
				result = o.createList(eFeature)
				o.properties[dynamicFeatureID] = result
			} else {
				result = o.AsStoreEObject().EStore().Get(o.AsEObject(), eFeature, NO_INDEX)
				if o.isCaching {
					o.properties[dynamicFeatureID] = result
				}
			}
		}
	}
	return result
}

func (o *EStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, value interface{}) {
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	if eFeature.IsTransient() {
		o.properties[dynamicFeatureID] = value
	} else {
		o.AsStoreEObject().EStore().Set(o.AsEObject(), eFeature, NO_INDEX, value)
		if o.isCaching {
			o.properties[dynamicFeatureID] = value
		}
	}
}

func (o *EStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	eFeature := o.eDynamicFeature(dynamicFeatureID)
	if eFeature.IsTransient() {
		o.properties[dynamicFeatureID] = nil
	} else {
		o.AsStoreEObject().EStore().UnSet(o.AsEObject(), eFeature)
		if o.isCaching {
			o.properties[dynamicFeatureID] = nil
		}
	}
}

func (o *EStoreEObjectImpl) eDynamicFeature(dynamicFeatureID int) EStructuralFeature {
	return o.EClass().GetEStructuralFeature(o.EStaticFeatureCount() + dynamicFeatureID)
}

func (o *EStoreEObjectImpl) createList(eFeature EStructuralFeature) EList {
	return nil
}
