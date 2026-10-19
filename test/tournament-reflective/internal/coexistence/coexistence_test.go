// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

// Package coexistence checks that two complete generated implementations of the
// same metamodel (typed fields in test/tournament, reflective delegation in
// test/tournament-reflective/internal/impls) coexist, chosen per resource set.
package coexistence

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/masagroup/soft.go/test/tournament"
	reflective "github.com/masagroup/soft.go/test/tournament-reflective"
	"github.com/masagroup/soft.go/test/tournament-reflective/internal/impls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tournamentXML = "../../../tournament/testdata/tournament.xml"

func typedPackage() ecore.EPackage {
	return tournament.GetPackage()
}

func reflectivePackage() ecore.EPackage {
	return impls.GetPackage()
}

func diagnosticError(errors ecore.EList) string {
	if errors.Empty() {
		return ""
	}
	return errors.Get(0).(ecore.EDiagnostic).GetMessage()
}

func loadXML(t *testing.T, ePackage ecore.EPackage) (ecore.EResourceSet, ecore.EResource) {
	t.Helper()
	resourceSet := ecore.CreateEResourceSet([]ecore.EPackage{ePackage})
	xmlProcessor := ecore.NewXMLProcessor(ecore.XMLProcessorResourceSet(resourceSet))
	resource := xmlProcessor.Load(ecore.CreateFileURI(tournamentXML))
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.Equal(t, 1, resource.GetContents().Size())
	return resourceSet, resource
}

func allContents(resource ecore.EResource) []ecore.EObject {
	objects := []ecore.EObject{}
	for it := resource.GetAllContents(); it.HasNext(); {
		objects = append(objects, it.Next().(ecore.EObject))
	}
	return objects
}

// requireImplementation checks that every object of resource belongs to ePackage
// and matches the reference resource object by object (classifier ID, name by feature ID).
func requireImplementation(t *testing.T, ePackage ecore.EPackage, reference ecore.EResource, resource ecore.EResource) {
	t.Helper()
	expected := allContents(reference)
	actual := allContents(resource)
	require.Equal(t, len(expected), len(actual))
	for i, object := range actual {
		require.Same(t, ePackage, object.EClass().GetEPackage())
		require.Equal(t, expected[i].EClass().GetClassifierID(), object.EClass().GetClassifierID())
		require.Equal(t, expected[i].EClass().GetName(), object.EClass().GetName())
		if object.EClass().GetEStructuralFeatureFromName("name") != nil {
			require.Equal(t, expected[i].(ecore.EObjectInternal).EGetFromID(tournament.NAMED_ELEMENT__NAME, true), object.(ecore.EObjectInternal).EGetFromID(tournament.NAMED_ELEMENT__NAME, true))
		}
	}
}

func TestCoexistencePackages(t *testing.T) {
	typed := typedPackage()
	reflected := reflectivePackage()
	require.Equal(t, tournament.NS_URI, reflective.NS_URI)
	require.Equal(t, typed.GetNsURI(), reflected.GetNsURI())
	require.NotSame(t, typed, reflected)
	require.NotSame(t, typed.GetEFactoryInstance(), reflected.GetEFactoryInstance())

	require.Equal(t, typed.GetEClassifiers().Size(), reflected.GetEClassifiers().Size())
	for it := typed.GetEClassifiers().Iterator(); it.HasNext(); {
		typedClassifier := it.Next().(ecore.EClassifier)
		reflectedClassifier := reflected.GetEClassifier(typedClassifier.GetName())
		require.NotNil(t, reflectedClassifier, typedClassifier.GetName())
		require.NotSame(t, typedClassifier, reflectedClassifier)
		assert.Equal(t, typedClassifier.GetClassifierID(), reflectedClassifier.GetClassifierID(), typedClassifier.GetName())
		if typedClass, isClass := typedClassifier.(ecore.EClass); isClass {
			reflectedClass := reflectedClassifier.(ecore.EClass)
			require.Equal(t, typedClass.GetFeatureCount(), reflectedClass.GetFeatureCount(), typedClass.GetName())
			for itFeature := typedClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
				typedFeature := itFeature.Next().(ecore.EStructuralFeature)
				reflectedFeature := reflectedClass.GetEStructuralFeatureFromName(typedFeature.GetName())
				require.NotNil(t, reflectedFeature, typedFeature.GetName())
				featureID := typedClass.GetFeatureID(typedFeature)
				assert.Equal(t, featureID, reflectedClass.GetFeatureID(reflectedFeature), typedClass.GetName()+"."+typedFeature.GetName())
				assert.Same(t, reflectedFeature, reflectedClass.GetEStructuralFeature(featureID))
			}
		}
	}

	// global registry stays neutral
	assert.Nil(t, ecore.GetPackageRegistry().GetPackage(tournament.NS_URI))
}

