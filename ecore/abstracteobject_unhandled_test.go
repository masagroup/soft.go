package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// A static feature that the generated ESetFromID or EUnsetFromID of a class does not handle, an
// unchangeable feature such as the package of a classifier or the derived attributes of a class,
// panics: it recursed between ESet and ESetFromID until the stack overflowed, a fatal error no
// recover catches. A malformed XMI document setting such a feature crashed the process decoding it.
func TestAbstractEObject_UnhandledStaticFeature(t *testing.T) {
	ePackage := GetPackage()
	for name, feature := range map[string]EStructuralFeature{
		"ePackage":       ePackage.GetEClassifier_EPackage(),
		"eAllAttributes": ePackage.GetEClass_EAllAttributes(),
	} {
		t.Run(name, func(t *testing.T) {
			eClass := GetFactory().CreateEClass()
			assert.PanicsWithValue(t, "The feature '"+feature.GetName()+"' cannot be set", func() {
				eClass.ESet(feature, nil)
			})
			assert.PanicsWithValue(t, "The feature '"+feature.GetName()+"' cannot be unset", func() {
				eClass.EUnset(feature)
			})
		})
	}
}
