package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type eAllContentsPackage struct {
	ePackage                      EPackage
	eRootClass                    EClass
	eRootPartiesReference         EReference
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
	eRootPartiesReference, _ := eRootClass.GetEStructuralFeatureFromName("parties").(EReference)
	require.NotNil(t, eRootPartiesReference)
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
		eRootPartiesReference:         eRootPartiesReference,
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

	table := newTransitionTable(p.eRootClass, p.eUnitClass)
	require.NotNil(t, table)
	assert.Equal(t, 4, len(table))
	{
		source := state{eClass: p.eRootClass}
		target := state{eClass: p.ePartyClass}
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eRootPartiesReference}, transitions[0])
	}
	{
		source := state{eClass: p.ePartyClass}
		target := state{eClass: p.eFormationClass}
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.ePartyFormationsReference}, transitions[0])
	}
	{
		source := state{eClass: p.eFormationClass}
		targetAutomat := state{eClass: p.eAutomatClass}
		targetFormation := state{eClass: p.eFormationClass}
		transitions := table.getTransitions(source)
		require.Equal(t, 2, len(transitions))
		assert.Equal(t, &transition{source: source, target: targetFormation, reference: p.eFormationFormationsReference}, transitions[0])
		assert.Equal(t, &transition{source: source, target: targetAutomat, reference: p.eFormationAutomatsReference}, transitions[1])
	}
	{
		source := state{eClass: p.eAutomatClass}
		target := state{eClass: p.eUnitClass, isEnd: true}
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eAutomatUnitsReference}, transitions[0])
	}

}

func TestTransitionTable_Integration_Cycle(t *testing.T) {
	p := loadEAllContentsPackage(t)

	table := newTransitionTable(p.eRootClass, p.eFormationClass)
	require.NotNil(t, table)
	assert.Equal(t, 3, len(table))

	{
		source := state{eClass: p.eRootClass}
		target := state{eClass: p.ePartyClass}
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.eRootPartiesReference}, transitions[0])
	}
	{
		source := state{eClass: p.ePartyClass}
		target := state{eClass: p.eFormationClass, isEnd: true}
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: p.ePartyFormationsReference}, transitions[0])
	}
	{
		source := state{eClass: p.eFormationClass, isEnd: true}
		target := state{eClass: p.eFormationClass, isEnd: true}
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

	table := newTransitionTable(p.eRootClass, p.eUnitClass)
	require.NotNil(t, table)

	it := newEAllContentsWithClassIterator(m.eRoot, table)
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

	table := newTransitionTable(p.eRootClass, p.eFormationClass)
	require.NotNil(t, table)

	it := newEAllContentsWithClassIterator(m.eRoot, table)
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
