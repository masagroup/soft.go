package ecore

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type EResourceInternal interface {
	EResource

	DoLoad(decoder EResourceDecoder)
	DoSave(encoder EResourceEncoder)
	DoUnload()

	IsAttachedDetachedRequired() bool
	DoAttached(o EObject)
	DoDetached(o EObject)

	BasicSetLoaded(bool, ENotificationChain) ENotificationChain
	BasicSetResourceSet(EResourceSet, ENotificationChain) ENotificationChain
}

type resourceNotification struct {
	AbstractNotification
	notifier  ENotifier
	featureID int
}

func (n *resourceNotification) GetNotifier() ENotifier {
	return n.notifier
}

func (n *resourceNotification) GetFeature() EStructuralFeature {
	return nil
}

func (n *resourceNotification) GetFeatureID() int {
	return n.featureID
}

func newResourceNotification(
	notifier ENotifier,
	featureID int,
	eventType EventType,
	oldValue interface{},
	newValue interface{},
	position int) *resourceNotification {
	n := new(resourceNotification)
	n.Initialize(n, eventType, oldValue, newValue, position)
	n.notifier = notifier
	n.featureID = featureID
	return n
}

type resourceContents struct {
	BasicENotifyingList
	resource *EResourceImpl
}

func newResourceContents(resource *EResourceImpl) *resourceContents {
	rc := new(resourceContents)
	rc.interfaces = rc
	rc.data = []interface{}{}
	rc.isUnique = true
	rc.resource = resource
	return rc
}

func (rc *resourceContents) GetNotifier() ENotifier {
	return rc.resource.AsENotifier()
}

func (rc *resourceContents) GetFeatureID() int {
	return RESOURCE__CONTENTS
}

func (rc *resourceContents) inverseAdd(object interface{}, notifications ENotificationChain) ENotificationChain {
	n := notifications
	if eObject, _ := object.(EObjectInternal); eObject != nil {
		eResource := rc.resource.AsEResource()
		n = eObject.ESetResource(eResource, n)
		eResource.Attached(eObject)
	}
	return n
}

func (rc *resourceContents) inverseRemove(object interface{}, notifications ENotificationChain) ENotificationChain {
	n := notifications
	if eObject, _ := object.(EObjectInternal); eObject != nil {
		eResource := rc.resource.AsEResource()
		eResource.Detached(eObject)
		n = eObject.ESetResource(nil, n)
	}
	return n
}

func (rc *resourceContents) didAdd(index int, elem interface{}) {
	rc.BasicENotifyingList.didAdd(index, elem)
	if index == rc.Size()-1 {
		rc.loaded()
	}
}

func (rc *resourceContents) didRemove(index int, old interface{}) {
	rc.BasicENotifyingList.didRemove(index, old)
	if rc.Size() == 0 {
		rc.unloaded()
	}
}

func (rc *resourceContents) didClear(oldObjects []interface{}) {
	rc.BasicENotifyingList.didClear(oldObjects)
	rc.unloaded()
}

func (rc *resourceContents) loaded() {
	if !rc.resource.IsLoaded() {
		n := rc.resource.BasicSetLoaded(true, nil)
		if n != nil {
			n.Dispatch()
		}
	}
}

func (rc *resourceContents) unloaded() {
	if rc.resource.IsLoaded() {
		n := rc.resource.BasicSetLoaded(false, nil)
		if n != nil {
			n.Dispatch()
		}
	}
}

type resourceDiagnostics struct {
	BasicENotifyingList
	resource  *EResourceImpl
	featureID int
}

func newResourceDiagnostics(resource *EResourceImpl, featureID int) *resourceDiagnostics {
	rd := new(resourceDiagnostics)
	rd.interfaces = rd
	rd.data = []interface{}{}
	rd.isUnique = true
	rd.resource = resource
	rd.featureID = featureID
	return rd
}

func (rd *resourceDiagnostics) GetNotifier() ENotifier {
	return rd.resource.AsENotifier()
}

func (rd *resourceDiagnostics) GetFeatureID() int {
	return rd.featureID
}

//EResource ...
type EResourceImpl struct {
	ENotifierImpl
	resourceSet     EResourceSet
	objectIDManager EObjectIDManager
	uri             *URI
	contents        EList
	errors          EList
	warnings        EList
	isLoaded        bool
	isLoading       bool
}

// NewBasicEObject is BasicEObject constructor
func NewEResourceImpl() *EResourceImpl {
	r := new(EResourceImpl)
	r.SetInterfaces(r)
	r.Initialize()
	return r
}

