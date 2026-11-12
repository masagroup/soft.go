package ecore

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type eAllContentsClassModel struct {
	eFactory EFactory
	root     EClass
	item     EClass // abstract super type of book
	book     EClass
	books    EReference // root -> book
	more     EReference // root -> book, the last transition of root
	items    EReference // root -> item
}

func newEAllContentsClassModel() *eAllContentsClassModel {
	f := GetFactory()
	ePackage := f.CreateEPackage()
	m := &eAllContentsClassModel{eFactory: f.CreateEFactory()}
	ePackage.SetEFactoryInstance(m.eFactory)
	newClass := func(name string, abstract bool) EClass {
		eClass := f.CreateEClass()
		eClass.SetName(name)
		eClass.SetAbstract(abstract)
		ePackage.GetEClassifiers().Add(eClass)
		return eClass
	}
	newContainment := func(eClass EClass, name string, eType EClass) EReference {
		eReference := f.CreateEReference()
		eReference.SetName(name)
		eReference.SetEType(eType)
		eReference.SetContainment(true)
		eReference.SetUpperBound(UNBOUNDED_MULTIPLICITY)
		eClass.GetEStructuralFeatures().Add(eReference)
		return eReference
	}
	m.root, m.item, m.book = newClass("Root", false), newClass("Item", true), newClass("Book", false)
	m.book.GetESuperTypes().Add(m.item)
	m.books = newContainment(m.root, "books", m.book)
	m.more = newContainment(m.root, "more", m.book)
	m.items = newContainment(m.root, "items", m.item)
	return m
}

func countWithClass(eObject EObject, eClass EClass) int {
	n := 0
	for it := EAllContentsWithClass(eObject, eClass); it.HasNext(); it.Next() {
		n++
	}
	return n
}

// an unset containment does not end the exploration of the object: its other containments are tried
func TestEAllContentsWithClass_UnsetTransition(t *testing.T) {
	m := newEAllContentsClassModel()
	root := m.eFactory.Create(m.root)
	root.EGet(m.books).(EList).Add(m.eFactory.Create(m.book))
	assert.Equal(t, 1, countWithClass(root, m.book))
}

// an instance of the class is found whatever the type of the reference holding it
func TestEAllContentsWithClass_SubClass(t *testing.T) {
	m := newEAllContentsClassModel()
	root := m.eFactory.Create(m.root)
	root.EGet(m.items).(EList).Add(m.eFactory.Create(m.book))
	root.EGet(m.more).(EList).Add(m.eFactory.Create(m.book))
	assert.Equal(t, 2, countWithClass(root, m.book))
	// the instances of a sub class are instances of the class
	assert.Equal(t, 2, countWithClass(root, m.item))
}

// instancesInContents returns the objects of EAllContents of eObject that are instances of eClass:
// the objects EAllContentsWithClass iterates, in the same order.
func instancesInContents(eObject EObject, eClass EClass) []EObject {
	objects := []EObject{}
	for it := eObject.EAllContents(); it.HasNext(); {
		if o := it.Next().(EObject); eClass.IsSuperTypeOf(o.EClass()) {
			objects = append(objects, o)
		}
	}
	return objects
}

// checkWithClass checks that EAllContentsWithClass of eObject iterates the instances of each class
// of ePackage in its contents, in the order of EAllContents.
func checkWithClass(t *testing.T, eObject EObject, ePackage EPackage) {
	for _, classifier := range ePackage.GetEClassifiers().ToArray() {
		eClass, _ := classifier.(EClass)
		if eClass == nil {
			continue
		}
		expected := instancesInContents(eObject, eClass)
		actual := collectIterator(EAllContentsWithClass(eObject, eClass))
		assert.Equal(t, positions(eObject, expected), positions(eObject, actual), "%s from %s", eClass.GetName(), eObject.EClass().GetName())
	}
}

