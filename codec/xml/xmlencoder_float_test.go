// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************
package xml

import (
	"math"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// floatModel is a dynamic model with single and many float attributes.
type floatModel struct {
	ePackage ecore.EPackage
	sample   ecore.EClass
	d1       ecore.EAttribute
	d2       ecore.EAttribute
	d3       ecore.EAttribute
	f1       ecore.EAttribute
	finite   ecore.EAttribute
	label    ecore.EAttribute
	ds       ecore.EAttribute
}

func newFloatModel() *floatModel {
	f := ecore.GetFactory()
	ecorePackage := ecore.GetPackage()
	m := &floatModel{}
	m.ePackage = f.CreateEPackage()
	m.ePackage.SetName("floats")
	m.ePackage.SetNsURI("http://floats")
	m.ePackage.SetNsPrefix("floats")
	m.ePackage.SetEFactoryInstance(f.CreateEFactory())
	m.sample = f.CreateEClass()
	m.sample.SetName("Sample")
	m.ePackage.GetEClassifiers().Add(m.sample)
	newAttribute := func(name string, eType ecore.EDataType) ecore.EAttribute {
		eAttribute := f.CreateEAttribute()
		eAttribute.SetName(name)
		eAttribute.SetEType(eType)
		m.sample.GetEStructuralFeatures().Add(eAttribute)
		return eAttribute
	}
	m.d1 = newAttribute("d1", ecorePackage.GetEDouble())
	m.d2 = newAttribute("d2", ecorePackage.GetEDouble())
	m.d3 = newAttribute("d3", ecorePackage.GetEDouble())
	m.f1 = newAttribute("f1", ecorePackage.GetEFloat())
	m.finite = newAttribute("finite", ecorePackage.GetEDouble())
	m.label = newAttribute("label", ecorePackage.GetEString())
	m.ds = newAttribute("ds", ecorePackage.GetEDouble())
	m.ds.SetUpperBound(ecore.UNBOUNDED_MULTIPLICITY)
	return m
}

func (m *floatModel) newResource() ecore.EResource {
	rs := ecore.NewEResourceSetImpl()
	rs.GetPackageRegistry().RegisterPackage(m.ePackage)
	return rs.CreateResource(ecore.NewURI("floats.xml"))
}

// The infinite and NaN floats, single or many, double or float, are written in the XSD lexical
// space (INF, -INF, NaN) and read back.
func TestXMLEncoder_SpecialFloats(t *testing.T) {
	m := newFloatModel()
	sample := m.ePackage.GetEFactoryInstance().Create(m.sample)
	sample.ESet(m.d1, math.Inf(1))
	sample.ESet(m.d2, math.Inf(-1))
	sample.ESet(m.d3, math.NaN())
	sample.ESet(m.f1, float32(math.Inf(1)))
	sample.ESet(m.finite, 1.5)
	sample.ESet(m.label, "+Inf")
	sample.EGet(m.ds).(ecore.EList).AddAll(ecore.NewImmutableEList([]any{math.Inf(1), math.Inf(-1), math.NaN(), 2.5}))
	eResource := m.newResource()
	eResource.GetContents().Add(sample)

	document := encodeResource(t, eResource, nil)
	assert.Contains(t, document, `d1="INF"`)
	assert.Contains(t, document, `d2="-INF"`)
	assert.Contains(t, document, `d3="NaN"`)
	assert.Contains(t, document, `f1="INF"`)
	assert.Contains(t, document, `finite="1.5"`)
	assert.Contains(t, document, `label="+Inf"`, "a string is not a float")
	assert.Contains(t, document, "<ds>INF</ds>")
	assert.Contains(t, document, "<ds>-INF</ds>")
	assert.Contains(t, document, "<ds>NaN</ds>")
	assert.Contains(t, document, "<ds>2.5</ds>")
	assert.NotContains(t, document, "Inf<")

	decoded := decodeResource(t, m.newResource(), document, nil)
	assert.True(t, math.IsInf(decoded.EGet(m.d1).(float64), 1))
	assert.True(t, math.IsInf(decoded.EGet(m.d2).(float64), -1))
	assert.True(t, math.IsNaN(decoded.EGet(m.d3).(float64)))
	assert.True(t, math.IsInf(float64(decoded.EGet(m.f1).(float32)), 1))
	assert.Equal(t, 1.5, decoded.EGet(m.finite))
	assert.Equal(t, "+Inf", decoded.EGet(m.label))
	ds := decoded.EGet(m.ds).(ecore.EList)
	require.Equal(t, 4, ds.Size())
	assert.True(t, math.IsInf(ds.Get(0).(float64), 1))
	assert.True(t, math.IsInf(ds.Get(1).(float64), -1))
	assert.True(t, math.IsNaN(ds.Get(2).(float64)))
	assert.Equal(t, 2.5, ds.Get(3))
}

// A document written before the XSD forms, with the Go forms +Inf and -Inf, is still read.
func TestXMLDecoder_GoSpecialFloats(t *testing.T) {
	m := newFloatModel()
	document := `<?xml version="1.0" encoding="UTF-8"?>
<floats:Sample xmlns:floats="http://floats" d1="+Inf" d2="-Inf" d3="NaN"/>`
	decoded := decodeResource(t, m.newResource(), document, nil)
	assert.True(t, math.IsInf(decoded.EGet(m.d1).(float64), 1))
	assert.True(t, math.IsInf(decoded.EGet(m.d2).(float64), -1))
	assert.True(t, math.IsNaN(decoded.EGet(m.d3).(float64)))
}

func TestConvertToString(t *testing.T) {
	eDouble := ecore.GetPackage().GetEDouble()
	eFloat := ecore.GetPackage().GetEFloat()
	eString := ecore.GetPackage().GetEString()
	for _, test := range []struct {
		name      string
		eDataType ecore.EDataType
		value     any
		expected  string
	}{
		{"double +inf", eDouble, math.Inf(1), "INF"},
		{"double -inf", eDouble, math.Inf(-1), "-INF"},
		{"double nan", eDouble, math.NaN(), "NaN"},
		{"double finite", eDouble, -0.25, "-0.25"},
		{"float +inf", eFloat, float32(math.Inf(1)), "INF"},
		{"float -inf", eFloat, float32(math.Inf(-1)), "-INF"},
		{"float nan", eFloat, float32(math.NaN()), "NaN"},
		{"float finite", eFloat, float32(0.5), "0.5"},
		{"string", eString, "+Inf", "+Inf"},
	} {
		t.Run(test.name, func(t *testing.T) {
			eFactory := test.eDataType.GetEPackage().GetEFactoryInstance()
			assert.Equal(t, test.expected, convertToString(eFactory, test.eDataType, test.value))
		})
	}
}
