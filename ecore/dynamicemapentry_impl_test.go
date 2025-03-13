package ecore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDynamicEMapEntryImpl_Accessors(t *testing.T) {
	me := NewDynamicEMapEntryImpl()
	mockClass := NewMockEClass(t)
	mockKeyFeature := NewMockEStructuralFeature(t)
	mockValueFeature := NewMockEStructuralFeature(t)
	mockClass.EXPECT().GetFeatureCount().Return(2).Once()
	mockClass.EXPECT().GetEStructuralFeatureFromName("key").Return(mockKeyFeature).Once()
	mockClass.EXPECT().GetEStructuralFeatureFromName("value").Return(mockValueFeature).Once()
	me.SetEClass(mockClass)
	require.Equal(t, mockClass, me.EClass())

	mockClass.EXPECT().GetFeatureID(mockKeyFeature).Return(0).Twice()
	mockClass.EXPECT().GetEStructuralFeature(0).Return(mockKeyFeature).Twice()
	me.SetKey("key")
	require.Equal(t, "key", me.GetKey())

	mockClass.EXPECT().GetFeatureID(mockValueFeature).Return(1).Twice()
	mockClass.EXPECT().GetEStructuralFeature(1).Return(mockValueFeature).Twice()
	me.SetValue("value")
	require.Equal(t, "value", me.GetValue())
}
