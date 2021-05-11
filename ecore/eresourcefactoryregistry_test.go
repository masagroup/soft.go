package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResoureFactoryRegistrySingleton(t *testing.T) {
	r := GetResourceFactoryRegistry()
	require.NotNil(t, r)
	assert.NotNil(t, r.GetExtensionToFactoryMap()["ecore"])
	assert.NotNil(t, r.GetExtensionToFactoryMap()["xml"])
	assert.NotNil(t, r.GetProtocolToFactoryMap()["memory"])
}

func TestResoureFactoryRegistrySingletonGetFactory(t *testing.T) {
	r := GetResourceFactoryRegistry()
	assert.NotNil(t, r.GetFactory(&URI{Path: "*.xml"}))
	assert.NotNil(t, r.GetFactory(&URI{Path: "*.ecore"}))
	assert.NotNil(t, r.GetFactory(&URI{Scheme: "memory", Path: "path"}))
}
