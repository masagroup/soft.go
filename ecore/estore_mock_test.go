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

func TestMockEStoreGet(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockResult := NewMockEObject(t)
	m := NewMockRun(t, mockObject, mockFeature, 0)
	mockEStore.EXPECT().Get(mockObject, mockFeature, 0).Return(mockResult).Run(func(object EObject, feature EStructuralFeature, index int) { m.Run(object, feature, index) }).Once()
	mockEStore.EXPECT().Get(mockObject, mockFeature, 0).Call.Return(func(EObject, EStructuralFeature, int) any {
		return mockResult
	}).Once()
	assert.Equal(t, mockResult, mockEStore.Get(mockObject, mockFeature, 0))
	assert.Equal(t, mockResult, mockEStore.Get(mockObject, mockFeature, 0))
}

func TestMockEStoreSet(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockValue := NewMockEObject(t)
	mockOld := NewMockEObject(t)
	m := NewMockRun(t, mockObject, mockFeature, 0, mockValue)
	mockEStore.EXPECT().Set(mockObject, mockFeature, 0, mockValue).Return(mockOld).Run(func(object EObject, feature EStructuralFeature, index int, value any) {
		m.Run(object, feature, index, value)
	}).Once()
	mockEStore.EXPECT().Set(mockObject, mockFeature, 0, mockValue).Call.Return(func(object EObject, feature EStructuralFeature, index int, value any) any {
		return mockOld
	}).Once()
	assert.Equal(t, mockOld, mockEStore.Set(mockObject, mockFeature, 0, mockValue))
	assert.Equal(t, mockOld, mockEStore.Set(mockObject, mockFeature, 0, mockValue))
}

func TestMockEStoreIsSet(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	m := NewMockRun(t, mockObject, mockFeature)
	mockEStore.EXPECT().IsSet(mockObject, mockFeature).Return(false).Run(func(object EObject, feature EStructuralFeature) {
		m.Run(object, feature)
	}).Once()
	mockEStore.EXPECT().IsSet(mockObject, mockFeature).Call.Return(func(EObject, EStructuralFeature) bool {
		return true
	}).Once()
	assert.False(t, mockEStore.IsSet(mockObject, mockFeature))
	assert.True(t, mockEStore.IsSet(mockObject, mockFeature))
}

func TestMockEStoreUnSet(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	m := NewMockRun(t, mockObject, mockFeature)
	mockEStore.EXPECT().UnSet(mockObject, mockFeature).Return().Run(func(object EObject, feature EStructuralFeature) {
		m.Run(object, feature)
	}).Once()
	mockEStore.UnSet(mockObject, mockFeature)
	mockEStore.AssertExpectations(t)
}

func TestMockEStoreIsEmpty(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	m := NewMockRun(t, mockObject, mockFeature)
	mockEStore.EXPECT().IsEmpty(mockObject, mockFeature).Return(false).Run(func(object EObject, feature EStructuralFeature) {
		m.Run(object, feature)
	}).Once()
	mockEStore.EXPECT().IsEmpty(mockObject, mockFeature).Call.Return(func(EObject, EStructuralFeature) bool {
		return true
	}).Once()
	assert.False(t, mockEStore.IsEmpty(mockObject, mockFeature))
	assert.True(t, mockEStore.IsEmpty(mockObject, mockFeature))
}

func TestMockEStoreContains(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockValue := NewMockEObject(t)
	m := NewMockRun(t, mockObject, mockFeature, mockValue)
	mockEStore.EXPECT().Contains(mockObject, mockFeature, mockValue).Return(false).Run(func(object EObject, feature EStructuralFeature, value any) {
		m.Run(object, feature, value)
	}).Once()
	mockEStore.EXPECT().Contains(mockObject, mockFeature, mockValue).Call.Return(func(EObject, EStructuralFeature, any) bool {
		return true
	}).Once()
	assert.False(t, mockEStore.Contains(mockObject, mockFeature, mockValue))
	assert.True(t, mockEStore.Contains(mockObject, mockFeature, mockValue))
}

func TestMockEStoreSize(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	m := NewMockRun(t, mockObject, mockFeature)
	mockEStore.EXPECT().Size(mockObject, mockFeature).Return(1).Run(func(object EObject, feature EStructuralFeature) {
		m.Run(object, feature)
	}).Once()
	mockEStore.EXPECT().Size(mockObject, mockFeature).Call.Return(func(EObject, EStructuralFeature) int {
		return 2
	}).Once()
	assert.Equal(t, 1, mockEStore.Size(mockObject, mockFeature))
	assert.Equal(t, 2, mockEStore.Size(mockObject, mockFeature))
}

func TestMockEStoreIndexOf(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockValue := NewMockEObject(t)
	m := NewMockRun(t, mockObject, mockFeature, mockValue)
	mockEStore.EXPECT().IndexOf(mockObject, mockFeature, mockValue).Return(1).Run(func(object EObject, feature EStructuralFeature, value any) {
		m.Run(object, feature, value)
	}).Once()
	mockEStore.EXPECT().IndexOf(mockObject, mockFeature, mockValue).Call.Return(func(EObject, EStructuralFeature, any) int {
		return 2
	}).Once()
	assert.Equal(t, 1, mockEStore.IndexOf(mockObject, mockFeature, mockValue))
	assert.Equal(t, 2, mockEStore.IndexOf(mockObject, mockFeature, mockValue))
}

