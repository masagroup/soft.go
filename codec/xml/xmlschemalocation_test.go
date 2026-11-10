package xml

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
)

// newUnqualifiedDocumentRootModel is the document root model of a schema without target namespace:
// its package is not qualified, and documents without namespace find it under "".
func newUnqualifiedDocumentRootModel() (*documentRootModel, ecore.EResourceSet) {
	m := newDocumentRootModel()
	eAnnotation := ecore.GetFactory().CreateEAnnotation()
	eAnnotation.SetSource(annotationURI)
	eAnnotation.GetDetails().Put("qualified", "false")
	m.ePackage.GetEAnnotations().Add(eAnnotation)
	rs := m.newResourceSet()
	rs.GetPackageRegistry().PutPackage("", m.ePackage)
	return m, rs
}

func TestExtendedMetaData_IsQualified(t *testing.T) {
	emd := NewExtendedMetaData()
	assert.True(t, emd.IsQualified(newLibraryModel().ePackage))
	m, _ := newUnqualifiedDocumentRootModel()
	assert.False(t, emd.IsQualified(m.ePackage))
}

// A document without namespace keeps its xsi:noNamespaceSchemaLocation, recorded in its document
// root under "", and is written back without namespace.
func TestXMLDocumentRootNoNamespaceSchemaLocation(t *testing.T) {
	m, rs := newUnqualifiedDocumentRootModel()
	options := map[string]any{OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}
	document := `<?xml version="1.0" encoding="UTF-8"?>
<library xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="schemas/1.0/lib.xsd" name="l">
  <books name="b"/>
</library>
`
	eResource := rs.CreateResource(ecore.NewURI("library.xml"))
	root := decodeResource(t, eResource, document, options)
	assert.Equal(t, "schemas/1.0/lib.xsd", root.EGet(m.schemaLocation).(ecore.EMap).GetValue(""))
	assert.Equal(t, document, encodeResource(t, eResource, options))
}

// A document of a package that is not qualified, and that names no schema without namespace, keeps
// its namespace: the package is written qualified, as before.
func TestXMLDocumentRootUnqualifiedNamespaced(t *testing.T) {
	_, rs := newUnqualifiedDocumentRootModel()
	options := map[string]any{OPTION_EXTENDED_META_DATA: NewExtendedMetaData()}
	document := `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://lib schemas/1.0/lib.xsd" name="l"/>
`
	eResource := rs.CreateResource(ecore.NewURI("library.xml"))
	decodeResource(t, eResource, document, options)
	assert.Equal(t, document, encodeResource(t, eResource, options))
}

// OPTION_SCHEMA_LOCATION names the schemas of a document a save writes, over the ones its document
// root records, and leaves the document root as it is.
func TestXMLSchemaLocationOption(t *testing.T) {
	m := newDocumentRootModel()
	emd := NewExtendedMetaData()
	document := `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://lib schemas/1.0/lib.xsd" name="l"/>
`
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("library.xml"))
	root := decodeResource(t, eResource, document, map[string]any{OPTION_EXTENDED_META_DATA: emd})
	encode := func(locations map[string]string) string {
		return encodeResource(t, eResource, map[string]any{OPTION_EXTENDED_META_DATA: emd, OPTION_SCHEMA_LOCATION: locations})
	}
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://lib schemas/2.0/lib.xsd http://other other.xsd" name="l"/>
`, encode(map[string]string{"http://other": "other.xsd", "http://lib": "schemas/2.0/lib.xsd"}))
	// an empty option writes no schema; xmlns:xsi is still declared by the prefix map of the document
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="l"/>
`, encode(map[string]string{}))
	assert.Equal(t, "schemas/1.0/lib.xsd", root.EGet(m.schemaLocation).(ecore.EMap).GetValue("http://lib"))
}

