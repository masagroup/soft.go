package ecore

import (
	"sync"
)

type resourcesList struct {
	BasicENotifyingList
	resourceSet *EResourceSetImpl
}

func newResourcesList(resourceSet *EResourceSetImpl) *resourcesList {
	l := new(resourcesList)
	l.interfaces = l
	l.data = []any{}
	l.isUnique = true
	l.resourceSet = resourceSet
	return l
}

func (l *resourcesList) GetNotifier() ENotifier {
	return l.resourceSet.AsENotifier()
}

func (l *resourcesList) GetFeatureID() int {
	return RESOURCE_SET__RESOURCES
}

func (l *resourcesList) InverseAdd(object any, notifications ENotificationChain) ENotificationChain {
	if eResource, _ := object.(EResourceInternal); eResource != nil {
		return eResource.BasicSetResourceSet(l.resourceSet.AsEResourceSet(), notifications)
	}
	return notifications
}

func (l *resourcesList) InverseRemove(object any, notifications ENotificationChain) ENotificationChain {
	if eResource, _ := object.(EResourceInternal); eResource != nil {
		return eResource.BasicSetResourceSet(nil, notifications)
	}
	return notifications
}

func (l *resourcesList) DidAdd(index int, elem any) {
	l.BasicENotifyingList.DidAdd(index, elem)
	if eResource, _ := elem.(EResource); eResource != nil {
		l.resourceSet.indexResource(eResource, eResource.GetURI())
	}
}

func (l *resourcesList) DidRemove(index int, old any) {
	l.BasicENotifyingList.DidRemove(index, old)
	if eResource, _ := old.(EResource); eResource != nil {
		l.resourceSet.unindexResource(eResource)
	}
}

func (l *resourcesList) DidSet(index int, newElem any, oldElem any) {
	l.BasicENotifyingList.DidSet(index, newElem, oldElem)
	if eResource, _ := oldElem.(EResource); eResource != nil {
		l.resourceSet.unindexResource(eResource)
	}
	if eResource, _ := newElem.(EResource); eResource != nil {
		l.resourceSet.indexResource(eResource, eResource.GetURI())
	}
}

func (l *resourcesList) DidClear(oldObjects []any) {
	l.BasicENotifyingList.DidClear(oldObjects)
	l.resourceSet.resourceIndex = nil
}

// EResourceSetImpl ...
// lazyMutex guards the lazily built registries and URI converter, never a resource load.
type EResourceSetImpl struct {
	ENotifierImpl
	lazyMutex             sync.Mutex
	resources             EList
	uriConverter          EURIConverter
	uriResourceMap        map[*URI]EResource
	resourceIndex         map[string]EResource
	resourceCodecRegistry ECodecRegistry
	packageRegistry       EPackageRegistry
}

func NewEResourceSetImpl() *EResourceSetImpl {
	rs := new(EResourceSetImpl)
	rs.SetInterfaces(rs)
	rs.Initialize()
	return rs
}

func (r *EResourceSetImpl) Initialize() {
	r.ENotifierImpl.Initialize()
	r.resources = newResourcesList(r)
}

func (r *EResourceSetImpl) AsEResourceSet() EResourceSet {
	return r.interfaces.(EResourceSet)
}

func (r *EResourceSetImpl) GetResources() EList {
	return r.resources
}

func (r *EResourceSetImpl) GetResource(uri *URI, loadOnDemand bool) EResource {
	if r.uriResourceMap != nil {
		resource := r.uriResourceMap[uri]
		if resource != nil {
			if loadOnDemand && !resource.IsLoaded() {
				resource.Load()
			}
			return resource
		}
	}

	ePackage := r.GetPackageRegistry().GetPackage(uri.String())
	if ePackage != nil {
		resource := ePackage.EResource()
		if resource != nil {
			if r.uriResourceMap != nil {
				r.uriResourceMap[uri] = resource
			}
			return resource
		}
	}

	resource := r.getIndexedResource(uri)
	if resource != nil {
		if loadOnDemand && !resource.IsLoaded() {
			resource.Load()
		}
		if r.uriResourceMap != nil {
			r.uriResourceMap[uri] = resource
		}
		return resource
	}

	if loadOnDemand {
		resource := r.AsEResourceSet().CreateResource(uri)
		if resource != nil {
			resource.Load()
			if r.uriResourceMap != nil {
				r.uriResourceMap[uri] = resource
			}
		}
		return resource
	}

	return nil
}