func (r *EResourceImpl) AsEResource() EResource {
	return r.GetInterfaces().(EResource)
}

func (r *EResourceImpl) AsEResourceInternal() EResourceInternal {
	return r.GetInterfaces().(EResourceInternal)
}

func (r *EResourceImpl) GetResourceSet() EResourceSet {
	return r.resourceSet
}

func (r *EResourceImpl) GetURI() *URI {
	return r.uri
}

func (r *EResourceImpl) SetURI(uri *URI) {
	oldURI := r.uri
	r.uri = uri
	if r.ENotificationRequired() {
		r.ENotify(newResourceNotification(r.AsENotifier(), RESOURCE__URI, SET, oldURI, uri, -1))
	}
}

func (r *EResourceImpl) GetContents() EList {
	if r.contents == nil {
		r.contents = newResourceContents(r)
	}
	return r.contents
}

func (r *EResourceImpl) GetAllContents() EIterator {
	return r.getAllContentsResolve(r.GetInterfaces(), true)
}

func (r *EResourceImpl) getAllContentsResolve(root interface{}, resolve bool) EIterator {
	return newTreeIterator(root, false, func(o interface{}) EIterator {
		if o == r.GetInterfaces() {
			return o.(EResource).GetContents().Iterator()
		}
		contents := o.(EObject).EContents()
		if !resolve {
			contents = contents.(EObjectList).GetUnResolvedList()
		}
		return contents.Iterator()
	})
}

func (r *EResourceImpl) GetEObject(uriFragment string) EObject {
	id := uriFragment
	size := len(uriFragment)
	if size > 0 {
		if uriFragment[0] == '/' {
			path := strings.Split(uriFragment, "/")
			path = path[1:]
			return r.getObjectByPath(path)
		} else if uriFragment[size-1] == '?' {
			if index := strings.LastIndex(uriFragment[:size-2], "?"); index != -1 {
				id = uriFragment[:index]
			}
		}
	}
	return r.getObjectByID(id)
}

func (r *EResourceImpl) GetURIFragment(eObject EObject) string {
	id := GetEObjectID(eObject)
	if len(id) > 0 {
		return id
	} else {
		internalEObject := eObject.(EObjectInternal)
		if internalEObject.EInternalResource() == r.interfaces {
			if id = r.getIDForObject(eObject); len(id) > 0 {
				return id
			} else {
				return "/" + r.getURIFragmentRootSegment(eObject)
			}
		} else {
			fragmentPath := []string{}
			isContained := false
			for eContainer, _ := internalEObject.EInternalContainer().(EObjectInternal); eContainer != nil; eContainer, _ = internalEObject.EInternalContainer().(EObjectInternal) {
				if id = r.getIDForObject(eObject); len(id) == 0 {
					fragmentPath = append([]string{eContainer.EURIFragmentSegment(internalEObject.EContainingFeature(), internalEObject)}, fragmentPath...)
				}
				internalEObject = eContainer
				if eContainer.EInternalResource() == r.interfaces {
					isContained = true
					break
				}
			}
			if !isContained {
				return "/-1"
			}
			if len(id) == 0 {
				fragmentPath = append([]string{r.getURIFragmentRootSegment(internalEObject)}, fragmentPath...)
				fragmentPath = append([]string{""}, fragmentPath...)
				return strings.Join(fragmentPath, "/")
			} else {
				return id
			}
		}
	}
}

func (r *EResourceImpl) getURIFragmentRootSegment(eObject EObject) string {
	contents := r.GetContents()
	if contents.Size() > 1 {
		return strconv.Itoa(contents.IndexOf(eObject))
	} else {
		return ""
	}
}

func (r *EResourceImpl) getIDForObject(eObject EObject) string {
	if r.objectIDManager != nil {
		if id := r.objectIDManager.GetID(eObject); id != nil {
			return fmt.Sprintf("%v", id)
		}
	}
	return GetEObjectID(eObject)
}

func (r *EResourceImpl) getObjectByID(id string) EObject {
	if r.objectIDManager != nil {
		return r.objectIDManager.GetEObject(id)
	}
	for it := r.getAllContentsResolve(r.GetInterfaces(), false); it.HasNext(); {
		eObject := it.Next().(EObject)
		objectID := GetEObjectID(eObject)
		if id == objectID {
			return eObject
		}
	}
	return nil
}

