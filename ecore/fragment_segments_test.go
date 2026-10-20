package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// the fragments of annotations and of named elements with escaped characters designate them back
func TestFragmentSegments_RoundTrip(t *testing.T) {
	f := GetFactory()
	ePackage := f.CreateEPackage()
	ePackage.SetName("p")
	eClass := f.CreateEClass()
	eClass.SetName("my class")
	ePackage.GetEClassifiers().Add(eClass)
	eAnnotation := f.CreateEAnnotation()
	eAnnotation.SetSource("http://x")
	ePackage.GetEAnnotations().Add(eAnnotation)
	eSecondAnnotation := f.CreateEAnnotation()
	eSecondAnnotation.SetSource("http://x")
	ePackage.GetEAnnotations().Add(eSecondAnnotation)
	eResource := NewEResourceImpl()
	eResource.SetURI(NewURI("file:///p.ecore"))
	eResource.GetContents().Add(ePackage)

	for _, eObject := range []EObject{eClass, eAnnotation, eSecondAnnotation} {
		fragment := eResource.GetURIFragment(eObject)
		require.NotPanics(t, func() {
			assert.Equal(t, eObject, eResource.GetEObject(fragment), fragment)
		})
	}
}
