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

func m(a, b any) []any {
	return []any{a, b}
}

func TestXmlNamespacesNoContext(t *testing.T) {
	n := newXMLNamespaces()
	assert.Equal(t, m("", false), m(n.getURI("prefix")))
	assert.Equal(t, m("", false), m(n.getPrefix("uri")))
}

func TestXmlNamespacesEmpty(t *testing.T) {
	n := newXMLNamespaces()
	n.pushContext()
	assert.Equal(t, m("", false), m(n.getURI("prefix")))
	assert.Equal(t, m("", false), m(n.getPrefix("uri")))
	c := n.popContext()
	assert.Equal(t, 0, len(c))
}

func TestXmlNamespacesContext(t *testing.T) {
	n := newXMLNamespaces()
	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri"))
	assert.Equal(t, m("uri", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri")))

	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri2"))
	assert.Equal(t, m("uri2", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri2")))

	n.popContext()
	assert.Equal(t, m("uri", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri")))

	n.popContext()
	assert.Equal(t, m("", false), m(n.getURI("prefix")))
	assert.Equal(t, m("", false), m(n.getPrefix("uri")))
}

func TestXmlNamespacesContextRemap(t *testing.T) {
	n := newXMLNamespaces()
	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri"))
	assert.Equal(t, m("uri", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri")))

	assert.True(t, n.declarePrefix("prefix", "uri2"))
	assert.Equal(t, m("uri2", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri2")))
}

func TestXmlNamespacesContextNoRemap(t *testing.T) {
	n := newXMLNamespaces()
	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri"))
	assert.Equal(t, m("uri", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri")))

	n.pushContext()
	assert.False(t, n.declarePrefix("prefix", "uri2"))
	assert.Equal(t, m("uri2", true), m(n.getURI("prefix")))
	assert.Equal(t, m("prefix", true), m(n.getPrefix("uri2")))
}