func (r *EResourceImpl) getObjectByPath(uriFragmentPath []string) EObject {
	var eObject EObject
	if uriFragmentPath == nil || len(uriFragmentPath) == 0 {
		eObject = r.getObjectForRootSegment("")
	} else {
		eObject = r.getObjectForRootSegment(uriFragmentPath[0])
	}
	for i := 1; i < len(uriFragmentPath) && eObject != nil; i++ {
		eObject = eObject.(EObjectInternal).EObjectForFragmentSegment(uriFragmentPath[i])
	}
	return eObject
}

func (r *EResourceImpl) getObjectForRootSegment(rootSegment string) EObject {
	position := 0
	if len(rootSegment) > 0 {
		if rootSegment[0] == '?' {
			return r.getObjectByID(rootSegment[1:])
		} else {
			position, _ = strconv.Atoi(rootSegment)
		}
	}
	if position >= 0 && position < r.GetContents().Size() {
		return r.GetContents().Get(position).(EObject)
	}
	return nil
}

func (r *EResourceImpl) IsAttachedDetachedRequired() bool {
	return r.objectIDManager != nil
}

func (r *EResourceImpl) Attached(object EObject) {
	if resourceInternal := r.AsEResourceInternal(); resourceInternal.IsAttachedDetachedRequired() {
		resourceInternal.DoAttached(object)
		for it := r.getAllContentsResolve(object, false); it.HasNext(); {
			if o, _ := it.Next().(EObject); o != nil {
				resourceInternal.DoAttached(o)
			}
		}
	}
}

func (r *EResourceImpl) DoAttached(object EObject) {
	if r.objectIDManager != nil {
		r.objectIDManager.Register(object)
	}
}

func (r *EResourceImpl) Detached(object EObject) {
	if resourceInternal := r.AsEResourceInternal(); resourceInternal.IsAttachedDetachedRequired() {
		resourceInternal.DoDetached(object)
		for it := r.getAllContentsResolve(object, false); it.HasNext(); {
			if o, _ := it.Next().(EObject); o != nil {
				resourceInternal.DoDetached(o)
			}
		}
	}
}

func (r *EResourceImpl) DoDetached(object EObject) {
	if r.objectIDManager != nil {
		r.objectIDManager.UnRegister(object)
	}
}

var defaultURIConverter EURIConverter = NewEURIConverterImpl()

func (r *EResourceImpl) getURIConverter() EURIConverter {
	if r.resourceSet != nil {
		return r.resourceSet.GetURIConverter()
	}
	return defaultURIConverter
}

func (r *EResourceImpl) getResourceCodecRegistry() EResourceCodecRegistry {
	if r.resourceSet != nil {
		return r.resourceSet.GetResourceCodecRegistry()
	}
	return GetResourceCodecRegistry()
}

func (r *EResourceImpl) Load() {
	r.LoadWithOptions(nil)
}

func (r *EResourceImpl) LoadWithOptions(options map[string]interface{}) {
	if !r.isLoaded {
		uriConverter := r.getURIConverter()
		if uriConverter != nil && r.uri != nil {
			rd, err := uriConverter.CreateReader(r.uri)
			if err != nil {
				errors := r.GetErrors()
				errors.Clear()
				errors.Add(NewEDiagnosticImpl("Unable to create reader for '"+r.uri.String()+"' :"+err.Error(), r.uri.String(), 0, 0))
			} else if rd != nil {
				r.LoadWithReader(rd, options)
				rd.Close()
			}
		}
	}
}

func (r *EResourceImpl) LoadWithReader(rd io.Reader, options map[string]interface{}) {
	if !r.isLoaded {
		codecs := r.getResourceCodecRegistry()
		if codec := codecs.GetCodec(r.uri); codec == nil {
			errors := r.GetErrors()
			errors.Clear()
			errors.Add(NewEDiagnosticImpl("Unable to find codec for '"+r.uri.String()+"'", r.uri.String(), 0, 0))
		} else if decoder := codec.NewDecoder(r.AsEResource(), rd, options); decoder == nil {
			errors := r.GetErrors()
			errors.Clear()
			errors.Add(NewEDiagnosticImpl("Unable to create decoder for '"+r.uri.String()+"'", r.uri.String(), 0, 0))
		} else {
			r.isLoading = true
			n := r.BasicSetLoaded(true, nil)
			if r.errors != nil {
				r.errors.Clear()
			}
			if r.warnings != nil {
				r.warnings.Clear()
			}
			r.GetInterfaces().(EResourceInternal).DoLoad(decoder)
			if n != nil {
				n.Dispatch()
			}
			r.isLoading = false
		}
	}
}

func (r *EResourceImpl) DoLoad(decoder EResourceDecoder) {
	decoder.Decode()
}

