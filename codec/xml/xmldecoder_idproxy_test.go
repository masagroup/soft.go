// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xml

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sharedFactory is the factory instance of a package shared by several implementations: it
// converts the data types, and leaves the objects to the factory of an implementation.
type sharedFactory struct {
	ecore.EFactory
}

func (f *sharedFactory) Create(eClass ecore.EClass) ecore.EObject {
	panic("no implementation chosen for " + eClass.GetName())
}

// implementationFactory is the factory of an implementation of a shared package, registered in
// the package registry of a resource set: it counts the objects it creates.
type implementationFactory struct {
	ecore.EFactory
	ePackage ecore.EPackage
	created  int
}

func (f *implementationFactory) GetEPackage() ecore.EPackage {
	return f.ePackage
}

func (f *implementationFactory) Create(eClass ecore.EClass) ecore.EObject {
	f.created++
	eObject := ecore.NewDynamicEObjectImpl()
	eObject.SetEClass(eClass)
	return eObject
}

// An id that resolves to nothing in the document is read as a proxy carrying it
// (OPTION_ID_REFERENCE_AS_PROXY), created as the objects of the document are: by the factory of
// the resource set for the package of the reference type, not by the factory instance of a
// shared package.
func TestXMLDecoder_IDReferenceAsProxy(t *testing.T) {
	m := newLibraryModel()
	m.ePackage.SetEFactoryInstance(&sharedFactory{EFactory: ecore.GetFactory().CreateEFactory()})
	implementation := &implementationFactory{EFactory: ecore.GetFactory().CreateEFactory(), ePackage: m.ePackage}
	rs := m.newResourceSet()
	rs.GetPackageRegistry().PutFactory(m.ePackage.GetNsURI(), implementation)
	eResource := rs.CreateResource(ecore.NewURI("file:///library.xml"))

	library := decodeResource(t, eResource, `<lib:Library xmlns:lib="http://lib" name="l" best="b1"/>`,
		map[string]any{OPTION_ID_REFERENCE_AS_PROXY: true})

	best, _ := library.EGetResolve(m.best, false).(ecore.EObjectInternal)
	require.NotNil(t, best)
	require.True(t, best.EIsProxy())
	assert.Equal(t, "file:///library.xml#b1", best.EProxyURI().String())
	assert.Equal(t, m.book, best.EClass())
	assert.Equal(t, 2, implementation.created, "the library and the proxy")
}