// getIndexedResource returns the resource of the set whose uri, once normalized
// by the uri converter, is uri normalized, or nil. The resources are indexed by
// uri and by normalized uri, and the list of the resources is only scanned when
// the index cannot answer.
func (r *EResourceSetImpl) getIndexedResource(uri *URI) EResource {
	// same uri : no normalization needed
	if resource := r.resourceIndex[uri.String()]; resource != nil && uri.Equals(resource.GetURI()) {
		return resource
	}
	uriConverter := r.GetURIConverter()
	normalizedURI := uriConverter.Normalize(uri)
	if resource := r.resourceIndex[normalizedURI.String()]; resource != nil {
		// the uri map of the converter may have changed since the resource was indexed
		if resourceURI := resource.GetURI(); resourceURI != nil && (resourceURI.Equals(uri) || uriConverter.Normalize(resourceURI).Equals(normalizedURI)) {
			return resource
		}
	}
	// a uri the index can not answer : scan the resources
	for it := r.resources.Iterator(); it.HasNext(); {
		resource := it.Next().(EResource)
		if resourceURI := resource.GetURI(); resourceURI != nil && uriConverter.Normalize(resourceURI).Equals(normalizedURI) {
			r.setIndexedResource(normalizedURI.String(), resource)
			return resource
		}
	}
	return nil
}

// indexResource registers resource, whose uri is uri, in the index of the
// resources of the set, by uri and by normalized uri. The first resource
// indexed for a uri is kept, as the first resource of the list is the one
// GetResource returns.
func (r *EResourceSetImpl) indexResource(resource EResource, uri *URI) {
	if uri == nil {
		return
	}
	if r.resourceIndex[uri.String()] == nil {
		r.setIndexedResource(uri.String(), resource)
	}
	normalizedURI := r.GetURIConverter().Normalize(uri)
	if r.resourceIndex[normalizedURI.String()] == nil {
		r.setIndexedResource(normalizedURI.String(), resource)
	}
}

func (r *EResourceSetImpl) setIndexedResource(key string, resource EResource) {
	if r.resourceIndex == nil {
		r.resourceIndex = map[string]EResource{}
	}
	r.resourceIndex[key] = resource
}

// unindexResource removes resource from the index of the resources of the set.
// Another resource of the set with the same uri, if any, is indexed again by the
// next GetResource that scans the list.
func (r *EResourceSetImpl) unindexResource(resource EResource) {
	for key, indexed := range r.resourceIndex {
		if indexed == resource {
			delete(r.resourceIndex, key)
		}
	}
}

// resourceURIChanged is called by a resource of the set when its uri changes.
func (r *EResourceSetImpl) resourceURIChanged(resource EResource, newURI *URI) {
	r.unindexResource(resource)
	r.indexResource(resource, newURI)
}

func (r *EResourceSetImpl) CreateResource(uri *URI) EResource {
	resource := NewEResourceImpl()
	resource.SetURI(uri)
	r.resources.Add(resource)
	return resource
}

func (r *EResourceSetImpl) GetEObject(uri *URI, loadOnDemand bool) EObject {
	trim := uri
	if uri.Fragment() != "" {
		trim = uri.TrimFragment()
	}
	resource := r.GetResource(trim, loadOnDemand)
	if resource != nil {
		return resource.GetEObject(uri.Fragment())
	}
	return nil
}

func (r *EResourceSetImpl) GetURIConverter() EURIConverter {
	r.lazyMutex.Lock()
	defer r.lazyMutex.Unlock()
	if r.uriConverter == nil {
		r.uriConverter = NewEURIConverterImpl()
	}
	return r.uriConverter
}

func (r *EResourceSetImpl) SetURIConverter(uriConverter EURIConverter) {
	r.lazyMutex.Lock()
	defer r.lazyMutex.Unlock()
	r.uriConverter = uriConverter
	// the index is keyed by normalized uri : rebuild it with the new converter
	r.resourceIndex = nil
	for it := r.resources.Iterator(); it.HasNext(); {
		resource := it.Next().(EResource)
		r.indexResource(resource, resource.GetURI())
	}
}

func (r *EResourceSetImpl) GetPackageRegistry() EPackageRegistry {
	r.lazyMutex.Lock()
	defer r.lazyMutex.Unlock()
	if r.packageRegistry == nil {
		r.packageRegistry = NewEPackageRegistryImplWithDelegate(GetPackageRegistry())
	}
	return r.packageRegistry
}

func (r *EResourceSetImpl) SetPackageRegistry(packageRegistry EPackageRegistry) {
	r.lazyMutex.Lock()
	defer r.lazyMutex.Unlock()
	r.packageRegistry = packageRegistry
}

func (r *EResourceSetImpl) GetCodecRegistry() ECodecRegistry {
	r.lazyMutex.Lock()
	defer r.lazyMutex.Unlock()
	if r.resourceCodecRegistry == nil {
		r.resourceCodecRegistry = NewECodecRegistryImplWithDelegate(GetCodecRegistry())
	}
	return r.resourceCodecRegistry
}

func (r *EResourceSetImpl) SetCodecRegistry(resourceCodecRegistry ECodecRegistry) {
	r.lazyMutex.Lock()
	defer r.lazyMutex.Unlock()
	r.resourceCodecRegistry = resourceCodecRegistry
}

func (r *EResourceSetImpl) SetURIResourceMap(uriResourceMap map[*URI]EResource) {
	r.uriResourceMap = uriResourceMap
}

func (r *EResourceSetImpl) GetURIResourceMap() map[*URI]EResource {
	return r.uriResourceMap
}
