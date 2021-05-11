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
}

func (r *memoryResourceImpl) LoadWithOptions(options map[string]interface{}) {
}

func (r *memoryResourceImpl) LoadWithReader(rd io.Reader, options map[string]interface{}) {
}

func (r *memoryResourceImpl) Unload() {
}

func (r *memoryResourceImpl) Save() {
}

func (r *memoryResourceImpl) SaveWithOptions(options map[string]interface{}) {
}

func (r *memoryResourceImpl) SaveWithWriter(w io.Writer, options map[string]interface{}) {
}
