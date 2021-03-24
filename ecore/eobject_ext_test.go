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

func TestEObjectEInvokeFromID(t *testing.T) {
	o := NewEObjectImpl()
	assert.Panics(t, func() { o.EInvokeFromID(-1, nil) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__EALL_CONTENTS, NewImmutableEList([]interface{}{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECLASS, NewImmutableEList([]interface{}{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECONTAINER, NewImmutableEList([]interface{}{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECONTAINING_FEATURE, NewImmutableEList([]interface{}{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECONTAINMENT_FEATURE, NewImmutableEList([]interface{}{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECONTENTS, NewImmutableEList([]interface{}{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECROSS_REFERENCES, NewImmutableEList([]interface{}{})) })
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.On("GetName").Return("name")
		o.EInvokeFromID(EOBJECT__EGET_ESTRUCTURALFEATURE, NewImmutableEList([]interface{}{mockFeature}))
	})
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.On("GetName").Return("name")
		o.EInvokeFromID(EOBJECT__EGET_ESTRUCTURALFEATURE_EBOOLEAN, NewImmutableEList([]interface{}{mockFeature, true}))
	})
	assert.Panics(t, func() {
		mockOperation := new(MockEOperation)
		mockOperation.On("GetName").Return("name")
		o.EInvokeFromID(EOBJECT__EINVOKE_EOPERATION_EELIST, NewImmutableEList([]interface{}{mockOperation, NewEmptyBasicEList()}))
	})
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__EIS_PROXY, NewImmutableEList([]interface{}{})) })
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.On("GetName").Return("name")
		o.EInvokeFromID(EOBJECT__EIS_SET_ESTRUCTURALFEATURE, NewImmutableEList([]interface{}{mockFeature}))
	})
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ERESOURCE, NewImmutableEList([]interface{}{})) })
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.On("GetName").Return("name")
		o.EInvokeFromID(EOBJECT__ESET_ESTRUCTURALFEATURE_EJAVAOBJECT, NewImmutableEList([]interface{}{mockFeature, interface{}(nil)}))
	})
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.On("GetName").Return("name")
		o.EInvokeFromID(EOBJECT__EUNSET_ESTRUCTURALFEATURE, NewImmutableEList([]interface{}{mockFeature}))
	})
}
