package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryResourceFactoryCreate(t *testing.T) {
	f := &MemoryResourceFactory{}
	uri := &URI{}
	r := f.CreateResource(uri)
	require.NotNil(t, r)
	assert.Equal(t, uri, r.GetURI())
}
