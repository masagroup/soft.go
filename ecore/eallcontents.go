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

type transitions []*transition

func (ts transitions) isEmpty() bool {
	return len(ts) == 0
}

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

func (ts transitions) getTransition(reference EReference) *transition {
	for _, t := range ts {
		if t.reference == reference {
			return t
		}
	}
	return nil
}

type transitionTable map[state]transitions

func (table transitionTable) union(other transitionTable) transitionTable {
	for s, ts := range other {
		for _, t := range ts {
			table[s] = table[s].addTransition(t)
		}
	}
	return table
}

func (table transitionTable) addTransition(t *transition) {
	table[t.source] = table[t.source].addTransition(t)
}

func (table transitionTable) removeTransition(t *transition) {
	if tableTransitions, isTransitions := table[t.source]; isTransitions {
		table[t.source] = tableTransitions.removeTransition(t)
	}
}

func (table transitionTable) getTransitions(source state) []*transition {
	if tableTransitions, isTransitions := table[source]; isTransitions {
		return tableTransitions
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
		if reference != nil {
			computeTransitionTableForReference(source, reference, endClass, currentTable, resultTable)
		}
	}
}

func computeTransitionTableForReference(source state, reference EReference, endClass EClass, currentTable transitionTable, resultTable transitionTable) {
	target := reference.GetEReferenceType()
	if target == endClass {
		transition := &transition{source: source, target: state{stateType: end, eClass: endClass}, reference: reference}
		resultTable.union(currentTable)
		resultTable.addTransition(transition)
	} else {
		state := state{stateType: active, eClass: target}
		transition := &transition{source: source, target: state, reference: reference}
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
			} else {
				d.transition = nil
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

		if d.iterator != nil && d.iterator.HasNext() {
			object := d.iterator.Next().(EObject)

			// push data to the stack before returning to look for children next iteration
			it.data = append(it.data, &data{object: object, state: d.transition.target})

			// a leaf is found
			if d.transition.target.stateType == end {
				return object
			}

		} else {
			// pop data from the stack
			it.data = it.data[:len(it.data)-1]
		}
	}
	return nil
}
