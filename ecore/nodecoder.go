package ecore

type NoDecoder struct {
}

func (de *NoDecoder) DecodeResource(resource EResource) {

}

func (de *NoDecoder) DecodeObject(object *EObject, resource EResource) error {
	return nil
}