func TestCoexistenceLoad(t *testing.T) {
	typedSet, typedResource := loadXML(t, typedPackage())
	reflectedSet, reflectedResource := loadXML(t, reflectivePackage())

	// each resource set resolves the nsURI to its own package and factory
	assert.Same(t, typedPackage(), typedSet.GetPackageRegistry().GetPackage(tournament.NS_URI))
	assert.Same(t, reflectivePackage(), reflectedSet.GetPackageRegistry().GetPackage(tournament.NS_URI))
	assert.Same(t, typedPackage().GetEFactoryInstance(), typedSet.GetPackageRegistry().GetFactory(tournament.NS_URI))
	assert.Same(t, reflectivePackage().GetEFactoryInstance(), reflectedSet.GetPackageRegistry().GetFactory(tournament.NS_URI))
	assert.Nil(t, ecore.GetPackageRegistry().GetPackage(tournament.NS_URI))

	typedRoot := typedResource.GetContents().Get(0)
	reflectedRoot := reflectedResource.GetContents().Get(0)
	assert.Equal(t, "*tournament.tournamentImpl", fmt.Sprintf("%T", typedRoot))
	assert.Equal(t, "*impls.tournamentImpl", fmt.Sprintf("%T", reflectedRoot))

	requireImplementation(t, typedPackage(), typedResource, typedResource)
	requireImplementation(t, reflectivePackage(), typedResource, reflectedResource)

	// typed accessors on each side
	typedTournament := typedRoot.(tournament.Tournament)
	reflectedTournament := reflectedRoot.(reflective.Tournament)
	assert.Equal(t, typedTournament.GetName(), reflectedTournament.GetName())
	require.Equal(t, typedTournament.GetGroups().Size(), reflectedTournament.GetGroups().Size())
	require.Equal(t, typedTournament.GetTeams().Size(), reflectedTournament.GetTeams().Size())
	for i := 0; i < typedTournament.GetTeams().Size(); i++ {
		typedTeam := typedTournament.GetTeams().Get(i).(tournament.Team)
		reflectedTeam := reflectedTournament.GetTeams().Get(i).(reflective.Team)
		assert.Equal(t, "*impls.teamImpl", fmt.Sprintf("%T", reflectedTeam))
		assert.Equal(t, typedTeam.GetName(), reflectedTeam.GetName())
		if typedTeam.GetGroup() != nil {
			require.NotNil(t, reflectedTeam.GetGroup())
			assert.Equal(t, typedTeam.GetGroup().GetName(), reflectedTeam.GetGroup().GetName())
			assert.Same(t, reflectivePackage(), reflectedTeam.GetGroup().EClass().GetEPackage())
		}
	}
}

func TestCoexistenceFactory(t *testing.T) {
	// cross Create by classifier ID: the result belongs to the factory's package
	reflectedTeam := impls.GetFactory().Create(tournament.GetPackage().GetTeam())
	require.NotNil(t, reflectedTeam)
	assert.Implements(t, (*reflective.Team)(nil), reflectedTeam)
	assert.Same(t, impls.GetPackage().GetTeam(), reflectedTeam.EClass())

	typedTeam := tournament.GetFactory().Create(impls.GetPackage().GetTeam())
	require.NotNil(t, typedTeam)
	assert.Implements(t, (*tournament.Team)(nil), typedTeam)
	assert.Same(t, tournament.GetPackage().GetTeam(), typedTeam.EClass())

	// cross CreateFromString / ConvertToString by classifier ID
	reflectedKind := impls.GetFactory().CreateFromString(tournament.GetPackage().GetMatchKind(), "SemiFinal")
	assert.Equal(t, reflective.SEMI_FINAL, reflectedKind)
	assert.Equal(t, "SemiFinal", impls.GetFactory().ConvertToString(tournament.GetPackage().GetMatchKind(), reflectedKind))
	typedKind := tournament.GetFactory().CreateFromString(impls.GetPackage().GetMatchKind(), "Final")
	assert.Equal(t, tournament.FINAL, typedKind)
	assert.Equal(t, "Final", tournament.GetFactory().ConvertToString(impls.GetPackage().GetMatchKind(), typedKind))
}

