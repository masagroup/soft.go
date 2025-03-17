package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONCodec_NewEncoder(t *testing.T) {
	mockResource := NewMockEResource(t)
	codec := &JSONCodec{}
	require.NotNil(t, codec.NewEncoder(mockResource, nil, nil))
}

func TestJSONCodec_NewDecoder(t *testing.T) {
	require.Panics(t, func() {
		mockResource := NewMockEResource(t)
		codec := &JSONCodec{}
		codec.NewDecoder(mockResource, nil, nil)
	})
}

func TestGetJSONCodecFeatureKind_Transient(t *testing.T) {
	mockFeature := NewMockEStructuralFeature(t)
	mockFeature.EXPECT().IsTransient().Return(true).Once()
	require.Equal(t, jfkTransient, getJSONCodecFeatureKind(mockFeature))
	mockFeature.EXPECT().IsTransient().Return(false).Once()
	require.Equal(t, jsonFeatureKind(-1), getJSONCodecFeatureKind(mockFeature))
}

func TestGetJSONCodecFeatureKind_Attribute(t *testing.T) {
	mockAttribute := NewMockEAttribute(t)
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(false).Once()
	require.Equal(t, jfkData, getJSONCodecFeatureKind(mockAttribute))
	mockAttribute.EXPECT().IsTransient().Return(false).Once()
	mockAttribute.EXPECT().IsMany().Return(true).Once()
	require.Equal(t, jfkDataList, getJSONCodecFeatureKind(mockAttribute))
}

func TestGetJSONCodecFeatureKind_Reference(t *testing.T) {
	mockReference := NewMockEReference(t)

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	require.Equal(t, jfkObject, getJSONCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	require.Equal(t, jfkObjectList, getJSONCodecFeatureKind(mockReference))

	mockOpposite := NewMockEReference(t)
	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(mockOpposite).Once()
	mockOpposite.EXPECT().IsContainment().Return(true).Once()
	require.Equal(t, jfkTransient, getJSONCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	require.Equal(t, jfkObjectReferenceList, getJSONCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(true).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	require.Equal(t, jfkObjectReference, getJSONCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsMany().Return(false).Once()
	require.Equal(t, jfkObject, getJSONCodecFeatureKind(mockReference))

	mockReference.EXPECT().IsTransient().Return(false).Once()
	mockReference.EXPECT().IsContainment().Return(false).Once()
	mockReference.EXPECT().GetEOpposite().Return(nil).Once()
	mockReference.EXPECT().IsResolveProxies().Return(false).Once()
	mockReference.EXPECT().IsMany().Return(true).Once()
	require.Equal(t, jfkObjectList, getJSONCodecFeatureKind(mockReference))
}
