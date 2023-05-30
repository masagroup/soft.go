package ecore

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type sqlEncoderFeatureData struct {
}

type sqlEncoderClassData struct {
	id int64
}

type sqlEncoderPackageData struct {
	id        int64
	classData []*sqlEncoderClassData
}

type SQLEncoder struct {
	resource          EResource
	writer            io.Writer
	driver            string
	db                *sql.DB
	packageDataMap    map[EPackage]*sqlEncoderPackageData
	classDataMap      map[EClass]*sqlEncoderClassData
	insertPackageStmt *sql.Stmt
	insertClassStmt   *sql.Stmt
}

func NewSQLEncoder(resource EResource, w io.Writer, options map[string]any) *SQLEncoder {
	e := &SQLEncoder{
		resource:       resource,
		writer:         w,
		driver:         "sqlite",
		packageDataMap: map[EPackage]*sqlEncoderPackageData{},
		classDataMap:   map[EClass]*sqlEncoderClassData{},
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

	// open db
	db, err := sql.Open(e.driver, dbPath)
	if err != nil {
		return nil, err
	}

	// version info
	version := fmt.Sprintf(`PRAGMA user_version = %v`, sqlCodecVersion)
	_, err = db.Exec(version)
	if err != nil {
		return nil, err
	}

	// properties infos
	properties := `
	PRAGMA synchronous = NORMAL;
	PRAGMA journal_mode = WAL;
	`
	_, err = db.Exec(properties)
	if err != nil {
		return nil, err
	}

	// common tables
	tables := `
	CREATE TABLE packages ( 
		packageID INTEGER PRIMARY KEY AUTOINCREMENT,
		namespace TEXT,
		uri TEXT
	);
	CREATE TABLE classes (
		classID INTEGER PRIMARY KEY AUTOINCREMENT,
		packageID INTEGER,
		name TEXT,
		FOREIGN KEY(packageID) REFERENCES packages(packageID)
	);
	`
	_, err = db.Exec(tables)
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

func (e *SQLEncoder) encodeObject(eObject EObject) error {
	// encode object class
	_, err := e.encodeClass(eObject.EClass())
	if err != nil {
		return err
	}
	return nil
}

func (e *SQLEncoder) encodeClass(eClass EClass) (*sqlEncoderClassData, error) {
	eClassData := e.classDataMap[eClass]
	if eClassData == nil {
		// encode package
		packageData, err := e.encodePackage(eClass.GetEPackage())
		if err != nil {
			return nil, err
		}
		// create statement if needed
		if e.insertClassStmt == nil {
			stmt, err := e.db.Prepare(`INSERT packageID,name VALUES (?,?)`)
			if err != nil {
				return nil, err
			}
			e.insertClassStmt = stmt
		}
		// insert new class
		sqlResult, err := e.insertClassStmt.Exec(packageData.id, eClass.GetName())
		if err != nil {
			return nil, err
		}
		// retrieve index
		id, err := sqlResult.LastInsertId()
		if err != nil {
			return nil, err
		}
		// create data
		eClassData = &sqlEncoderClassData{
			id: id,
		}
		e.classDataMap[eClass] = eClassData

	}
	return eClassData, nil
}

func (e *SQLEncoder) encodePackage(ePackage EPackage) (*sqlEncoderPackageData, error) {
	ePackageData := e.packageDataMap[ePackage]
	if ePackageData == nil {
		// create statement if needed
		if e.insertPackageStmt == nil {
			stmt, err := e.db.Prepare(`INSERT namespace,uri VALUES (?,?)`)
			if err != nil {
				return nil, err
			}
			e.insertPackageStmt = stmt
		}
		// insert new package
		sqlResult, err := e.insertPackageStmt.Exec(ePackage.GetNsURI(), GetURI(ePackage))
		if err != nil {
			return nil, err
		}
		// retrieve package index
		id, err := sqlResult.LastInsertId()
		if err != nil {
			return nil, err
		}
		// create data
		ePackageData = &sqlEncoderPackageData{
			id:        id,
			classData: make([]*sqlEncoderClassData, ePackage.GetEClassifiers().Size()),
		}
		e.packageDataMap[ePackage] = ePackageData
	}
	return ePackageData, nil
}
