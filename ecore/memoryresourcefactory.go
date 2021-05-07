package ecore

type MemoryResourceFactory struct {
}

func (f *MemoryResourceFactory) CreateResource(uri *URI) EResource {
	r := newMemoryResourceImpl()
	r.SetURI(uri)
	return r
}
