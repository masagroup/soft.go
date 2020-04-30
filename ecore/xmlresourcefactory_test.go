package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLResourceFactoryCreate(t *testing.T) {
	f := &XMLResourceFactory{}
	uri := &url.URL{}
	r := f.CreateResource(uri)
	assert.NotNil(t, r)
	assert.Equal(t, uri, r.GetURI())
}
