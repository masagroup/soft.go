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
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
)

// the root element of an object whose class is in a cycle of first super types, an invalid model,
// is searched in each class once
func TestXMLEncoder_RootFeatureSuperTypesCycle(t *testing.T) {
	m := newKindsModel()
	m.root.GetESuperTypes().Add(m.target)
	m.target.GetESuperTypes().Add(m.root)
	processor := NewProcessor(ProcessorPackages([]ecore.EPackage{m.pkg}))
	resource := processor.GetResourceSet().CreateResource(ecore.NewURI("root.xml"))
	resource.GetContents().Add(m.pkg.GetEFactoryInstance().Create(m.root))
	encoded := processor.SaveToString(resource, map[string]any{OPTION_EXTENDED_META_DATA: NewExtendedMetaData()})
	assert.Contains(t, encoded, "<kinds:Root")
}
