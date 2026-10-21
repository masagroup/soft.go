package binary

import (
	"bytes"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a data type whose Go type is given by the GenGo annotation is encoded as that type
func TestBinaryCodec_GenGoInstanceType(t *testing.T) {
	f := ecore.GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetName("p")
	ePackage.SetNsURI("http://p")
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eDataType := f.CreateEDataType()
	eDataType.SetName("Small")
	eDataType.SetInstanceTypeName("java.lang.Integer")
	eAnnotation := f.CreateEAnnotation()
	eAnnotation.SetSource("http://net.masagroup/soft/2019/GenGo")
	eAnnotation.GetDetails().Put("instanceTypeName", "int32")
	eDataType.GetEAnnotations().Add(eAnnotation)
	ePackage.GetEClassifiers().Add(eDataType)
	eClass := f.CreateEClass()
	eClass.SetName("C")
	ePackage.GetEClassifiers().Add(eClass)
	eAttribute := f.CreateEAttribute()
	eAttribute.SetName("v")
	eAttribute.SetEType(eDataType)
	eClass.GetEStructuralFeatures().Add(eAttribute)

	eObject := ePackage.GetEFactoryInstance().Create(eClass)
	eObject.ESet(eAttribute, ePackage.GetEFactoryInstance().CreateFromString(eDataType, "7"))
	eResourceSet := ecore.NewEResourceSetImpl()
	eResourceSet.GetPackageRegistry().RegisterPackage(ePackage)
	eResource := eResourceSet.CreateResource(ecore.NewURI("file:///o.bin"))
	eResource.GetContents().Add(eObject)

	var buffer bytes.Buffer
	require.NotPanics(t, func() { NewEncoder(eResource, &buffer, nil).EncodeResource() })
	require.True(t, eResource.GetErrors().Empty())

	eDecoded := eResourceSet.CreateResource(ecore.NewURI("file:///o-decoded.bin"))
	NewDecoder(eDecoded, bytes.NewReader(buffer.Bytes()), nil).DecodeResource()
	require.True(t, eDecoded.GetErrors().Empty())
	assert.Equal(t, int32(7), eDecoded.GetContents().Get(0).(ecore.EObject).EGet(eAttribute))
}
