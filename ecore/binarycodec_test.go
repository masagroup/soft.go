package ecore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBinaryCodec_NewEncoder(t *testing.T) {
	c := &BinaryCodec{}
	mockResource := &MockEResource{}
	mockResource.On("GetURI").Return(nil).Once()
	e := c.NewEncoder(mockResource, nil, nil)
	require.NotNil(t, e)
	mock.AssertExpectationsForObjects(t, mockResource)
}

func TestBinaryCodec_NewDecoder(t *testing.T) {
	c := &BinaryCodec{}
	mockResource := &MockEResource{}
	mockResource.On("GetURI").Return(nil).Once()
	e := c.NewDecoder(mockResource, nil, nil)
	require.NotNil(t, e)
	mock.AssertExpectationsForObjects(t, mockResource)
}

func TestBinaryCodec_GetFeatureKind_Reference(t *testing.T) {
	mockReference := &MockEReference{}
	mockReference.On("IsContainment").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(true).Once()
	mockReference.On("IsMany").Return(true).Once()
	assert.Equal(t, bfkObjectContainmentListProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(true).Once()
	mockReference.On("IsMany").Return(false).Once()
	assert.Equal(t, bfkObjectContainmentProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(false).Once()
	mockReference.On("IsMany").Return(true).Once()
	assert.Equal(t, bfkObjectContainmentList, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(false).Once()
	mockReference.On("IsMany").Return(false).Once()
	assert.Equal(t, bfkObjectContainment, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(false).Once()
	mockReference.On("IsContainer").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(true).Once()
	assert.Equal(t, bfkObjectContainerProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(false).Once()
	mockReference.On("IsContainer").Return(true).Once()
	mockReference.On("IsResolveProxies").Return(false).Once()
	assert.Equal(t, bfkObjectContainer, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(false).Once()
	mockReference.On("IsContainer").Return(false).Once()
	mockReference.On("IsResolveProxies").Return(true).Once()
	mockReference.On("IsMany").Return(true).Once()
	assert.Equal(t, bfkObjectListProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(false).Once()
	mockReference.On("IsContainer").Return(false).Once()
	mockReference.On("IsResolveProxies").Return(true).Once()
	mockReference.On("IsMany").Return(false).Once()
	assert.Equal(t, bfkObjectProxy, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(false).Once()
	mockReference.On("IsContainer").Return(false).Once()
	mockReference.On("IsResolveProxies").Return(false).Once()
	mockReference.On("IsMany").Return(true).Once()
	assert.Equal(t, bfkObjectList, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)

	mockReference.On("IsContainment").Return(false).Once()
	mockReference.On("IsContainer").Return(false).Once()
	mockReference.On("IsResolveProxies").Return(false).Once()
	mockReference.On("IsMany").Return(false).Once()
	assert.Equal(t, bfkObject, getBinaryCodecFeatureKind(mockReference))
	mockReference.AssertExpectations(t)
}

func TestBinaryCodec_GetFeatureKind_Attribute(t *testing.T) {
	mockAttribute := &MockEAttribute{}
	mockAttribute.On("IsMany").Return(true).Once()
	assert.Equal(t, bfkDataList, getBinaryCodecFeatureKind(mockAttribute))
	mockAttribute.AssertExpectations(t)

	mockEnum := &MockEEnum{}
	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockEnum).Once()
	assert.Equal(t, bfkEnum, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockEnum)

	mockDataType := &MockEDataType{}
	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("float64").Once()
	assert.Equal(t, bfkFloat64, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("float32").Once()
	assert.Equal(t, bfkFloat32, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("int").Once()
	assert.Equal(t, bfkInt, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("int64").Once()
	assert.Equal(t, bfkInt64, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("int32").Once()
	assert.Equal(t, bfkInt32, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("int16").Once()
	assert.Equal(t, bfkInt16, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("byte").Once()
	assert.Equal(t, bfkByte, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("bool").Once()
	assert.Equal(t, bfkBool, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("string").Once()
	assert.Equal(t, bfkString, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("[]byte").Once()
	assert.Equal(t, bfkByteArray, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("*time.Time").Once()
	assert.Equal(t, bfkDate, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)

	mockAttribute.On("IsMany").Return(false).Once()
	mockAttribute.On("GetEAttributeType").Return(mockDataType).Once()
	mockDataType.On("GetInstanceTypeName").Return("data").Once()
	assert.Equal(t, bfkData, getBinaryCodecFeatureKind(mockAttribute))
	mock.AssertExpectationsForObjects(t, mockAttribute, mockDataType)
}

func TestBinaryCodec_EncodeDecodeEcore(t *testing.T) {
	// load package
	ePackage := loadPackage("library.simple.ecore")
	require.NotNil(t, ePackage)

	// encode package resource in binary format
	buffer := bytes.Buffer{}
	c := &BinaryCodec{}
	encoder := c.NewEncoder(ePackage.EResource(), &buffer, nil)
	encoder.Encode()

	// decode buffer into another resource
	eNewResource := NewEResourceImpl()
	decoder := c.NewDecoder(eNewResource, &buffer, nil)
	decoder.Decode()
	require.True(t, eNewResource.GetErrors().Empty(), diagnosticError(eNewResource.GetErrors()))

	eNewPackage, _ := eNewResource.GetContents().Get(0).(EPackage)
	require.NotNil(t, eNewPackage)

	// retrieve document root class , library class & library name attribute
	eLibraryClass, _ := eNewPackage.GetEClassifier("Library").(EClass)
	require.NotNil(t, eLibraryClass)
	eLibraryOwnerAttribute, _ := eLibraryClass.GetEStructuralFeatureFromName("owner").(EAttribute)
	require.NotNil(t, eLibraryOwnerAttribute)
	eDataType := eLibraryOwnerAttribute.GetEAttributeType()
	require.NotNil(t, eDataType)
	assert.Equal(t, "EString", eDataType.GetName())

}
