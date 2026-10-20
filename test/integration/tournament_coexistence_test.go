// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

// The tests of tournament_coexistence_test.go check that two implementations of the tournament
// model (typed fields in test/tournament/typed, reflective delegation in test/tournament/reflective)
// share its package (test/tournament): one EPackage, the same EClass for the objects of both, each
// implementation with a factory of its own chosen per resource set, and objects that do not mix.

package integration

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/masagroup/soft.go/codec/binary"
	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
	"github.com/masagroup/soft.go/test/tournament"
	"github.com/masagroup/soft.go/test/tournament/reflective"
	"github.com/masagroup/soft.go/test/tournament/typed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tournamentXML = "testdata/tournament/tournament.xml"

// newTournamentResourceSet returns a resource set creating the objects of the tournament model with factory.
func newTournamentResourceSet(factory ecore.EFactory) ecore.EResourceSet {
	resourceSet := ecore.CreateEResourceSet([]ecore.EPackage{tournament.GetPackage()})
	resourceSet.GetPackageRegistry().RegisterFactory(factory)
	return resourceSet
}

func loadXML(t *testing.T, factory ecore.EFactory) ecore.EResource {
	t.Helper()
	xmlProcessor := xml.NewProcessor(xml.ProcessorResourceSet(newTournamentResourceSet(factory)))
	resource := xmlProcessor.Load(ecore.CreateFileURI(tournamentXML))
	require.True(t, resource.IsLoaded())
	require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
	require.Equal(t, 1, resource.GetContents().Size())
	return resource
}

func allContents(resource ecore.EResource) []ecore.EObject {
	objects := []ecore.EObject{}
	for it := resource.GetAllContents(); it.HasNext(); {
		objects = append(objects, it.Next().(ecore.EObject))
	}
	return objects
}

// requireImplementation checks that every object of resource is an object of the implementation of
// factory, of the shared class of the reference object at the same place, with the same name.
func requireImplementation(t *testing.T, factory ecore.EFactory, reference ecore.EResource, resource ecore.EResource) {
	t.Helper()
	nameFeature := tournament.GetPackage().GetNamedElement_Name()
	expected := allContents(reference)
	actual := allContents(resource)
	require.Equal(t, len(expected), len(actual))
	for i, object := range actual {
		require.Same(t, expected[i].EClass(), object.EClass())
		require.Same(t, factory, ecore.GetObjectFactory(object))
		if nameFeature.GetEContainingClass().IsSuperTypeOf(object.EClass()) {
			require.Equal(t, expected[i].EGet(nameFeature), object.EGet(nameFeature))
		}
	}
}

func TestSharedPackage(t *testing.T) {
	pack := tournament.GetPackage()
	require.Same(t, pack, typed.GetFactory().GetEPackage())
	require.Same(t, pack, reflective.GetFactory().GetEPackage())
	require.NotSame(t, typed.GetFactory(), reflective.GetFactory())
	// the factory instance of the package converts data types and creates no object
	assert.Same(t, pack, pack.GetEFactoryInstance().GetEPackage())
	assert.NotEqual(t, typed.GetFactory(), pack.GetEFactoryInstance())
	assert.NotEqual(t, reflective.GetFactory(), pack.GetEFactoryInstance())
	assert.Panics(t, func() { pack.GetEFactoryInstance().Create(pack.GetTeam()) })
	assert.Equal(t, tournament.SEMI_FINAL, pack.GetEFactoryInstance().CreateFromString(pack.GetMatchKind(), "SemiFinal"))
	assert.Equal(t, "Final", pack.GetEFactoryInstance().ConvertToString(pack.GetMatchKind(), tournament.FINAL))

	typedTeam := typed.GetFactory().CreateTeam()
	reflectiveTeam := reflective.GetFactory().CreateTeam()
	assert.Equal(t, "*typed.teamImpl", fmt.Sprintf("%T", typedTeam))
	assert.Equal(t, "*reflective.teamImpl", fmt.Sprintf("%T", reflectiveTeam))
	// the same class
	assert.Same(t, pack.GetTeam(), typedTeam.EClass())
	assert.Same(t, pack.GetTeam(), reflectiveTeam.EClass())
	assert.True(t, pack.GetNamedElement().IsSuperTypeOf(typedTeam.EClass()))
	assert.True(t, pack.GetNamedElement().IsSuperTypeOf(reflectiveTeam.EClass()))
	// the factory of their implementation
	assert.Same(t, typed.GetFactory(), ecore.GetObjectFactory(typedTeam))
	assert.Same(t, reflective.GetFactory(), ecore.GetObjectFactory(reflectiveTeam))
	// created by class
	assert.Equal(t, "*typed.teamImpl", fmt.Sprintf("%T", typed.GetFactory().Create(pack.GetTeam())))
	assert.Equal(t, "*reflective.teamImpl", fmt.Sprintf("%T", reflective.GetFactory().Create(pack.GetTeam())))

	// global registry stays neutral
	assert.Nil(t, ecore.GetPackageRegistry().GetPackage(tournament.NS_URI))
}

