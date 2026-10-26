package xml

import (
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
)

// with a line width, an attribute starts a new line once the line is longer, as EMF wraps them:
// indented by two levels from its element, the namespaces of the root element included, which
// the width of the other attributes does not count, as they are added at the end
func TestXMLEncoder_LineWidth(t *testing.T) {
	m := newLibraryModel()
	library := m.create(m.library, "library")
	library.EGet(m.books).(ecore.EList).Add(m.create(m.book, "a book whose name is long enough"))
	library.EGet(m.books).(ecore.EList).Get(0).(ecore.EObject).ESet(m.pages, 10)
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.GetContents().Add(library)

	for _, options := range []map[string]any{nil, {OPTION_LINE_WIDTH: 0}, {OPTION_LINE_WIDTH: -1}} {
		document := encodeResource(t, eResource, options)
		assert.Equal(t, 4, strings.Count(document, "\n"), options)
	}

	document := encodeResource(t, eResource, map[string]any{OPTION_LINE_WIDTH: 10})
	assert.Equal(t, strings.Join([]string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<lib:Library`,
		`    xmlns:lib="http://lib"`,
		`    name="library">`,
		`  <books name="a book whose name is long enough"`,
		`      pages="10"/>`,
		`</lib:Library>`,
		``,
	}, "\n"), strings.ReplaceAll(document, "\r\n", "\n"))
}
