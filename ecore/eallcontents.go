package ecore

import (
	"iter"
	"slices"
)

// EAllContentsWithClass returns an iterator on the instances of eClass in the contents of eObject,
// in the order of EAllContents: the walk only follows the containments which may lead to them.
func EAllContentsWithClass(eObject EObject, eClass EClass) EIterator {
	return EAllContentsWithTable(eObject, NewEClassTransitionsTable(eObject.EClass(), eClass))
}

// EAllContentsWithTable returns an iterator on the searched objects of table in the contents of
// eObject, in the order of EAllContents: the walk only follows the containments of the table.
func EAllContentsWithTable(eObject EObject, table *EClassTransitionsTable) EIterator {
	return newEAllContentsIterator(newEAllContentsFrame(eObject, table), table)
}

// AllContents returns the contents of notifier, recursively, depth first, without the allocations
// of an iterator:
//   - of an object, its contents, not the object itself: the objects EAllContents iterates
//   - of a resource, each root then its contents: the objects GetAllContents iterates
//   - of a resource set, the roots and contents of each of its resources
//
// Any other notifier has no contents.
func AllContents(notifier ENotifier) iter.Seq[EObject] {
	return allContents(notifier, nil, true)
}

// AllContentsWithClass returns the instances of eClass in the contents of notifier, as AllContents,
// in the same order: the walk only follows the containments which may lead to them.
func AllContentsWithClass(notifier ENotifier, eClass EClass) iter.Seq[EObject] {
	return func(yield func(EObject) bool) {
		startClasses := []EClass{}
		for root := range roots(notifier) {
			if rootClass := root.EClass(); !slices.Contains(startClasses, rootClass) {
				startClasses = append(startClasses, rootClass)
			}
		}
		AllContentsWithTable(notifier, newEClassTransitionsTable(startClasses, eClass))(yield)
	}
}

// AllContentsWithTable returns the searched objects of table in the contents of notifier, as
// AllContents, in the same order: the walk only follows the containments of the table. The table
// is the one of the class of notifier, or of the classes of the roots of a resource.
func AllContentsWithTable(notifier ENotifier, table *EClassTransitionsTable) iter.Seq[EObject] {
	return allContents(notifier, table, true)
}

// allContents returns the searched objects of table in the contents of notifier, all of them when
// table is nil, their proxies resolved when resolve is true.
func allContents(notifier ENotifier, table *EClassTransitionsTable, resolve bool) iter.Seq[EObject] {
	return func(yield func(EObject) bool) {
		if eObject, _ := notifier.(EObject); eObject != nil {
			yieldContents(eObject, table.transitions(eObject.EClass()).features, table, resolve, yield)
			return
		}
		for root := range roots(notifier) {
			if !yieldObject(root, table, resolve, yield) {
				return
			}
		}
	}
}

// roots returns the object notifier, or the roots of a resource or of the resources of a resource
// set.
func roots(notifier ENotifier) iter.Seq[EObject] {
	return func(yield func(EObject) bool) {
		yieldRoots := func(resource EResource) bool {
			for it := resource.GetContents().Iterator(); it.HasNext(); {
				if !yield(it.Next().(EObject)) {
					return false
				}
			}
			return true
		}
		switch n := notifier.(type) {
		case EObject:
			yield(n)
		case EResource:
			yieldRoots(n)
		case EResourceSet:
			for it := n.GetResources().Iterator(); it.HasNext(); {
				if !yieldRoots(it.Next().(EResource)) {
					return
				}
			}
		}
	}
}

// yieldObject yields o if it is searched, then the searched objects of its contents, false when
// the iteration is stopped.
func yieldObject(o EObject, table *EClassTransitionsTable, resolve bool, yield func(EObject) bool) bool {
	state := table.transitions(o.EClass())
	return (!state.end || yield(o)) && yieldContents(o, state.features, table, resolve, yield)
}

