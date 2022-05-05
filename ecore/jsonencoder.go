package ecore

import (
	"io"

	"github.com/karlseguin/jsonwriter"
)

type JSONEncoder struct {
	w                    *jsonwriter.Writer
	resource             EResource
	objectRoot           EObject
	baseURI              *URI
	isIDAttributeEncoded bool
	errorFn              func(diagnostic EDiagnostic)
}

func NewJSONEncoder(resource EResource, w io.Writer, options map[string]interface{}) *JSONEncoder {
	e := &JSONEncoder{
		w:        jsonwriter.New(w),
		resource: resource,
	}
	if uri := resource.GetURI(); uri != nil {
		e.baseURI = uri
	}
	if options != nil {
		e.isIDAttributeEncoded = options[JSON_OPTION_ID_ATTRIBUTE] == true
	}
	return e
}

func (e *JSONEncoder) Encode() {
	e.errorFn = func(diagnostic EDiagnostic) {
		e.resource.GetErrors().Add(diagnostic)
	}
	if contents := e.resource.GetContents(); !contents.Empty() {
		object := contents.Get(0).(EObject)
		e.encodeTopObject(object)
	}
}

func (e *JSONEncoder) EncodeObject(object EObject) (err error) {
	e.errorFn = func(diagnostic EDiagnostic) {
		if err == nil {
			err = diagnostic
		}
	}
	e.encodeTopObject(object)
	return
}

func (e *JSONEncoder) encodeTopObject(eObject EObject) {
	e.w.RootObject(func() {
		e.encodeObject(eObject, checkContainer)
	})
}

func (e *JSONEncoder) encodeObject(eObject EObject, check checkType) {
	eObjectInternal := eObject.(EObjectInternal)
	eClass := eObject.EClass()
	ePackage := eClass.GetEPackage()
	e.w.KeyString("eClass", ePackage.GetNsURI()+"#//"+eClass.GetName())
	saveFeatureValues := true
	switch check {
	case checkDirectResource:
		if eObjectInternal.EIsProxy() {
			e.w.KeyString("eRef", eObjectInternal.EProxyURI().String())
			saveFeatureValues = false
		} else if eResource := eObjectInternal.EInternalResource(); eResource != nil {
			uri := eResource.GetURI().Copy()
			uri.Fragment = eResource.GetURIFragment(eObject)
			e.w.KeyString("eRef", uri.String())
			saveFeatureValues = false
		}
	case checkResource:
		if eObjectInternal.EIsProxy() {
			e.w.KeyString("eRef", eObjectInternal.EProxyURI().String())
			saveFeatureValues = false
		} else if eResource := eObjectInternal.EResource(); eResource != nil &&
			(eResource != e.resource ||
				(e.objectRoot != nil && !IsAncestor(e.objectRoot, eObjectInternal))) {
			// encode object as uri and fragment if object is in a different resource
			// or if in the same resource and root object is not its ancestor
			uri := eResource.GetURI().Copy()
			uri.Fragment = eResource.GetURIFragment(eObject)
			e.w.KeyString("eRef", uri.String())
			saveFeatureValues = false
		}
	case checkNothing:
	case checkContainer:
	}

	if saveFeatureValues {

		// id attribute
		if objectIDManager := e.resource.GetObjectIDManager(); e.isIDAttributeEncoded && objectIDManager != nil {
			if id := objectIDManager.GetID(eObject); id != nil {
				e.w.KeyValue("eID", id)
			}
		}

		// features
		for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
			eFeature := itFeature.Next().(EStructuralFeature)

			if !eFeature.IsTransient() {
				e.encodeFeatureValue(eObjectInternal, eFeature)
			}
		}
	}

}

func (e *JSONEncoder) encodeFeatureValue(eObject EObject, check checkType) {

}
