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

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewXmlProcessorWithOption(t *testing.T) {
	mockOption := newMockXmlProcessorOption(t)
	mockOption.EXPECT().apply(mock.AnythingOfType("*xml.Processor")).Once()
	p := NewProcessor(mockOption)
	require.NotNil(t, p)
}
