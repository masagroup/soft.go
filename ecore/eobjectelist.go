package ecore

type eObjectEList struct {
	*ENotifyingListImpl
	owner            EObjectInternal
	featureID        int
	inverseFeatureID int
	containment      bool
	inverse          bool
	opposite         bool
	proxies          bool
	unset            bool
}

func NewEObjectEList(owner EObjectInternal, featureID int, inverseFeatureID int, containment, inverse, opposite, proxies, unset bool) *eObjectEList {
	l := new(eObjectEList)
	l.ENotifyingListImpl = NewENotifyingListImpl()
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
func (list *eObjectEList) GetNotifier() ENotifier {
	return list.owner
}

// GetFeature ...
func (list *eObjectEList) GetFeature() EStructuralFeature {
	if list.owner != nil {
		return list.owner.EClass().GetEStructuralFeature(list.featureID)
	}
	return nil
}

// GetFeatureID ...
func (list *eObjectEList) GetFeatureID() int {
	return list.featureID
}

func (list *eObjectEList) doGet(index int) interface{} {
	return list.resolve(index, list.ENotifyingListImpl.doGet(index))
}

func (list *eObjectEList) resolve(index int, object interface{}) interface{} {
	resolved := list.resolveProxy(object.(EObject))
	if resolved != object {
		list.basicEList.doSet(index, object)
		var notifications ENotificationChain
		if list.containment {
			notifications = list.interfaces.(eNotifyingListInternal).inverseRemove(object, notifications)
			notifications = list.interfaces.(eNotifyingListInternal).inverseAdd(resolved, notifications)
		}
		list.createAndDispatchNotification(notifications, RESOLVE, object, resolved, index)
	}
	return resolved
}

func (list *eObjectEList) resolveProxy(eObject EObject) EObject {
	if list.proxies && eObject.EIsProxy() {
		return list.owner.(EObjectInternal).EResolveProxy(eObject)
	}
	return eObject
}

func (list *eObjectEList) inverseAdd(object interface{}, notifications ENotificationChain) ENotificationChain {
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

func (list *eObjectEList) inverseRemove(object interface{}, notifications ENotificationChain) ENotificationChain {
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
