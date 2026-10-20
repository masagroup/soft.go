// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// The packages and models of the core tests are built in code: the codecs reading .ecore, .xmi
// and .xml files are in other modules (codec/xml, codec/xmi).

func newTestPackage(name string, nsURI string, nsPrefix string) EPackage {
	ePackage := GetFactory().CreateEPackage()
	ePackage.SetName(name)
	ePackage.SetNsURI(nsURI)
	ePackage.SetNsPrefix(nsPrefix)
	return ePackage
}

func newTestClass(ePackage EPackage, name string) EClass {
	eClass := GetFactory().CreateEClass()
	eClass.SetName(name)
	ePackage.GetEClassifiers().Add(eClass)
	return eClass
}

func newTestAttribute(eClass EClass, name string, eType EDataType) EAttribute {
	eAttribute := GetFactory().CreateEAttribute()
	eAttribute.SetName(name)
	eAttribute.SetEType(eType)
	eClass.GetEStructuralFeatures().Add(eAttribute)
	return eAttribute
}

func newTestReference(eClass EClass, name string, eType EClass, many bool, containment bool, resolveProxies bool) EReference {
	eReference := GetFactory().CreateEReference()
	eReference.SetName(name)
	eReference.SetEType(eType)
	if many {
		eReference.SetUpperBound(UNBOUNDED_MULTIPLICITY)
	}
	eReference.SetContainment(containment)
	eReference.SetResolveProxies(resolveProxies)
	eClass.GetEStructuralFeatures().Add(eReference)
	return eReference
}

// newTestObject creates an object of the class named className in ePackage, with the attribute
// name set when it is not empty.
func newTestObject(ePackage EPackage, className string, name string) EObject {
	eClass := ePackage.GetEClassifier(className).(EClass)
	eObject := ePackage.GetEFactoryInstance().Create(eClass)
	if name != "" {
		eObject.ESet(eClass.GetEStructuralFeatureFromName("name"), name)
	}
	return eObject
}

// addTestValue adds value to the many-valued feature of eObject named featureName.
func addTestValue(eObject EObject, featureName string, value EObject) {
	eObject.EGet(eObject.EClass().GetEStructuralFeatureFromName(featureName)).(EList).Add(value)
}

// newDeletePackage is the package of testdata/delete.ecore in soft.go/codec/xml.
func newDeletePackage() EPackage {
	ePackage := newTestPackage("delete", "http:///org/eclipse/emf/examples/delete/delete.ecore/1.0.0", "delete")
	eRoot := newTestClass(ePackage, "Root")
	eA := newTestClass(ePackage, "A")
	eB := newTestClass(ePackage, "B")
	eO := newTestClass(ePackage, "O")
	eC := newTestClass(ePackage, "C")
	eString := GetPackage().GetEString()
	newTestReference(eRoot, "a", eA, true, true, false)
	newTestReference(eRoot, "b", eB, true, false, true)
	newTestReference(eRoot, "o", eO, true, true, false)
	newTestReference(eRoot, "c", eC, true, true, false)
	newTestAttribute(eA, "name", eString)
	newTestReference(eA, "b", eB, false, true, false)
	newTestAttribute(eB, "name", eString)
	newTestAttribute(eO, "name", eString)
	newTestReference(eO, "a", eA, true, false, true)
	newTestAttribute(eC, "name", eString)
	newTestReference(eC, "m", eC, false, false, true)
	return ePackage
}

// newDeleteResource is the resource of testdata/delete.xmi in soft.go/codec/xml.
func newDeleteResource(ePackage EPackage) EResource {
	root := newTestObject(ePackage, "Root", "")
	as := []EObject{}
	for _, n := range []string{"1", "2", "3"} {
		a := newTestObject(ePackage, "A", "a"+n)
		b := newTestObject(ePackage, "B", "b"+n)
		a.ESet(a.EClass().GetEStructuralFeatureFromName("b"), b)
		addTestValue(root, "a", a)
		as = append(as, a)
	}
	for _, a := range as {
		addTestValue(root, "b", a.EGet(a.EClass().GetEStructuralFeatureFromName("b")).(EObject))
	}
	o1 := newTestObject(ePackage, "O", "o1")
	addTestValue(o1, "a", as[0])
	addTestValue(o1, "a", as[1])
	addTestValue(root, "o", o1)
	c1 := newTestObject(ePackage, "C", "c1")
	c2 := newTestObject(ePackage, "C", "c2")
	c3 := newTestObject(ePackage, "C", "c3")
	cM := c1.EClass().GetEStructuralFeatureFromName("m")
	c2.ESet(cM, c1)
	c3.ESet(cM, c1)
	addTestValue(root, "c", c1)
	addTestValue(root, "c", c2)
	addTestValue(root, "c", c3)

	eResourceSet := NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(NewURI("memory://delete.xmi"))
	eResource.GetContents().Add(root)
	return eResource
}

