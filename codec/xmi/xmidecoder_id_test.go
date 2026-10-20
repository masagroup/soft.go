package xmi

import (
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// an xmi resource whose objects have ids (EMF's xmi:id) references them by id
func TestXMIDecoder_XMIIDs(t *testing.T) {
	doc := `<?xml version="1.0" encoding="UTF-8"?>
<ecore:EPackage xmi:version="2.0" xmlns:xmi="http://www.omg.org/XMI" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:ecore="http://www.eclipse.org/emf/2002/Ecore" xmi:id="_p" name="p">
  <eClassifiers xsi:type="ecore:EClass" xmi:id="_a" name="A" eSuperTypes="_b"/>
  <eClassifiers xsi:type="ecore:EClass" xmi:id="_b" name="B"/>
</ecore:EPackage>`
	// an id manager of string ids, which registers the xmi ids
	newStringIDManager := func() ecore.EObjectIDManager {
		return ecore.NewUniqueIDManager(
			func() string { return "" },
			func(id string) bool { return id != "" },
			func(id any) (string, error) { s, _ := id.(string); return s, nil },
			func(string) {})
	}
	for name, withIDManager := range map[string]bool{"without id manager": false, "with id manager": true} {
		r := ecore.NewEResourceImpl()
		r.SetURI(ecore.NewURI("file:///p.ecore"))
		if withIDManager {
			r.SetObjectIDManager(newStringIDManager())
		}
		NewDecoder(r, strings.NewReader(doc), nil).DecodeResource()
		require.True(t, r.GetErrors().Empty(), name)
		ePackage := r.GetContents().Get(0).(ecore.EPackage)
		a := ePackage.GetEClassifier("A").(ecore.EClass)
		b := ePackage.GetEClassifier("B").(ecore.EClass)
		assert.Equal(t, []any{b}, a.GetESuperTypes().ToArray(), name)
		if manager := r.GetObjectIDManager(); manager != nil {
			assert.Equal(t, "_b", manager.GetID(b), name)
		}
	}
}
