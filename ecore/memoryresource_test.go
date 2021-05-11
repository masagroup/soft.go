package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryResourceIsLoaded(t *testing.T) {
	r := newMemoryResourceImpl()
	assert.True(t, r.IsLoaded())
}

func TestMemoryResourceLoad(t *testing.T) {
	r := newMemoryResourceImpl()
	assert.Panics(t, func() { r.Load() })
}

func TestMemoryResourceLoadWithOptions(t *testing.T) {
	r := newMemoryResourceImpl()
	assert.Panics(t, func() { r.LoadWithOptions(nil) })
}

func TestMemoryResourceLoadWithWriter(t *testing.T) {
	r := newMemoryResourceImpl()
	assert.Panics(t, func() { r.LoadWithReader(nil, nil) })
}

func TestMemoryResourceSave(t *testing.T) {
	r := newMemoryResourceImpl()
	assert.Panics(t, func() { r.Save() })
}

func TestMemoryResourceSaveWithOptions(t *testing.T) {
	r := newMemoryResourceImpl()
	assert.Panics(t, func() { r.SaveWithOptions(nil) })
}

func TestMemoryResourceSaveWithWriter(t *testing.T) {
	r := newMemoryResourceImpl()
	assert.Panics(t, func() { r.SaveWithWriter(nil, nil) })
}
