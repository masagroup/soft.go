package ecore

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	SQL_OPTION_ERROR_HANDLER      = "ERROR_HANDLER"      // error handler funtion
	SQL_OPTION_KEEP_DEFAULTS      = "KEEP_DEFAULTS"      // keep default values
	SQL_OPTION_CODEC_VERSION      = "CODEC_VERSION"      // codec version ( int64 )
	SQL_OPTION_OBJECT_ID          = "OBJECT_ID"          // object id column name ( string )
	SQL_OPTION_CONTAINER_ID       = "CONTAINER_ID"       // container id ( boolean ) if true, encode object's container and container feature id
	SQL_OPTION_SQL_ID_MANAGER     = "SQL_ID_MANAGER"     // SQL id manager
	SQL_OPTION_IN_MEMORY_DATABASE = "IN_MEMORY_DATABASE" // uses a memory database
	SQL_OPTION_SCHEDULED_QUERIES  = "SCHEDULED_QUERIES"  // sqlite scheduled queries if true
	SQL_OPTION_LOGGER             = "LOGGER"
)

type SQLCodec struct {
}

const sqlCodecVersion int64 = 1

func sqlTmpDB(prefix string) (string, error) {
	try := 0
	for {
		randBytes := make([]byte, 16)
		_, err := rand.Read(randBytes)
		if err != nil {
			return "", err
		}
		f := filepath.Join(os.TempDir(), prefix+"."+hex.EncodeToString(randBytes)+".sqlite")
		_, err = os.Stat(f)
		if os.IsExist(err) {
			if try++; try < 10000 {
				continue
			}
			return "", &fs.PathError{Op: "sqlTmpDB", Path: prefix, Err: fs.ErrExist}
		}
		return f, nil
	}
}

func (d *SQLCodec) NewEncoder(resource EResource, w io.Writer, options map[string]any) EEncoder {
	return NewSQLWriterEncoder(w, resource, options)
}
func (d *SQLCodec) NewDecoder(resource EResource, r io.Reader, options map[string]any) EDecoder {
	return NewSQLReaderDecoder(r, resource, options)
}

type sqlFeatureKind int

const (
	sfkTransient sqlFeatureKind = iota
	sfkFloat64
	sfkFloat32
	sfkInt
	sfkInt64
	sfkInt32
	sfkInt16
	sfkByte
	sfkBool
	sfkString
	sfkByteArray
	sfkEnum
	sfkDate
	sfkData
	sfkDataList
	sfkObject
	sfkObjectList
	sfkObjectReference
	sfkObjectReferenceList
)

func getSQLCodecFeatureKind(eFeature EStructuralFeature) sqlFeatureKind {
	if eFeature.IsTransient() {
		return sfkTransient
	} else if eReference, _ := eFeature.(EReference); eReference != nil {
		if eReference.IsContainment() {
			if eReference.IsMany() {
				return sfkObjectList
			} else {
				return sfkObject
			}
		}
		opposite := eReference.GetEOpposite()
		if opposite != nil && opposite.IsContainment() {
			return sfkTransient
		}
		if eReference.IsResolveProxies() {
			if eReference.IsMany() {
				return sfkObjectReferenceList
			} else {
				return sfkObjectReference
			}
		}
		if eReference.IsContainer() {
			return sfkTransient
		}
		if eReference.IsMany() {
			return sfkObjectList
		} else {
			return sfkObject
		}
	} else if eAttribute, _ := eFeature.(EAttribute); eAttribute != nil {
		if eAttribute.IsMany() {
			return sfkDataList
		} else {
			eDataType := eAttribute.GetEAttributeType()
			if eEnum, _ := eDataType.(EEnum); eEnum != nil {
				return sfkEnum
			}

			switch eDataType.GetInstanceTypeName() {
			case "float64", "java.lang.Double", "double":
				return sfkFloat64
			case "float32", "java.lang.Float", "float":
				return sfkFloat32
			case "int", "java.lang.Integer":
				return sfkInt
			case "int64", "java.lang.Long", "java.math.BigInteger", "long":
				return sfkInt64
			case "int32":
				return sfkInt32
			case "int16", "java.lang.Short", "short":
				return sfkInt16
			case "byte":
				return sfkByte
			case "bool", "java.lang.Boolean", "boolean":
				return sfkBool
			case "string", "java.lang.String":
				return sfkString
			case "[]byte", "java.util.ByteArray":
				return sfkByteArray
			case "*time/time.Time", "java.util.Date":
				return sfkDate
			}

			return sfkData
		}
	}
	return -1
}

type SQLCodecIDManager interface {
	SetPackageID(EPackage, int64)
	SetObjectID(EObject, int64)
	SetClassID(EClass, int64)
	SetEnumLiteralID(EEnumLiteral, int64)
}
