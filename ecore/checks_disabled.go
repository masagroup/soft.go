//go:build !ecore_check

// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// Checks is false when the program is built without the tag ecore_check: the checks that only
// detect programming errors are compiled out (see checks_enabled.go).
const Checks = false
