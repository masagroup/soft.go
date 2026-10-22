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

func TestEEnumLiteralExt_Literal(t *testing.T) {
	l := newEEnumLiteralExt()
	l.SetName("name")

	// unset : the name
	assert.Equal(t, "name", l.GetLiteral())
	assert.False(t, l.EIsSetFromID(EENUM_LITERAL__LITERAL))

	// set to "" : stays ""
	l.SetLiteral("")
	assert.Equal(t, "", l.GetLiteral())
	assert.True(t, l.EIsSetFromID(EENUM_LITERAL__LITERAL))
	assert.Equal(t, "", l.EGetFromID(EENUM_LITERAL__LITERAL, false))

	l.SetLiteral("literal")
	assert.Equal(t, "literal", l.GetLiteral())
	assert.True(t, l.EIsSetFromID(EENUM_LITERAL__LITERAL))

	// unset again
	l.EUnsetFromID(EENUM_LITERAL__LITERAL)
	assert.Equal(t, "name", l.GetLiteral())
	assert.False(t, l.EIsSetFromID(EENUM_LITERAL__LITERAL))

	// reflective set
	l.ESetFromID(EENUM_LITERAL__LITERAL, "")
	assert.Equal(t, "", l.GetLiteral())
	assert.True(t, l.EIsSetFromID(EENUM_LITERAL__LITERAL))

	// the other features
	assert.False(t, l.EIsSetFromID(EENUM_LITERAL__VALUE))
	l.SetValue(1)
	assert.True(t, l.EIsSetFromID(EENUM_LITERAL__VALUE))
}

func TestEEnumExt_GetEEnumLiteralByLiteral_Empty(t *testing.T) {
	e := newEEnumExt()
	empty := newEEnumLiteralExt()
	empty.SetName("_")
	empty.SetLiteral("")
	unset := newEEnumLiteralExt()
	unset.SetName("Unset")
	e.GetELiterals().AddAll(NewImmutableEList([]any{empty, unset}))

	assert.Equal(t, empty, e.GetEEnumLiteralByLiteral(""))
	assert.Nil(t, e.GetEEnumLiteralByLiteral("_"))
	assert.Equal(t, unset, e.GetEEnumLiteralByLiteral("Unset"))
}
