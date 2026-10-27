// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

type EEnumExt struct {
	EEnumImpl
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
	if ePackage := eEnum.GetEPackage(); ePackage != nil {
		if eFactory := ePackage.GetEFactoryInstance(); eFactory != nil {
			return eFactory.CreateFromString(eEnum, eLiteral.GetLiteral())
		}
	}
	return eLiteral.GetValue()
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
