package ecore

type basicEDataTypeList struct {
	BasicENotifyingList
	owner     EObjectInternal
	featureID int
}

func NewBasicEDataTypeList(owner EObjectInternal, featureID int, isUnique bool) *basicEDataTypeList {
	l := new(basicEDataTypeList)
	l.interfaces = l
	l.data = []interface{}{}
	l.owner = owner
	l.featureID = featureID
	l.isUnique = isUnique
	return l
}

// GetNotifier ...
func (list *basicEDataTypeList) GetNotifier() ENotifier {
	return list.owner
}

// GetFeature ...
func (list *basicEDataTypeList) GetFeature() EStructuralFeature {
	if list.owner != nil {
		return list.owner.EClass().GetEStructuralFeature(list.featureID)
	}
	return nil
}

// GetFeatureID ...
func (list *basicEDataTypeList) GetFeatureID() int {
	return list.featureID
}
