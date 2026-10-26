package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBasicEDataTypeListAccessors(t *testing.T) {
	{
		list := NewBasicEDataTypeList(nil, 1, false, false)
		assert.Equal(t, nil, list.GetNotifier())
		assert.Equal(t, nil, list.GetFeature())
		assert.Equal(t, 1, list.GetFeatureID())
	}
	{
		mockOwner := NewMockEObjectInternal(t)
		list := NewBasicEDataTypeList(mockOwner, 1, true, true)
		assert.Equal(t, mockOwner, list.GetNotifier())
		assert.Equal(t, 1, list.GetFeatureID())
		mockClass := NewMockEClass(t)
		mockFeature := NewMockEStructuralFeature(t)
		mockClass.EXPECT().GetEStructuralFeature(1).Return(mockFeature)
		mockOwner.EXPECT().EClass().Return(mockClass)
		assert.Equal(t, mockFeature, list.GetFeature())
		mock.AssertExpectationsForObjects(t, mockOwner, mockClass, mockFeature, mockClass)
	}
	{
		// not unsettable: set when not empty
		list := NewBasicEDataTypeList(nil, 1, true, false)
		list.Clear()
		assert.False(t, list.IsSet())
		// unsettable: set from its first change
		list = NewBasicEDataTypeList(nil, 1, true, true)
		assert.False(t, list.IsSet())
		list.Clear()
		assert.True(t, list.IsSet())
	}
}