func TestSharedPackageReflectiveAccess(t *testing.T) {
	nameFeature := tournament.GetPackage().GetNamedElement_Name()
	for _, team := range []tournament.Team{typed.GetFactory().CreateTeam(), reflective.GetFactory().CreateTeam()} {
		team.ESet(nameFeature, "team")
		assert.Equal(t, "team", team.GetName())
		assert.Equal(t, "team", team.EGet(nameFeature))
		assert.True(t, team.EIsSet(nameFeature))
		team.EUnset(nameFeature)
		assert.False(t, team.EIsSet(nameFeature))
	}
}

func TestSharedPackageLoad(t *testing.T) {
	typedResource := loadXML(t, typed.GetFactory())
	reflectiveResource := loadXML(t, reflective.GetFactory())

	typedRoot := typedResource.GetContents().Get(0).(tournament.Tournament)
	reflectiveRoot := reflectiveResource.GetContents().Get(0).(tournament.Tournament)
	assert.Equal(t, "*typed.tournamentImpl", fmt.Sprintf("%T", typedRoot))
	assert.Equal(t, "*reflective.tournamentImpl", fmt.Sprintf("%T", reflectiveRoot))

	requireImplementation(t, typed.GetFactory(), typedResource, typedResource)
	requireImplementation(t, reflective.GetFactory(), typedResource, reflectiveResource)

	// the same accessors on both sides
	assert.Equal(t, typedRoot.GetName(), reflectiveRoot.GetName())
	require.Equal(t, typedRoot.GetGroups().Size(), reflectiveRoot.GetGroups().Size())
	require.Equal(t, typedRoot.GetTeams().Size(), reflectiveRoot.GetTeams().Size())
	for i := 0; i < typedRoot.GetTeams().Size(); i++ {
		typedTeam := typedRoot.GetTeams().Get(i).(tournament.Team)
		reflectiveTeam := reflectiveRoot.GetTeams().Get(i).(tournament.Team)
		assert.Equal(t, typedTeam.GetName(), reflectiveTeam.GetName())
		if typedTeam.GetGroup() != nil {
			require.NotNil(t, reflectiveTeam.GetGroup())
			assert.Equal(t, typedTeam.GetGroup().GetName(), reflectiveTeam.GetGroup().GetName())
			assert.Same(t, reflective.GetFactory(), ecore.GetObjectFactory(reflectiveTeam.GetGroup()))
		}
	}
}

func TestSharedPackageLoadWithoutFactory(t *testing.T) {
	// a resource set with the shared package and no factory has no implementation to create objects with
	xmlProcessor := xml.NewProcessor(xml.ProcessorPackages([]ecore.EPackage{tournament.GetPackage()}))
	assert.Panics(t, func() { xmlProcessor.Load(ecore.CreateFileURI(tournamentXML)) })
}

func TestSharedPackageInstall(t *testing.T) {
	// Install of an implementation registers the shared package and the factory of the implementation
	for _, test := range []struct {
		install func(ecore.EPackageRegistry)
		factory ecore.EFactory
	}{
		{typed.Install, typed.GetFactory()},
		{reflective.Install, reflective.GetFactory()},
	} {
		resourceSet := ecore.NewEResourceSetImpl()
		test.install(resourceSet.GetPackageRegistry())
		assert.Same(t, tournament.GetPackage(), resourceSet.GetPackageRegistry().GetPackage(tournament.NS_URI))
		assert.Same(t, test.factory, resourceSet.GetPackageRegistry().GetFactory(tournament.NS_URI))

		xmlProcessor := xml.NewProcessor(xml.ProcessorResourceSet(resourceSet))
		resource := xmlProcessor.Load(ecore.CreateFileURI(tournamentXML))
		require.True(t, resource.GetErrors().Empty(), diagnosticError(resource.GetErrors()))
		requireImplementation(t, test.factory, resource, resource)
	}
	// no global default: the factory instance of the package creates no object
	assert.Panics(t, func() { tournament.GetPackage().GetEFactoryInstance().Create(tournament.GetPackage().GetTeam()) })
	// the global registry stays neutral
	assert.Nil(t, ecore.GetPackageRegistry().GetPackage(tournament.NS_URI))
}

func TestSharedPackageCopy(t *testing.T) {
	for _, factory := range []ecore.EFactory{typed.GetFactory(), reflective.GetFactory()} {
		resource := loadXML(t, factory)
		root := resource.GetContents().Get(0).(ecore.EObject)
		copied := ecore.Copy(root)
		require.NotNil(t, copied)
		assert.Same(t, root.EClass(), copied.EClass())
		assert.Same(t, factory, ecore.GetObjectFactory(copied))
		assert.True(t, ecore.Equals(root, copied))
	}
}

type roundTrip struct {
	name   string
	source ecore.EResource
	target ecore.EFactory
}

func roundTrips(t *testing.T) []roundTrip {
	return []roundTrip{
		{"typed to reflective", loadXML(t, typed.GetFactory()), reflective.GetFactory()},
		{"reflective to typed", loadXML(t, reflective.GetFactory()), typed.GetFactory()},
	}
}

