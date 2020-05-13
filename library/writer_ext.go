// *****************************************************************************
//
// This file is part of a MASA library or program.
// Refer to the included end-user license agreement for restrictions.
//
// Copyright (c) 2020 MASA Group
//
// *****************************************************************************

package library

import (
	"strings"
)

type writerExt struct {
	*writerImpl
}

func newWriterExt() *writerExt {
	w := new(writerExt)
	w.writerImpl = newWriterImpl()
	w.SetInterfaces(w)
	return w
}

// GetName get the value of name
func (writer *writerExt) GetName() string {
	return writer.GetFirstName() + "--" + writer.GetLastName()
}

// SetName set the value of name
func (writer *writerExt) SetName(newName string) {
	index := strings.Index(newName, "--")
	if index != -1 {
		writer.SetFirstName(newName[:index])
		writer.SetLastName(newName[index+2:])
	}
}
