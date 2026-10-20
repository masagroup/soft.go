// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xmlcodec

import (
	"testing"

	"github.com/masagroup/soft.go/ecore"
	"github.com/stretchr/testify/mock"
)

func TestMockXmlProcessorOption_Apply(t *testing.T) {
	mockOption := newMockXmlProcessorOption(t)
	mockRun := ecore.NewMockRun(t, mock.AnythingOfType("*xmlcodec.XMLProcessor"))
	mockOption.EXPECT().apply(mock.AnythingOfType("*xmlcodec.XMLProcessor")).Return().Run(func(_a0 *XMLProcessor) { mockRun.Run(_a0) })
	mockOption.apply(nil)
}
