// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "fmt"

// EObjectFactory is implemented by the objects of an implementation of a shared package: a
// package generated once, with the interfaces of its model (soft.generator.go property
// packageInstance=interfaces), and implemented by several implementations, each with a factory of
// its own whose package is the shared one (EFactoryExt.SetSharedEPackage). The objects of all
// implementations have the same EClass; EFactory tells their implementation apart.
type EObjectFactory interface {
	// EFactory returns the factory of the implementation of the object.
	EFactory() EFactory
}

// GetObjectFactory returns the factory that creates objects of the implementation of eObject:
// its own factory when it is an object of an implementation of a shared package (EObjectFactory),
// the factory instance of the package of its class otherwise.
func GetObjectFactory(eObject EObject) EFactory {
	if o, _ := eObject.(EObjectFactory); o != nil {
		return o.EFactory()
	}
	if eClass := eObject.EClass(); eClass != nil {
		if ePackage := eClass.GetEPackage(); ePackage != nil {
			return ePackage.GetEFactoryInstance()
		}
	}
	return nil
}

// getClassFactory returns the factory creating the objects of eClass for owner: the factory of
// the implementation of owner when eClass is a class of its package (an object of an implementation
// of a shared package), the factory instance of the package of eClass otherwise.
func getClassFactory(owner EObject, eClass EClass) EFactory {
	ePackage := eClass.GetEPackage()
	if o, _ := owner.(EObjectFactory); o != nil {
		if factory := o.EFactory(); factory.GetEPackage() == ePackage {
			return factory
		}
	}
	return ePackage.GetEFactoryInstance()
}

// getRegistryFactory returns the factory creating the objects of the classes of ePackage in the
// resource set of packageRegistry: the factory registered for its nsURI (PutFactory), else the
// factory instance of ePackage.
func getRegistryFactory(packageRegistry EPackageRegistry, ePackage EPackage) EFactory {
	if packageRegistry != nil {
		if factory := packageRegistry.GetFactory(ePackage.GetNsURI()); factory != nil {
			return factory
		}
	}
	return ePackage.GetEFactoryInstance()
}

// EDefaultFactoryHolder is implemented by the factory instance of a shared package: it converts
// the data types of the package and creates its objects with its default factory, the factory of
// the default implementation of the package.
type EDefaultFactoryHolder interface {
	SetDefaultFactory(defaultFactory EFactory)
	GetDefaultFactory() EFactory
}

// SetDefaultFactory makes factory, the factory of an implementation of the shared package ePackage,
// the default implementation of ePackage: the factory instance of ePackage creates its objects with
// it, as do the resource sets without a factory registered for its nsURI (PutFactory). A nil factory
// removes the default implementation.
func SetDefaultFactory(ePackage EPackage, factory EFactory) {
	holder, _ := ePackage.GetEFactoryInstance().(EDefaultFactoryHolder)
	if holder == nil {
		panic(fmt.Sprintf("the package '%s' is not a shared package", ePackage.GetNsURI()))
	}
	if factory != nil && factory.GetEPackage() != ePackage {
		panic(fmt.Sprintf("the factory is not a factory of the package '%s'", ePackage.GetNsURI()))
	}
	holder.SetDefaultFactory(factory)
}
