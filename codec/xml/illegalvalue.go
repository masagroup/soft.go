package xml

import "github.com/masagroup/soft.go/ecore"

// IllegalValueDiagnostic reports a value its data type does not accept, as EMF's
// IllegalValueException: the decoder leaves the feature unset and goes on, so a reader that tells
// a refused value from the other diagnostics can fail the load instead of losing the value.
type IllegalValueDiagnostic struct {
	*ecore.EDiagnosticImpl
	// Feature is the feature the value was read for.
	Feature ecore.EStructuralFeature
	// DataType is the type that refused the value.
	DataType ecore.EDataType
	// Value is the text read.
	Value string
}
