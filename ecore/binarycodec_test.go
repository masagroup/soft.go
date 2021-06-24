package ecore

import (
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
