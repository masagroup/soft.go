//go:build ecore_check

// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

// Checks is true when the program is built with the tag ecore_check: the checks that only detect
// programming errors run, such as a unique list set to an element it holds at another index, or
// objects of two implementations of a shared package mixed (CheckSameImplementation). Without the
// tag, Checks is false and the code guarded by it is compiled out (see checks_disabled.go).
const Checks = true