func TestSharedPackageXMLRoundTrip(t *testing.T) {
	for _, test := range roundTrips(t) {
		t.Run(test.name, func(t *testing.T) {
			sourceProcessor := xml.NewProcessor(xml.ProcessorResourceSet(test.source.GetResourceSet()))
			saved := sourceProcessor.SaveToString(test.source, nil)
			require.NotEmpty(t, saved)

			targetProcessor := xml.NewProcessor(xml.ProcessorResourceSet(newTournamentResourceSet(test.target)))
			loaded := targetProcessor.LoadWithReader(strings.NewReader(saved), nil)
			require.True(t, loaded.IsLoaded())
			require.True(t, loaded.GetErrors().Empty(), diagnosticError(loaded.GetErrors()))
			requireImplementation(t, test.target, test.source, loaded)
			assert.Equal(t, saved, targetProcessor.SaveToString(loaded, nil))
		})
	}
}

func TestSharedPackageBinaryRoundTrip(t *testing.T) {
	for _, test := range roundTrips(t) {
		t.Run(test.name, func(t *testing.T) {
			var saved bytes.Buffer
			binary.NewEncoder(test.source, &saved, nil).EncodeResource()
			require.True(t, test.source.GetErrors().Empty(), diagnosticError(test.source.GetErrors()))
			require.NotZero(t, saved.Len())

			loaded := newTournamentResourceSet(test.target).CreateResource(ecore.NewURI("tournament.bin"))
			loaded.LoadWithReader(bytes.NewReader(saved.Bytes()), nil)
			require.True(t, loaded.IsLoaded())
			require.True(t, loaded.GetErrors().Empty(), diagnosticError(loaded.GetErrors()))
			requireImplementation(t, test.target, test.source, loaded)

			var resaved bytes.Buffer
			binary.NewEncoder(loaded, &resaved, nil).EncodeResource()
			assert.Equal(t, saved.Bytes(), resaved.Bytes())
		})
	}
}

func TestSharedPackageNoMixing(t *testing.T) {
	if !ecore.Checks {
		t.Skip("the implementations are checked with the build tag ecore_check")
	}
	const mixing = "cannot mix implementations of the package '" + tournament.NS_URI + "'"
	assertMixing := func(t *testing.T, f func()) {
		t.Helper()
		defer func() {
			r := recover()
			require.NotNil(t, r, "no panic")
			assert.Contains(t, fmt.Sprint(r), mixing)
		}()
		f()
	}
	factories := map[string]tournament.TournamentFactory{"typed": typed.GetFactory(), "reflective": reflective.GetFactory()}
	for name, factory := range factories {
		for otherName, other := range factories {
			if name == otherName {
				continue
			}
			t.Run(name+" with "+otherName, func(t *testing.T) {
				t.Run("containment", func(t *testing.T) {
					assertMixing(t, func() { factory.CreateTournament().GetTeams().Add(other.CreateTeam()) })
				})
				t.Run("reference list", func(t *testing.T) {
					assertMixing(t, func() { factory.CreateGroup().GetTeams().Add(other.CreateTeam()) })
				})
				t.Run("bidirectional reference", func(t *testing.T) {
					assertMixing(t, func() { factory.CreateTeam().SetGroup(other.CreateGroup()) })
				})
				t.Run("single reference", func(t *testing.T) {
					assertMixing(t, func() { factory.CreateMatch().SetHomeTeam(other.CreateTeam()) })
				})
				t.Run("reflective set", func(t *testing.T) {
					assertMixing(t, func() {
						factory.CreateMatch().ESet(tournament.GetPackage().GetMatch_GuestTeam(), other.CreateTeam())
					})
				})
				t.Run("resource roots", func(t *testing.T) {
					resource := newTournamentResourceSet(factory).CreateResource(ecore.NewURI("tournament.xml"))
					resource.GetContents().Add(factory.CreateTournament())
					assertMixing(t, func() { resource.GetContents().Add(other.CreateTournament()) })
				})
			})
		}
		t.Run(name+" alone", func(t *testing.T) {
			group := factory.CreateGroup()
			team := factory.CreateTeam()
			tournamentObject := factory.CreateTournament()
			tournamentObject.GetTeams().Add(team)
			tournamentObject.GetGroups().Add(group)
			team.SetGroup(group)
			match := factory.CreateMatch()
			match.SetHomeTeam(team)
			match.ESet(tournament.GetPackage().GetMatch_GuestTeam(), team)
			assert.Same(t, group, team.GetGroup())
			assert.True(t, group.GetTeams().Contains(team))
			assert.Same(t, team, match.GetGuestTeam())
			resource := newTournamentResourceSet(factory).CreateResource(ecore.NewURI("tournament.xml"))
			resource.GetContents().Add(tournamentObject)
			resource.GetContents().Add(factory.CreateTournament())
			assert.Equal(t, 2, resource.GetContents().Size())
		})
	}
}
