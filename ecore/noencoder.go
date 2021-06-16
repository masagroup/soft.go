package ecore

type NoEncoder struct {
}

func (ne *NoEncoder) EncodeResource(resource EResource) {

}

func (ne *NoEncoder) EncodeObject(object EObject, context EResource) error {
	return nil
}
