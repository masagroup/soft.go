package ecore

func diagnosticError(errors EList) string {
	if errors.Empty() {
		return ""
	} else {
		return errors.Get(0).(EDiagnostic).GetMessage()
	}
}
