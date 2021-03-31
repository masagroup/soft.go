package ecore

type EDiagnostic interface {
	error
	GetMessage() string
	GetLocation() string
	GetLine() int
	GetColumn() int
}
