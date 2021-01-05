package ecore

type CompactEObjectContainer struct {
	CompactEObjectImpl
	container EObject
}

func (o *CompactEObjectContainer) ESetInternalContainer(newContainer EObject, newContainerFeatureID int) {
	o.container = newContainer
	o.flags = newContainerFeatureID<<16 | (o.flags & 0x00FF)
}

func (o *CompactEObjectContainer) EInternalContainer() EObject {
	return o.container
}
