package ecore

type BinaryCodec struct {
}

func (bc *BinaryCodec) NewEncoder(options map[string]interface{}) EResourceEncoder {
	return NewBinaryEncoder(options)
}
func (bc *BinaryCodec) NewDecoder(options map[string]interface{}) EResourceDecoder {
	return NewBinaryDecoder(options)
}
