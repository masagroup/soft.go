package ecore

type XMICodec struct {
}

func (d *XMICodec) NewEncoder(options map[string]interface{}) EResourceEncoder {
	return NewXMIEncoder(options)
}
func (d *XMICodec) NewDecoder(options map[string]interface{}) EResourceDecoder {
	return NewXMIDecoder(options)
}
