package ecore

type eAllContentIterator = treeIterator

func newEAllContentsIterator(object EObject) *eAllContentIterator {
	return &eAllContentIterator{object: object, root: false, getChildren: func(o interface{}) EIterator {
		return o.(EObject).EContents().Iterator()
	}}
}

type stateType int

const (
	start stateType = iota
	active
	end
)

type state struct {
	stateType stateType
	eClass    EClass
}

type transition struct {
	reference EReference
	source    state
	target    state
}

type transitions struct {
	array []*transition
}

func (ts transitions) isEmpty() bool {
	return len(ts.array) == 0
}

func (ts transitions) addTransition(transition *transition) {
	for _, t := range ts.array {
		if *t == *transition {
			return
		}
	}
	ts.array = append(ts.array, transition)
}

func (ts transitions) removeTransition(transition *transition) {
	for i, t := range ts.array {
		if *t == *transition {
			ts.array[i] = ts.array[len(ts.array)-1]
			ts.array = ts.array[:len(ts.array)-1]
			return
		}
	}
}

func (ts transitions) getTransition(reference EReference) *transition {
	for _, t := range ts.array {
		if t.reference == reference {
			return t
		}
	}
	return nil
}

type transitionTable map[state]transitions

func (table transitionTable) union(other transitionTable) transitionTable {
	for s, ts := range other {
		for _, t := range ts.array {
			table[s].addTransition(t)
		}
	}
	return table
}

func (table transitionTable) addTransition(t *transition) {
	table[t.source].addTransition(t)
}

func (table transitionTable) removeTransition(t *transition) {
	if tableTransitions, isTransitions := table[t.source]; isTransitions {
		tableTransitions.removeTransition(t)
	}
}

func (table transitionTable) getTransition(source state, reference EReference) *transition {
	if tableTransitions, isTransitions := table[source]; isTransitions {
		return tableTransitions.getTransition(reference)
	}
	return nil
}

func (table transitionTable) getTransitions(source state) []*transition {
	if tableTransitions, isTransitions := table[source]; isTransitions {
		return tableTransitions.array
	}
	return nil
}

func (table transitionTable) contains(source state) bool {
	_, isTransitions := table[source]
	return isTransitions
}

func newTransitionTable(startClass EClass, endClass EClass) transitionTable {
	resultTable := transitionTable{}
	if startClass != endClass {
		computeTransitionTableForState(state{stateType: start, eClass: startClass}, endClass, transitionTable{}, resultTable)
	}
	return resultTable
}

func computeTransitionTableForState(source state, endClass EClass, currentTable transitionTable, resultTable transitionTable) {
	for itFeature := source.eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
		reference, _ := itFeature.Next().(EReference)
		computeTransitionTableForReference(source, reference, endClass, currentTable, resultTable)
	}
}

func computeTransitionTableForReference(source state, reference EReference, endClass EClass, currentTable transitionTable, resultTable transitionTable) {
	target := reference.GetEReferenceType()
	transition := &transition{source: source, target: state{stateType: end, eClass: endClass}, reference: reference}
	if target == endClass {
		resultTable.union(currentTable)
		resultTable.addTransition(transition)
	} else {
		state := state{stateType: active, eClass: target}
		if currentTable.contains(state) {
			// cycle
			// check if target is in result and add the current
			// transition table to keep track of this cycle
			if resultTable.contains(state) {
				resultTable.union(currentTable)
				resultTable.addTransition(transition)
			}
			return
		}
		currentTable.addTransition(transition)
		computeTransitionTableForState(state, endClass, currentTable, resultTable)
		currentTable.removeTransition(transition)
	}
}

type data struct {
	object      EObject
	state       state
	transition  *transition   // current transition
	transitions []*transition // remaining transitions
	iterator    EIterator
}

type eObjectIterator struct {
	object EObject
	next   bool
}

func newEObjectIterator(object EObject) *eObjectIterator {
	return &eObjectIterator{
		object: object,
		next:   true,
	}
}

func (it *eObjectIterator) HasNext() bool {
	return it.next
}

func (it *eObjectIterator) Next() interface{} {
	if it.next {
		it.next = false
		return it.object
	}
	panic("Not such an element")
}

type eAllContentsWithClassIterator struct {
	table transitionTable
	data  []*data
	next  interface{}
}

func newEAllContentsWithClassIterator(object EObject, table transitionTable) *eAllContentsWithClassIterator {
	it := &eAllContentsWithClassIterator{
		table: table,
		data:  []*data{{object: object, state: state{stateType: start, eClass: object.EClass()}}},
	}
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

func (it *eAllContentsWithClassIterator) findNext() interface{} {
	notransitions := []*transition{}
	for len(it.data) != 0 {
		d := it.data[len(it.data)-1]

		iterator := d.iterator
		if iterator == nil || !iterator.HasNext() {
			// retrieve state transitions
			if d.transitions == nil {
				if transitions := it.table.getTransitions(d.state); transitions != nil {
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
			}

			// compute iterator
			if d.transition != nil && d.object.EIsSet(d.transition.reference) {
				value := d.object.EGet(d.transition.reference)
				switch v := value.(type) {
				case EList:
					d.iterator = v.Iterator()
				case EObject:
					d.iterator = newEObjectIterator(v)
				}
			}
		}

		if iterator != nil && iterator.HasNext() {
			object := iterator.Next().(EObject)

			// a leaf is found
			if d.state.stateType == end {
				return object
			}

			// push data to the stack
			it.data = append(it.data, &data{object: object, state: d.transition.target})
		} else {
			// pop data from the stack
			it.data = it.data[:len(it.data)-1]
		}
	}
	return nil
}
