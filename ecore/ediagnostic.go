package ecore

type EDiagnostic interface {
	Error() string
	GetMessage() string
	GetLocation() string
	GetLine() int
	GetColumn() int
}
