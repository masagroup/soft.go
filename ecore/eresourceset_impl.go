package ecore

import "net/url"

type resourcesList struct {
	BasicENotifyingList
	resourceSet *EResourceSetImpl
}

func newResourcesList(resourceSet *EResourceSetImpl) *resourcesList {
	l := new(resourcesList)
	l.interfaces = l
	l.data = []interface{}{}
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

func (l *resourcesList) inverseAdd(object interface{}, notifications ENotificationChain) ENotificationChain {
	eResource := object.(EResourceInternal)
	n := notifications
	n = eResource.basicSetResourceSet(l.resourceSet, n)
	return n
}

func (l *resourcesList) inverseRemove(object interface{}, notifications ENotificationChain) ENotificationChain {
	eResource := object.(EResourceInternal)
	n := notifications
	n = eResource.basicSetResourceSet(nil, n)
	return n
}

//EResourceSetImpl ...
type EResourceSetImpl struct {
	ENotifierImpl
	resources               EList
	uriConverter            EURIConverter
	uriResourceMap          map[*url.URL]EResource
	resourceFactoryRegistry EResourceFactoryRegistry
	packageRegistry         EPackageRegistry
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

func (r *EResourceSetImpl) GetResource(uri *url.URL, loadOnDemand bool) EResource {
	if r.uriResourceMap != nil {
		resource := r.uriResourceMap[uri]
		if resource != nil {
			if loadOnDemand && !resource.IsLoaded() {
				resource.Load()
			}
			return resource
		}
	}

	normalizedURI := r.GetURIConverter().Normalize(uri)
	for it := r.resources.Iterator(); it.HasNext(); {
		resource := it.Next().(EResource)
		resourceURI := r.GetURIConverter().Normalize(resource.GetURI())
		if *resourceURI == *normalizedURI {
			if loadOnDemand && !resource.IsLoaded() {
				resource.Load()
			}
			if r.uriResourceMap != nil {
				r.uriResourceMap[uri] = resource
			}
			return resource
		}
	}

	ePackage := r.GetPackageRegistry().GetPackage(uri.String())
	if ePackage != nil {
		return ePackage.EResource()
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

func (r *EResourceSetImpl) CreateResource(uri *url.URL) EResource {
	resourceFactory := r.GetResourceFactoryRegistry().GetFactory(uri)
	if resourceFactory != nil {
		resource := resourceFactory.CreateResource(uri)
		r.resources.Add(resource)
		return resource
	}
	return nil
}

func (r *EResourceSetImpl) GetEObject(uri *url.URL, loadOnDemand bool) EObject {
	trim := TrimURIFragment(uri)
	resource := r.GetResource(trim, loadOnDemand)
	if resource != nil {
		return resource.GetEObject(uri.Fragment)
	}
	return nil
}

func (r *EResourceSetImpl) GetURIConverter() EURIConverter {
	if r.uriConverter == nil {
		r.uriConverter = NewEURIConverterImpl()
	}
	return r.uriConverter
}

func (r *EResourceSetImpl) SetURIConverter(uriConverter EURIConverter) {
	r.uriConverter = uriConverter
}

func (r *EResourceSetImpl) GetPackageRegistry() EPackageRegistry {
	if r.packageRegistry == nil {
		r.packageRegistry = NewEPackageRegistryImplWithDelegate(GetPackageRegistry())
	}
	return r.packageRegistry
}

func (r *EResourceSetImpl) SetPackageRegistry(packageRegistry EPackageRegistry) {
	r.packageRegistry = packageRegistry
}

func (r *EResourceSetImpl) GetResourceFactoryRegistry() EResourceFactoryRegistry {
	if r.resourceFactoryRegistry == nil {
		r.resourceFactoryRegistry = NewEResourceFactoryRegistryImplWithDelegate(GetResourceFactoryRegistry())
	}
	return r.resourceFactoryRegistry
}

func (r *EResourceSetImpl) SetResourceFactoryRegistry(resourceFactoryRegistry EResourceFactoryRegistry) {
	r.resourceFactoryRegistry = resourceFactoryRegistry
}

func (r *EResourceSetImpl) SetURIResourceMap(uriResourceMap map[*url.URL]EResource) {
	r.uriResourceMap = uriResourceMap
}

func (r *EResourceSetImpl) GetURIResourceMap() map[*url.URL]EResource {
	return r.uriResourceMap
}
