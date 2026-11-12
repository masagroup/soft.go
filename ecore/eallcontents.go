package ecore

import (
	"iter"
	"slices"
)

// containmentFeatures returns the containment features of the class of o, those EContents walks:
// the contents are walked from them, without the view and its iterator. The class holds them in
// an immutable list, whose array is not a copy.
func containmentFeatures(o EObject) []any {
	return o.EClass().GetEContainmentFeatures().ToArray()
}

// eAllContentsFrame is the walk of the contents of an object: its containment features from
// cursor, then values, the remaining values of the current many feature.
type eAllContentsFrame struct {
	object   EObject
	features []any
	cursor   int
	values   EIterator
}

func newEAllContentsFrame(o EObject) eAllContentsFrame {
	return eAllContentsFrame{object: o, features: containmentFeatures(o)}
}

// next returns the next content of the object, nil when there is none.
func (f *eAllContentsFrame) next() any {
	if f.values != nil {
		if f.values.HasNext() {
			return f.values.Next()
		}
		f.values = nil
	}
	o := f.object
	for f.cursor < len(f.features) {
		feature := f.features[f.cursor].(EStructuralFeature)
		f.cursor++
		if !o.EIsSet(feature) {
			continue
		}
		value := o.EGetResolve(feature, true)
		if feature.IsMany() {
			if values := value.(EList).Iterator(); values.HasNext() {
				f.values = values
				return values.Next()
			}
		} else if value != nil {
			return value
		}
	}
	return nil
}

// eAllContentsIterator iterates the contents of an object, depth first. The frame of an object is
// pushed when it is returned and walked when the next content is searched, so that Prune skips its
// contents without reading them.
type eAllContentsIterator struct {
	stack []eAllContentsFrame
	next  any
	// prunable is the index of the frame of the last object returned, -1 when there is none
	prunable int
}

func newEAllContentsIterator(object EObject) *eAllContentsIterator {
	return newEAllContentsIteratorFrom(newEAllContentsFrame(object))
}

// newResourceAllContentsIterator returns an iterator on the roots of resource and their contents:
// its first frame has no object, its values are the roots.
func newResourceAllContentsIterator(resource EResource) *eAllContentsIterator {
	return newEAllContentsIteratorFrom(eAllContentsFrame{values: resource.GetContents().Iterator()})
}

func newEAllContentsIteratorFrom(frame eAllContentsFrame) *eAllContentsIterator {
	stack := make([]eAllContentsFrame, 1, 8)
	stack[0] = frame
	return &eAllContentsIterator{stack: stack, prunable: -1}
}

func (it *eAllContentsIterator) HasNext() bool {
	if it.next != nil {
		return true
	}
	for len(it.stack) != 0 {
		last := len(it.stack) - 1
		if next := it.stack[last].next(); next != nil {
			it.next = next
			return true
		}
		it.stack[last] = eAllContentsFrame{}
		it.stack = it.stack[:last]
	}
	return false
}

func (it *eAllContentsIterator) Next() any {
	if !it.HasNext() {
		panic("Not such an element")
	}
	next := it.next
	it.next = nil
	it.prunable = len(it.stack)
	it.stack = append(it.stack, newEAllContentsFrame(next.(EObject)))
	return next
}

// Prune skips the contents of the last object returned by Next.
func (it *eAllContentsIterator) Prune() {
	if it.prunable >= 0 && it.prunable < len(it.stack) {
		// the frame of the object is in the stack: the next object, if prepared, is one of its
		// contents
		clear(it.stack[it.prunable:])
		it.stack = it.stack[:it.prunable]
		it.next = nil
	}
	it.prunable = -1
}

// AllContents returns the contents of notifier, recursively, depth first, without the allocations
// of an iterator:
//   - of an object, its contents, not the object itself: the objects EAllContents iterates
//   - of a resource, each root then its contents: the objects GetAllContents iterates
//   - of a resource set, the roots and contents of each of its resources
//
// Any other notifier has no contents.
func AllContents(notifier ENotifier) iter.Seq[EObject] {
	return allContentsResolve(notifier, true)
}

// allContentsResolve returns the contents of notifier as AllContents, the proxies of the contents
// resolved when resolve is true.
func allContentsResolve(notifier ENotifier, resolve bool) iter.Seq[EObject] {
	return func(yield func(EObject) bool) {
		switch n := notifier.(type) {
		case EObject:
			yieldAllContents(n, resolve, yield)
		case EResource:
			yieldRoots(n, resolve, yield)
		case EResourceSet:
			for it := n.GetResources().Iterator(); it.HasNext(); {
				if !yieldRoots(it.Next().(EResource), resolve, yield) {
					return
				}
			}
		}
	}
}

// yieldRoots yields the roots of resource and their contents, false when the iteration is
// stopped.
func yieldRoots(resource EResource, resolve bool, yield func(EObject) bool) bool {
	for it := resource.GetContents().Iterator(); it.HasNext(); {
		if root := it.Next().(EObject); !yield(root) || !yieldAllContents(root, resolve, yield) {
			return false
		}
	}
	return true
}

// yieldAllContents yields the contents of o and their contents, false when the iteration is
// stopped.
func yieldAllContents(o EObject, resolve bool, yield func(EObject) bool) bool {
	for _, f := range containmentFeatures(o) {
		feature := f.(EStructuralFeature)
		if !o.EIsSet(feature) {
			continue
		}
		value := o.EGetResolve(feature, resolve)
		if feature.IsMany() {
			values := value.(EList)
			if objects, _ := values.(EObjectList); objects != nil && !resolve {
				values = objects.GetUnResolvedList()
			}
			for it := values.Iterator(); it.HasNext(); {
				if child := it.Next().(EObject); !yield(child) || !yieldAllContents(child, resolve, yield) {
					return false
				}
			}
		} else if child, _ := value.(EObject); child != nil {
			if !yield(child) || !yieldAllContents(child, resolve, yield) {
				return false
			}
		}
	}
	return true
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
		if internal, _ := classes[i].(eClassSubClasses); internal != nil {
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
