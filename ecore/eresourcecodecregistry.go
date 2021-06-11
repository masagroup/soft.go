package ecore

const (
	DEFAULT_EXTENSION = "*"
)

type EResourceCodecRegistry interface {
	GetCodec(uri *URI) EResourceCodec
	GetProtocolToCodecMap() map[string]EResourceCodec
	GetExtensionToCodecMap() map[string]EResourceCodec
}

var resourceCodecRegistryInstance EResourceCodecRegistry

func GetResourceCodecRegistry() EResourceCodecRegistry {
	if resourceCodecRegistryInstance == nil {
		resourceCodecRegistryInstance = NewEResourceCodecRegistryImpl()
		// initialize with default factories
		extensionToCodecs := resourceCodecRegistryInstance.GetExtensionToCodecMap()
		extensionToCodecs["ecore"] = &XMICodec{}
		extensionToCodecs["xml"] = &XMLCodec{}
	}
	return resourceCodecRegistryInstance
}