// positions returns the positions of objects in EAllContents of eObject, -1 for an object not in
// its contents: a readable comparison of objects.
func positions(eObject EObject, objects []EObject) []int {
	contents := collectIterator(eObject.EAllContents())
	result := []int{}
	for _, o := range objects {
		result = append(result, slices.Index(contents, o))
	}
	return result
}

func TestEAllContentsWithClass_ContentsOrder(t *testing.T) {
	p := loadEAllContentsPackage(t)
	m := loadEAllContentsModel(t, p.ePackage)
	checkWithClass(t, m.eRoot, p.ePackage)
	bench := newEAllContentsBenchPackage()
	root, _ := bench.newTree(3, 3, 2)
	checkWithClass(t, root, bench.nodeClass.GetEPackage())
}

// the instances of the start class in its contents are found
func TestEAllContentsWithClass_StartClass(t *testing.T) {
	p := loadEAllContentsPackage(t)
	m := loadEAllContentsModel(t, p.ePackage)
	for it := m.eRoot.EAllContents(); it.HasNext(); {
		o := it.Next().(EObject)
		checkWithClass(t, o, p.ePackage)
	}
}

// a cycle of containments explored before the path to the class is followed: A holds a B, first,
// which holds an A, and an End
func TestEAllContentsWithClass_CycleBeforeEnd(t *testing.T) {
	f := GetFactory()
	ePackage := f.CreateEPackage()
	eFactory := f.CreateEFactory()
	ePackage.SetEFactoryInstance(eFactory)
	newClass := func(name string) EClass {
		eClass := f.CreateEClass()
		eClass.SetName(name)
		ePackage.GetEClassifiers().Add(eClass)
		return eClass
	}
	newContainment := func(eClass EClass, name string, eType EClass) EReference {
		eReference := f.CreateEReference()
		eReference.SetName(name)
		eReference.SetEType(eType)
		eReference.SetContainment(true)
		eClass.GetEStructuralFeatures().Add(eReference)
		return eReference
	}
	a, b, end := newClass("A"), newClass("B"), newClass("End")
	ab := newContainment(a, "b", b)
	ae := newContainment(a, "end", end)
	ba := newContainment(b, "a", a)

	root := eFactory.Create(a)
	root.ESet(ae, eFactory.Create(end))
	child := eFactory.Create(b)
	root.ESet(ab, child)
	nested := eFactory.Create(a)
	child.ESet(ba, nested)
	nested.ESet(ae, eFactory.Create(end))

	require.Len(t, instancesInContents(root, end), 2)
	checkWithClass(t, root, ePackage)
}

// collectSeq returns the objects of seq, an empty slice for none.
func collectSeq(seq iter.Seq[EObject]) []EObject {
	return append([]EObject{}, slices.Collect(seq)...)
}

// filterInstances returns the objects of seq which are instances of eClass.
func filterInstances(seq iter.Seq[EObject], eClass EClass) []EObject {
	objects := []EObject{}
	for o := range seq {
		if isEObjectClass(eClass) || eClass.IsSuperTypeOf(o.EClass()) {
			objects = append(objects, o)
		}
	}
	return objects
}

func TestAllContentsWithClass(t *testing.T) {
	p := loadEAllContentsPackage(t)
	m := loadEAllContentsModel(t, p.ePackage)
	resourceSet := NewEResourceSetImpl()
	resourceSet.GetResources().Add(m.eResource)
	other, _ := newAllContentsResource(newEAllContentsBenchPackage(), 2, 2, 1)
	resourceSet.GetResources().Add(other)
	for _, notifier := range []ENotifier{m.eRoot, m.eResource, resourceSet} {
		for _, classifier := range p.ePackage.GetEClassifiers().ToArray() {
			eClass := classifier.(EClass)
			expected := filterInstances(AllContents(notifier), eClass)
			assert.Equal(t, expected, collectSeq(AllContentsWithClass(notifier, eClass)), "%s in %T", eClass.GetName(), notifier)
			if eObject, _ := notifier.(EObject); eObject != nil {
				table := NewEClassTransitionsTable(eObject.EClass(), eClass)
				assert.Equal(t, expected, collectSeq(AllContentsWithTable(eObject, table)), "%s in %T", eClass.GetName(), notifier)
			}
		}
	}
	// early stop
	units := filterInstances(AllContents(m.eRoot), p.eUnitClass)
	objects := []EObject{}
	for o := range AllContentsWithClass(m.eRoot, p.eUnitClass) {
		objects = append(objects, o)
		if len(objects) == 2 {
			break
		}
	}
	assert.Equal(t, units[:2], objects)
}

