package ecore

type EPackageRegistryImpl struct {
	packages map[string]interface{}
	delegate EPackageRegistry
}

func NewEPackageRegistryImpl() *EPackageRegistryImpl {
	r := &EPackageRegistryImpl{
		packages: map[string]interface{}{},
	}
	return r
}

func NewEPackageRegistryImplWithDelegate(delegate EPackageRegistry) *EPackageRegistryImpl {
	r := &EPackageRegistryImpl{
		packages: map[string]interface{}{},
		delegate: delegate,
	}
	return r
}

func (r *EPackageRegistryImpl) RegisterPackage(pack EPackage) {
	r.packages[pack.GetNsURI()] = pack
}

func (r *EPackageRegistryImpl) UnregisterPackage(pack EPackage) {
	delete(r.packages, pack.GetNsURI())
}

func (r *EPackageRegistryImpl) PutPackage(nsURI string, pack EPackage) {
	r.packages[nsURI] = pack
}

func (r *EPackageRegistryImpl) PutSupplier(nsURI string, supplier func() EPackage) {
	r.packages[nsURI] = supplier
}

func (r *EPackageRegistryImpl) Remove(nsURI string) {
	delete(r.packages, nsURI)
}

func (r *EPackageRegistryImpl) doGetPackage(nsURI string) EPackage {
	p, _ := r.packages[nsURI]
	if p != nil {
		if pack, _ := p.(EPackage); pack != nil {
			return pack
		} else if f, _ := p.(func() EPackage); f != nil {
			return f()
		}
	}
	return nil
}

func (r *EPackageRegistryImpl) GetPackage(nsURI string) EPackage {
	p := r.doGetPackage(nsURI)
	if p != nil {
		return p
	} else {
		if r.delegate != nil {
			return r.delegate.GetPackage(nsURI)
		}
	}
	return nil
}

func (r *EPackageRegistryImpl) GetFactory(nsURI string) EFactory {
	p := r.doGetPackage(nsURI)
	if p != nil {
		return p.GetEFactoryInstance()
	} else {
		if r.delegate != nil {
			return r.delegate.GetFactory(nsURI)
		}
	}
	return nil
}
