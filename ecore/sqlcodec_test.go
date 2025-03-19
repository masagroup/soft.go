package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSQLCodecFeatureKind_Transient(t *testing.T) {
	mockFeature := NewMockEStructuralFeature(t)
	mockFeature.EXPECT().IsTransient().Return(true).Once()
	require.Equal(t, sfkTransient, getSQLCodecFeatureKind(mockFeature))
	mockFeature.EXPECT().IsTransient().Return(false).Once()
	require.Equal(t, sqlFeatureKind(-1), getSQLCodecFeatureKind(mockFeature))
}

func TestGetSQLCodecFeatureKind_Attribute(t *testing.T) {
	mockAttribute := NewMockEAttribute(t)
	mockDataType := NewMockEDataType(t)
	mockEnumType := NewMockEEnum(t)

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(true).Once()
	require.Equal(t, sfkDataList, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("").Once()
	require.Equal(t, sfkData, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("float64").Once()
	require.Equal(t, sfkFloat64, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("float32").Once()
	require.Equal(t, sfkFloat32, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("int").Once()
	require.Equal(t, sfkInt, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("int64").Once()
	require.Equal(t, sfkInt64, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("int32").Once()
	require.Equal(t, sfkInt32, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("int16").Once()
	require.Equal(t, sfkInt16, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("byte").Once()
	require.Equal(t, sfkByte, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("bool").Once()
	require.Equal(t, sfkBool, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("string").Once()
	require.Equal(t, sfkString, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("[]byte").Once()
	require.Equal(t, sfkByteArray, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockDataType).Once()
	mockDataType.EXPECT().GetInstanceTypeName().Return("java.util.Date").Once()
	require.Equal(t, sfkDate, getSQLCodecFeatureKind(mockAttribute))

	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	mockAttribute.EXPECT().GetEAttributeType().Return(mockEnumType).Once()
	require.Equal(t, sfkEnum, getSQLCodecFeatureKind(mockAttribute))
}

func TestGetSQLCodecFeatureKind_Reference(t *testing.T) {
	mockReference := NewMockEReference(t)

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	require.Equal(t, sfkObject, getSQLCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	require.Equal(t, sfkObjectList, getSQLCodecFeatureKind(mockReference))

	mockOpposite := NewMockEReference(t)
	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	require.Equal(t, sfkTransient, getSQLCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	require.Equal(t, sfkObjectReferenceList, getSQLCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	require.Equal(t, sfkObjectReference, getSQLCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(true).Once()
	require.Equal(t, sfkTransient, getSQLCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(false).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	require.Equal(t, sfkObject, getSQLCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsContainer().Return(false).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	require.Equal(t, sfkObjectList, getSQLCodecFeatureKind(mockReference))
}

func TestNewSQLEncoder(t *testing.T) {
	codec := &SQLCodec{}
	mockResource := NewMockEResource(t)
	mockResource.EXPECT().GetURI().Return(NewURI("")).Twice()
	mockResource.EXPECT().GetObjectIDManager().Return(nil).Once()
	require.NotNil(t, codec.NewEncoder(mockResource, nil, nil))
}

func TestNewSQLDecoder(t *testing.T) {
	codec := &SQLCodec{}
	mockResource := NewMockEResource(t)
	mockResource.EXPECT().GetResourceSet().Return(nil).Once()
	mockResource.EXPECT().GetURI().Return(NewURI("")).Once()
	mockResource.EXPECT().GetObjectIDManager().Return(nil).Once()
	require.NotNil(t, codec.NewDecoder(mockResource, nil, nil))
}
