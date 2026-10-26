package json

import (
	"fmt"
	"io"

	"github.com/masagroup/soft.go/ecore"
)

type Encoder struct {
	resource        ecore.EResource
	w               *jsonWriter
	featureKinds    map[ecore.EStructuralFeature]jsonFeatureKind
	errorFn         func(diagnostic ecore.EDiagnostic)
	idAttributeName string
	keepDefaults    bool
}

func NewEncoder(resource ecore.EResource, w io.Writer, options map[string]any) *Encoder {
	e := &Encoder{
		w:            newJSONWriter(w),
		resource:     resource,
		featureKinds: map[ecore.EStructuralFeature]jsonFeatureKind{},
	}
	if options != nil {
		e.idAttributeName, _ = options[OPTION_ID_ATTRIBUTE_NAME].(string)
		e.keepDefaults = options[OPTION_KEEP_DEFAULTS] == true
	}
	return e
}

func (e *Encoder) EncodeResource() {
	e.errorFn = func(diagnostic ecore.EDiagnostic) {
		e.resource.GetErrors().Add(diagnostic)
	}
	if contents := e.resource.GetContents(); !contents.Empty() {
		object := contents.Get(0).(ecore.EObject)
		e.encodeTopObject(object)
	}
}

func (e *Encoder) EncodeObject(object ecore.EObject) (err error) {
	e.errorFn = func(diagnostic ecore.EDiagnostic) {
		if err == nil {
			err = diagnostic
		}
	}
	e.encodeTopObject(object)
	return
}

func (e *Encoder) encodeTopObject(eObject ecore.EObject) {
	e.w.beginObject()
	e.encodeObject(eObject)
	e.w.endObject()
	if err := e.w.flush(); err != nil {
		location := ""
		if uri := e.resource.GetURI(); uri != nil {
			location = uri.String()
		}
		e.errorFn(ecore.NewEDiagnosticImpl(err.Error(), location, 0, 0))
	}
}

func (e *Encoder) encodeObject(eObject ecore.EObject) {
	eClass := eObject.EClass()
	// class
	e.w.key("eClass")
	e.w.stringValue(e.getClassName(eClass))
	// id
	if idManager := e.resource.GetObjectIDManager(); len(e.idAttributeName) > 0 && idManager != nil {
		if id := idManager.GetID(eObject); id != nil {
			e.w.key(e.idAttributeName)
			e.w.stringValue(fmt.Sprintf("%v", id))
		}
	}
	// features
	for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
		eFeature := itFeature.Next().(ecore.EStructuralFeature)
		e.encodeFeatureValue(eObject, eFeature)
	}
}

func (e *Encoder) encodeObjectReference(eObject ecore.EObject) {
	e.w.key("eClass")
	e.w.stringValue(e.getClassName(eObject.EClass()))
	e.w.key("eRef")
	e.w.stringValue(e.getReference(eObject))
}

func (e *Encoder) encodeFeatureValue(eObject ecore.EObject, eFeature ecore.EStructuralFeature) {
	if !e.shouldSaveFeature(eObject, eFeature) {
		return
	}

	// compute feature kind
	kind, ok := e.featureKinds[eFeature]
	if !ok {
		kind = getJSONCodecFeatureKind(eFeature)
		e.featureKinds[eFeature] = kind
	}

	value := eObject.EGetResolve(eFeature, false)
	switch kind {
	case jfkTransient:
	case jfkData:
		if str, ok := e.getData(value, eFeature); ok {
			e.w.key(eFeature.GetName())
			e.w.stringValue(str)
		}
	case jfkDataList:
		l := value.(ecore.EList)
		e.w.key(eFeature.GetName())
		e.w.beginArray()
		for it := l.Iterator(); it.HasNext(); {
			if str, ok := e.getData(it.Next(), eFeature); ok {
				e.w.stringValue(str)
			}
		}
		e.w.endArray()
	case jfkObject:
		// an unsettable feature can be set to nil : it is not written
		if eValue, _ := value.(ecore.EObject); eValue != nil {
			e.w.key(eFeature.GetName())
			e.w.beginObject()
			e.encodeObject(eValue)
			e.w.endObject()
		}
	case jfkObjectList:
		l := value.(ecore.EList)
		e.w.key(eFeature.GetName())
		e.w.beginArray()
		for it := l.Iterator(); it.HasNext(); {
			e.w.beginObject()
			e.encodeObject(it.Next().(ecore.EObject))
			e.w.endObject()
		}
		e.w.endArray()
	case jfkObjectReference:
		if eValue, _ := value.(ecore.EObject); eValue != nil {
			e.w.key(eFeature.GetName())
			e.w.beginObject()
			e.encodeObjectReference(eValue)
			e.w.endObject()
		}
	case jfkObjectReferenceList:
		l := value.(ecore.EList)
		e.w.key(eFeature.GetName())
		e.w.beginArray()
		for it := l.Iterator(); it.HasNext(); {
			e.w.beginObject()
			e.encodeObjectReference(it.Next().(ecore.EObject))
			e.w.endObject()
		}
		e.w.endArray()
	}
}

func (e *Encoder) shouldSaveFeature(o ecore.EObject, f ecore.EStructuralFeature) bool {
	return o.EIsSet(f) || (e.keepDefaults && f.GetDefaultValueLiteral() != "")
}

func (e *Encoder) getClassName(eClass ecore.EClass) string {
	ePackage := eClass.GetEPackage()
	return ePackage.GetNsURI() + "#//" + eClass.GetName()
}

func (e *Encoder) getData(value any, f ecore.EStructuralFeature) (string, bool) {
	if value == nil {
		return "", false
	} else {
		d := f.GetEType().(ecore.EDataType)
		p := d.GetEPackage()
		f := p.GetEFactoryInstance()
		s := f.ConvertToString(d, value)
		return s, true
	}
}

func (e *Encoder) getReference(eObject ecore.EObject) string {
	eInternal, _ := eObject.(ecore.EObjectInternal)
	if eInternal != nil {
		objectURI := eInternal.EProxyURI()
		if objectURI == nil {
			eOtherResource := eObject.EResource()
			if eOtherResource == nil {
				if e.resource != nil && e.resource.GetObjectIDManager() != nil && e.resource.GetObjectIDManager().GetID(eObject) != nil {
					objectURI = e.getResourceReference(e.resource, eObject)
				} else {
					return ""
				}
			} else {
				objectURI = e.getResourceReference(eOtherResource, eObject)
			}
		}
		if resourceURI := e.resource.GetURI(); resourceURI != nil {
			objectURI = resourceURI.Relativize(objectURI)
		}
		return objectURI.String()
	}
	return ""
}

func (e *Encoder) getResourceReference(resource ecore.EResource, object ecore.EObject) *ecore.URI {
	return ecore.NewURIBuilder(resource.GetURI()).SetFragment(resource.GetURIFragment(object)).URI()
}
