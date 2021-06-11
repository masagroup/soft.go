package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResoureCodecRegistrySingleton(t *testing.T) {
	r := GetResourceCodecRegistry()
	require.NotNil(t, r)
	assert.NotNil(t, r.GetExtensionToCodecMap()["ecore"])
	assert.NotNil(t, r.GetExtensionToCodecMap()["xml"])
}

func TestResoureCodecRegistrySingletonGetCodec(t *testing.T) {
	r := GetResourceCodecRegistry()
	assert.NotNil(t, r.GetCodec(&URI{Path: "*.xml"}))
	assert.NotNil(t, r.GetCodec(&URI{Path: "*.ecore"}))
}
