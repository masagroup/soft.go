// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EPackageRegistryImpl struct {
	packages  map[string]any
	factories map[string]EFactory
	delegate  EPackageRegistry
}

func NewEPackageRegistryImpl() *EPackageRegistryImpl {
	r := &EPackageRegistryImpl{
		packages: map[string]any{},
	}
	return r
}

func NewEPackageRegistryImplWithDelegate(delegate EPackageRegistry) *EPackageRegistryImpl {
	r := &EPackageRegistryImpl{
		packages: map[string]any{},
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
	if p := r.packages[nsURI]; p != nil {
		if pack, _ := p.(EPackage); pack != nil {
			return pack
		} else if f, _ := p.(func() EPackage); f != nil {
			return f()
		}
	}
	return nil
}

func (r *EPackageRegistryImpl) GetPackage(nsURI string) EPackage {
	if p := r.doGetPackage(nsURI); p != nil {
		return p
	} else if r.delegate != nil {
		return r.delegate.GetPackage(nsURI)
	}
	return nil
}

func (r *EPackageRegistryImpl) PutFactory(nsURI string, factory EFactory) {
	if r.factories == nil {
		r.factories = map[string]EFactory{}
	}
	r.factories[nsURI] = factory
}

func (r *EPackageRegistryImpl) RemoveFactory(nsURI string) {
	delete(r.factories, nsURI)
}

func (r *EPackageRegistryImpl) RegisterFactory(factory EFactory) {
	r.PutFactory(getFactoryNsURI(factory), factory)
}

func (r *EPackageRegistryImpl) UnregisterFactory(factory EFactory) {
	r.RemoveFactory(getFactoryNsURI(factory))
}

// getFactoryNsURI returns the nsURI of the package of factory: a factory registered by its package
// has one.
func getFactoryNsURI(factory EFactory) string {
	pack := factory.GetEPackage()
	if pack == nil {
		panic("the factory has no package")
	}
	return pack.GetNsURI()
}

// GetFactory returns the factory registered for nsURI (PutFactory), else the factory instance of
// the package of nsURI, else the factory of the delegate registry.
func (r *EPackageRegistryImpl) GetFactory(nsURI string) EFactory {
	if f := r.factories[nsURI]; f != nil {
		return f
	} else if p := r.doGetPackage(nsURI); p != nil {
		return p.GetEFactoryInstance()
	} else if r.delegate != nil {
		return r.delegate.GetFactory(nsURI)
	}
	return nil
}
