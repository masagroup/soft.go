// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"iter"
	"strings"
)

func GetEObjectID(eObject EObject) string {
	eClass := eObject.EClass()
	eIDAttribute := eClass.GetEIDAttribute()
	if eIDAttribute == nil || !eObject.EIsSet(eIDAttribute) {
		return ""
	} else {
		return ConvertToString(eIDAttribute.GetEAttributeType(), eObject.EGet(eIDAttribute))
	}
}

func SetEObjectID(eObject EObject, id string) {
	eClass := eObject.EClass()
	eIDAttribute := eClass.GetEIDAttribute()
	if eIDAttribute == nil {
		panic("The object doesn't have an ID feature.")
	} else if len(id) == 0 {
		eObject.EUnset(eIDAttribute)
	} else {
		eObject.ESet(eIDAttribute, CreateFromString(eIDAttribute.GetEAttributeType(), id))
	}
}

func ConvertToString(eDataType EDataType, value any) string {
	eFactory := eDataType.GetEPackage().GetEFactoryInstance()
	return eFactory.ConvertToString(eDataType, value)
}

func CreateFromString(eDataType EDataType, literal string) any {
	eFactory := eDataType.GetEPackage().GetEFactoryInstance()
	return eFactory.CreateFromString(eDataType, literal)
}

func GetURI(eObject EObject) *URI {
	if eObject.EIsProxy() {
		return eObject.(EObjectInternal).EProxyURI()
	} else {
		resource := eObject.EResource()
		if resource != nil {
			return NewURIBuilder(resource.GetURI()).SetFragment(resource.GetURIFragment(eObject)).URI()
		} else {
			id := GetEObjectID(eObject)
			if len(id) == 0 {
				return NewURIBuilder(nil).SetFragment("//" + getRelativeURIFragmentPath(nil, eObject, false)).URI()
			} else {
				return NewURIBuilder(nil).SetFragment(id).URI()
			}
		}
	}
}

func getRelativeURIFragmentPath(ancestor EObject, descendant EObject, _ bool) string {
	if ancestor == descendant {
		return ""
	}
	eObject := descendant
	eContainer := eObject.EContainer()
	visited := make(map[EObject]struct{})
	fragmentPath := []string{}
	for {
		if eContainer == nil {
			break
		}
		if _, isVisited := visited[eObject]; isVisited {
			break
		}
		visited[eObject] = struct{}{}
		fragmentPath = append([]string{eContainer.(EObjectInternal).EURIFragmentSegment(eObject.EContainingFeature(), eObject)}, fragmentPath...)
		eObject = eContainer
		if eContainer == ancestor {
			break
		}
		eContainer = eObject.EContainer()
	}
	if eObject != ancestor && ancestor != nil {
		panic("The ancestor not found")
	}

	return strings.Join(fragmentPath, "/")
}

func GetEObject(rootEObject EObject, relativeFragmentPath string) EObject {
	segments := strings.Split(relativeFragmentPath, "/")
	eObject := rootEObject.(EObjectInternal)
	for i := 0; i < len(segments) && eObject != nil; i++ {
		eObject = eObject.EObjectForFragmentSegment(segments[i]).(EObjectInternal)
	}
	return eObject
}

func ResolveInObject(proxy EObject, context EObject) EObject {
	var resource EResource
	if context != nil {
		resource = context.EResource()
	}
	if resource != nil {
		return ResolveInResourceSet(proxy, resource.GetResourceSet())
	} else {
		return ResolveInResourceSet(proxy, nil)
	}

}

func ResolveInResource(proxy EObject, resource EResource) EObject {
	if resource != nil {
		return ResolveInResourceSet(proxy, resource.GetResourceSet())
	} else {
		return ResolveInResourceSet(proxy, nil)
	}
}

func ResolveInResourceSet(proxy EObject, resourceSet EResourceSet) EObject {
	if proxyInternal, _ := proxy.(EObjectInternal); proxyInternal != nil && proxyInternal.EProxyURI() != nil {
		proxyURI := proxyInternal.EProxyURI()
		var resolved EObject
		if resourceSet != nil {
			resolved = resourceSet.GetEObject(proxyURI, true)
		} else {
			trim := proxyURI.TrimFragment()
			ePackage := GetPackageRegistry().GetPackage(trim.String())
			if ePackage != nil {
				eResource := ePackage.EResource()
				if eResource != nil {
					resolved = eResource.GetEObject(proxyURI.fragment)
				}
			}
		}
		if resolved != nil && resolved != proxy {
			return ResolveInResourceSet(resolved, resourceSet)
		}
	}
	return proxy
}

func Copy(eObject EObject) EObject {
	dC := newDeepCopy(true, true)
	c := dC.copy(eObject)
	dC.copyReferences()
	return c
}

func CopyAll(l EList) EList {
	dC := newDeepCopy(true, true)
	c := dC.copyAll(l)
	dC.copyReferences()
	return c
}

