// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2019 MASA Group
//
// *****************************************************************************

package ecore

import (
	"time"
)

type ecoreFactoryExt struct {
	ecoreFactoryImpl
}

func newEcoreFactoryExt() *ecoreFactoryExt {
	factory := new(ecoreFactoryExt)
	factory.SetInterfaces(factory)
	factory.Initialize()
	return factory
}

func (factory *ecoreFactoryExt) createEDateFromString(dataType EDataType, literalValue string) interface{} {
	t, _ := time.Parse(dateFormat, literalValue)
	return &t
}

const (
	dateFormat string = "2006-01-02T15:04:05.999Z"
)

func (factory *ecoreFactoryExt) convertEDateToString(dataType EDataType, instanceValue interface{}) string {
	t := instanceValue.(*time.Time)
	return t.Format(dateFormat)
}
