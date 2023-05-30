package ecore

import (
	"database/sql"
	"io"
	"path/filepath"

	_ "modernc.org/sqlite"
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

	db, err := sql.Open(e.driver, dbPath)
	if err != nil {
		return nil, err
	}

	_, err = e.db.Exec(`
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_mode = WAL;
	CREATE TABLE packages ( 
		packageID INTEGER PRIMARY KEY AUTOINCREMENT,
		uri TEXT,
	);
	CREATE TABLE classes (
		classID INTEGER PRIMARY KEY AUTOINCREMENT,
		packageID INTEGER,
		name TEXT,
		FOREIGN KEY(packageID) REFERENCES packages(packageID)
	);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (e *SQLEncoder) EncodeResource() {
	var err error
	e.db, err = e.createDB()
	if err != nil {
		e.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), e.resource.GetURI().String(), 0, 0))
		return
	}
	defer func() {
		_ = e.db.Close()
	}()

}

func (e *SQLEncoder) EncodeObject(object EObject) error {
	return nil
}
