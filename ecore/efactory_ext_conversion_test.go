package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newConversionDataType(instanceTypeName string) (EFactory, EDataType) {
	f := GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eDataType := f.CreateEDataType()
	eDataType.SetName("T")
	eDataType.SetInstanceTypeName(instanceTypeName)
	ePackage.GetEClassifiers().Add(eDataType)
	return ePackage.GetEFactoryInstance(), eDataType
}

func TestEFactoryExt_ByteArrayRoundTrip(t *testing.T) {
	eFactory, eDataType := newConversionDataType("[]byte")
	value := eFactory.CreateFromString(eDataType, "abc")
	assert.Equal(t, []byte("abc"), value)
	assert.Equal(t, "abc", eFactory.ConvertToString(eDataType, value))
}

// an invalid literal panics, as EMF throws, instead of becoming a zero value
func TestEFactoryExt_InvalidLiteral(t *testing.T) {
	for _, instanceTypeName := range []string{"int", "int64", "int32", "int16", "int8", "uint64", "float64", "float32", "bool"} {
		eFactory, eDataType := newConversionDataType(instanceTypeName)
		assert.Panics(t, func() { eFactory.CreateFromString(eDataType, "12x") }, instanceTypeName)
	}
	eFactory, eDataType := newConversionDataType("int")
	assert.Equal(t, 12, eFactory.CreateFromString(eDataType, "12"))
}
