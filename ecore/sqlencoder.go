package ecore

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"

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
	insertObjectStms  map[EClass]*sql.Stmt
	idAttributeName   string
}

func NewSQLEncoder(resource EResource, w io.Writer, options map[string]any) *SQLEncoder {
	e := &SQLEncoder{
		resource:         resource,
		writer:           w,
		driver:           "sqlite",
		packageDataMap:   map[EPackage]*sqlEncoderPackageData{},
		classDataMap:     map[EClass]*sqlEncoderClassData{},
		insertObjectStms: map[EClass]*sql.Stmt{},
	}
	if options != nil {
		if driver, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			e.driver = driver.(string)
		}
		e.idAttributeName, _ = options[JSON_OPTION_ID_ATTRIBUTE_NAME].(string)
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
	eClass := eObject.EClass()
	_, err := e.encodeClass(eClass)
	if err != nil {
		return err
	}

	idManager := e.resource.GetObjectIDManager()
	if len(e.idAttributeName) == 0 {
		idManager = nil
	}

	// create table
	insertObjectStmt, isObjectInsertStmt := e.insertObjectStms[eClass]
	if !isObjectInsertStmt {
		// create table
		{
			var query strings.Builder
			query.WriteString("CREATE TABLE ")
			query.WriteString(eClass.GetName())
			query.WriteString("( objectID INTEGER PRIMARY KEY AUTOINCREMENT")
			if idManager != nil {
				query.WriteString(",")
				query.WriteString(e.idAttributeName)
			}
			query.WriteString(" )")

			if _, err := e.db.Exec(query.String()); err != nil {
				return err
			}
		}
		// create stmt
		{
			var query strings.Builder
			query.WriteString("INSERT INTO ")
			query.WriteString(eClass.GetName())
			if idManager != nil {
				query.WriteString("( ")
				query.WriteString(e.idAttributeName)
				query.WriteString(")")
				query.WriteString("VALUES (?)")
			}

			insertObjectStmt, err = e.db.Prepare(`INSERT packageID,name VALUES (?,?)`)
			if err != nil {
				return err
			}

		}
	}

	// insert object in table
	args := []any{}
	if idManager != nil {
		args = append(args, idManager.GetID(eObject))
	}

	sqlResult, err := insertObjectStmt.Exec(args...)
	if err != nil {
		return err
	}

	// retrieve object id
	objectID, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	// features
	for itFeature := eClass.GetEAllStructuralFeatures().Iterator(); itFeature.HasNext(); {
		eFeature := itFeature.Next().(EStructuralFeature)
		e.encodeFeatureValue(eObject, objectID, eFeature)
	}

	return nil
}

func (e *SQLEncoder) encodeObjectReference(eObject EObject) error {
	// encode object class
	_, err := e.encodeClass(eObject.EClass())
	if err != nil {
		return err
	}

	return nil
}

func (e *SQLEncoder) encodeFeatureValue(eObject EObject, objectID int64, eFeature EStructuralFeature) {
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
