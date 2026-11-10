// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import "sync/atomic"

type EEnumExt struct {
	EEnumImpl
	defaultValue atomic.Pointer[enumDefaultValue]
}

// enumDefaultValue is the default value of an enum, read by factory from its first literal, of
// literal name and value ordinal: kept as long as they do not change, as the default value is read
// at each access to an unset attribute of the enum. The value depends on nothing else: another
// first literal of the same literal and value gives the same.
type enumDefaultValue struct {
	factory EFactory
	name    string
	ordinal int
	value   any
}

func newEEnumExt() *EEnumExt {
	eEnum := new(EEnumExt)
	eEnum.SetInterfaces(eEnum)
	eEnum.Initialize()
	return eEnum
}

func (eEnum *EEnumExt) Initialize() {
	eEnum.EEnumImpl.Initialize()
	// created now, not at the first access: the literals are read without lock, at each literal
	// lookup, and a lazy creation would race between readers
	eEnum.GetELiterals()
}

// GetDefaultValue returns the first literal as the factory of the enum's package reads it: a value
// of the enum's Go type for a generated package, the literal's int value otherwise.
func (eEnum *EEnumExt) GetDefaultValue() any {
	eLiterals := eEnum.GetELiterals()
	if eLiterals.Empty() {
		return nil
	}
	eLiteral := eLiterals.Get(0).(EEnumLiteral)
	var eFactory EFactory
	if ePackage := eEnum.GetEPackage(); ePackage != nil {
		eFactory = ePackage.GetEFactoryInstance()
	}
	if eFactory == nil {
		return eLiteral.GetValue()
	}
	name, ordinal := eLiteral.GetLiteral(), eLiteral.GetValue()
	if current := eEnum.defaultValue.Load(); current != nil && current.factory == eFactory && current.name == name && current.ordinal == ordinal {
		return current.value
	}
	value := eFactory.CreateFromString(eEnum, name)
	eEnum.defaultValue.Store(&enumDefaultValue{factory: eFactory, name: name, ordinal: ordinal, value: value})
	return value
}

// GetEEnumLiteralByName default implementation
func (eEnum *EEnumExt) GetEEnumLiteralByName(name string) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetName() == name {
			return eLiteral
		}
	}
	return nil
}

// GetEEnumLiteralByValue default implementation
func (eEnum *EEnumExt) GetEEnumLiteralByValue(value int) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetValue() == value {
			return eLiteral
		}
	}
	return nil
}

// GetEEnumLiteralByLiteral default implementation
func (eEnum *EEnumExt) GetEEnumLiteralByLiteral(literal string) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetLiteral() == literal {
			return eLiteral
		}
	}
	return nil
}