// OPTION_SCHEMA_LOCATION with a schema without namespace writes a document without namespace from a
// document with one: its xsi:noNamespaceSchemaLocation, or none for an empty location.
func TestXMLSchemaLocationOptionNoNamespace(t *testing.T) {
	_, rs := newUnqualifiedDocumentRootModel()
	emd := NewExtendedMetaData()
	eResource := rs.CreateResource(ecore.NewURI("library.xml"))
	decodeResource(t, eResource, `<lib:library xmlns:lib="http://lib" name="l"><books name="b"/></lib:library>`, map[string]any{OPTION_EXTENDED_META_DATA: emd})
	encode := func(locations map[string]string) string {
		return encodeResource(t, eResource, map[string]any{OPTION_EXTENDED_META_DATA: emd, OPTION_SCHEMA_LOCATION: locations})
	}
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<library xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="lib.xsd" name="l">
  <books name="b"/>
</library>
`, encode(map[string]string{"": "lib.xsd"}))
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<library name="l">
  <books name="b"/>
</library>
`, encode(map[string]string{"": ""}))
}

// putInResource puts the package of a model in a resource at uri, as a package loaded from a file.
func putInResource(ePackage ecore.EPackage, uri string) {
	eResource := ecore.NewEResourceImpl()
	eResource.SetURI(ecore.NewURI(uri))
	eResource.GetContents().Add(ePackage)
}

// OPTION_SCHEMA_LOCATION true writes the location of the packages of the document, as EMF: the uri of
// the resource of a package, relative to the document, and none for a package without resource or
// registered at its namespace.
func TestXMLSchemaLocationOptionPackages(t *testing.T) {
	m := newDocumentRootModel()
	emd := NewExtendedMetaData()
	eResource := m.newResourceSet().CreateResource(ecore.NewURI("file:///data/library.xml"))
	decodeResource(t, eResource, `<lib:library xmlns:lib="http://lib" name="l"/>`, map[string]any{OPTION_EXTENDED_META_DATA: emd})
	options := map[string]any{OPTION_EXTENDED_META_DATA: emd, OPTION_SCHEMA_LOCATION: true}

	// a package without resource
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" name="l"/>
`, encodeResource(t, eResource, options))

	// a package loaded from a file
	putInResource(m.ePackage, "file:///data/schemas/lib.ecore")
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://lib schemas/lib.ecore" name="l"/>
`, encodeResource(t, eResource, options))
	// false writes none
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" name="l"/>
`, encodeResource(t, eResource, map[string]any{OPTION_EXTENDED_META_DATA: emd, OPTION_SCHEMA_LOCATION: false}))

	// a package registered at its namespace, as a generated package is
	putInResource(m.ePackage, "http://lib")
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<lib:library xmlns:lib="http://lib" name="l"/>
`, encodeResource(t, eResource, options))
}

// OPTION_SCHEMA_LOCATION true keeps the schemas a document names, and gives a location to the schema
// without namespace it names without one.
func TestXMLSchemaLocationOptionPackagesNoNamespace(t *testing.T) {
	m, rs := newUnqualifiedDocumentRootModel()
	emd := NewExtendedMetaData()
	putInResource(m.ePackage, "file:///data/schemas/lib.ecore")
	eResource := rs.CreateResource(ecore.NewURI("file:///data/library.xml"))
	root := decodeResource(t, eResource, `<library xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="lib.xsd" name="l"/>`, map[string]any{OPTION_EXTENDED_META_DATA: emd})
	options := map[string]any{OPTION_EXTENDED_META_DATA: emd, OPTION_SCHEMA_LOCATION: true}
	kept := `<?xml version="1.0" encoding="UTF-8"?>
<library xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="lib.xsd" name="l"/>
`
	assert.Equal(t, kept, encodeResource(t, eResource, options))

	root.EGet(m.schemaLocation).(ecore.EMap).Put("", "")
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<library xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="schemas/lib.ecore" name="l"/>
`, encodeResource(t, eResource, options))
}
