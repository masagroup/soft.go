// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package sqlite

import (
	"path/filepath"
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnumLiteralByStoredLiteral(t *testing.T) {
	eEnum := ecore.GetFactory().CreateEEnum()
	empty := ecore.GetFactory().CreateEEnumLiteralFromContainer(eEnum)
	empty.SetName("_")
	empty.SetLiteral("")
	named := ecore.GetFactory().CreateEEnumLiteralFromContainer(eEnum)
	named.SetName("Named")
	named.SetLiteral("named")
	unset := ecore.GetFactory().CreateEEnumLiteralFromContainer(eEnum)
	unset.SetName("Unset")

	assert.Equal(t, empty, enumLiteralByStoredLiteral(eEnum, ""))
	assert.Equal(t, empty, enumLiteralByStoredLiteral(eEnum, "_"))
	assert.Equal(t, named, enumLiteralByStoredLiteral(eEnum, "named"))
	assert.Nil(t, enumLiteralByStoredLiteral(eEnum, "Named"))
	assert.Equal(t, unset, enumLiteralByStoredLiteral(eEnum, "Unset"))
	assert.Nil(t, enumLiteralByStoredLiteral(eEnum, "missing"))

	assert.Equal(t, []any{"", "_"}, storedEnumLiterals(empty))
	assert.Equal(t, []any{"named"}, storedEnumLiterals(named))
	assert.Equal(t, []any{"Unset"}, storedEnumLiterals(unset))
}

// a literal set to "" was stored as its name: the stored name still decodes
// to the literal, and a lookup of the literal finds the stored row
func TestSQLStore_Enum_EmptyLiteralStoredAsName(t *testing.T) {
	ePackage := loadPackage("library.complex.ecore")
	require.NotNil(t, ePackage)
	eClass, _ := ePackage.GetEClassifier("Book").(ecore.EClass)
	require.NotNil(t, eClass)
	eFeature := eClass.GetEStructuralFeatureFromName("category")
	require.NotNil(t, eFeature)
	eEnum, _ := eFeature.GetEType().(ecore.EEnum)
	require.NotNil(t, eEnum)
	biography := eEnum.GetEEnumLiteralByName("Biography")
	require.NotNil(t, biography)
	biography.SetLiteral("")

	dbPath := filepath.Join(t.TempDir(), "library.store.sqlite")
	require.NoError(t, copyFile("testdata/library.store.sqlite", dbPath))
	mockPackageRegitry := newSQLMockEPackageRegistry(t)
	s, err := NewStore(dbPath, ecore.NewURI(""), nil, mockPackageRegitry, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	defer s.Close()

	// lookup : the row stored as "Biography"
	categorySchema, err := s.getFeatureSchema(newLendable(t, eClass, 6), eFeature)
	require.NoError(t, err)
	encoded, isEncoded, err := s.encodeLookupValue(categorySchema, 2)
	require.NoError(t, err)
	require.True(t, isEncoded)
	require.Equal(t, int64(1), encoded)

	// decoding : Biography == 2
	mockObject := NewMockObject(t)
	mockObject.EXPECT().GetSQLID().Return(int64(6)).Once()
	mockObject.EXPECT().SetSQLID(int64(6)).Once()
	mockObject.EXPECT().EClass().Return(eClass).Once()
	mockPackageRegitry.EXPECT().GetPackage(libraryNSURI).Return(ePackage).Maybe()
	assert.Equal(t, 2, s.Get(mockObject, eFeature, -1))
}
