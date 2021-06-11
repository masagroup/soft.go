package ecore

type XMLDriver struct {
}

func (d *XMLDriver) NewEncoder(options map[string]interface{}) EResourceEncoder {
	return NewXMLEncoder(options)
}
func (d *XMLDriver) NewDecoder(options map[string]interface{}) EResourceDecoder {
	return NewXMLDecoder(options)
}
