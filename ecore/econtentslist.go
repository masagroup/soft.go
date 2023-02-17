package ecore

type eContentsList struct {
	emptyImmutableEList
	o        EObject
	features EList
	resolve  bool
}

type eContentsListIterator struct {
	next          any
	values        EIterator
	l             *eContentsList
	prepared      int
	featureCursor int
}

func (it *eContentsListIterator) Next() any {
	if it.prepared > 1 || it.HasNext() {
		it.prepared = 0
		return it.next
	}
	panic("not such element")
}

func (it *eContentsListIterator) HasNext() bool {
	switch it.prepared {
	case 2:
		return true
	case 1:
		return false
	default:
		features := it.l.features
		resolve := it.l.resolve
		o := it.l.o
		if it.values == nil || !it.values.HasNext() {
			// no feature list or feature list iterator is finished
			for it.featureCursor < features.Size() {
				feature := features.Get(it.featureCursor).(EStructuralFeature)
				it.featureCursor++
				if o.EIsSet(feature) {
					value := o.EGetResolve(feature, resolve)
					if feature.IsMany() {
						// list of values
						values := value.(EList)
						// get unresolved list if object list and not resolved iterator
						if objectList, _ := value.(EObjectList); objectList != nil && !resolve {
							values = objectList.GetUnResolvedList()
						}
						if itValues := values.Iterator(); itValues.HasNext() {
							// we have a value
							it.values = itValues
							it.next = itValues.Next()
							it.prepared = 2
							return true
						}
					} else if value != nil {
						// we have a value
						it.values = nil
						it.next = value
						it.prepared = 2
						return true
					}
				}
			}
			it.values = nil
			it.next = nil
			it.prepared = 1
			return false
		} else {
			it.next = it.values.Next()
			it.prepared = 2
			return true
		}
	}
}

func newEContentsList(o EObject, features EList, resolve bool) *eContentsList {
	return &eContentsList{
		o:        o,
		features: features,
		resolve:  resolve,
	}
}

// Get an element of the array
func (l *eContentsList) Get(index int) any {
	it := l.features.Iterator()
	for i := 0; i < index; i++ {
		it.Next()
	}
	return it.Next()
}

// Size count the number of element in the array
func (l *eContentsList) Size() int {
	size := 0
	for it := l.features.Iterator(); it.HasNext(); {
		eFeature := it.Next().(EStructuralFeature)
		if l.o.EIsSet(eFeature) {
			value := l.o.EGetResolve(eFeature, false)
			if eFeature.IsMany() {
				list := value.(EList)
				size += list.Size()
			} else if value != nil {
				size++
			}
		}
	}
	return size
}

// Empty return true if the array contains 0 element
func (l *eContentsList) Empty() bool {
	for it := l.features.Iterator(); it.HasNext(); {
		eFeature := it.Next().(EStructuralFeature)
		if l.o.EIsSet(eFeature) {
			value := l.o.EGetResolve(eFeature, false)
			if eFeature.IsMany() {
				list := value.(EList)
				if !list.Empty() {
					return false
				}
			} else if value != nil {
				return false
			}
		}
	}
	return true
}

// Contains return if an array contains or not an element
func (l *eContentsList) Contains(elem any) bool {
	for it := l.Iterator(); it.HasNext(); {
		if e := it.Next(); e == elem {
			return true
		}
	}
	return false
}

// IndexOf return the index on an element in an array, else return -1
func (l *eContentsList) IndexOf(elem any) int {
	index := 0
	for it := l.Iterator(); it.HasNext(); {
		if e := it.Next(); e == elem {
			return index
		}
		index++
	}
	return -1
}

// Iterator through the array
func (l *eContentsList) Iterator() EIterator {
	return &eContentsListIterator{l: l}
}

// ToArray convert to array
func (l *eContentsList) ToArray() []any {
	arr := []any{}
	for it := l.Iterator(); it.HasNext(); {
		arr = append(arr, it.Next())
	}
	return arr
}

func (l *eContentsList) GetUnResolvedList() EList {
	if l.resolve {
		return newEContentsList(l.o, l.features, false)
	}
	return l
}
