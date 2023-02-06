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
)

func TestEFactoryExtCreate(t *testing.T) {
	f := newEFactoryExt()
	mockClass := NewMockEClass(t)
	mockPackage := NewMockEPackage(t)
	mockClass.EXPECT().GetEPackage().Return(mockPackage).Once()
	mockClass.EXPECT().GetName().Return("mockClass").Once()
	assert.Panics(t, func() { f.Create(mockClass) })
}

func TestEFactoryExt_CreateFromString_IncompatiblePackage(t *testing.T) {
	f := newEFactoryExt()
	mockEPackage := NewMockEPackage(t)
	mockEDataType := NewMockEDataType(t)
	mockEDataType.EXPECT().GetEPackage().Return(mockEPackage).Once()
	mockEDataType.EXPECT().GetName().Return("name").Once()
	assert.Panics(t, func() {
		f.CreateFromString(mockEDataType, "")
	})
}

func TestEFactoryExt_CreateFromString_TypeNotSupported(t *testing.T) {
	f := newEFactoryExt()
	mockEDataType := NewMockEDataType(t)
	mockEDataType.EXPECT().GetEPackage().Return(nil).Once()
	mockEDataType.EXPECT().GetEAnnotation("http://net.masagroup/soft/2019/GenGo").Return(nil).Once()
	mockEDataType.EXPECT().GetInstanceTypeName().Return("unknown").Once()
	assert.Panics(t, func() {
		f.CreateFromString(mockEDataType, "literalValue")
	})
}

func TestEFactoryExt_CreateFromString_EEnum(t *testing.T) {
	f := newEFactoryExt()
	mockEEnum := NewMockEEnum(t)
	mockEEnum.EXPECT().GetEPackage().Return(nil).Once()
	mockEEnum.EXPECT().GetEEnumLiteralByLiteral("literalValue").Return(nil).Once()
	mockEEnum.EXPECT().GetName().Return("name").Once()
	assert.Panics(t, func() {
		f.CreateFromString(mockEEnum, "literalValue")
	})

	mockEEnumLiteral := NewMockEEnumLiteral(t)
	mockEEnum.EXPECT().GetEPackage().Return(nil).Once()
	mockEEnum.EXPECT().GetEEnumLiteralByLiteral("literalValue").Return(mockEEnumLiteral).Once()
	mockEEnumLiteral.EXPECT().GetValue().Return(1).Once()
	assert.Equal(t, 1, f.CreateFromString(mockEEnum, "literalValue"))
}

func TestEFactoryExt_CreateFromString_Primitives(t *testing.T) {
	f := newEFactoryExt()
	mockEDataType := NewMockEDataType(t)
	for _, test := range []struct {
		instanceTypeName string
		literalValue     string
		expectedValue    any
	}{
		{"float64", "1.2", float64(1.2)},
		{"float32", "1.2", float32(1.2)},
		{"int", "1", int(1)},
		{"uint64", "1", uint64(1)},
		{"int64", "1", int64(1)},
		{"int32", "1", int32(1)},
		{"int16", "1", int16(1)},
		{"int8", "1", int8(1)},
		{"bool", "true", true},
		{"string", "string", "string"},
	} {
		mockEDataType.EXPECT().GetEPackage().Return(nil).Once()
		mockEDataType.EXPECT().GetEAnnotation("http://net.masagroup/soft/2019/GenGo").Return(nil).Once()
		mockEDataType.EXPECT().GetInstanceTypeName().Return(test.instanceTypeName).Once()
		assert.Equal(t, test.expectedValue, f.CreateFromString(mockEDataType, test.literalValue))
		mockEDataType.AssertExpectations(t)
	}
}

func TestEFactoryExt_ConvertToString_IncompatiblePackage(t *testing.T) {
	f := newEFactoryExt()
	mockEPackage := NewMockEPackage(t)
	mockEDataType := NewMockEDataType(t)
	mockEDataType.EXPECT().GetEPackage().Return(mockEPackage).Once()
	mockEDataType.EXPECT().GetName().Return("name").Once()
	assert.Panics(t, func() {
		f.ConvertToString(mockEDataType, "")
	})
}

func TestEFactoryExt_ConvertToString_TypeNotSupported(t *testing.T) {
	f := newEFactoryExt()
	mockEDataType := NewMockEDataType(t)
	mockEDataType.EXPECT().GetEPackage().Return(nil).Once()
	mockEDataType.EXPECT().GetEAnnotation("http://net.masagroup/soft/2019/GenGo").Return(nil).Once()
	mockEDataType.EXPECT().GetInstanceTypeName().Return("unknown").Once()
	assert.Panics(t, func() {
		f.ConvertToString(mockEDataType, "")
	})
}

func TestEFactoryExt_ConvertToString_EEnum(t *testing.T) {
	f := newEFactoryExt()
	mockEEnum := NewMockEEnum(t)
	mockEEnum.EXPECT().GetEPackage().Return(nil).Once()
	mockEEnum.EXPECT().GetEEnumLiteralByValue(0).Return(nil).Once()
	mockEEnum.EXPECT().GetName().Return("name").Once()
	assert.Panics(t, func() {
		f.ConvertToString(mockEEnum, 0)
	})

	mockEEnumLiteral := NewMockEEnumLiteral(t)
	mockEEnum.EXPECT().GetEPackage().Return(nil).Once()
	mockEEnum.EXPECT().GetEEnumLiteralByValue(0).Return(mockEEnumLiteral).Once()
	mockEEnumLiteral.EXPECT().GetLiteral().Return("enumValue").Once()
	assert.Equal(t, "enumValue", f.ConvertToString(mockEEnum, 0))
}

func TestEFactoryExt_ConvertToString_Primitives(t *testing.T) {
	f := newEFactoryExt()
	mockEDataType := NewMockEDataType(t)
	for _, test := range []struct {
		instanceTypeName string
		stringValue      string
		inputValue       any
	}{
		{"float64", "1.2", float64(1.2)},
		{"float32", "1.2", float32(1.2)},
		{"int", "1", int(1)},
		{"uint64", "1", uint64(1)},
		{"int64", "1", int64(1)},
		{"int32", "1", int32(1)},
		{"int16", "1", int16(1)},
		{"int8", "1", int8(1)},
		{"bool", "true", true},
		{"string", "string", "string"},
	} {
		mockEDataType.EXPECT().GetEPackage().Return(nil).Once()
		mockEDataType.EXPECT().GetEAnnotation("http://net.masagroup/soft/2019/GenGo").Return(nil).Once()
		mockEDataType.EXPECT().GetInstanceTypeName().Return(test.instanceTypeName).Once()
		assert.Equal(t, test.stringValue, f.ConvertToString(mockEDataType, test.inputValue))
		mockEDataType.AssertExpectations(t)
	}
}
