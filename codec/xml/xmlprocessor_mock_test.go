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
	"github.com/stretchr/testify/mock"
)

func TestMockXmlProcessorOption_Apply(t *testing.T) {
	mockOption := newMockXmlProcessorOption(t)
	mockRun := ecore.NewMockRun(t, mock.AnythingOfType("*xml.Processor"))
	mockOption.EXPECT().apply(mock.AnythingOfType("*xml.Processor")).Return().Run(func(_a0 *Processor) { mockRun.Run(_a0) })
	mockOption.apply(nil)
}
