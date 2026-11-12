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
			if node.leads {
				continue
			}
			for _, targets := range node.targets {
				if featureLeads(targets) {
					node.leads = true
					changed = true
					break
				}
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

// transitions returns the state of eClass. A class the table does not know, created after it or
// not in the contents of its start classes, may lead anywhere: all its containment features are
// followed.
func (table *EClassTransitionsTable) transitions(eClass EClass) eClassTransitions {
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

// eAllContentsWithTableIterator iterates the searched objects of the contents of an object, depth
// first, in the order of EAllContents: it walks the features of the table from a stack of frames,
// pushed for the objects whose features may lead to searched objects.
type eAllContentsWithTableIterator struct {
	table *EClassTransitionsTable
	stack []eAllContentsFrame
	next  EObject
}

func newEAllContentsWithClassIterator(eObject EObject, eClass EClass) *eAllContentsWithTableIterator {
	return newEAllContentsWithTableIterator(eObject, NewEClassTransitionsTable(eObject.EClass(), eClass))
}

func newEAllContentsWithTableIterator(eObject EObject, table *EClassTransitionsTable) *eAllContentsWithTableIterator {
	stack := make([]eAllContentsFrame, 1, 8)
	stack[0] = eAllContentsFrame{object: eObject, features: table.transitions(eObject.EClass()).features}
	return &eAllContentsWithTableIterator{table: table, stack: stack}
}

func (it *eAllContentsWithTableIterator) HasNext() bool {
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
		eObject := next.(EObject)
		state := it.table.transitions(eObject.EClass())
		if len(state.features) != 0 {
			it.stack = append(it.stack, eAllContentsFrame{object: eObject, features: state.features})
		}
		if state.end {
			it.next = eObject
			return true
		}
	}
	return false
}

func (it *eAllContentsWithTableIterator) Next() any {
	if !it.HasNext() {
		panic("Not such an element")
	}
	next := it.next
	it.next = nil
	return next
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
		allContentsWithTable(notifier, newEClassTransitionsTable(startClasses, eClass))(yield)
	}
}

// AllContentsWithTable returns the searched objects of table in the contents of notifier, as
// AllContents, in the same order: the walk only follows the containments of the table. The table
// is the one of the class of notifier, or of the classes of the roots of a resource.
func AllContentsWithTable(notifier ENotifier, table *EClassTransitionsTable) iter.Seq[EObject] {
	return allContentsWithTable(notifier, table)
}

func allContentsWithTable(notifier ENotifier, table *EClassTransitionsTable) iter.Seq[EObject] {
	return func(yield func(EObject) bool) {
		if eObject, _ := notifier.(EObject); eObject != nil {
			yieldContentsWithTable(eObject, table.transitions(eObject.EClass()).features, table, yield)
			return
		}
		for root := range roots(notifier) {
			if !yieldWithTable(root, table, yield) {
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

// yieldWithTable yields o if it is searched, then the searched objects of its contents, false when
// the iteration is stopped.
func yieldWithTable(o EObject, table *EClassTransitionsTable, yield func(EObject) bool) bool {
	state := table.transitions(o.EClass())
	if state.end && !yield(o) {
		return false
	}
	return len(state.features) == 0 || yieldContentsWithTable(o, state.features, table, yield)
}

// yieldContentsWithTable yields the searched objects of the values of the features of o and of
// their contents, false when the iteration is stopped.
func yieldContentsWithTable(o EObject, features []any, table *EClassTransitionsTable, yield func(EObject) bool) bool {
	for _, f := range features {
		feature := f.(EStructuralFeature)
		if !o.EIsSet(feature) {
			continue
		}
		value := o.EGetResolve(feature, true)
		if feature.IsMany() {
			for it := value.(EList).Iterator(); it.HasNext(); {
				if !yieldWithTable(it.Next().(EObject), table, yield) {
					return false
				}
			}
		} else if child, _ := value.(EObject); child != nil {
			if !yieldWithTable(child, table, yield) {
				return false
			}
		}
	}
	return true
}
