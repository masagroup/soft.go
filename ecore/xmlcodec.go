package ecore

type XMLCodec struct {
}

func (d *XMLCodec) NewEncoder(options map[string]interface{}) EResourceEncoder {
	return NewXMLEncoder(options)
}
func (d *XMLCodec) NewDecoder(options map[string]interface{}) EResourceDecoder {
	return NewXMLDecoder(options)
}
