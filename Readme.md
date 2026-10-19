[![](https://img.shields.io/github/license/masagroup/soft.go.svg)](https://github.com/masagroup/soft.go/blob/master/LICENSE)
![](https://github.com/masagroup/soft.go/actions/workflows/build_and_test.yaml/badge.svg)
[![Go](https://img.shields.io/github/go-mod/go-version/masagroup/soft.go)](https://github.com/masagroup/soft.go)
![](https://img.shields.io/github/v/release/masagroup/soft.go)
# Soft Go #

Soft Go is an implementation of the EMF Ecore library in Golang. This library is partially generated and referenced by the [Soft Go Generator](https://github.com/masagroup/soft.gen)
 
Soft Go is part of [Soft](https://github.com/masagroup/soft) project

# Installation #
To install Ecore, use go get:
```shell
go get github.com/masagroup/soft.go
```

This will then make the following packages available to you:
```
github.com/masagroup/soft.go/ecore
```

Import the masagroup/soft.go/ecore package into your code using this template:
```Golang
package yours

import (
  "github.com/masagroup/soft.go/ecore"
)
```

# Several implementations of a metamodel #
A program may link several complete generated implementations of the same metamodel (for instance one with typed fields and one with reflective delegation). They share the nsURI and the classifier and feature IDs, but each has its own `EPackage` instance and its own `EFactory`.

A package registry holds one package per nsURI, so the implementation is chosen per resource set: register the package of the implementation a resource set uses in that resource set's registry, with `ecore.CreateEResourceSet(packages)`, `ecore.XMLProcessorPackages(packages)` or `resourceSet.GetPackageRegistry().RegisterPackage(pack)`. The global registry (`ecore.GetPackageRegistry()`) holds `ecore` only and stays neutral: generated packages never register themselves in it.

```Golang
fieldsSet := ecore.CreateEResourceSet([]ecore.EPackage{fields.GetPackage()})
storeSet := ecore.CreateEResourceSet([]ecore.EPackage{store.GetPackage()})
```

What works within one implementation and across two implementations of the same metamodel:

| Operation | Same implementation | Across (same nsURI, same IDs) |
| --- | --- | --- |
| generated factory `Create` / `CreateFromString` / `ConvertToString` | works | works, by classifier ID; the result belongs to the factory's package |
| dynamic `EFactoryExt` | works | panics (package mismatch), unchanged |
| `EGetFromID` / `ESetFromID`, `GetEStructuralFeature(id)`, `...FromName` | works | works |
| `EGet` / `ESet` / `EIsSet` / `EUnset(feature)` | works | panics "not a valid feature": look the feature up by ID or name on `obj.EClass()` first |
| load into a resource set with one implementation registered | works | works; the whole resource is that implementation |
| save a resource mixing implementations | n/a | unsupported |
| `Copy` | works, same implementation | not a cross copy |
| `Equals`, `IsSuperTypeOf`, `IsInstance`, proxies across resource sets | works | false / unsupported |

## One shared package, several implementations ##
Generated with the property `packageInstance=interfaces` of soft.generator.go, the package of a metamodel (`GetPackage`) is generated once, with its interfaces, and shared by its implementations: each implementation (`generateImplementations`, one run each) holds its classes and a factory of its own (`GetFactory`), whose package is the shared one. The objects of all implementations have the same `EClass`. The factory instance of the shared package only converts data types: its `Create` panics.

The implementation a resource set creates its objects with (load, XML, binary and SQL) is the factory registered for the nsURI in its package registry:

```Golang
typedSet := ecore.CreateEResourceSet([]ecore.EPackage{model.GetPackage()})
typedSet.GetPackageRegistry().PutFactory(model.NS_URI, typed.GetFactory())
reflectiveSet := ecore.CreateEResourceSet([]ecore.EPackage{model.GetPackage()})
reflectiveSet.GetPackageRegistry().PutFactory(model.NS_URI, reflective.GetFactory())
```

An implementation can be made the default implementation of the package with its generated `SetDefault()` (`ecore.SetDefaultFactory(pack, factory)`): the factory instance of the package then creates its objects with it, as do the resource sets without a factory registered for the nsURI.

`ecore.GetObjectFactory(obj)` returns the factory of the implementation of an object: `Copy` and the entries of maps are of the implementation of the object they come from. The objects of different implementations of a shared package do not mix: containing, referencing (lists, single references, `ESet`) or adding as a root of the same resource an object of another implementation panics (`ecore.CheckSameImplementation`) when the program is built with the tag `ecore_check` (`ecore.Checks`). Without it, the checks are compiled out and cost nothing. The same tag enables the other checks of programming errors, such as setting an element of a unique list to an element it holds at another index. The property `checkImplementation=false` of soft.generator.go also leaves them out of the generated setters. Within one implementation, everything works as with a complete package, including `EGet` / `ESet` with the features of `GetPackage()`, `Equals` and `IsSuperTypeOf`.

# Supported go versions #
We support Go v1.12

