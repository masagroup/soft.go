package ecore

type MemoryResourceFactory struct {
}

func (f *MemoryResourceFactory) CreateResource(uri *URI) EResource {
	r := NewMemoryResourceImpl()
	r.SetURI(uri)
	return r
}
