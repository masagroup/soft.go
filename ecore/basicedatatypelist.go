package ecore

type basicEDataTypeList struct {
	BasicENotifyingList
	owner     EObjectInternal
	featureID int
}

// NewBasicEDataTypeList returns the list of the values of the many-valued attribute of ID featureID
// of owner. The list of an unsettable attribute is set from its first change until Unset
// (EUnsettableList).
func NewBasicEDataTypeList(owner EObjectInternal, featureID int, isUnique, unsettable bool) *basicEDataTypeList {
	l := new(basicEDataTypeList)
	l.interfaces = l
	l.data = []any{}
	l.owner = owner
	l.featureID = featureID
	l.isUnique = isUnique
	l.unsettable = unsettable
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
