package ecore

import (
	"testing"
)

func TestMemoryResourceLoad(t *testing.T) {
	r := newMemoryResourceImpl()
	r.Load()
}

func TestMemoryResourceLoadWithOptions(t *testing.T) {
	r := newMemoryResourceImpl()
	r.LoadWithOptions(nil)
}

func TestMemoryResourceLoadWithWriter(t *testing.T) {
	r := newMemoryResourceImpl()
	r.LoadWithReader(nil, nil)
}

func TestMemoryResourceSave(t *testing.T) {
	r := newMemoryResourceImpl()
	r.Save()
}

func TestMemoryResourceSaveWithOptions(t *testing.T) {
	r := newMemoryResourceImpl()
	r.SaveWithOptions(nil)
}

func TestMemoryResourceSaveWithWriter(t *testing.T) {
	r := newMemoryResourceImpl()
	r.SaveWithWriter(nil, nil)
}
