package ecore

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