func TestCoexistenceReflectiveAccess(t *testing.T) {
	reflectedTeam := impls.GetFactory().CreateTeam()
	typedTeam := tournament.GetFactory().CreateTeam()

	// a feature of the other implementation is not a valid feature
	assert.PanicsWithValue(t, "The feature 'name' is not a valid feature", func() { reflectedTeam.ESet(tournament.GetPackage().GetNamedElement_Name(), "team") })
	assert.PanicsWithValue(t, "The feature 'name' is not a valid feature", func() { typedTeam.ESet(impls.GetPackage().GetNamedElement_Name(), "team") })
	assert.PanicsWithValue(t, "The feature 'name' is not a valid feature", func() { reflectedTeam.EGet(tournament.GetPackage().GetNamedElement_Name()) })

	// by ID works
	reflectedTeam.(ecore.EObjectInternal).ESetFromID(tournament.TEAM__NAME, "reflected")
	typedTeam.(ecore.EObjectInternal).ESetFromID(reflective.TEAM__NAME, "typed")
	assert.Equal(t, "reflected", reflectedTeam.GetName())
	assert.Equal(t, "typed", typedTeam.GetName())

	// the feature looked up by ID or name on obj.EClass() works
	nameFeature := tournament.GetPackage().GetNamedElement_Name()
	byID := reflectedTeam.EClass().GetEStructuralFeature(tournament.GetPackage().GetTeam().GetFeatureID(nameFeature))
	byName := reflectedTeam.EClass().GetEStructuralFeatureFromName(nameFeature.GetName())
	require.NotNil(t, byID)
	assert.Same(t, byID, byName)
	reflectedTeam.ESet(byID, "renamed")
	assert.Equal(t, "renamed", reflectedTeam.EGet(byName))
}

type roundTrip struct {
	name   string
	source ecore.EResource
	target ecore.EPackage
}

func roundTrips(t *testing.T) []roundTrip {
	_, typedResource := loadXML(t, typedPackage())
	_, reflectedResource := loadXML(t, reflectivePackage())
	return []roundTrip{
		{"typed to reflective", typedResource, reflectivePackage()},
		{"reflective to typed", reflectedResource, typedPackage()},
	}
}

func TestCoexistenceXMLRoundTrip(t *testing.T) {
	for _, test := range roundTrips(t) {
		t.Run(test.name, func(t *testing.T) {
			sourceProcessor := ecore.NewXMLProcessor(ecore.XMLProcessorResourceSet(test.source.GetResourceSet()))
			saved := sourceProcessor.SaveToString(test.source, nil)
			require.NotEmpty(t, saved)

			targetProcessor := ecore.NewXMLProcessor(ecore.XMLProcessorResourceSet(ecore.CreateEResourceSet([]ecore.EPackage{test.target})))
			loaded := targetProcessor.LoadWithReader(strings.NewReader(saved), nil)
			require.True(t, loaded.IsLoaded())
			require.True(t, loaded.GetErrors().Empty(), diagnosticError(loaded.GetErrors()))
			requireImplementation(t, test.target, test.source, loaded)
			assert.Equal(t, saved, targetProcessor.SaveToString(loaded, nil))
		})
	}
}

func TestCoexistenceBinaryRoundTrip(t *testing.T) {
	for _, test := range roundTrips(t) {
		t.Run(test.name, func(t *testing.T) {
			var saved bytes.Buffer
			ecore.NewBinaryEncoder(test.source, &saved, nil).EncodeResource()
			require.True(t, test.source.GetErrors().Empty(), diagnosticError(test.source.GetErrors()))
			require.NotZero(t, saved.Len())

			targetSet := ecore.CreateEResourceSet([]ecore.EPackage{test.target})
			loaded := targetSet.CreateResource(ecore.NewURI("tournament.bin"))
			loaded.LoadWithReader(bytes.NewReader(saved.Bytes()), nil)
			require.True(t, loaded.IsLoaded())
			require.True(t, loaded.GetErrors().Empty(), diagnosticError(loaded.GetErrors()))
			requireImplementation(t, test.target, test.source, loaded)

			var resaved bytes.Buffer
			ecore.NewBinaryEncoder(loaded, &resaved, nil).EncodeResource()
			assert.Equal(t, saved.Bytes(), resaved.Bytes())
		})
	}
}