func TestMockEStoreLastIndexOf(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockValue := NewMockEObject(t)
	m := NewMockRun(t, mockObject, mockFeature, mockValue)
	mockEStore.EXPECT().LastIndexOf(mockObject, mockFeature, mockValue).Return(1).Run(func(object EObject, feature EStructuralFeature, value any) {
		m.Run(object, feature, value)
	}).Once()
	mockEStore.EXPECT().LastIndexOf(mockObject, mockFeature, mockValue).Call.Return(func(EObject, EStructuralFeature, any) int {
		return 2
	}).Once()
	assert.Equal(t, 1, mockEStore.LastIndexOf(mockObject, mockFeature, mockValue))
	assert.Equal(t, 2, mockEStore.LastIndexOf(mockObject, mockFeature, mockValue))
}

func TestMockEStoreAdd(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockValue := NewMockEObject(t)
	m := NewMockRun(t, mockObject, mockFeature, 0, mockValue)
	mockEStore.EXPECT().Add(mockObject, mockFeature, 0, mockValue).Return().Run(func(object EObject, feature EStructuralFeature, index int, value interface{}) {
		m.Run(object, feature, index, value)
	}).Once()
	mockEStore.Add(mockObject, mockFeature, 0, mockValue)
}

func TestMockEStoreRemove(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockOld := NewMockEObject(t)
	m := NewMockRun(t, mockObject, mockFeature, 0)
	mockEStore.EXPECT().Remove(mockObject, mockFeature, 0).Return(mockOld).Run(func(object EObject, feature EStructuralFeature, index int) {
		m.Run(object, feature, index)
	}).Once()
	mockEStore.EXPECT().Remove(mockObject, mockFeature, 0).Call.Return(func(object EObject, feature EStructuralFeature, index int) any {
		return mockOld
	}).Once()
	assert.Equal(t, mockOld, mockEStore.Remove(mockObject, mockFeature, 0))
	assert.Equal(t, mockOld, mockEStore.Remove(mockObject, mockFeature, 0))
}

func TestMockEStoreMove(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockOld := NewMockEObject(t)
	m := NewMockRun(t, mockObject, mockFeature, 0, 1)
	mockEStore.EXPECT().Move(mockObject, mockFeature, 0, 1).Return(mockOld).Run(func(object EObject, feature EStructuralFeature, targetIndex, sourceIndex int) {
		m.Run(object, feature, targetIndex, sourceIndex)
	}).Once()
	mockEStore.EXPECT().Move(mockObject, mockFeature, 0, 1).Call.Return(func(object EObject, feature EStructuralFeature, index int, old int) any {
		return mockOld
	}).Once()
	assert.Equal(t, mockOld, mockEStore.Move(mockObject, mockFeature, 0, 1))
	assert.Equal(t, mockOld, mockEStore.Move(mockObject, mockFeature, 0, 1))
}

func TestMockEStoreClear(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	m := NewMockRun(t, mockObject, mockFeature)
	mockEStore.EXPECT().Clear(mockObject, mockFeature).Return().Run(func(object EObject, feature EStructuralFeature) {
		m.Run(object, feature)
	}).Once()
	mockEStore.Clear(mockObject, mockFeature)
	mockEStore.AssertExpectations(t)
}

func TestMockEStoreToArray(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockFeature := NewMockEStructuralFeature(t)
	mockResult := []any{}
	m := NewMockRun(t, mockObject, mockFeature)
	mockEStore.EXPECT().ToArray(mockObject, mockFeature).Return(mockResult).Run(func(object EObject, feature EStructuralFeature) {
		m.Run(object, feature)
	}).Once()
	mockEStore.EXPECT().ToArray(mockObject, mockFeature).Call.Return(func(EObject, EStructuralFeature) []any {
		return mockResult
	}).Once()
	assert.Equal(t, mockResult, mockEStore.ToArray(mockObject, mockFeature))
	assert.Equal(t, mockResult, mockEStore.ToArray(mockObject, mockFeature))
}

func TestMockEStoreGetContainer(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockResult := NewMockEObject(t)
	m := NewMockRun(t, mockObject)
	mockEStore.EXPECT().GetContainer(mockObject).Return(mockResult).Run(func(object EObject) {
		m.Run(object)
	}).Once()
	mockEStore.EXPECT().GetContainer(mockObject).Call.Return(func(EObject) EObject {
		return mockResult
	}).Once()
	assert.Equal(t, mockResult, mockEStore.GetContainer(mockObject))
	assert.Equal(t, mockResult, mockEStore.GetContainer(mockObject))
}

func TestMockEStoreGetContainingFeature(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockObject := NewMockEObject(t)
	mockResult := NewMockEStructuralFeature(t)
	m := NewMockRun(t, mockObject)
	mockEStore.EXPECT().GetContainingFeature(mockObject).Return(mockResult).Run(func(object EObject) {
		m.Run(object)
	}).Once()
	mockEStore.EXPECT().GetContainingFeature(mockObject).Call.Return(func(EObject) EStructuralFeature {
		return mockResult
	}).Once()
	assert.Equal(t, mockResult, mockEStore.GetContainingFeature(mockObject))
	assert.Equal(t, mockResult, mockEStore.GetContainingFeature(mockObject))
}

func TestMockEStoreCreate(t *testing.T) {
	mockEStore := NewMockEStore(t)
	mockClass := NewMockEClass(t)
	mockResult := NewMockEObject(t)
	m := NewMockRun(t, mockClass)
	mockEStore.EXPECT().Create(mockClass).Return(mockResult).Run(func(eClass EClass) { m.Run(eClass) }).Once()
	mockEStore.EXPECT().Create(mockClass).Call.Return(func(EClass) EObject {
		return mockResult
	}).Once()
	assert.Equal(t, mockResult, mockEStore.Create(mockClass))
	assert.Equal(t, mockResult, mockEStore.Create(mockClass))
}
