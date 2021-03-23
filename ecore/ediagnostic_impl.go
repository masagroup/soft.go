package ecore

import "fmt"

type EDiagnosticImpl struct {
	message  string
	location string
	line     int
	column   int
}

func NewEDiagnosticImpl(message string, location string, line int, column int) *EDiagnosticImpl {
	return &EDiagnosticImpl{
		message:  message,
		location: location,
		line:     line,
		column:   column,
	}
}

func (d *EDiagnosticImpl) GetMessage() string {
	return d.message
}

func (d *EDiagnosticImpl) GetLocation() string {
	return d.location
}

func (d *EDiagnosticImpl) GetLine() int {
	return d.line
}

func (d *EDiagnosticImpl) GetColumn() int {
	return d.column
}

func (d *EDiagnosticImpl) Error() string {
	return fmt.Sprintf("%v(%v,%v):%v", d.location, d.line, d.column, d.message)
}
