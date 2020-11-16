package ecore

type EStoreEObjectImpl struct {
	*ReflectiveEObjectImpl
}

func NewEStoreEObjectImpl() *EStoreEObjectImpl {
	o := new(EStoreEObjectImpl)
	o.SetInterfaces(o)
	return o
}

func (o *EStoreEObjectImpl) AsStoreEObject() EStoreEObject {
	return o.GetInterfaces().(EStoreEObject)
}

func (o *EStoreEObjectImpl) EDynamicGet(dynamicFeatureID int) interface{} {
	return o.properties[dynamicFeatureID]
}

func (o *EStoreEObjectImpl) EDynamicSet(dynamicFeatureID int, newValue interface{}) {
	o.properties[dynamicFeatureID] = newValue
}

func (o *EStoreEObjectImpl) EDynamicUnset(dynamicFeatureID int) {
	o.properties[dynamicFeatureID] = nil
}
