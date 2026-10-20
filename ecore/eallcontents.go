package ecore

import "slices"

type eAllContentIterator = treeIterator

func newEAllContentsIterator(object EObject) *eAllContentIterator {
	return &eAllContentIterator{object: object, root: false, getChildren: func(o any) EIterator {
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
	// endClass is the class of the objects searched, nil for an empty table: an object is an end
	// when it is an instance of it, whatever the type of the reference holding it
	endClass EClass
}

func newEClassTransitionsTable() *EClassTransitionsTable {
	return &EClassTransitionsTable{
		transitionsMap: map[EClass]transitions{},
		endSet:         map[EClass]struct{}{},
	}
}

func (table *EClassTransitionsTable) union(other *EClassTransitionsTable) *EClassTransitionsTable {
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

// isInstanceEnd returns true if an object of class eClass is searched: its class is the end class or
// one of its sub classes
func (table *EClassTransitionsTable) isInstanceEnd(eClass EClass) bool {
	return table.isEnd(eClass) || (table.endClass != nil && table.endClass.IsSuperTypeOf(eClass))
}

func NewEClassTransitionsTable(startClass EClass, endClass EClass) *EClassTransitionsTable {
	resultTable := newEClassTransitionsTable()
	if startClass != endClass {
		resultTable.endClass = endClass
		computeTransitionTableForState(startClass, endClass, newEClassTransitionsTable(), resultTable)
	}
	return resultTable
}

// classesOf returns eClass and its sub classes, recursively: the classes of the objects a reference
// typed eClass may hold
func classesOf(eClass EClass) []EClass {
	classes := []EClass{eClass}
	for i := 0; i < len(classes); i++ {
		if internal, _ := classes[i].(EClassInternal); internal != nil {
			for _, subClass := range internal.getSubClasses() {
				if subEClass, _ := subClass.(EClass); subEClass != nil && !slices.Contains(classes, subEClass) {
					classes = append(classes, subEClass)
				}
			}
		}
	}
	return classes
}

func computeTransitionTableForState(eClass EClass, endClass EClass, currentTable *EClassTransitionsTable, resultTable *EClassTransitionsTable) {
	for itFeature := eClass.GetEContainmentFeatures().Iterator(); itFeature.HasNext(); {
		reference, _ := itFeature.Next().(EReference)
		if reference != nil {
			computeTransitionTableForReference(eClass, reference, endClass, currentTable, resultTable)
		}
	}
}

func computeTransitionTableForReference(sourceClass EClass, reference EReference, endClass EClass, currentTable *EClassTransitionsTable, resultTable *EClassTransitionsTable) {
	targetClass := reference.GetEReferenceType()
	transition := &transition{source: sourceClass, target: targetClass, reference: reference}

	// the reference may hold instances of the sub classes of its type: each one may be an end, or
	// lead to one through its own containments
	targetClasses := classesOf(targetClass)
	for _, eClass := range targetClasses {
		if endClass.IsSuperTypeOf(eClass) {
			// end
			// add current to result table
			resultTable.setIsEnd(eClass)
			resultTable.union(currentTable)
			resultTable.addTransition(transition)
		}
	}

	for _, eClass := range targetClasses {
		if currentTable.contains(eClass) {
			// cycle
			// check if target is in result and add the current
			// transition table to keep track of this cycle
			if resultTable.contains(eClass) {
				resultTable.union(currentTable)
				resultTable.addTransition(transition)
			}
			continue
		}

		currentTable.addTransition(transition)
		computeTransitionTableForState(eClass, endClass, currentTable, resultTable)
		currentTable.removeTransition(transition)
	}
}

type data struct {
	eObject     EObject
	eClass      EClass
	iterator    EIterator
	transition  *transition   // current transition
	transitions []*transition // remaining transitions
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

func (it *eObjectIterator) Next() any {
	if it.next {
		it.next = false
		return it.eObject
	}
	panic("Not such an element")
}

type eAllContentsWithClassIterator struct {
	next  any
	table *EClassTransitionsTable
	data  []*data
}

func newEAllContentsWithClassIterator(eObject EObject, eClass EClass) *eAllContentsWithClassIterator {
	return newEAllContentsWithTableIterator(eObject, NewEClassTransitionsTable(eObject.EClass(), eClass))
}

func newEAllContentsWithTableIterator(eObject EObject, table *EClassTransitionsTable) *eAllContentsWithClassIterator {
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

func (it *eAllContentsWithClassIterator) Next() any {
	next := it.next
	it.next = it.findNext()
	return next
}

var notransitions = []*transition{}

func (it *eAllContentsWithClassIterator) findNext() any {
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
			// the next transition with values: an unset or empty feature does not end the
			// exploration of the object, its other transitions are tried
			d.iterator = nil
			for (d.iterator == nil || !d.iterator.HasNext()) && len(d.transitions) != 0 {
				last := len(d.transitions) - 1
				d.transition = d.transitions[last]
				d.transitions = d.transitions[:last]
				d.iterator = nil
				if d.eObject.EIsSet(d.transition.reference) {
					value := d.eObject.EGet(d.transition.reference)
					switch v := value.(type) {
					case EList:
						d.iterator = v.Iterator()
					case EObject:
						d.iterator = newEObjectIterator(v)
					}
				}
			}
		}

		if d.iterator != nil && d.iterator.HasNext() {
			eObject := d.iterator.Next().(EObject)
			// the class of the object, not the type of the reference: it may be a sub class
			eClass := eObject.EClass()
			// push data to the stack before returning to look for children next iteration
			it.data = append(it.data, &data{eObject: eObject, eClass: eClass})

			// a leaf is found
			if it.table.isInstanceEnd(eClass) {
				return eObject
			}

		} else {
			// pop data from the stack
			it.data = it.data[:len(it.data)-1]
		}
	}
	return nil
}
