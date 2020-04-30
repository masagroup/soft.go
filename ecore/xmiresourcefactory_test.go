package ecore

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMIResourceFactoryCreate(t *testing.T) {
	f := &XMIResourceFactory{}
	uri := &url.URL{}
	r := f.CreateResource(uri)
	assert.NotNil(t, r)
	assert.Equal(t, uri, r.GetURI())
}
