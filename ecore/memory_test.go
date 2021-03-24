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
	"testing"

	"github.com/OneOfOne/go-utils/memory"
)

func TestMemorySizes(t *testing.T) {
	t.Log(memory.Sizeof(GetPackage()))
	t.Log(memory.Sizeof(newEClassExt()))
	t.Log(memory.Sizeof(&CompactEObjectImpl{}))
	t.Log(memory.Sizeof(&CompactEObjectContainer{}))
	t.Log(memory.Sizeof(&BasicEObjectImpl{}))
}
