package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResoureDriverRegistrySingleton(t *testing.T) {
	r := GetResourceDriverRegistry()
	require.NotNil(t, r)
	assert.NotNil(t, r.GetExtensionToDriverMap()["ecore"])
	assert.NotNil(t, r.GetExtensionToDriverMap()["xml"])
	assert.NotNil(t, r.GetProtocolToDriverMap()["memory"])
}

func TestResoureDriverRegistrySingletonGetDriver(t *testing.T) {
	r := GetResourceDriverRegistry()
	assert.NotNil(t, r.GetDriver(&URI{Path: "*.xml"}))
	assert.NotNil(t, r.GetDriver(&URI{Path: "*.ecore"}))
	assert.NotNil(t, r.GetDriver(&URI{Scheme: "memory", Path: "path"}))
}
