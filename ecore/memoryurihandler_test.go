package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryURIHandler_CanHandle(t *testing.T) {
	m := &MemoryURIHandler{}
	assert.True(t, m.CanHandle(&URI{scheme: "memory"}))
	assert.False(t, m.CanHandle(&URI{scheme: "file"}))
}

func TestMemoryURIHandler_CreateReader(t *testing.T) {
	m := &MemoryURIHandler{}
	r, err := m.CreateReader(nil)
	assert.Nil(t, r)
	assert.Nil(t, err)
}

func TestMemoryURIHandler_CreateWriter(t *testing.T) {
	m := &MemoryURIHandler{}
	r, err := m.CreateWriter(nil)
	assert.Nil(t, r)
	assert.Nil(t, err)
}