// yieldContents yields the searched objects of the values of the features of o and of their
// contents, false when the iteration is stopped. The contents of an object without features to
// walk are not walked.
func yieldContents(o EObject, features []any, table *EClassTransitionsTable, resolve bool, yield func(EObject) bool) bool {
	for _, f := range features {
		feature := f.(EStructuralFeature)
		if !o.EIsSet(feature) {
			continue
		}
		value := o.EGetResolve(feature, resolve)
		// the values of the feature: those of values when it is many, single otherwise
		var values EIterator
		var single EObject
		if feature.IsMany() {
			list := value.(EList)
			if objects, _ := list.(EObjectList); objects != nil && !resolve {
				list = objects.GetUnResolvedList()
			}
			values = list.Iterator()
		} else if single, _ = value.(EObject); single == nil {
			continue
		}
		for {
			var child EObject
			if values != nil {
				if !values.HasNext() {
					break
				}
				child = values.Next().(EObject)
			} else if child, single = single, nil; child == nil {
				break
			}
			state := table.transitions(child.EClass())
			if state.end && !yield(child) {
				return false
			}
			if len(state.features) != 0 && !yieldContents(child, state.features, table, resolve, yield) {
				return false
			}
		}
	}
	return true
}

// eAllContentsFrame is the walk of the contents of an object: its features from cursor, then
// values, the remaining values of the current many feature.
type eAllContentsFrame struct {
	object   EObject
	features []any
	cursor   int
	values   EIterator
}

// newEAllContentsFrame returns the frame of the walk of the features of o of table, all its
// containment features when table is nil.
func newEAllContentsFrame(o EObject, table *EClassTransitionsTable) eAllContentsFrame {
	return eAllContentsFrame{object: o, features: table.transitions(o.EClass()).features}
}

// next returns the next content of the object, nil when there is none.
func (f *eAllContentsFrame) next() EObject {
	if f.values != nil {
		if f.values.HasNext() {
			return f.values.Next().(EObject)
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
				return values.Next().(EObject)
			}
		} else if child, _ := value.(EObject); child != nil {
			return child
		}
	}
	return nil
}

// eAllContentsIterator iterates the searched objects of table in the contents of an object, all
// of them when table is nil, depth first. The frame of a searched object is pushed when it is
// returned and walked when the next object is searched, so that Prune skips its contents without
// reading them; the frames of the other objects are pushed when they are walked through.
type eAllContentsIterator struct {
	table *EClassTransitionsTable
	stack []eAllContentsFrame
	next  EObject
	// nextFeatures are the features of next to walk
	nextFeatures []any
	// prunable is the index of the frame of the last object returned, -1 when there is none
	prunable int
}

// newEAllContentsIterator returns an iterator on the searched objects of table, all of them when
// table is nil, walked from frame.
func newEAllContentsIterator(frame eAllContentsFrame, table *EClassTransitionsTable) *eAllContentsIterator {
	return &eAllContentsIterator{table: table, stack: append(make([]eAllContentsFrame, 0, 8), frame), prunable: -1}
}

func (it *eAllContentsIterator) HasNext() bool {
	if it.next != nil {
		return true
	}
	for len(it.stack) != 0 {
		last := len(it.stack) - 1
		next := it.stack[last].next()
		if next == nil {
			it.stack[last] = eAllContentsFrame{}
			it.stack = it.stack[:last]
			continue
		}
		state := it.table.transitions(next.EClass())
		if state.end {
			it.next, it.nextFeatures = next, state.features
			return true
		}
		if len(state.features) != 0 {
			it.stack = append(it.stack, eAllContentsFrame{object: next, features: state.features})
		}
	}
	return false
}

func (it *eAllContentsIterator) Next() any {
	if !it.HasNext() {
		panic("Not such an element")
	}
	next := it.next
	// an object without features to walk has no frame: there is nothing to prune
	it.prunable = -1
	if len(it.nextFeatures) != 0 {
		it.prunable = len(it.stack)
		it.stack = append(it.stack, eAllContentsFrame{object: next, features: it.nextFeatures})
	}
	it.next, it.nextFeatures = nil, nil
	return next
}

// Prune skips the contents of the last object returned by Next.
func (it *eAllContentsIterator) Prune() {
	if it.prunable >= 0 && it.prunable < len(it.stack) {
		// the frame of the object is in the stack: the next object, if prepared, is one of its
		// contents
		clear(it.stack[it.prunable:])
		it.stack = it.stack[:it.prunable]
		it.next, it.nextFeatures = nil, nil
	}
	it.prunable = -1
}

