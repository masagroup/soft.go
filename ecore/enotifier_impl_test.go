package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockENotifierImpl_Adapters(t *testing.T) {
	n := &ENotifierImpl{}
	n.SetInterfaces(n)
	assert.Nil(t, n.EBasicAdapters())
	assert.False(t, n.EBasicHasAdapters())

	assert.NotNil(t, n.EAdapters())
	assert.NotNil(t, n.EBasicAdapters())
	assert.False(t, n.EBasicHasAdapters())
}

func TestMockENotifierImpl_Deliver(t *testing.T) {
	n := &ENotifierImpl{}

	assert.False(t, n.EDeliver())

	n.ESetDeliver(true)
	assert.True(t, n.EDeliver())
}
