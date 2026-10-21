package xml

import (
	"bytes"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// the values with XML special characters and whitespaces are escaped: the document is read back
// with the same values
func TestXMLEncoder_Escape(t *testing.T) {
	f := ecore.GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetName("p")
	ePackage.SetNsURI("http://p")
	ePackage.SetNsPrefix("p")
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eClass := f.CreateEClass()
	eClass.SetName("C")
	ePackage.GetEClassifiers().Add(eClass)
	single := f.CreateEAttribute()
	single.SetName("single")
	single.SetEType(ecore.GetPackage().GetEString())
	eClass.GetEStructuralFeatures().Add(single)
	many := f.CreateEAttribute()
	many.SetName("many")
	many.SetUpperBound(ecore.UNBOUNDED_MULTIPLICITY)
	many.SetEType(ecore.GetPackage().GetEString())
	eClass.GetEStructuralFeatures().Add(many)

	eObject := ePackage.GetEFactoryInstance().Create(eClass)
	eObject.ESet(single, "a\nb\tc\r<&>\"'")
	eObject.EGet(many).(ecore.EList).Add("x<y&z")

	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(ecore.NewURI("file:///escape.xml"))
	eResource.GetContents().Add(eObject)
	var buffer bytes.Buffer
	NewEncoder(eResource, &buffer, nil).EncodeResource()
	require.True(t, eResource.GetErrors().Empty())
	assert.Contains(t, buffer.String(), `single="a&#xA;b&#x9;c&#13;&lt;&amp;&gt;&#34;&#39;"`)
	assert.Contains(t, buffer.String(), `<many>x&lt;y&amp;z</many>`)

	eDecoded := eResourceSet.CreateResource(ecore.NewURI("file:///escape-decoded.xml"))
	NewDecoder(eDecoded, bytes.NewReader(buffer.Bytes()), nil).DecodeResource()
	require.True(t, eDecoded.GetErrors().Empty())
	require.Equal(t, 1, eDecoded.GetContents().Size())
	eDecodedObject := eDecoded.GetContents().Get(0).(ecore.EObject)
	assert.Equal(t, "a\nb\tc\r<&>\"'", eDecodedObject.EGet(single))
	assert.Equal(t, []any{"x<y&z"}, eDecodedObject.EGet(many).(ecore.EList).ToArray())
}