func (r *EResourceImpl) Unload() {
	if r.isLoaded {
		n := r.BasicSetLoaded(false, nil)
		r.GetInterfaces().(EResourceInternal).DoUnload()
		if n != nil {
			n.Dispatch()
		}
	}
}

func (r *EResourceImpl) DoUnload() {
	r.contents = nil
	r.errors = nil
	r.warnings = nil
	if r.objectIDManager != nil {
		r.objectIDManager.Clear()
	}
}

func (r *EResourceImpl) IsLoaded() bool {
	return r.isLoaded
}

func (r *EResourceImpl) IsLoading() bool {
	return r.isLoading
}

func (r *EResourceImpl) SetLoading(isLoading bool) {
	r.isLoading = isLoading
}

func (r *EResourceImpl) Save() {
	r.SaveWithOptions(nil)
}

func (r *EResourceImpl) SaveWithOptions(options map[string]interface{}) {
	uriConverter := r.getURIConverter()
	if uriConverter != nil && r.uri != nil {
		w, err := uriConverter.CreateWriter(r.uri)
		if err != nil {
			errors := r.GetErrors()
			errors.Clear()
			errors.Add(NewEDiagnosticImpl("Unable to create writer for '"+r.uri.String()+"' :"+err.Error(), r.uri.String(), 0, 0))
		} else if w != nil {
			r.SaveWithWriter(w, options)
			w.Close()
		}
	}
}

func (r *EResourceImpl) SaveWithWriter(w io.Writer, options map[string]interface{}) {
	codecs := r.getResourceCodecRegistry()
	if codec := codecs.GetCodec(r.uri); codec == nil {
		errors := r.GetErrors()
		errors.Clear()
		errors.Add(NewEDiagnosticImpl("Unable to find codec for '"+r.uri.String()+"'", r.uri.String(), 0, 0))
	} else if encoder := codec.NewEncoder(r.AsEResource(), w, options); encoder == nil {
		errors := r.GetErrors()
		errors.Clear()
		errors.Add(NewEDiagnosticImpl("Unable to create encoder for '"+r.uri.String()+"'", r.uri.String(), 0, 0))
	} else {
		if r.errors != nil {
			r.errors.Clear()
		}
		if r.warnings != nil {
			r.warnings.Clear()
		}
		r.GetInterfaces().(EResourceInternal).DoSave(encoder)
	}
}

func (r *EResourceImpl) DoSave(encoder EResourceEncoder) {
	encoder.Encode()
}

func (r *EResourceImpl) GetErrors() EList {
	if r.errors == nil {
		r.errors = newResourceDiagnostics(r, RESOURCE__ERRORS)
	}
	return r.errors
}

func (r *EResourceImpl) GetWarnings() EList {
	if r.warnings == nil {
		r.warnings = newResourceDiagnostics(r, RESOURCE__WARNINGS)
	}
	return r.warnings
}

func (r *EResourceImpl) BasicSetLoaded(isLoaded bool, msgs ENotificationChain) ENotificationChain {
	notifications := msgs
	oldLoaded := r.isLoaded
	r.isLoaded = isLoaded
	if r.ENotificationRequired() {
		if notifications == nil {
			notifications = NewNotificationChain()
		}
		notifications.Add(newResourceNotification(r.GetInterfaces().(ENotifier), RESOURCE__IS_LOADED, SET, oldLoaded, r.isLoaded, -1))
	}
	return notifications
}

func (r *EResourceImpl) BasicSetResourceSet(resourceSet EResourceSet, msgs ENotificationChain) ENotificationChain {
	notifications := msgs
	oldAbstractResourceSet := r.resourceSet
	if oldAbstractResourceSet != nil {
		l := oldAbstractResourceSet.GetResources().(ENotifyingList)
		notifications = l.RemoveWithNotification(r.GetInterfaces().(ENotifier), notifications)
	}
	r.resourceSet = resourceSet
	if r.ENotificationRequired() {
		if notifications == nil {
			notifications = NewNotificationChain()
		}
		notifications.Add(newResourceNotification(r.GetInterfaces().(ENotifier), RESOURCE__RESOURCE_SET, SET, oldAbstractResourceSet, resourceSet, -1))
	}
	return notifications
}

func (r *EResourceImpl) SetObjectIDManager(objectIDManager EObjectIDManager) {
	r.objectIDManager = objectIDManager
}

func (r *EResourceImpl) GetObjectIDManager() EObjectIDManager {
	return r.objectIDManager
}
