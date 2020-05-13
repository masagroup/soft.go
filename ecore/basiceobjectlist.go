package ecore

import "strconv"

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

// GetUnResolvedList ...
func (list *basicEObjectList) GetUnResolvedList() EList {
	if list.proxies {
		u := new(unResolvedEList)
		u.delegate = list
		return u
	}
	return list
}

func (list *basicEObjectList) IndexOf(elem interface{}) int {
	if list.proxies {
		for i, value := range list.data {
			if value == elem || list.resolve(i, value) == elem {
				return i
			}
		}
		return -1
	}
	return list.basicEList.IndexOf(elem)
}

func (list *basicEObjectList) doGet(index int) interface{} {
	return list.resolve(index, list.basicEList.doGet(index))
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

type unResolvedEList struct {
	delegate *basicEObjectList
}

func (l *unResolvedEList) Add(elem interface{}) bool {
	if l.delegate.isUnique && l.Contains(elem) {
		return false
	}
	l.delegate.interfaces.(abstractEList).doAdd(elem)
	return true
}

// AddAll elements of an list in the current one
func (l *unResolvedEList) AddAll(list EList) bool {
	if l.delegate.isUnique {
		list = getNonDuplicates(list, l)
		if list.Size() == 0 {
			return false
		}
	}
	l.delegate.interfaces.(abstractEList).doAddAll(list)
	return true
}

// Insert an element in the list
func (l *unResolvedEList) Insert(index int, elem interface{}) bool {
	if index < 0 || index > l.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(l.Size()))
	}
	if l.delegate.isUnique && l.Contains(elem) {
		return false
	}
	l.delegate.interfaces.(abstractEList).doInsert(index, elem)
	return true
}

// InsertAll element of an list at a given position
func (l *unResolvedEList) InsertAll(index int, list EList) bool {
	if index < 0 || index > l.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(l.Size()))
	}
	if l.delegate.isUnique {
		list = getNonDuplicates(list, l)
		if list.Size() == 0 {
			return false
		}
	}
	l.delegate.interfaces.(abstractEList).doInsertAll(index, list)
	return true
}

// Move an element to the given index
func (l *unResolvedEList) MoveObject(newIndex int, elem interface{}) {
	oldIndex := l.IndexOf(elem)
	if oldIndex == -1 {
		panic("Object not found")
	}
	l.delegate.Move(oldIndex, newIndex)
}

// Swap move an element from oldIndex to newIndex
func (l *unResolvedEList) Move(oldIndex, newIndex int) interface{} {
	return l.delegate.Move(oldIndex, newIndex)
}

// RemoveAt remove an element at a given position
func (l *unResolvedEList) RemoveAt(index int) interface{} {
	return l.delegate.RemoveAt(index)
}

// Remove an element in an list
func (l *unResolvedEList) Remove(elem interface{}) bool {
	index := l.IndexOf(elem)
	if index == -1 {
		return false
	}
	l.delegate.RemoveAt(index)
	return true
}

// Get an element of the list
func (l *unResolvedEList) Get(index int) interface{} {
	if index < 0 || index >= l.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(l.Size()))
	}
	return l.delegate.data[index]
}

// Set an element of the list
func (l *unResolvedEList) Set(index int, elem interface{}) interface{} {
	if index < 0 || index >= l.Size() {
		panic("Index out of bounds: index=" + strconv.Itoa(index) + " size=" + strconv.Itoa(l.Size()))
	}
	if l.delegate.isUnique {
		currIndex := l.IndexOf(elem)
		if currIndex >= 0 && currIndex != index {
			panic("element already in list")
		}
	}
	return l.delegate.interfaces.(abstractEList).doSet(index, elem)
}

// Size count the number of element in the list
func (l *unResolvedEList) Size() int {
	return l.delegate.Size()
}

// Clear remove all elements of the list
func (l *unResolvedEList) Clear() {
	l.delegate.Clear()
}

// Empty return true if the list contains 0 element
func (l *unResolvedEList) Empty() bool {
	return l.delegate.Empty()
}

// Contains return if an list contains or not an element
func (l *unResolvedEList) Contains(elem interface{}) bool {
	return l.IndexOf(elem) != -1
}

// IndexOf return the index on an element in an list, else return -1
func (l *unResolvedEList) IndexOf(elem interface{}) int {
	return l.delegate.basicEList.IndexOf(elem)
}

// Iterator through the list
func (l *unResolvedEList) Iterator() EIterator {
	return &listIterator{list: l}
}

func (l *unResolvedEList) ToArray() []interface{} {
	return l.delegate.ToArray()
}

func (l *unResolvedEList) GetNotifier() ENotifier {
	return l.delegate.GetNotifier()
}

func (l *unResolvedEList) GetFeature() EStructuralFeature {
	return l.delegate.GetFeature()
}

func (l *unResolvedEList) GetFeatureID() int {
	return l.delegate.GetFeatureID()
}

func (l *unResolvedEList) AddWithNotification(object interface{}, notifications ENotificationChain) ENotificationChain {
	return l.delegate.AddWithNotification(object, notifications)
}

func (l *unResolvedEList) RemoveWithNotification(object interface{}, notifications ENotificationChain) ENotificationChain {
	return l.delegate.RemoveWithNotification(object, notifications)
}

func (l *unResolvedEList) SetWithNotification(index int, object interface{}, notifications ENotificationChain) ENotificationChain {
	return l.delegate.SetWithNotification(index, object, notifications)
}
