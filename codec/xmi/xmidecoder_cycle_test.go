// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmi

import (
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// cyclePackage is a package of the classes A, B and C, with the super types given by class
func cyclePackage(superTypes map[string]string) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<ecore:EPackage xmi:version="2.0" xmlns:xmi="http://www.omg.org/XMI" xmlns:ecore="http://www.eclipse.org/emf/2002/Ecore" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="cycle" nsURI="http://example.com/cycle" nsPrefix="cycle">
`)
	for _, name := range []string{"A", "B", "C"} {
		b.WriteString(`  <eClassifiers xsi:type="ecore:EClass" name="` + name + `"`)
		if supers := superTypes[name]; supers != "" {
			b.WriteString(` eSuperTypes="` + supers + `"`)
		}
		b.WriteString(`>
    <eStructuralFeatures xsi:type="ecore:EAttribute" name="` + strings.ToLower(name) + `">
      <eType xsi:type="ecore:EDataType" href="http://www.eclipse.org/emf/2002/Ecore#//EString"/>
    </eStructuralFeatures>
  </eClassifiers>
`)
	}
	b.WriteString("</ecore:EPackage>\n")
	return b.String()
}

// a package whose super types close a cycle loads with an error, without recursing without end
func TestXMIDecoder_SuperTypesCycle(t *testing.T) {
	for name, superTypes := range map[string]map[string]string{
		"self":             {"A": "#//A"},
		"two classes":      {"A": "#//B", "B": "#//A"},
		"forward and back": {"B": "#//C", "C": "#//B"},
		"three classes":    {"A": "#//B", "B": "#//C", "C": "#//A"},
		"above a cycle":    {"A": "#//B", "B": "#//A", "C": "#//A"},
	} {
		t.Run(name, func(t *testing.T) {
			resource := NewProcessor().LoadWithReader(strings.NewReader(cyclePackage(superTypes)), nil)
			require.NotNil(t, resource)
			require.Equal(t, 1, resource.GetErrors().Size(), diagnosticError(resource.GetErrors()))
			assert.Contains(t, resource.GetErrors().Get(0).(ecore.EDiagnostic).GetMessage(), "cycle of super types")
		})
	}
}

// a hierarchy without cycle, diamond included, loads without error
func TestXMIDecoder_SuperTypesDiamond(t *testing.T) {
	resource := NewProcessor().LoadWithReader(strings.NewReader(cyclePackage(map[string]string{"B": "#//A", "C": "#//A #//B"})), nil)
	require.NotNil(t, resource)
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	c := resource.GetContents().Get(0).(ecore.EPackage).GetEClassifier("C").(ecore.EClass)
	assert.Equal(t, 3, c.GetEAllAttributes().Size())
}
