// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEStructuralFeatureExtGetDefaultValue(t *testing.T) {
	{
		o := newEStructuralFeatureExt()
		assert.Nil(t, o.GetDefaultValue())
	}
	{
		// with a type & no literal
		o := newEStructuralFeatureExt()
		mockType := NewMockEDataType(t)
		o.SetEType(mockType)

		mockDefaultValue := NewMockEObject(t)
		mockType.EXPECT().EIsProxy().Return(false).Once()
		mockType.EXPECT().GetDefaultValue().Return(mockDefaultValue).Once()
		assert.Equal(t, mockDefaultValue, o.GetDefaultValue())
		mock.AssertExpectationsForObjects(t, mockType, mockDefaultValue)
	}
	{
		// with a type & no literal & many
		o := newEStructuralFeatureExt()
		mockType := NewMockEDataType(t)
		o.SetEType(mockType)
		o.SetUpperBound(UNBOUNDED_MULTIPLICITY)

		mockType.EXPECT().EIsProxy().Return(false).Once()
		assert.Nil(t, o.GetDefaultValue())
		mock.AssertExpectationsForObjects(t, mockType)
	}
	{
		// with a type & literal
		o := newEStructuralFeatureExt()
		mockType := NewMockEDataType(t)
		mockDefaultValue := NewMockEObject(t)
		mockPackage := NewMockEPackage(t)
		mockFactory := NewMockEFactory(t)
		o.SetEType(mockType)
		o.SetDefaultValueLiteral("defaultLiteralValue")
		mockType.EXPECT().EIsProxy().Return(false).Once()
		mockType.EXPECT().GetEPackage().Return(mockPackage).Once()
		mockType.EXPECT().IsSerializable().Return(true).Once()
		mockPackage.EXPECT().GetEFactoryInstance().Return(mockFactory).Once()
		mockFactory.EXPECT().CreateFromString(mockType, "defaultLiteralValue").Return(mockDefaultValue).Once()
		assert.Equal(t, mockDefaultValue, o.GetDefaultValue())
		mock.AssertExpectationsForObjects(t, mockType, mockDefaultValue, mockPackage, mockFactory)
	}
}

func TestEStructuralFeatureExtSetDefaultValue(t *testing.T) {
	{
		o := newEStructuralFeatureExt()
		assert.Panics(t, func() {
			o.SetDefaultValue(nil)
		})
	}
	{
		o := newEStructuralFeatureExt()
		mockType := NewMockEDataType(t)
		mockDefaultValue := NewMockEObject(t)
		mockPackage := NewMockEPackage(t)
		mockFactory := NewMockEFactory(t)
		o.SetEType(mockType)
		mockType.EXPECT().EIsProxy().Return(false).Once()
		mockType.EXPECT().GetEPackage().Return(mockPackage).Once()
		mockPackage.EXPECT().GetEFactoryInstance().Return(mockFactory).Once()
		mockFactory.EXPECT().ConvertToString(mockType, mockDefaultValue).Return("defaultValueLiteral").Once()
		o.SetDefaultValue(mockDefaultValue)
		mock.AssertExpectationsForObjects(t, mockType, mockDefaultValue, mockPackage, mockFactory)
	}
}