func TestEAllContentsWithClass_Exhausted(t *testing.T) {
	p := loadEAllContentsPackage(t)
	m := loadEAllContentsModel(t, p.ePackage)
	it := EAllContentsWithClass(m.eRoot, p.eRootClass)
	assert.False(t, it.HasNext())
	assert.Panics(t, func() { it.Next() })
}

// all objects are instances of EObject, and a reference typed EObject may hold objects of any class
func TestEAllContentsWithClass_EObject(t *testing.T) {
	m := newEAllContentsClassModel()
	anything := GetFactory().CreateEReference()
	anything.SetName("anything")
	anything.SetEType(GetPackage().GetEObject())
	anything.SetContainment(true)
	anything.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	m.root.GetEStructuralFeatures().Add(anything)

	root := m.eFactory.Create(m.root)
	nested := m.eFactory.Create(m.root)
	nested.EGet(m.books).(EList).Add(m.eFactory.Create(m.book))
	root.EGet(anything).(EList).Add(nested)
	root.EGet(m.items).(EList).Add(m.eFactory.Create(m.book))

	assert.Equal(t, 3, countWithClass(root, GetPackage().GetEObject()))
	assert.Equal(t, 2, countWithClass(root, m.book))
	assert.Equal(t, filterInstances(AllContents(root), m.book), slices.Collect(AllContentsWithClass(root, m.book)))
}

// a class created after a table is walked with all its containments
func TestEAllContentsWithTable_UnknownClass(t *testing.T) {
	m := newEAllContentsClassModel()
	table := NewEClassTransitionsTable(m.root, m.book)
	// a sub class of item holding books, unknown to the table
	shelf := GetFactory().CreateEClass()
	shelf.SetName("Shelf")
	shelf.GetESuperTypes().Add(m.item)
	m.root.GetEPackage().GetEClassifiers().Add(shelf)
	shelfBooks := GetFactory().CreateEReference()
	shelfBooks.SetName("books")
	shelfBooks.SetEType(m.book)
	shelfBooks.SetContainment(true)
	shelfBooks.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	shelf.GetEStructuralFeatures().Add(shelfBooks)

	root := m.eFactory.Create(m.root)
	s := m.eFactory.Create(shelf)
	s.EGet(shelfBooks).(EList).Add(m.eFactory.Create(m.book))
	root.EGet(m.items).(EList).Add(s)

	expected := filterInstances(AllContents(root), m.book)
	require.Len(t, expected, 1)
	assert.Equal(t, expected, collectIterator(EAllContentsWithTable(root, table)))
	assert.Equal(t, expected, slices.Collect(AllContentsWithTable(root, table)))
}

// Prune skips the contents of the last instance returned: the top formations only
func TestEAllContentsWithClass_Prune(t *testing.T) {
	p := loadEAllContentsPackage(t)
	m := loadEAllContentsModel(t, p.ePackage)
	names := []string{}
	for it := EAllContentsWithClass(m.eRoot, p.eFormationClass).(*eAllContentsIterator); it.HasNext(); {
		names = append(names, it.Next().(EObject).EGet(p.eFormationNameAttribute).(string))
		it.Prune()
	}
	assert.Equal(t, []string{"formation 1", "formation 3", "formation 4"}, names)
}
