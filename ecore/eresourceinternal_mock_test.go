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

func TestMockEResourceInternalDoLoad(t *testing.T) {
	mockResource := NewMockEResourceInternal(t)
	mockDecoder := NewMockEResourceDecoder(t)
	m := NewMockRun(t, mockDecoder)
	mockResource.EXPECT().DoLoad(mockDecoder).Return().Run(func(decoder EResourceDecoder) { m.Run(decoder) }).Once()
	mockResource.DoLoad(mockDecoder)
}

func TestMockEResourceInternalDoSave(t *testing.T) {
	mockResource := NewMockEResourceInternal(t)
	mockEncoder := NewMockEEncoder(t)
	m := NewMockRun(t, mockEncoder)
	mockResource.EXPECT().DoSave(mockEncoder).Return().Run(func(encoder EEncoder) { m.Run(encoder) }).Once()
	mockResource.DoSave(mockEncoder)
}

func TestMockEResourceInternalDoUnLoad(t *testing.T) {
	m := NewMockRun(t)
	mockResource := NewMockEResourceInternal(t)
	mockResource.EXPECT().DoUnload().Return().Run(func() { m.Run() }).Once()
	mockResource.DoUnload()
}

func TestMockEResourceInternalBasicSetLoaded(t *testing.T) {
	mockResource := NewMockEResourceInternal(t)
	n1 := NewMockENotificationChain(t)
	n2 := NewMockENotificationChain(t)
	m := NewMockRun(t, false, n1)
	mockResource.EXPECT().BasicSetLoaded(false, n1).Return(n2).Run(func(_a0 bool, _a1 ENotificationChain) { m.Run(_a0, _a1) }).Once()
	mockResource.EXPECT().BasicSetLoaded(false, n1).Call.Return(func(bool, ENotificationChain) ENotificationChain {
		return n2
	}).Once()
	assert.Equal(t, n2, mockResource.BasicSetLoaded(false, n1))
	assert.Equal(t, n2, mockResource.BasicSetLoaded(false, n1))
}

func TestMockEResourceInternalBasicSetResourceSet(t *testing.T) {
	mockResource := NewMockEResourceInternal(t)
	rs := NewMockEResourceSet(t)
	n1 := NewMockENotificationChain(t)
	n2 := NewMockENotificationChain(t)
	m := NewMockRun(t, rs, n1)
	mockResource.EXPECT().BasicSetResourceSet(rs, n1).Return(n2).Run(func(_a0 EResourceSet, _a1 ENotificationChain) { m.Run(_a0, _a1) }).Once()
	mockResource.EXPECT().BasicSetResourceSet(rs, n1).Call.Return(func(EResourceSet, ENotificationChain) ENotificationChain {
		return n2
	}).Once()
	assert.Equal(t, n2, mockResource.BasicSetResourceSet(rs, n1))
	assert.Equal(t, n2, mockResource.BasicSetResourceSet(rs, n1))
}

func TestMockEResourceInternalDoAttached(t *testing.T) {
	mockResource := NewMockEResourceInternal(t)
	mockObject := NewMockEObject(t)
	m := NewMockRun(t, mockObject)
	mockResource.EXPECT().DoAttached(mockObject).Return().Run(func(o EObject) { m.Run(o) }).Once()
	mockResource.DoAttached(mockObject)
}

func TestMockEResourceInternalDetached(t *testing.T) {
	mockResource := NewMockEResourceInternal(t)
	mockObject := NewMockEObject(t)
	m := NewMockRun(t, mockObject)
	mockResource.EXPECT().DoDetached(mockObject).Return().Run(func(o EObject) { m.Run(o) }).Once()
	mockResource.DoDetached(mockObject)
}

func TestMockEResourceInternalIsAttachedDetachedRequired(t *testing.T) {
	mockResource := NewMockEResourceInternal(t)
	m := NewMockRun(t)
	mockResource.EXPECT().IsAttachedDetachedRequired().Return(true).Run(func() { m.Run() }).Once()
	mockResource.EXPECT().IsAttachedDetachedRequired().Call.Return(func() bool {
		return false
	}).Once()
	assert.True(t, mockResource.IsAttachedDetachedRequired())
	assert.False(t, mockResource.IsAttachedDetachedRequired())
}
