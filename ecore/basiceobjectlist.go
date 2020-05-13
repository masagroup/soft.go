package ecore

type basicEObjectList struct {
	*BasicENotifyingList
	owner            EObjectInternal
	featureID        int
	inverseFeatureID int
	containment      bool
	inverse          bool
	opposite         bool
	proxies          bool
	unset            bool
}

func NewBasicEObjectList(owner EObjectInternal, featureID int, inverseFeatureID int, containment, inverse, opposite, proxies, unset bool) *basicEObjectList {
	l := new(basicEObjectList)
	l.BasicENotifyingList = NewBasicENotifyingList()
	l.interfaces = l
	l.owner = owner
	l.featureID = featureID
	l.inverseFeatureID = inverseFeatureID
	l.containment = containment
	l.inverse = inverse
	l.opposite = opposite
	l.proxies = proxies
	l.unset = unset
	return l
}

// GetNotifier ...
func (list *basicEObjectList) GetNotifier() ENotifier {
	return list.owner
}

// GetFeature ...
func (list *basicEObjectList) GetFeature() EStructuralFeature {
	if list.owner != nil {
		return list.owner.EClass().GetEStructuralFeature(list.featureID)
	}
	return nil
}

// GetFeatureID ...
func (list *basicEObjectList) GetFeatureID() int {
	return list.featureID
}

func (list *basicEObjectList) doGet(index int) interface{} {
	return list.resolve(index, list.BasicENotifyingList.doGet(index))
}

func (list *basicEObjectList) resolve(index int, object interface{}) interface{} {
	resolved := list.resolveProxy(object.(EObject))
	if resolved != object {
		list.basicEList.doSet(index, object)
		var notifications ENotificationChain
		if list.containment {
			notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(object, notifications)
			if resolvedInternal, _ := resolved.(EObjectInternal); resolvedInternal != nil && resolvedInternal.EInternalContainer() == nil {
				notifications = list.interfaces.(eNotifyingListInternal).inverseAdd(resolved, notifications)
			}
		}
		list.createAndDispatchNotification(notifications, RESOLVE, object, resolved, index)
	}
	return resolved
}

func (list *basicEObjectList) resolveProxy(eObject EObject) EObject {
	if list.proxies && eObject.EIsProxy() {
		return list.owner.(EObjectInternal).EResolveProxy(eObject)
	}
	return eObject
}

func (list *basicEObjectList) inverseAdd(object interface{}, notifications ENotificationChain) ENotificationChain {
	internal, _ := object.(EObjectInternal)
	if internal != nil && list.inverse {
		if list.opposite {
			return internal.EInverseAdd(list.owner, list.inverseFeatureID, notifications)
		} else {
			return internal.EInverseAdd(list.owner, EOPPOSITE_FEATURE_BASE-list.featureID, notifications)
		}
	}
	return notifications
}

func (list *basicEObjectList) inverseRemove(object interface{}, notifications ENotificationChain) ENotificationChain {
	internal, _ := object.(EObjectInternal)
	if internal != nil && list.inverse {
		if list.opposite {
			return internal.EInverseRemove(list.owner, list.inverseFeatureID, notifications)
		} else {
			return internal.EInverseRemove(list.owner, EOPPOSITE_FEATURE_BASE-list.featureID, notifications)
		}
	}
	return notifications
}