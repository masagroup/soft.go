package ecore

type EResourceDriver interface {
	NewEncoder(options map[string]interface{})
	NewDecoder(options map[string]interface{})
}