// newEAllContentsPackage is the package of testdata/eallcontents.ecore in soft.go/codec/xml.
func newEAllContentsPackage() EPackage {
	ePackage := newTestPackage("eallcontents", "http:///org/eclipse/emf/examples/eallcontents/eallcontents.ecore/1.0.0", "eac")
	eRoot := newTestClass(ePackage, "Root")
	eParty := newTestClass(ePackage, "Party")
	eFormation := newTestClass(ePackage, "Formation")
	eAutomat := newTestClass(ePackage, "Automat")
	eUnit := newTestClass(ePackage, "Unit")
	eTheater := newTestClass(ePackage, "Theater")
	eString := GetPackage().GetEString()
	newTestReference(eRoot, "theater", eTheater, false, true, true)
	newTestReference(eParty, "formations", eFormation, true, true, true)
	newTestAttribute(eParty, "name", eString)
	newTestAttribute(eFormation, "name", eString)
	newTestReference(eFormation, "formations", eFormation, true, true, true)
	newTestReference(eFormation, "automats", eAutomat, true, true, true)
	newTestAttribute(eAutomat, "name", eString)
	newTestReference(eAutomat, "units", eUnit, true, true, true)
	newTestAttribute(eUnit, "name", eString)
	newTestReference(eTheater, "parties", eParty, true, true, true)
	return ePackage
}

// newEAllContentsResource is the resource of testdata/eallcontents.xml in soft.go/codec/xml.
func newEAllContentsResource(ePackage EPackage) EResource {
	automat := func(name string, units ...string) EObject {
		a := newTestObject(ePackage, "Automat", name)
		for _, u := range units {
			addTestValue(a, "units", newTestObject(ePackage, "Unit", u))
		}
		return a
	}
	formation := func(name string, children ...EObject) EObject {
		f := newTestObject(ePackage, "Formation", name)
		for _, c := range children {
			if c.EClass().GetName() == "Formation" {
				addTestValue(f, "formations", c)
			} else {
				addTestValue(f, "automats", c)
			}
		}
		return f
	}
	party := func(name string, formations ...EObject) EObject {
		p := newTestObject(ePackage, "Party", name)
		for _, f := range formations {
			addTestValue(p, "formations", f)
		}
		return p
	}
	theater := newTestObject(ePackage, "Theater", "")
	addTestValue(theater, "parties", party("party 1",
		formation("formation 1",
			formation("formation 2", automat("automat 1", "unit 1"), automat("automat 2", "unit 2")),
			automat("automat 3", "unit 3"),
			automat("automat 4")),
		formation("formation 3", automat("automat 5", "unit 4"))))
	addTestValue(theater, "parties", party("party 2"))
	addTestValue(theater, "parties", party("party 3",
		formation("formation 4", automat("automat 6", "unit 5"))))
	root := newTestObject(ePackage, "Root", "")
	root.ESet(root.EClass().GetEStructuralFeatureFromName("theater"), theater)

	eResource := NewEResourceSetImpl().CreateResource(NewURI("memory://eallcontents.xml"))
	eResource.GetContents().Add(root)
	return eResource
}

// newLibrarySimplePackage is the package of testdata/library.simple.ecore in
// soft.go/codec/xml.
func newLibrarySimplePackage() EPackage {
	ePackage := newTestPackage("library", "http:///org/eclipse/emf/examples/library/library.simple.ecore/1.0.0", "lib")
	eLibrary := newTestClass(ePackage, "Library")
	eBook := newTestClass(ePackage, "Book")
	newTestAttribute(eLibrary, "owner", GetPackage().GetEString())
	newTestAttribute(eLibrary, "location", GetPackage().GetEString())
	newTestReference(eLibrary, "books", eBook, true, true, true)
	newTestAttribute(eBook, "name", GetPackage().GetEString())
	newTestAttribute(eBook, "isbn", GetPackage().GetEInt())
	return ePackage
}
