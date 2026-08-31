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
	"strings"

	"github.com/masagroup/soft.go/ecore"
)

// xmlWhitespace is XML's S production - the characters XSD's whiteSpace facet acts on.
const xmlWhitespace = " \t\r\n"

// xsdCollapsesWhitespace reports whether XSD gives a datatype whiteSpace="collapse": every built-in
// type but xs:string and its derivatives, judged on the Go instance type since Ecore's data types
// carry no XSD annotation. An enumeration collapses too.
func xsdCollapsesWhitespace(eDataType ecore.EDataType) bool {
	if _, isEnum := eDataType.(ecore.EEnum); isEnum {
		return true
	}
	switch instanceTypeName(eDataType) {
	case "float64", "java.lang.Double", "double",
		"float32", "java.lang.Float", "float",
		"int", "java.lang.Integer",
		"uint64", "com.google.common.primitives.UnsignedLong",
		"int64", "java.lang.Long", "java.math.BigInteger", "long",
		"int32", "int16", "java.lang.Short", "short", "int8",
		"uint", "uint32", "uint16", "uint8",
		"byte", "java.lang.Byte",
		"bool", "java.lang.Boolean", "boolean",
		"java.math.BigDecimal",
		"time.Time", "*time.Time", "java.util.Date":
		return true
	}
	return false
}

// xsdCollapse applies that facet to a literal on its way into a factory: the factories answer zero
// on a parse error, so y1=" 10" was read as 0. Only the lexical form changes.
func xsdCollapse(eDataType ecore.EDataType, literalValue string) string {
	if !xsdCollapsesWhitespace(eDataType) {
		return literalValue
	}
	return strings.Trim(literalValue, xmlWhitespace)
}

// instanceTypeName is the Go type of a datatype: its GenGo instanceTypeName annotation, or its
// instance type name, as ecore's own conversions read it.
func instanceTypeName(eDataType ecore.EDataType) string {
	if eAnnotation := eDataType.GetEAnnotation("http://net.masagroup/soft/2019/GenGo"); eAnnotation != nil {
		if name, _ := eAnnotation.GetDetails().GetValue("instanceTypeName").(string); name != "" {
			return name
		}
	}
	return eDataType.GetInstanceTypeName()
}
