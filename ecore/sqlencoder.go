package ecore

import (
	"database/sql"
	"io"
	"path/filepath"
)

type SQLEncoder struct {
	resource EResource
	writer   io.Writer
	driver   string
	db       *sql.DB
}

func NewSQLEncoder(resource EResource, w io.Writer, options map[string]any) *SQLEncoder {
	e := &SQLEncoder{
		resource: resource,
		writer:   w,
		driver:   "sqlite",
	}
	if options != nil {
		if driver, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			e.driver = driver.(string)
		}
	}
	return e
}

func (e *SQLEncoder) createDB() (*sql.DB, error) {
	fileName := filepath.Base(e.resource.GetURI().Path())
	dbPath, err := sqlTmpDB(fileName)
	if err != nil {
		return nil, err
	}

	return sql.Open(e.driver, dbPath)
}

func (e *SQLEncoder) EncodeResource() {
	_, err := e.createDB()
	if err != nil {
		e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), e.resource.GetURI().String(), 0, 0))
		return
	}
}

func (e *SQLEncoder) EncodeObject(object EObject) error {
	return nil
}
