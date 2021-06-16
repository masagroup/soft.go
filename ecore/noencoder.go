package ecore

type NoEncoder struct {
}

func (ne *NoEncoder) Encode() {

}

func (ne *NoEncoder) EncodeObject(object EObject) error {
	return nil
}
