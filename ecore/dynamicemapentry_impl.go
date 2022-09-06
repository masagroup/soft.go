package ecore

type DynamicEMapEntryImpl struct {
	DynamicEObjectImpl
	keyFeature   EStructuralFeature
	valueFeature EStructuralFeature
}

func NewDynamicEMapEntryImpl() *DynamicEMapEntryImpl {
	o := new(DynamicEMapEntryImpl)
	o.SetInterfaces(o)
	o.Initialize()
	return o
}

func (o *DynamicEMapEntryImpl) SetEClass(class EClass) {
	o.DynamicEObjectImpl.SetEClass(class)
	o.keyFeature = class.GetEStructuralFeatureFromName("key")
	o.valueFeature = class.GetEStructuralFeatureFromName("value")
}

func (o *DynamicEMapEntryImpl) GetKey() any {
	return o.EGet(o.keyFeature)
}

func (o *DynamicEMapEntryImpl) SetKey(key any) {
	o.ESet(o.keyFeature, key)
}

func (o *DynamicEMapEntryImpl) GetValue() any {
	return o.EGet(o.valueFeature)
}

func (o *DynamicEMapEntryImpl) SetValue(value any) {
	o.ESet(o.valueFeature, value)
}
