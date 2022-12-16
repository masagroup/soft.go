package ecore

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

func (l *resourcesList) inverseAdd(object any, notifications ENotificationChain) ENotificationChain {
	if eResource, _ := object.(EResourceInternal); eResource != nil {
		return eResource.BasicSetResourceSet(l.resourceSet.AsEResourceSet(), notifications)
	}
	return notifications
}

func (l *resourcesList) inverseRemove(object any, notifications ENotificationChain) ENotificationChain {
	if eResource, _ := object.(EResourceInternal); eResource != nil {
		return eResource.BasicSetResourceSet(nil, notifications)
	}
	return notifications
}

// EResourceSetImpl ...
type EResourceSetImpl struct {
	ENotifierImpl
	resources             EList
	uriConverter          EURIConverter
	uriResourceMap        map[*URI]EResource
	resourceCodecRegistry EResourceCodecRegistry
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

	normalizedURI := r.GetURIConverter().Normalize(uri)
	for it := r.resources.Iterator(); it.HasNext(); {
		resource := it.Next().(EResource)
		resourceURI := r.GetURIConverter().Normalize(resource.GetURI())
		if resourceURI.Equals(normalizedURI) {
			if loadOnDemand && !resource.IsLoaded() {
				resource.Load()
			}
			if r.uriResourceMap != nil {
				r.uriResourceMap[uri] = resource
			}
			return resource
		}
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

func (r *EResourceSetImpl) CreateResource(uri *URI) EResource {
	resource := NewEResourceImpl()
	resource.SetURI(uri)
	r.resources.Add(resource)
	return resource
}

func (r *EResourceSetImpl) GetEObject(uri *URI, loadOnDemand bool) EObject {
	trim := uri.TrimFragment()
	resource := r.GetResource(trim, loadOnDemand)
	if resource != nil {
		return resource.GetEObject(uri.Fragment())
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

func (r *EResourceSetImpl) GetResourceCodecRegistry() EResourceCodecRegistry {
	if r.resourceCodecRegistry == nil {
		r.resourceCodecRegistry = NewEResourceCodecRegistryImplWithDelegate(GetResourceCodecRegistry())
	}
	return r.resourceCodecRegistry
}

func (r *EResourceSetImpl) SetResourceCodecRegistry(resourceCodecRegistry EResourceCodecRegistry) {
	r.resourceCodecRegistry = resourceCodecRegistry
}

func (r *EResourceSetImpl) SetURIResourceMap(uriResourceMap map[*URI]EResource) {
	r.uriResourceMap = uriResourceMap
}

func (r *EResourceSetImpl) GetURIResourceMap() map[*URI]EResource {
	return r.uriResourceMap
}
