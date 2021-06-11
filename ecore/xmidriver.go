package ecore

type XMIDriver struct {
}

func (d *XMIDriver) NewEncoder(options map[string]interface{}) EResourceEncoder {
	return NewXMIEncoder(options)
}
func (d *XMIDriver) NewDecoder(options map[string]interface{}) EResourceDecoder {
	return NewXMIDecoder(options)
}
