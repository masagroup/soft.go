package ecore

const (
	OPTION_EXTENDED_META_DATA        = "EXTENDED_META_DATA"        // ExtendedMetaData pointer
	OPTION_SUPPRESS_DOCUMENT_ROOT    = "SUPPRESS_DOCUMENT_ROOT"    // if true , suppress document root if found
	OPTION_IDREF_RESOLUTION_DEFERRED = "IDREF_RESOLUTION_DEFERRED" // if true , defer id ref resolution
	OPTION_ID_ATTRIBUTE_NAME         = "ID_ATTRIBUTE_NAME"         // value of the id attribute
	OPTION_ROOT_OBJECTS              = "ROOT_OBJECTS"              // list of root objects to save
)

type EResourceDriver interface {
	NewEncoder(options map[string]interface{}) EResourceEncoder
	NewDecoder(options map[string]interface{}) EResourceDecoder
}