// EClassTransitionsTable tells, for the classes of the objects in the contents of an object of
// its start classes, whether their instances are searched, those of its end class, and which of
// their containment features may hold searched objects or objects containing some: the contents
// walks with a table only follow these features.
type EClassTransitionsTable struct {
	states map[EClass]eClassTransitions
	// endClass is the class of the searched objects
	endClass EClass
}

// eClassTransitions is the state of a class in a table: whether its instances are searched, and
// its containment features which may lead to searched objects, in the order of
// GetEContainmentFeatures.
type eClassTransitions struct {
	features []any
	end      bool
}

// NewEClassTransitionsTable returns the table of the walks of the contents of the objects of
// startClass searching the instances of endClass.
func NewEClassTransitionsTable(startClass EClass, endClass EClass) *EClassTransitionsTable {
	return newEClassTransitionsTable([]EClass{startClass}, endClass)
}

// newEClassTransitionsTable returns the table of the walks of the contents of the objects of
// startClasses searching the instances of endClass. A class may lead to a searched object when one
// of its containment features may hold an instance of endClass, or an object of a class which
// may lead to one: computed as a fixed point, whatever the cycles of containments.
func newEClassTransitionsTable(startClasses []EClass, endClass EClass) *EClassTransitionsTable {
	table := &EClassTransitionsTable{states: map[EClass]eClassTransitions{}, endClass: endClass}
	// the classes of the objects in the contents of the start objects, and for each one the
	// classes of the objects each of its containment features may hold
	type classNode struct {
		eClass   EClass
		features []any
		targets  [][]EClass // nil for a feature holding objects of any class
		end      bool
		leads    bool
	}
	nodes := []*classNode{}
	index := map[EClass]*classNode{}
	add := func(eClass EClass) {
		if _, known := index[eClass]; !known {
			node := &classNode{eClass: eClass, end: table.isInstance(eClass)}
			index[eClass] = node
			nodes = append(nodes, node)
		}
	}
	for _, startClass := range startClasses {
		add(startClass)
	}
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		node.features = node.eClass.GetEContainmentFeatures().ToArray()
		node.targets = make([][]EClass, len(node.features))
		for j, feature := range node.features {
			targetClass := feature.(EReference).GetEReferenceType()
			if isEObjectClass(targetClass) {
				// a feature typed EObject holds objects of any class
				continue
			}
			targets := classesOf(targetClass)
			node.targets[j] = targets
			for _, target := range targets {
				add(target)
			}
		}
	}
	// a feature leads to searched objects when it may hold a searched object or an object leading
	// to one
	featureLeads := func(targets []EClass) bool {
		if targets == nil {
			return true
		}
		for _, target := range targets {
			if node := index[target]; node.end || node.leads {
				return true
			}
		}
		return false
	}
	for changed := true; changed; {
		changed = false
		for _, node := range nodes {
			if !node.leads && slices.ContainsFunc(node.targets, featureLeads) {
				node.leads = true
				changed = true
			}
		}
	}
	for _, node := range nodes {
		state := eClassTransitions{end: node.end}
		if node.leads {
			for j, feature := range node.features {
				if featureLeads(node.targets[j]) {
					state.features = append(state.features, feature)
				}
			}
		}
		table.states[node.eClass] = state
	}
	return table
}

// isEObjectClass returns true for the class EObject of the Ecore package: all objects are its
// instances, as in EMF.
func isEObjectClass(eClass EClass) bool {
	return eClass == GetPackage().GetEObject()
}

// isInstance returns true if the instances of eClass are searched: it is the end class or one of
// its sub classes.
func (table *EClassTransitionsTable) isInstance(eClass EClass) bool {
	return table.endClass != nil && (isEObjectClass(table.endClass) || table.endClass.IsSuperTypeOf(eClass))
}

// transitions returns the state of eClass. Without table, all the objects are searched and all
// their containment features are followed. A class the table does not know, created after it or
// not in the contents of its start classes, may lead anywhere: all its containment features are
// followed.
func (table *EClassTransitionsTable) transitions(eClass EClass) eClassTransitions {
	if table == nil {
		return eClassTransitions{features: eClass.GetEContainmentFeatures().ToArray(), end: true}
	}
	if state, known := table.states[eClass]; known {
		return state
	}
	return eClassTransitions{features: eClass.GetEContainmentFeatures().ToArray(), end: table.isInstance(eClass)}
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
