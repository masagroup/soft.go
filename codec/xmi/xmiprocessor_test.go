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

	"github.com/masagroup/soft.go/codec/xml"
	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXMIProcessor_New(t *testing.T) {
	resourceSet := ecore.NewEResourceSetImpl()
	processor := NewProcessor(xml.ProcessorResourceSet(resourceSet))
	require.NotNil(t, processor)
	assert.Same(t, resourceSet, processor.GetResourceSet())
}

func TestXMIProcessor_LoadWithReader(t *testing.T) {
	processor := NewProcessor()
	reader := strings.NewReader("")
	resource := processor.LoadWithReader(reader, map[string]any{xml.OPTION_EXTENDED_META_DATA: xml.NewExtendedMetaData()})
	require.NotNil(t, resource)
}

func TestXMIRegistration(t *testing.T) {
	for _, extension := range []string{"ecore", "xmi"} {
		assert.IsType(t, &Codec{}, ecore.GetCodecRegistry().GetCodec(ecore.NewURI("model."+extension)), extension)
	}
}
