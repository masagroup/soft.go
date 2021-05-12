package ecore

type MemoryResource interface {
	EResource
}

type MemoryResourceImpl struct {
	EResourceImpl
}

func NewMemoryResourceImpl() *MemoryResourceImpl {
	r := &MemoryResourceImpl{}
	r.SetInterfaces(r)
	r.Initialize()
	return r
}
