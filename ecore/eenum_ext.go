package ecore

type eEnumExt struct {
	eEnumImpl
}

func newEEnumExt() *eEnumExt {
	eEnum := new(eEnumExt)
	eEnum.SetInterfaces(eEnum)
	eEnum.Initialize()
	return eEnum
}

// GetEEnumLiteralByName default implementation
func (eEnum *eEnumExt) GetEEnumLiteralByName(name string) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetName() == name {
			return eLiteral
		}
	}
	return nil
}

// GetEEnumLiteralByValue default implementation
func (eEnum *eEnumExt) GetEEnumLiteralByValue(value int) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetValue() == value {
			return eLiteral
		}
	}
	return nil
}

// GetEEnumLiteralByLiteral default implementation
func (eEnum *eEnumExt) GetEEnumLiteralByLiteral(literal string) EEnumLiteral {
	for it := eEnum.GetELiterals().Iterator(); it.HasNext(); {
		eLiteral := it.Next().(EEnumLiteral)
		if eLiteral.GetLiteral() == literal {
			return eLiteral
		}
	}
	return nil
}