func Equals(eObj1 EObject, eObj2 EObject) bool {
	dE := newDeepEqual()
	return dE.equals(eObj1, eObj2)
}

func EqualsAll(l1 EList, l2 EList) bool {
	dE := newDeepEqual()
	return dE.equalsObjectList(l1, l2)
}

// Removes the object from its containing resource and/or its containing object.
func Remove(eObject EObject) {
	if eObjectInternal, _ := eObject.(EObjectInternal); eObjectInternal != nil {
		if eContainer := eObjectInternal.EInternalContainer(); eContainer != nil {
			if eFeature := eObject.EContainmentFeature(); eFeature != nil {
				if eFeature.IsMany() {
					l := eContainer.EGet(eFeature).(EList)
					l.Remove(eObject)
				} else {
					eContainer.EUnset(eFeature)
				}
			}
		}
		if eResource := eObjectInternal.EInternalResource(); eResource != nil {
			eResource.GetContents().Remove(eObject)
		}
	}
}

/**
 * Deletes the object from its containing resource and/or its containing object
 * as well as from any other feature that references it within the enclosing resource set,
 * resource, or root object.
 */
func Delete(eObject EObject) {
	rootEObject := GetRootContainer(eObject)
	resource := rootEObject.EResource()
	var usages iter.Seq2[EObject, EStructuralFeature]
	if resource == nil {
		usages = findUsages(eObject, rootEObject)
	} else {
		resourceSet := resource.GetResourceSet()
		if resourceSet == nil {
			usages = findUsages(eObject, resource)
		} else {
			usages = findUsages(eObject, resourceSet)
		}
	}
	for object, feature := range usages {
		if feature.IsChangeable() {
			if feature.IsMany() {
				l := object.EGet(feature).(EList)
				l.Remove(eObject)
			} else {
				object.EUnset(feature)
			}
		}
	}
	Remove(eObject)
}

func findUsages(eObjectOfInterest EObject, object any) iter.Seq2[EObject, EStructuralFeature] {
	return func(yield func(EObject, EStructuralFeature) bool) {
		iterator := &eAllContentIterator{
			object: object,
			root:   false,
			getChildren: func(o any) EIterator {
				switch t := o.(type) {
				case EObject:
					return t.EContents().Iterator()
				case EResource:
					return t.GetContents().Iterator()
				case EResourceSet:
					return t.GetResources().Iterator()
				default:
					return nil
				}
			}}
		for iterator.HasNext() {
			eObject, _ := iterator.Next().(EObject)
			if eObject != nil {

			}
		}
	}
}

/**
 * Deletes the object from its containing resource and/or its containing object
 * as well as from any other feature that references it
 * within the enclosing resource set, resource, or root object.
 * If recursive true, contained children of the object that are in the same resource
 * are similarly removed from any features that reference them.
 */
func DeleteRecursive(eObject EObject, recursive bool) {
}

/**
 * Returns the root container
 * it may be this object itself and it will have a nil container
 */
func GetRootContainer(eObject EObject) EObject {
	eCurrent := eObject
	if eCurrent != nil {
		parent := eCurrent.EContainer()
		for parent != nil {
			eCurrent = parent
			parent = eCurrent.EContainer()
		}
	}
	return eCurrent
}

func GetAncestor(eObject EObject, eClass EClass) EObject {
	eCurrent := eObject
	for eCurrent != nil && eCurrent.EClass() != eClass {
		eCurrent = eCurrent.EContainer()
	}
	return eCurrent
}

func IsAncestor(eAncestor EObject, eObject EObject) bool {
	eCurrent := eObject
	for eCurrent != nil && eCurrent != eAncestor {
		eCurrent = eCurrent.EContainer()
	}
	return eCurrent == eAncestor
}

func ResolveAllInResourceSet(resourceSet EResourceSet) {
	for it := resourceSet.GetResources().Iterator(); it.HasNext(); {
		resource := it.Next().(EResource)
		ResolveAllInResource(resource)
	}
}

func ResolveAllInResource(resource EResource) {
	for it := resource.GetContents().Iterator(); it.HasNext(); {
		object := it.Next().(EObject)
		ResolveAll(object)
	}
}

func ResolveAll(eObject EObject) {
	resolveCrossReferences(eObject)
	for it := eObject.EAllContents(); it.HasNext(); {
		childEObject := it.Next().(EObject)
		resolveCrossReferences(childEObject)
	}
}

func resolveCrossReferences(eObject EObject) {
	for it := eObject.ECrossReferences().Iterator(); it.HasNext(); it.Next() {
		// The loop resolves the cross references by visiting them.
	}
}

func EAllContentsWithClass(eObject EObject, eClass EClass) EIterator {
	return newEAllContentsWithClassIterator(eObject, eClass)
}

func EAllContentsWithTable(eObject EObject, table *EClassTransitionsTable) EIterator {
	return newEAllContentsWithTableIterator(eObject, table)
}
