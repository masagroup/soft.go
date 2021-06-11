package ecore

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXMIResourceLoadLibrarySimple(t *testing.T) {
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(&URI{Path: "testdata/library.simple.ecore"})
	require.NotNil(t, resource)
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))

	contents := resource.GetContents()
	assert.Equal(t, 1, contents.Size())

	ePackage, _ := contents.Get(0).(EPackage)
	assert.NotNil(t, ePackage)
	assert.Equal(t, "library", ePackage.GetName())
	assert.Equal(t, "lib", ePackage.GetNsPrefix())
	assert.Equal(t, "http:///org/eclipse/emf/examples/library/library.simple.ecore/1.0.0", ePackage.GetNsURI())

	eClassifiers := ePackage.GetEClassifiers()
	assert.Equal(t, 2, eClassifiers.Size())

	eLibrary, _ := eClassifiers.Get(0).(EClassifier)
	assert.NotNil(t, eLibrary)
	assert.Equal(t, "Library", eLibrary.GetName())

	eLibraryClass, _ := eLibrary.(EClass)
	assert.NotNil(t, eLibraryClass)
	assert.Equal(t, 3, eLibraryClass.GetFeatureCount())

	eOwnerFeature := eLibraryClass.GetEStructuralFeature(0)
	assert.Equal(t, "owner", eOwnerFeature.GetName())
	eOwnerAttribute, _ := eOwnerFeature.(EAttribute)
	assert.NotNil(t, eOwnerAttribute)
	eOwnerAttributeType := eOwnerAttribute.GetEAttributeType()
	assert.Equal(t, "EString", eOwnerAttributeType.GetName())

	eLocationFeature := eLibraryClass.GetEStructuralFeature(1)
	assert.Equal(t, "location", eLocationFeature.GetName())
	eLocationAttribute, _ := eLocationFeature.(EAttribute)
	assert.NotNil(t, eLocationAttribute)
	eLocationType := eLocationAttribute.GetEAttributeType()
	assert.NotNil(t, eLocationType)

	eBooksFeature := eLibraryClass.GetEStructuralFeature(2)
	assert.Equal(t, "books", eBooksFeature.GetName())
	eBooksReference, _ := eBooksFeature.(EReference)
	assert.NotNil(t, eBooksReference)

	eBook := eClassifiers.Get(1).(EClassifier)
	assert.NotNil(t, eBook)
	assert.Equal(t, "Book", eBook.GetName())
	eBookClass, _ := eBook.(EClass)
	assert.NotNil(t, eBookClass)
	assert.Equal(t, 2, eBookClass.GetFeatureCount())

	eNameFeature := eBookClass.GetEStructuralFeature(0)
	assert.Equal(t, "name", eNameFeature.GetName())
	eNameAttribute, _ := eNameFeature.(EAttribute)
	assert.NotNil(t, eNameAttribute)

	eISBNFeature := eBookClass.GetEStructuralFeature(1)
	assert.Equal(t, "isbn", eISBNFeature.GetName())
	eISBNAttribute, _ := eISBNFeature.(EAttribute)
	assert.NotNil(t, eISBNAttribute)

	// check resolved reference
	assert.Equal(t, eBookClass, eBooksReference.GetEReferenceType())
}

