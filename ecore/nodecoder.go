package ecore

type NoDecoder struct {
}

func (de *NoDecoder) Decode() {

}

func (de *NoDecoder) DecodeObject() (EObject, error) {
	return nil, nil
}
