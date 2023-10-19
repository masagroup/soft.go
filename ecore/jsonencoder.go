package ecore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/karlseguin/jsonwriter"
)

type JSONEncoder struct {
	resource        EResource
	w               *jsonwriter.Writer
	featureKinds    map[EStructuralFeature]jsonFeatureKind
	errorFn         func(diagnostic EDiagnostic)
	idAttributeName string
	keepDefaults    bool
}

func NewJSONEncoder(resource EResource, w io.Writer, options map[string]interface{}) *JSONEncoder {
	e := &JSONEncoder{
		w:            jsonwriter.New(w),
		resource:     resource,
		featureKinds: map[EStructuralFeature]jsonFeatureKind{},
	}
	if options != nil {
		e.idAttributeName, _ = options[JSON_OPTION_ID_ATTRIBUTE_NAME].(string)
	}
	return e
}

func (e *JSONEncoder) EncodeResource() {
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
		e.encodeObject(eObject)
	})
}

func (e *JSONEncoder) encodeObject(eObject EObject) {
	eClass := eObject.EClass()
	// class
	e.w.KeyString("eClass", e.getClassName(eClass))
	// id
	if idManager := e.resource.GetObjectIDManager(); len(e.idAttributeName) > 0 && idManager != nil {
		if id := idManager.GetID(eObject); id != nil {
			e.w.KeyValue(e.idAttributeName, fmt.Sprintf("%v", id))
		}
	}
	// features
	for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
		eFeature := itFeature.Next().(EStructuralFeature)
		e.encodeFeatureValue(eObject, eFeature)
	}
}

func (e *JSONEncoder) encodeObjectReference(eObject EObject) {
	e.w.KeyString("eClass", e.getClassName(eObject.EClass()))
	e.w.KeyString("eRef", e.getReference(eObject))
}

func jsonEscape(i string) string {
	w := &bytes.Buffer{}
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	if err := e.Encode(i); err != nil {
		panic(err)
	}
	// last \n
	b := w.Bytes()
	return string(b[:len(b)-1])
}

func (e *JSONEncoder) encodeFeatureValue(eObject EObject, eFeature EStructuralFeature) {
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
		str, ok := e.getData(value, eFeature)
		if ok {
			e.w.Key(eFeature.GetName())
			e.w.Raw([]byte(jsonEscape(str)))
		}
	case jfkDataList:
		l := value.(EList)
		e.w.Array(eFeature.GetName(), func() {
			for it := l.Iterator(); it.HasNext(); {
				str, ok := e.getData(it.Next(), eFeature)
				if ok {
					e.w.Value(str)
				}
			}
		})
	case jfkObject:
		e.w.Object(eFeature.GetName(), func() {
			e.encodeObject(value.(EObject))
		})
	case jfkObjectList:
		l := value.(EList)
		e.w.Array(eFeature.GetName(), func() {
			for it := l.Iterator(); it.HasNext(); {
				e.w.ArrayObject(func() {
					e.encodeObject(it.Next().(EObject))
				})
			}
		})
	case jfkObjectReference:
		e.w.Object(eFeature.GetName(), func() {
			e.encodeObjectReference(value.(EObject))
		})
	case jfkObjectReferenceList:
		l := value.(EList)
		e.w.Array(eFeature.GetName(), func() {
			for it := l.Iterator(); it.HasNext(); {
				e.w.ArrayObject(func() {
					e.encodeObjectReference(it.Next().(EObject))
				})
			}
		})
	}
}

func (e *JSONEncoder) shouldSaveFeature(o EObject, f EStructuralFeature) bool {
	return o.EIsSet(f) || (e.keepDefaults && f.GetDefaultValueLiteral() != "")
}

func (e *JSONEncoder) getClassName(eClass EClass) string {
	ePackage := eClass.GetEPackage()
	return ePackage.GetNsURI() + "#//" + eClass.GetName()
}

func (e *JSONEncoder) getData(value interface{}, f EStructuralFeature) (string, bool) {
	if value == nil {
		return "", false
	} else {
		d := f.GetEType().(EDataType)
		p := d.GetEPackage()
		f := p.GetEFactoryInstance()
		s := f.ConvertToString(d, value)
		return s, true
	}
}

func (e *JSONEncoder) getReference(eObject EObject) string {
	eInternal, _ := eObject.(EObjectInternal)
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
		objectURI = e.resource.GetURI().Relativize(objectURI)
		return objectURI.String()
	}
	return ""
}

func (e *JSONEncoder) getResourceReference(resource EResource, object EObject) *URI {
	return NewURIBuilder(resource.GetURI()).SetFragment(resource.GetURIFragment(object)).URI()
}
