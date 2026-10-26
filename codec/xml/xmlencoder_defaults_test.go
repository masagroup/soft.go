package xml

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
)

// an unset feature with a default value is saved only with the keep defaults option
func TestXMLEncoder_KeepDefaults(t *testing.T) {
	m := newLibraryModel()
	m.pages.SetDefaultValueLiteral("100")
	library := m.create(m.library, "l")
	library.EGet(m.books).(ecore.EList).Add(m.create(m.book, "b"))
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///library.xml"))
	eResource.GetContents().Add(library)

	assert.NotContains(t, encodeResource(t, eResource, nil), `pages=`)
	assert.NotContains(t, encodeResource(t, eResource, map[string]any{OPTION_KEEP_DEFAULTS: false}), `pages=`)
	assert.Contains(t, encodeResource(t, eResource, map[string]any{OPTION_KEEP_DEFAULTS: true}), `pages="100"`)
	// a feature without a default value is still not saved
	assert.NotContains(t, encodeResource(t, eResource, map[string]any{OPTION_KEEP_DEFAULTS: true}), `nick=`)
}
