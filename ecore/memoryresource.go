package ecore

import "io"

type memoryResource interface {
	EResource
}

type memoryResourceImpl struct {
	EResourceImpl
}

func newMemoryResourceImpl() *memoryResourceImpl {
	r := &memoryResourceImpl{}
	r.SetInterfaces(r)
	r.Initialize()
	return r
}

func (r *memoryResourceImpl) Load() {
	panic("operation not supported")
}

func (r *memoryResourceImpl) LoadWithOptions(options map[string]interface{}) {
	panic("operation not supported")
}

func (r *memoryResourceImpl) LoadWithReader(rd io.Reader, options map[string]interface{}) {
	panic("operation not supported")
}

func (r *memoryResourceImpl) Unload() {
	panic("operation not supported")
}

func (r *memoryResourceImpl) IsLoaded() bool {
	return true
}

func (r *memoryResourceImpl) Save() {
	panic("operation not supported")
}

func (r *memoryResourceImpl) SaveWithOptions(options map[string]interface{}) {
	panic("operation not supported")
}

func (r *memoryResourceImpl) SaveWithWriter(w io.Writer, options map[string]interface{}) {
	panic("operation not supported")
}
