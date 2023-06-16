package ecore

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type sqlDecoderClassData struct {
	ePackage EPackage
	eClass   EClass
}

type SQLDecoder struct {
	resource    EResource
	reader      io.Reader
	driver      string
	db          *sql.DB
	classesData map[int]*sqlDecoderClassData
}

func NewSQLDecoder(resource EResource, r io.Reader, options map[string]any) *SQLDecoder {
	d := &SQLDecoder{
		resource:    resource,
		reader:      r,
		driver:      "sqlite",
		classesData: map[int]*sqlDecoderClassData{},
	}
	if options != nil {
		if driver, isDriver := options[SQL_OPTION_DRIVER]; isDriver {
			d.driver = driver.(string)
		}
	}
	return d
}

func (d *SQLDecoder) createDB() (*sql.DB, error) {
	fileName := filepath.Base(d.resource.GetURI().Path())
	dbPath, err := sqlTmpDB(fileName)
	if err != nil {
		return nil, err
	}

	dbFile, err := os.Create(dbPath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(dbFile, d.reader)
	if err != nil {
		dbFile.Close()
		return nil, err
	}
	dbFile.Close()

	return sql.Open(d.driver, dbPath)
}

func (d *SQLDecoder) DecodeResource() {
	var err error
	d.db, err = d.createDB()
	if err != nil {
		d.addError(err)
		return
	}

	if err := d.decodeVersion(); err != nil {
		d.addError(err)
		return
	}

	if err := d.decodeClasses(); err != nil {
		d.addError(err)
		return
	}

	if err := d.decodeContents(); err != nil {
		d.addError(err)
		return
	}
}

func (d *SQLDecoder) DecodeObject() (EObject, error) {
	panic("SQLDecoder doesn't support object decoding")
}

func (d *SQLDecoder) addError(err error) {
	d.resource.GetErrors().Add(NewEDiagnosticImpl(err.Error(), d.resource.GetURI().String(), 0, 0))
}

func (d *SQLDecoder) decodeVersion() error {
	if row := d.db.QueryRow("PRAGMA user_version;"); row == nil {
		return fmt.Errorf("unable to retrieve user version")
	} else {
		var v int
		if err := row.Scan(&v); err != nil {
			return err
		}
		if v != sqlCodecVersion {
			return fmt.Errorf("history version %v is not supported", v)
		}
		return nil
	}
}

func (d *SQLDecoder) decodeContents() error {
	return nil
}

func (d *SQLDecoder) decodeObject() (EObject, error) {
	return nil, nil
}

func (d *SQLDecoder) decodeClasses() error {
	// read packages
	packagesData, err := d.decodePackages()
	if err != nil {
		return err
	}

	// read classes
	rows, err := d.db.Query("SELECT classID,packageID,name FROM classes")
	if err != nil {
		return err
	}
	defer rows.Close()

	classesData := map[int]*sqlDecoderClassData{}
	rawBuffer := make([]sql.RawBytes, 3)
	scanCallArgs := make([]any, len(rawBuffer))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanCallArgs...); err != nil {
			return err
		}
		classID, _ := strconv.Atoi(string(rawBuffer[0]))
		packageID, _ := strconv.Atoi(string(rawBuffer[1]))
		className := string(rawBuffer[2])
		ePackage, _ := packagesData[packageID]
		if ePackage == nil {
			return fmt.Errorf("unable to find package with id '%d'", packageID)
		}
		eClass, _ := ePackage.GetEClassifier(className).(EClass)
		if eClass == nil {
			return fmt.Errorf("unable to find class '%s' in package '%s'", className, ePackage.GetNsURI())
		}
		classesData[classID] = &sqlDecoderClassData{
			ePackage: ePackage,
			eClass:   eClass,
		}

	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func (d *SQLDecoder) decodePackages() (map[int]EPackage, error) {
	rows, err := d.db.Query("SELECT packageID,uri FROM packages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packagesData := map[int]EPackage{}
	rawBuffer := make([]sql.RawBytes, 2)
	scanCallArgs := make([]any, len(rawBuffer))
	for i := range rawBuffer {
		scanCallArgs[i] = &rawBuffer[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanCallArgs...); err != nil {
			return nil, err
		}
		packageID, _ := strconv.Atoi(string(rawBuffer[0]))
		packageURI := string(rawBuffer[1])
		packageRegistry := GetPackageRegistry()
		resourceSet := d.resource.GetResourceSet()
		if resourceSet != nil {
			packageRegistry = resourceSet.GetPackageRegistry()
		}
		ePackage := packageRegistry.GetPackage(packageURI)
		if ePackage == nil {
			return nil, fmt.Errorf("unable to find package '%s'", packageURI)
		}
		packagesData[packageID] = ePackage
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return packagesData, nil
}
