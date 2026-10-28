package xml

import (
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sibling elements declaring namespaces reuse the slots of the namespaces of their previous sibling
func TestXMLDecoder_SiblingNamespaces(t *testing.T) {
	r := ecore.NewEResourceImpl()
	r.SetURI(ecore.NewURI("file:///a.xml"))
	require.NotPanics(t, func() {
		NewDecoder(r, strings.NewReader(`<a:x xmlns:a="u"><a:y xmlns:b="v"/><a:z xmlns:c="w"/></a:x>`), nil).DecodeResource()
	})
}

// an invalid literal is reported as a diagnostic, the other values are decoded
func TestXMLDecoder_InvalidLiteral(t *testing.T) {
	f := ecore.GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetName("p")
	ePackage.SetNsURI("http://p")
	ePackage.SetNsPrefix("p")
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eInt := f.CreateEDataType()
	eInt.SetName("Int")
	eInt.SetInstanceTypeName("int")
	ePackage.GetEClassifiers().Add(eInt)
	eClass := f.CreateEClass()
	eClass.SetName("C")
	ePackage.GetEClassifiers().Add(eClass)
	count := f.CreateEAttribute()
	count.SetName("count")
	count.SetEType(eInt)
	eClass.GetEStructuralFeatures().Add(count)
	name := f.CreateEAttribute()
	name.SetName("name")
	name.SetEType(ecore.GetPackage().GetEString())
	eClass.GetEStructuralFeatures().Add(name)

	rs := ecore.NewEResourceSetImpl()
	rs.GetPackageRegistry().RegisterPackage(ePackage)
	r := rs.CreateResource(ecore.NewURI("file:///c.xml"))
	require.NotPanics(t, func() {
		NewDecoder(r, strings.NewReader(`<p:C xmlns:p="http://p" count="12x" name="n"/>`), nil).DecodeResource()
	})
	require.Equal(t, 1, r.GetErrors().Size())
	illegal, isIllegal := r.GetErrors().Get(0).(*IllegalValueDiagnostic)
	require.True(t, isIllegal, "a refused value is told from the other diagnostics: %T", r.GetErrors().Get(0))
	assert.Equal(t, count, illegal.Feature)
	assert.Equal(t, eInt, illegal.DataType)
	assert.Equal(t, "12x", illegal.Value)
	assert.Equal(t, "file:///c.xml", illegal.GetLocation())
	require.Equal(t, 1, r.GetContents().Size())
	eObject := r.GetContents().Get(0).(ecore.EObject)
	assert.False(t, eObject.EIsSet(count))
	assert.Equal(t, "n", eObject.EGet(name))
}

// each refused token of a list written as one attribute is reported, the accepted ones are kept
func TestXMLDecoder_InvalidLiteralInAList(t *testing.T) {
	f := ecore.GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetName("p")
	ePackage.SetNsURI("http://p")
	ePackage.SetNsPrefix("p")
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eInt := f.CreateEDataType()
	eInt.SetName("Int")
	eInt.SetInstanceTypeName("int")
	ePackage.GetEClassifiers().Add(eInt)
	eClass := f.CreateEClass()
	eClass.SetName("C")
	ePackage.GetEClassifiers().Add(eClass)
	counts := f.CreateEAttribute()
	counts.SetName("counts")
	counts.SetEType(eInt)
	counts.SetUpperBound(ecore.UNBOUNDED_MULTIPLICITY)
	eClass.GetEStructuralFeatures().Add(counts)

	rs := ecore.NewEResourceSetImpl()
	rs.GetPackageRegistry().RegisterPackage(ePackage)
	r := rs.CreateResource(ecore.NewURI("file:///c.xml"))
	require.NotPanics(t, func() {
		NewDecoder(r, strings.NewReader(`<p:C xmlns:p="http://p" counts="1 x 3 y"/>`), nil).DecodeResource()
	})
	var refused []string
	for it := r.GetErrors().Iterator(); it.HasNext(); {
		illegal, isIllegal := it.Next().(*IllegalValueDiagnostic)
		require.True(t, isIllegal)
		assert.Equal(t, counts, illegal.Feature)
		refused = append(refused, illegal.Value)
	}
	assert.Equal(t, []string{"x", "y"}, refused)
	eObject := r.GetContents().Get(0).(ecore.EObject)
	assert.Equal(t, []any{1, 3}, eObject.EGet(counts).(ecore.EList).ToArray())
}
