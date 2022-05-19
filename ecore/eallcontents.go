package ecore

type eAllContentIterator = treeIterator

func newEAllContentsIterator(object EObject) *eAllContentIterator {
	return &eAllContentIterator{object: object, root: false, getChildren: func(o interface{}) EIterator {
		return o.(EObject).EContents().Iterator()
	}}
}

type transition struct {
	reference EReference
	source    EClass
	target    EClass
}

type transitions []*transition

func (ts transitions) addTransition(transition *transition) transitions {
	for _, t := range ts {
		if t == transition || *t == *transition {
			return ts
		}
	}
	return append(ts, transition)
}

func (ts transitions) removeTransition(transition *transition) transitions {
	for i, t := range ts {
		if t == transition || *t == *transition {
			ts[i] = ts[len(ts)-1]
			return ts[:len(ts)-1]
		}
	}
	return ts
}

type EClassTransitionsTable struct {
	transitionsMap map[EClass]transitions
	endSet         map[EClass]struct{}
}

func newEClassTransitionsTable() *EClassTransitionsTable {
	return &EClassTransitionsTable{
		transitionsMap: map[EClass]transitions{},
		endSet:         map[EClass]struct{}{},
	}
}

func (table *EClassTransitionsTable) Union(other *EClassTransitionsTable) *EClassTransitionsTable {
	for _, transitions := range other.transitionsMap {
		for _, transition := range transitions {
			table.addTransition(transition)
		}
	}
	for eClass := range other.endSet {
		table.setIsEnd(eClass)
	}
	return table
}

func (table *EClassTransitionsTable) addTransition(t *transition) {
	table.transitionsMap[t.source] = table.transitionsMap[t.source].addTransition(t)
}

func (table *EClassTransitionsTable) removeTransition(t *transition) {
	if tableTransitions, isTransitions := table.transitionsMap[t.source]; isTransitions {
		table.transitionsMap[t.source] = tableTransitions.removeTransition(t)
	}
}

func (table *EClassTransitionsTable) getTransitions(eClass EClass) []*transition {
	if tableTransitions, isTransitions := table.transitionsMap[eClass]; isTransitions {
		return tableTransitions
	}
	return nil
}

func (table *EClassTransitionsTable) contains(source EClass) bool {
	_, isTransitions := table.transitionsMap[source]
	return isTransitions
}

func (table *EClassTransitionsTable) setIsEnd(eClass EClass) {
	table.endSet[eClass] = struct{}{}
}

func (table *EClassTransitionsTable) isEnd(eClass EClass) bool {
	_, isEnd := table.endSet[eClass]
	return isEnd
}

func NewEClassTransitionsTable(startClass EClass, endClass EClass) *EClassTransitionsTable {
	resultTable := newEClassTransitionsTable()
	if startClass != endClass {
		computeTransitionTableForState(startClass, endClass, newEClassTransitionsTable(), resultTable)
	}
	return resultTable
}

func computeTransitionTableForState(eClass EClass, endClass EClass, currentTable *EClassTransitionsTable, resultTable *EClassTransitionsTable) {
	for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
		reference, _ := itFeature.Next().(EReference)
		if reference != nil {
			computeTransitionTableForReference(eClass, reference, endClass, currentTable, resultTable)
		}
	}
}

func computeTransitionTableForReference(sourceClass EClass, reference EReference, endClass EClass, currentTable *EClassTransitionsTable, resultTable *EClassTransitionsTable) {
	targetClass := reference.GetEReferenceType()
	transition := &transition{source: sourceClass, target: targetClass, reference: reference}

	if targetClass == endClass {
		// end
		// add current to result table
		resultTable.setIsEnd(targetClass)
		resultTable.Union(currentTable)
		resultTable.addTransition(transition)
	}

	if currentTable.contains(targetClass) {
		// cycle
		// check if target is in result and add the current
		// transition table to keep track of this cycle
		if resultTable.contains(targetClass) {
			resultTable.Union(currentTable)
			resultTable.addTransition(transition)
		}
		return
	}

	currentTable.addTransition(transition)
	computeTransitionTableForState(targetClass, endClass, currentTable, resultTable)
	currentTable.removeTransition(transition)

}

type data struct {
	eObject     EObject
	eClass      EClass
	transition  *transition   // current transition
	transitions []*transition // remaining transitions
	iterator    EIterator
}

type eObjectIterator struct {
	eObject EObject
	next    bool
}

func newEObjectIterator(eObject EObject) *eObjectIterator {
	return &eObjectIterator{
		eObject: eObject,
		next:    true,
	}
}

func (it *eObjectIterator) HasNext() bool {
	return it.next
}

func (it *eObjectIterator) Next() interface{} {
	if it.next {
		it.next = false
		return it.eObject
	}
	panic("Not such an element")
}

type eAllContentsWithClassIterator struct {
	table *EClassTransitionsTable
	data  []*data
	next  interface{}
}

func newEAllContentsWithClassIterator(eObject EObject, table *EClassTransitionsTable) *eAllContentsWithClassIterator {
	it := &eAllContentsWithClassIterator{
		table: table,
		data:  []*data{{eObject: eObject, eClass: eObject.EClass()}},
	}
	it.next = it.findNext()
	return it
}

func (it *eAllContentsWithClassIterator) HasNext() bool {
	return it.next != nil
}

func (it *eAllContentsWithClassIterator) Next() interface{} {
	next := it.next
	it.next = it.findNext()
	return next
}

var notransitions []*transition = []*transition{}

func (it *eAllContentsWithClassIterator) findNext() interface{} {
	for len(it.data) != 0 {
		d := it.data[len(it.data)-1]

		if d.iterator == nil || !d.iterator.HasNext() {
			// retrieve state transitions
			if d.transitions == nil {
				if transitions := it.table.getTransitions(d.eClass); transitions != nil {
					d.transitions = transitions
				} else {
					d.transitions = notransitions
				}
			}
			// compute current and remaining transition
			if len(d.transitions) != 0 {
				last := len(d.transitions) - 1
				d.transition = d.transitions[last]
				d.transitions = d.transitions[:last]
			} else {
				d.transition = nil
			}

			// compute iterator
			if d.transition != nil && d.eObject.EIsSet(d.transition.reference) {
				value := d.eObject.EGet(d.transition.reference)
				switch v := value.(type) {
				case EList:
					d.iterator = v.Iterator()
				case EObject:
					d.iterator = newEObjectIterator(v)
				}
			}
		}

		if d.iterator != nil && d.iterator.HasNext() {
			eObject := d.iterator.Next().(EObject)
			eClass := d.transition.target
			// push data to the stack before returning to look for children next iteration
			it.data = append(it.data, &data{eObject: eObject, eClass: eClass})

			// a leaf is found
			if it.table.isEnd(eClass) {
				return eObject
			}

		} else {
			// pop data from the stack
			it.data = it.data[:len(it.data)-1]
		}
	}
	return nil
}
