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
	o := newEObjectImpl()
	assert.Panics(t, func() { o.EInvokeFromID(-1, nil) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__EALL_CONTENTS, NewImmutableEList([]any{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECLASS, NewImmutableEList([]any{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECONTAINER, NewImmutableEList([]any{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECONTAINING_FEATURE, NewImmutableEList([]any{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECONTAINMENT_FEATURE, NewImmutableEList([]any{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECONTENTS, NewImmutableEList([]any{})) })
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ECROSS_REFERENCES, NewImmutableEList([]any{})) })
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.EXPECT().GetName().Return("name")
		o.EInvokeFromID(EOBJECT__EGET_ESTRUCTURALFEATURE, NewImmutableEList([]any{mockFeature}))
	})
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.EXPECT().GetName().Return("name")
		o.EInvokeFromID(EOBJECT__EGET_ESTRUCTURALFEATURE_EBOOLEAN, NewImmutableEList([]any{mockFeature, true}))
	})
	assert.Panics(t, func() {
		mockOperation := new(MockEOperation)
		mockOperation.EXPECT().GetName().Return("name")
		o.EInvokeFromID(EOBJECT__EINVOKE_EOPERATION_EELIST, NewImmutableEList([]any{mockOperation, NewEmptyBasicEList()}))
	})
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__EIS_PROXY, NewImmutableEList([]any{})) })
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.EXPECT().GetName().Return("name")
		o.EInvokeFromID(EOBJECT__EIS_SET_ESTRUCTURALFEATURE, NewImmutableEList([]any{mockFeature}))
	})
	assert.NotPanics(t, func() { o.EInvokeFromID(EOBJECT__ERESOURCE, NewImmutableEList([]any{})) })
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.EXPECT().GetName().Return("name")
		o.EInvokeFromID(EOBJECT__ESET_ESTRUCTURALFEATURE_EJAVAOBJECT, NewImmutableEList([]any{mockFeature, any(nil)}))
	})
	assert.Panics(t, func() {
		mockFeature := new(MockEStructuralFeature)
		mockFeature.EXPECT().GetName().Return("name")
		o.EInvokeFromID(EOBJECT__EUNSET_ESTRUCTURALFEATURE, NewImmutableEList([]any{mockFeature}))
	})
}
