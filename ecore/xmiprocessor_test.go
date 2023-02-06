// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestXMIProcessor_New(t *testing.T) {
	mockOption := newMockXmlProcessorOption(t)
	mockOption.EXPECT().apply(mock.AnythingOfType("*ecore.XMLProcessor")).Once()
	processor := NewXMIProcessor(mockOption)
	require.NotNil(t, processor)
}

func TestXMIProcessor_LoadWithReader(t *testing.T) {
	processor := NewXMIProcessor()
	reader := strings.NewReader("")
	resource := processor.LoadWithReader(reader, map[string]any{XML_OPTION_EXTENDED_META_DATA: NewExtendedMetaData()})
	require.NotNil(t, resource)
}
