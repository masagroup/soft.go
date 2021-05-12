package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryResourceConstructor(t *testing.T) {
	require.NotNil(t, NewMemoryResourceImpl())
}

func TestMemoryLoad(t *testing.T) {
	r := NewMemoryResourceImpl()
	r.SetURI(CreateMemoryURI("memory"))
	r.Load()
	assert.True(t, r.GetErrors().Empty(), diagnosticError(r.GetErrors()))
}

func TestMemorySave(t *testing.T) {
	r := NewMemoryResourceImpl()
	r.SetURI(CreateMemoryURI("memory"))
	r.Save()
	assert.True(t, r.GetErrors().Empty(), diagnosticError(r.GetErrors()))
}
