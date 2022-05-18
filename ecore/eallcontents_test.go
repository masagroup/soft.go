package ecore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransitionTable_Leaf(t *testing.T) {
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
	eAutomatClass, _ := ePackage.GetEClassifier("Automat").(EClass)
	require.NotNil(t, eAutomatClass)
	eAutomatUnitsReference, _ := eAutomatClass.GetEStructuralFeatureFromName("units").(EReference)
	require.NotNil(t, eAutomatUnitsReference)
	eUnitClass, _ := ePackage.GetEClassifier("Unit").(EClass)
	require.NotNil(t, eUnitClass)

	table := newTransitionTable(eRootClass, eUnitClass)
	require.NotNil(t, table)
	assert.Equal(t, 4, len(table))
	{
		source := state{stateType: start, eClass: eRootClass}
		target := state{stateType: active, eClass: ePartyClass}
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: eRootPartiesReference}, transitions[0])
	}
	{
		source := state{stateType: active, eClass: ePartyClass}
		target := state{stateType: active, eClass: eFormationClass}
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: ePartyFormationsReference}, transitions[0])
	}
	{
		source := state{stateType: active, eClass: eFormationClass}
		targetAutomat := state{stateType: active, eClass: eAutomatClass}
		targetFormation := state{stateType: active, eClass: eFormationClass}
		transitions := table.getTransitions(source)
		require.Equal(t, 2, len(transitions))
		assert.Equal(t, &transition{source: source, target: targetFormation, reference: eFormationFormationsReference}, transitions[0])
		assert.Equal(t, &transition{source: source, target: targetAutomat, reference: eFormationAutomatsReference}, transitions[1])
	}
	{
		source := state{stateType: active, eClass: eAutomatClass}
		target := state{stateType: end, eClass: eUnitClass}
		transitions := table.getTransitions(source)
		require.Equal(t, 1, len(transitions))
		assert.Equal(t, &transition{source: source, target: target, reference: eAutomatUnitsReference}, transitions[0])
	}

}

func TestEAllContentsWithClass_Leaf(t *testing.T) {
	// load package
	ePackage := loadPackage("eallcontents.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/eallcontents.xml"})
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	eRoot, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eRoot)

	eRootClass, _ := ePackage.GetEClassifier("Root").(EClass)
	require.NotNil(t, eRootClass)
	eUnitClass, _ := ePackage.GetEClassifier("Unit").(EClass)
	require.NotNil(t, eUnitClass)
	eUnitNameAttribute, _ := eUnitClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eUnitNameAttribute)

	table := newTransitionTable(eRootClass, eUnitClass)
	require.NotNil(t, table)

	it := newEAllContentsWithClassIterator(eRoot, table)
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
		require.Equal(t, eUnitClass, unit.EClass())
		name := unit.EGet(eUnitNameAttribute).(string)
		result = append(result, name)
	}
	assert.Equal(t, expected, result)
}

func TestEAllContentsWithClass_Cycle(t *testing.T) {
	// load package
	ePackage := loadPackage("eallcontents.ecore")
	require.NotNil(t, ePackage)

	// load resource
	xmlProcessor := NewXMLProcessor([]EPackage{ePackage})
	eResource := xmlProcessor.Load(&URI{Path: "testdata/eallcontents.xml"})
	require.NotNil(t, eResource)
	require.True(t, eResource.IsLoaded())
	require.True(t, eResource.GetErrors().Empty(), diagnosticError(eResource.GetErrors()))
	require.True(t, eResource.GetWarnings().Empty(), diagnosticError(eResource.GetWarnings()))
	eRoot, _ := eResource.GetContents().Get(0).(EObject)
	require.NotNil(t, eRoot)

	eRootClass, _ := ePackage.GetEClassifier("Root").(EClass)
	require.NotNil(t, eRootClass)
	eFormationClass, _ := ePackage.GetEClassifier("Formation").(EClass)
	require.NotNil(t, eFormationClass)
	eFormationNameAttribute, _ := eFormationClass.GetEStructuralFeatureFromName("name").(EAttribute)
	require.NotNil(t, eFormationNameAttribute)

	table := newTransitionTable(eRootClass, eFormationClass)
	require.NotNil(t, table)

	it := newEAllContentsWithClassIterator(eRoot, table)
	require.NotNil(t, it)
	result := []string{}
	expected := []string{}
	for it.HasNext() {
		formation, _ := it.Next().(EObject)
		require.NotNil(t, formation)
		require.Equal(t, eFormationClass, formation.EClass())
		name := formation.EGet(eFormationNameAttribute).(string)
		result = append(result, name)
	}
	assert.Equal(t, expected, result)
}
