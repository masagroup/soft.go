package sqlite

import (
	"bytes"
	"os"
	"testing"

	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// classes of the same name in different packages have their own tables
func TestSQLEncoder_SameClassNameInTwoPackages(t *testing.T) {
	f := ecore.GetFactory()
	newPackage := func(name, attribute string) (ecore.EPackage, ecore.EClass, ecore.EAttribute) {
		ePackage := f.CreateEPackage()
		ePackage.SetName(name)
		ePackage.SetNsURI("http://" + name)
		ePackage.SetEFactoryInstance(f.CreateEFactory())
		eClass := f.CreateEClass()
		eClass.SetName("Item")
		ePackage.GetEClassifiers().Add(eClass)
		eAttribute := f.CreateEAttribute()
		eAttribute.SetName(attribute)
		eAttribute.SetEType(ecore.GetPackage().GetEString())
		eClass.GetEStructuralFeatures().Add(eAttribute)
		return ePackage, eClass, eAttribute
	}
	p1, c1, title := newPackage("p1", "title")
	p2, c2, weight := newPackage("p2", "weight")
	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(p1)
	eResourceSet.GetPackageRegistry().RegisterPackage(p2)
	eResource := eResourceSet.CreateResource(ecore.NewURI("file:///items.sqlite"))
	o1 := p1.GetEFactoryInstance().Create(c1)
	o1.ESet(title, "t")
	o2 := p2.GetEFactoryInstance().Create(c2)
	o2.ESet(weight, "w")
	eResource.GetContents().Add(o1)
	eResource.GetContents().Add(o2)

	var buffer bytes.Buffer
	NewWriterEncoder(&buffer, eResource, map[string]any{OPTION_IN_MEMORY_DATABASE: true}).EncodeResource()
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))

	eDecoded := eResourceSet.CreateResource(ecore.NewURI("file:///items-decoded.sqlite"))
	NewReaderDecoder(bytes.NewReader(buffer.Bytes()), eDecoded, map[string]any{OPTION_IN_MEMORY_DATABASE: true}).DecodeResource()
	require.True(t, eDecoded.GetErrors().Empty(), diagnosticError(eDecoded.GetErrors()))
	require.Equal(t, 2, eDecoded.GetContents().Size())
	assert.Equal(t, "t", eDecoded.GetContents().Get(0).(ecore.EObject).EGet(title))
	assert.Equal(t, "w", eDecoded.GetContents().Get(1).(ecore.EObject).EGet(weight))
}

// a database written before the qualified tables, without their property, is decoded with the
// tables named by class: as the same model written with qualified tables and decoded
func TestSQLDecoder_LegacyTableNames(t *testing.T) {
	for _, name := range []string{"library.simple", "library.complex"} {
		ePackage := loadPackage(name + ".ecore")
		require.NotNil(t, ePackage)
		eResourceSet := ecore.NewEResourceSetImpl()
		eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)

		// the model written with qualified tables, decoded
		xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{ePackage}))
		model := xmlProcessor.LoadWithOptions(ecore.NewURI("testdata/"+name+".xml"), nil)
		require.True(t, model.GetErrors().Empty(), diagnosticError(model.GetErrors()))
		var buffer bytes.Buffer
		NewWriterEncoder(&buffer, model, map[string]any{OPTION_IN_MEMORY_DATABASE: true}).EncodeResource()
		require.True(t, model.GetErrors().Empty(), diagnosticError(model.GetErrors()))
		expected := eResourceSet.CreateResource(ecore.NewURI("file:///" + name + ".qualified.sqlite"))
		NewReaderDecoder(bytes.NewReader(buffer.Bytes()), expected, map[string]any{OPTION_IN_MEMORY_DATABASE: true}).DecodeResource()
		require.True(t, expected.GetErrors().Empty(), diagnosticError(expected.GetErrors()))

		// the legacy database
		legacy := eResourceSet.CreateResource(ecore.NewURI("testdata/legacy/" + name + ".sqlite"))
		r, err := os.Open("testdata/legacy/" + name + ".sqlite")
		require.NoError(t, err)
		NewReaderDecoder(r, legacy, nil).DecodeResource()
		require.NoError(t, r.Close())
		require.True(t, legacy.GetErrors().Empty(), diagnosticError(legacy.GetErrors()))
		if name == "library.complex" {
			// written before the forward references were encoded by id: they are proxies, so the
			// models are only compared by their objects
			assert.Equal(t, countContents(expected), countContents(legacy), name)
			continue
		}
		assert.True(t, ecore.EqualsAll(expected.GetContents(), legacy.GetContents()), name)
	}
}

func countContents(eResource ecore.EResource) int {
	count := 0
	for it := eResource.GetAllContents(); it.HasNext(); it.Next() {
		count++
	}
	return count
}
