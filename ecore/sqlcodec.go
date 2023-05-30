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
	SQL_OPTION_DRIVER            = "DRIVER_NAME"       // value of the sql driver
	SQL_OPTION_ID_ATTRIBUTE_NAME = "ID_ATTRIBUTE_NAME" // value of the id attribute
)

type SQLCodec struct {
}

func sqlTmpDB(prefix string) (string, error) {
	try := 0
	for {
		randBytes := make([]byte, 16)
		rand.Read(randBytes)
		f := filepath.Join(os.TempDir(), prefix+"."+hex.EncodeToString(randBytes)+".sqlite")
		_, err := os.Stat(f)
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
	return NewSQLEncoder(resource, w, options)
}
func (d *SQLCodec) NewDecoder(resource EResource, r io.Reader, options map[string]any) EDecoder {
	return NewSQLDecoder(resource, r, options)
}
