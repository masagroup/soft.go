package ecore

// eAttributeExt is the extension of the model object 'EAttribute'
type eAttributeExt struct {
	eAttributeImpl
}

func newEAttributeExt() *eAttributeExt {
	eAttribute := new(eAttributeExt)
	eAttribute.SetInterfaces(eAttribute)
	eAttribute.Initialize()
	return eAttribute
}

func (eAttribute *eAttributeExt) GetEAttributeType() EDataType {
	return eAttribute.GetEType().(EDataType)
}

func (eAttribute *eAttributeExt) basicGetEAttributeType() EDataType {
	return eAttribute.basicGetEType().(EDataType)
}

func (eAttribute *eAttributeExt) SetID(newIsID bool) {
	eAttribute.eAttributeImpl.SetID(newIsID)
	eClass := eAttribute.GetEContainingClass()
	if eClass != nil {
		classExt := eClass.(*eClassExt)
		classExt.setModified(ECLASS__EATTRIBUTES)
	}

}
