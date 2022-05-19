package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type eAllContentsPackage struct {
	ePackage                      EPackage
	eRootClass                    EClass
	eRootTheaterReference         EReference
	eTheaterClass                 EClass
	eTheaterPartiesReference      EReference
	ePartyClass                   EClass
	ePartyFormationsReference     EReference
	eFormationClass               EClass
	eFormationFormationsReference EReference
	eFormationAutomatsReference   EReference
	eFormationNameAttribute       EAttribute
	eAutomatClass                 EClass
	eAutomatUnitsReference        EReference
	eUnitClass                    EClass
	eUnitNameAttribute            EAttribute
}

func loadEAllContentsPackage(t *testing.T) *eAllContentsPackage {
	// load package
	ePackage := loadPackage("eallcontents.ecore")
	require.NotNil(t, ePackage)
	eRootClass, _ := ePackage.GetEClassifier("Root").(EClass)
	require.NotNil(t, eRootClass)
	eRootTheaterReference, _ := eRootClass.GetEStructuralFeatureFromName("theater").(EReference)
	require.NotNil(t, eRootTheaterReference)
	eTheaterClass, _ := ePackage.GetEClassifier("Theater").(EClass)
	require.NotNil(t, eTheaterClass)
	eTheaterPartiesReference, _ := eTheaterClass.GetEStructuralFeatureFromName("parties").(EReference)
	require.NotNil(t, eTheaterPartiesReference)
	ePartyClass, _ := ePackage.GetEClassifier("Party").(EClass)
	require.NotNil(t, ePartyClass)
	ePartyFormationsReference, _ := ePartyClass.GetEStructuralFeatureFromName("formations").(EReference)
	require.NotNil(t, ePartyFormationsReference)
	eFormationClass, _ := ePackage.GetEClassifier("Formation").(EClass)
	require.NotNil(t, eFormationClass)
	eFormationFormationsReference, _ := eFormationClass.GetEStructuralFeatureFromName("formations").(EReference)
	require.NotNil(t, eFormationFormationsReference)
	eFormationAutomatsReference, _ := eFormationClass.GetEStructuralFeatureFromName("automats").(EReference)
	require.NotNil(t, eFormationAutomatsReference)
	eFormationNameAttribute, _ := eFormationClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eFormationNameAttribute)
	eAutomatClass, _ := ePackage.GetEClassifier("Automat").(EClass)
	require.NotNil(t, eAutomatClass)
	eAutomatUnitsReference, _ := eAutomatClass.GetEStructuralFeatureFromName("units").(EReference)
	require.NotNil(t, eAutomatUnitsReference)
	eUnitClass, _ := ePackage.GetEClassifier("Unit").(EClass)
	require.NotNil(t, eUnitClass)
	eUnitNameAttribute, _ := eUnitClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eUnitNameAttribute)

	return &eAllContentsPackage{
		ePackage:                      ePackage,
		eRootClass:                    eRootClass,
		eRootTheaterReference:         eRootTheaterReference,
		eTheaterClass:                 eTheaterClass,
		eTheaterPartiesReference:      eTheaterPartiesReference,
		ePartyClass:                   ePartyClass,
		ePartyFormationsReference:     ePartyFormationsReference,
		eFormationClass:               eFormationClass,
		eFormationFormationsReference: eFormationFormationsReference,
		eFormationAutomatsReference:   eFormationAutomatsReference,
		eFormationNameAttribute:       eFormationNameAttribute,
		eAutomatClass:                 eAutomatClass,
		eAutomatUnitsReference:        eAutomatUnitsReference,
		eUnitClass:                    eUnitClass,
		eUnitNameAttribute:            eUnitNameAttribute,
	}
}

func TestTransitionTable_Integration_Leaf(t *testing.T) {
	p := loadEAllContentsPackage(t)

	table := NewEClassTransitionsTable(p.eRootClass, p.eUnitClass)
	require.NotNil(t, table)
	assert.True(t, table.isEnd(p.eUnitClass))
	{
		source := p.eRootClass
		target := p.eTheaterClass
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eRootTheaterReference}, transitions[0])
	}
	{
		source := p.eTheaterClass
		target := p.ePartyClass
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eTheaterPartiesReference}, transitions[0])
	}
	{
		source := p.ePartyClass
		target := p.eFormationClass
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.ePartyFormationsReference}, transitions[0])
	}
	{
		source := p.eFormationClass
		targetAutomat := p.eAutomatClass
		targetFormation := p.eFormationClass
		transitions := table.getTransitions(source)
		require.Equal(t, 2, len(transitions))
		assert.Equal(t, &transition{source: source, target: targetFormation, reference: p.eFormationFormationsReference}, transitions[0])
		assert.Equal(t, &transition{source: source, target: targetAutomat, reference: p.eFormationAutomatsReference}, transitions[1])
	}
	{
		source := p.eAutomatClass
		target := p.eUnitClass
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eAutomatUnitsReference}, transitions[0])
	}

}

func TestTransitionTable_Integration_Cycle(t *testing.T) {
	p := loadEAllContentsPackage(t)

	table := NewEClassTransitionsTable(p.eRootClass, p.eFormationClass)
	require.NotNil(t, table)
	assert.True(t, table.isEnd(p.eFormationClass))
	{
		source := p.eRootClass
		target := p.eTheaterClass
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eRootTheaterReference}, transitions[0])
	}
	{
		source := p.eTheaterClass
		target := p.ePartyClass
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eTheaterPartiesReference}, transitions[0])
	}
	{
		source := p.ePartyClass
		target := p.eFormationClass
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.ePartyFormationsReference}, transitions[0])
	}
	{
		source := p.eFormationClass
		target := p.eFormationClass
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eFormationFormationsReference}, transitions[0])
	}

}

type eAllContentsModel struct {
	eResource EResource
	eRoot     EObject
}

func loadEAllContentsModel(t *testing.T, ePackage EPackage) *eAllContentsModel {
	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/eallcontents.xml"})
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	require.Equal(t, 1, eResource.GetContents().Size())
	eRoot, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eRoot)
	return &eAllContentsModel{
		eResource: eResource,
		eRoot:     eRoot,
	}
}

func TestEAllContentsWithClass_Leaf(t *testing.T) {
	p := loadEAllContentsPackage(t)
	m := loadEAllContentsModel(t, p.ePackage)
	it := newEAllContentsWithClassIterator(m.eRoot, p.eUnitClass)
	require.NotNil(t, it)
	result := []string{}
	expected := []string{
		"unit 3",
		"unit 1",
		"unit 2",
		"unit 4",
		"unit 5",
	}
	for it.HasNext() {
		unit, _ := it.Next().(EObject)
		require.NotNil(t, unit)
		require.Equal(t, p.eUnitClass, unit.EClass())
		name := unit.EGet(p.eUnitNameAttribute).(string)
		result = append(result, name)
	}
	assert.Equal(t, expected, result)
}

func TestEAllContentsWithClass_Cycle(t *testing.T) {
	p := loadEAllContentsPackage(t)
	m := loadEAllContentsModel(t, p.ePackage)
	it := newEAllContentsWithClassIterator(m.eRoot, p.eFormationClass)
	require.NotNil(t, it)
	result := []string{}
	expected := []string{
		"formation 1",
		"formation 2",
		"formation 3",
		"formation 4",
	}
	for it.HasNext() {
		formation, _ := it.Next().(EObject)
		require.NotNil(t, formation)
		require.Equal(t, p.eFormationClass, formation.EClass())
		name := formation.EGet(p.eFormationNameAttribute).(string)
		result = append(result, name)
	}
	assert.Equal(t, expected, result)
}
