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
	"maps"
	"slices"
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

func crossReferencesInObject(eObject EObject, resolve bool) iter.Seq2[EStructuralFeature, EObject] {
	return func(yield func(EStructuralFeature, EObject) bool) {
		for crossReferenceAny := range eObject.EClass().GetEAllCrossReferences().All() {
			crossReference := crossReferenceAny.(EStructuralFeature)
			if eObject.EIsSet(crossReference) {
				value := eObject.EGetResolve(crossReference, resolve)
				if crossReference.IsMany() {
					for it := value.(EList).Iterator(); it.HasNext(); {
						eObject := it.Next().(EObject)
						if !yield(crossReference, eObject) {
							return
						}
					}
				} else if value != nil {
					eObject := value.(EObject)
					if !yield(crossReference, eObject) {
						return
					}
				}
			}
		}
	}

}

type usage struct {
	object  EObject
	feature EStructuralFeature
}

func findUsages(root any, lookups []EObject) map[EObject][]usage {
	usages := make(map[EObject][]usage)
	iterator := &eAllContentIterator{
		object: root,
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
		if eObject, _ := iterator.Next().(EObject); eObject != nil {
			for crossReference, referencedObject := range crossReferencesInObject(eObject, true) {
				if slices.Contains(lookups, referencedObject) {
					usages[referencedObject] = append(
						usages[referencedObject],
						usage{object: eObject, feature: crossReference},
					)
				}
			}
		}
	}
	return usages
}

func removeValue(eObject EObject, feature EStructuralFeature, value any) {
	if feature.IsMany() {
		l := eObject.EGet(feature).(EList)
		l.Remove(value)
	} else {
		eObject.EUnset(feature)
	}
}

// Removes the object from its containing resource and/or its containing object.
func Remove(eObject EObject) {
	if eObjectInternal, _ := eObject.(EObjectInternal); eObjectInternal != nil {
		if eContainer := eObjectInternal.EInternalContainer(); eContainer != nil {
			if eFeature := eObject.EContainmentFeature(); eFeature != nil {
				removeValue(eContainer, eFeature, eObject)
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
	// compute the usages of the object to be deleted
	rootEObject := GetRootContainer(eObject)
	resource := rootEObject.EResource()
	var usages map[EObject][]usage
	if resource == nil {
		usages = findUsages(rootEObject, []EObject{eObject})
	} else {
		resourceSet := resource.GetResourceSet()
		if resourceSet == nil {
			usages = findUsages(resource, []EObject{eObject})
		} else {
			usages = findUsages(resourceSet, []EObject{eObject})
		}
	}
	// remove the object to be deleted from any referencing features
	for _, usage := range usages[eObject] {
		feature := usage.feature
		object := usage.object
		if feature.IsChangeable() {
			removeValue(object, feature, eObject)
		}
	}
	// remove the object from its containing resource and/or its containing object
	Remove(eObject)
}

/**
 * Deletes the object from its containing resource and/or its containing object
 * as well as from any other feature that references it
 * within the enclosing resource set, resource, or root object.
 * If recursive true, contained children of the object that are in the same resource
 * are similarly removed from any features that reference them.
 */
func DeleteRecursive(eObject EObject, recursive bool) {
	if recursive {
		rootEObject := GetRootContainer(eObject)
		resource := rootEObject.EResource()
		// compute the set of objects to be deleted and the set of directly
		// contained objects in the same resource
		objectSet := map[EObject]struct{}{}
		directSet := map[EObject]struct{}{}
		objectSet[eObject] = struct{}{}
		for it := eObject.EAllContents(); it.HasNext(); {
			childEObject := it.Next().(EObjectInternal)
			if childEObject.EInternalResource() != nil {
				directSet[childEObject] = struct{}{}
				type prunableEIterator interface {
					EIterator
					Prune()
				}
				it.(prunableEIterator).Prune()
			} else {
				objectSet[childEObject] = struct{}{}
			}
		}
		objects := slices.Collect(maps.Keys(objectSet))
		// compute the usages of the objects to be deleted
		var usages map[EObject][]usage
		if resource == nil {
			usages = findUsages(rootEObject, objects)
		} else {
			resourceSet := resource.GetResourceSet()
			if resourceSet == nil {
				usages = findUsages(resource, objects)
			} else {
				usages = findUsages(resourceSet, objects)
			}
		}
		// remove the objects to be deleted from any referencing features
		for deletedObject, usages := range usages {
			for _, usage := range usages {
				feature := usage.feature
				object := usage.object
				if _, contains := objectSet[object]; !contains && feature.IsChangeable() {
					removeValue(object, feature, deletedObject)
				}
			}
		}
		// remove the object from its containing resource and/or its containing object
		Remove(eObject)
		// remove the directly contained objects in the same resource from any referencing features
		for direct := range directSet {
			removeValue(direct.EContainer(), direct.EContainingFeature(), direct)
		}
	} else {
		// non-recursive delete is sufficient
		Delete(eObject)
	}
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