func TestXMIResourceLoadLibraryNoRoot(t *testing.T) {
	xmiProcessor := NewXMIProcessor()
	resource, _ := xmiProcessor.Load(&URI{Path: "testdata/library.noroot.ecore"}).(XMIResource)
	require.NotNil(t, resource)
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
	assert.Equal(t, "2.0", resource.GetXMIVersion())
	assert.Equal(t, "1.0", resource.GetXMLVersion())
	assert.Equal(t, "UTF-8", resource.GetEncoding())

	contents := resource.GetContents()
	assert.Equal(t, 1, contents.Size())

	ePackage, _ := contents.Get(0).(EPackage)
	assert.NotNil(t, ePackage)

	eClassifiers := ePackage.GetEClassifiers()
	eBook, _ := eClassifiers.Get(0).(EClassifier)
	assert.NotNil(t, eBook)
	assert.Equal(t, "Book", eBook.GetName())

	eBookClass, _ := eBook.(EClass)
	assert.NotNil(t, eBookClass)
	superTypes := eBookClass.GetESuperTypes()
	assert.Equal(t, 1, superTypes.Size())
	eCirculationItemClass := superTypes.Get(0).(EClass)
	assert.Equal(t, "CirculatingItem", eCirculationItemClass.GetName())

	eWriter, _ := eClassifiers.Get(2).(EClass)
	assert.NotNil(t, eWriter)
	assert.False(t, eWriter.GetEAnnotations().Empty())
	eAnnotation := eWriter.GetEAnnotation("http://net.masagroup/soft/2019/GenGo")
	assert.NotNil(t, eAnnotation)
	assert.Equal(t, "true", eAnnotation.GetDetails().GetValue("extension"))
}

func TestXMIResourceLoadLibraryComplex(t *testing.T) {
	xmiProcessor := NewXMIProcessor()
	resource, _ := xmiProcessor.Load(&URI{Path: "testdata/library.complex.ecore"}).(XMIResource)
	require.NotNil(t, resource)
	assert.True(t, resource.IsLoaded())
	assert.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	assert.True(t, resource.GetWarnings().Empty(), diagnosticError(resource.GetWarnings()))
	assert.Equal(t, "2.0", resource.GetXMIVersion())

	contents := resource.GetContents()
	assert.Equal(t, 1, contents.Size())

	ePackage, _ := contents.Get(0).(EPackage)
	assert.NotNil(t, ePackage)

	eClassifiers := ePackage.GetEClassifiers()
	eDocumentRootClass, _ := eClassifiers.Get(0).(EClass)
	assert.NotNil(t, eDocumentRootClass)

	eXMNLSPrefixFeature, _ := eDocumentRootClass.GetEStructuralFeatureFromName("xMLNSPrefixMap").(EReference)
	assert.NotNil(t, eXMNLSPrefixFeature)

	eType := eXMNLSPrefixFeature.GetEType()
	assert.NotNil(t, eType)
	assert.Equal(t, "EStringToStringMapEntry", eType.GetName())
	assert.False(t, eType.EIsProxy())
}

func TestXMIResourceSaveLibrarySimple(t *testing.T) {
	// load/save
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(&URI{Path: "testdata/library.simple.ecore"})
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, nil)
	// check
	bytes, err := ioutil.ReadFile("testdata/library.simple.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMIResourceSaveLibraryNoRoot(t *testing.T) {
	// load/save
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(&URI{Path: "testdata/library.noroot.ecore"})
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, nil)
	// check
	bytes, err := ioutil.ReadFile("testdata/library.noroot.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func TestXMIResourceSaveLibraryComplex(t *testing.T) {
	// load/save
	xmiProcessor := NewXMIProcessor()
	resource := xmiProcessor.Load(&URI{Path: "testdata/library.complex.ecore"})
	require.NotNil(t, resource)
	result := xmiProcessor.SaveToString(resource, nil)
	// check
	bytes, err := ioutil.ReadFile("testdata/library.complex.ecore")
	assert.Nil(t, err)
	assert.Equal(t, strings.ReplaceAll(string(bytes), "\r\n", "\n"), strings.ReplaceAll(result, "\r\n", "\n"))
}

func BenchmarkXMIResourceLoadSaveLibrarySimple(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resource := NewXMIResourceImpl()
		resource.SetURI(&URI{Path: "testdata/library.simple.ecore"})
		resource.Load()

		var strbuff strings.Builder
		resource.SaveWithWriter(&strbuff, nil)
		resource = nil
	}
}

func BenchmarkXMIResourceLoadSaveLibraryNoRoot(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resource := NewXMIResourceImpl()
		resource.SetURI(&URI{Path: "testdata/library.noroot.ecore"})
		resource.Load()

		var strbuff strings.Builder
		resource.SaveWithWriter(&strbuff, nil)
		resource = nil
	}
}
