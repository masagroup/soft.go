package xml

import (
	"github.com/masagroup/soft.go/ecore"
)

// libraryModel is a dynamic model covering the feature kinds of the codec: containments, single
// and many references, with and without opposite, unsettable features and an abstract type.
type libraryModel struct {
	ePackage  ecore.EPackage
	named     ecore.EClass
	library   ecore.EClass
	book      ecore.EClass
	novel     ecore.EClass
	writer    ecore.EClass
	name      ecore.EAttribute
	pages     ecore.EAttribute
	tags      ecore.EAttribute
	nick      ecore.EAttribute
	books     ecore.EReference
	writers   ecore.EReference
	main      ecore.EReference
	optional  ecore.EReference
	best      ecore.EReference
	favorites ecore.EReference
	any       ecore.EReference
	maybe     ecore.EReference
	others    ecore.EReference
	author    ecore.EReference
	written   ecore.EReference
}

func newLibraryModel() *libraryModel {
	f := ecore.GetFactory()
	ecorePackage := ecore.GetPackage()
	m := &libraryModel{}
	m.ePackage = f.CreateEPackage()
	m.ePackage.SetName("lib")
	m.ePackage.SetNsURI("http://lib")
	m.ePackage.SetNsPrefix("lib")
	m.ePackage.SetEFactoryInstance(f.CreateEFactory())

	newClass := func(name string, superTypes ...ecore.EClass) ecore.EClass {
		eClass := f.CreateEClass()
		eClass.SetName(name)
		for _, superType := range superTypes {
			eClass.GetESuperTypes().Add(superType)
		}
		m.ePackage.GetEClassifiers().Add(eClass)
		return eClass
	}
	newAttribute := func(eClass ecore.EClass, name string, eType ecore.EDataType) ecore.EAttribute {
		eAttribute := f.CreateEAttribute()
		eAttribute.SetName(name)
		eAttribute.SetEType(eType)
		eClass.GetEStructuralFeatures().Add(eAttribute)
		return eAttribute
	}
	newReference := func(eClass ecore.EClass, name string, eType ecore.EClass, containment bool, many bool) ecore.EReference {
		eReference := f.CreateEReference()
		eReference.SetName(name)
		eReference.SetEType(eType)
		eReference.SetContainment(containment)
		if many {
			eReference.SetUpperBound(ecore.UNBOUNDED_MULTIPLICITY)
		}
		eClass.GetEStructuralFeatures().Add(eReference)
		return eReference
	}

	m.named = newClass("Named")
	m.named.SetAbstract(true)
	m.name = newAttribute(m.named, "name", ecorePackage.GetEString())
	m.book = newClass("Book", m.named)
	m.pages = newAttribute(m.book, "pages", ecorePackage.GetEInt())
	m.novel = newClass("Novel", m.book)
	m.writer = newClass("Writer", m.named)
	m.library = newClass("Library", m.named)
	m.tags = newAttribute(m.library, "tags", ecorePackage.GetEString())
	m.tags.SetUpperBound(ecore.UNBOUNDED_MULTIPLICITY)
	m.nick = newAttribute(m.library, "nick", ecorePackage.GetEString())
	m.nick.SetUnsettable(true)
	m.books = newReference(m.library, "books", m.book, true, true)
	m.writers = newReference(m.library, "writers", m.writer, true, true)
	m.main = newReference(m.library, "main", m.writer, true, false)
	m.optional = newReference(m.library, "optional", m.writer, true, false)
	m.optional.SetUnsettable(true)
	m.best = newReference(m.library, "best", m.book, false, false)
	m.favorites = newReference(m.library, "favorites", m.book, false, true)
	m.any = newReference(m.library, "any", m.named, false, false)
	m.maybe = newReference(m.library, "maybe", m.book, false, false)
	m.maybe.SetUnsettable(true)
	m.others = newReference(m.library, "others", m.book, false, true)
	m.others.SetUnsettable(true)
	m.author = newReference(m.book, "author", m.writer, false, false)
	m.written = newReference(m.writer, "books", m.book, false, true)
	m.author.SetEOpposite(m.written)
	m.written.SetEOpposite(m.author)
	return m
}

func (m *libraryModel) create(eClass ecore.EClass, name string) ecore.EObject {
	eObject := m.ePackage.GetEFactoryInstance().Create(eClass)
	eObject.ESet(m.name, name)
	return eObject
}

func (m *libraryModel) newResourceSet() ecore.EResourceSet {
	rs := ecore.NewEResourceSetImpl()
	rs.GetPackageRegistry().RegisterPackage(m.ePackage)
	return rs
}

// newStringIDManager is an object id manager of string ids
func newStringIDManager() ecore.EObjectIDManager {
	return ecore.NewUniqueIDManager(
		func() string { return "" },
		func(id string) bool { return id != "" },
		func(id any) (string, error) { s, _ := id.(string); return s, nil },
		func(string) {})
}
