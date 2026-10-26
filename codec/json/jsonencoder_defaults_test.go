package json

import (
	"bytes"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// an unset feature with a default value is saved only with the keep defaults option
func TestJSONEncoder_KeepDefaults(t *testing.T) {
	f := ecore.GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetName("lib")
	ePackage.SetNsURI("http://lib")
	ePackage.SetNsPrefix("lib")
	ePackage.SetEFactoryInstance(f.CreateEFactory())
	eBook := f.CreateEClass()
	eBook.SetName("Book")
	ePackage.GetEClassifiers().Add(eBook)
	ePages := f.CreateEAttribute()
	ePages.SetName("pages")
	ePages.SetEType(ecore.GetPackage().GetEInt())
	ePages.SetDefaultValueLiteral("100")
	eBook.GetEStructuralFeatures().Add(ePages)
	eTitle := f.CreateEAttribute()
	eTitle.SetName("title")
	eTitle.SetEType(ecore.GetPackage().GetEString())
	eBook.GetEStructuralFeatures().Add(eTitle)

	eResource := ecore.NewEResourceImpl()
	eResource.SetURI(ecore.NewURI("file:///library.json"))
	eResource.GetContents().Add(ePackage.GetEFactoryInstance().Create(eBook))

	encode := func(options map[string]any) string {
		var buffer bytes.Buffer
		NewEncoder(eResource, &buffer, options).EncodeResource()
		require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
		return buffer.String()
	}
	assert.NotContains(t, encode(nil), `"pages"`)
	assert.NotContains(t, encode(map[string]any{OPTION_KEEP_DEFAULTS: false}), `"pages"`)
	document := encode(map[string]any{OPTION_KEEP_DEFAULTS: true})
	assert.Contains(t, document, `"pages":"100"`)
	// a feature without a default value is still not saved
	assert.NotContains(t, document, `"title"`)
}
