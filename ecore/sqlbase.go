package ecore

import (
	"context"
	"database/sql"
	"sync"
)

type sqlSafeDB struct {
	delegate *sql.DB
	sync.Mutex
}

func (db *sqlSafeDB) Conn(ctx context.Context) (*sql.Conn, error) {
	db.Lock()
	defer db.Unlock()
	return db.delegate.Conn(ctx)
}

func (db *sqlSafeDB) Prepare(query string) (*sqlSafeStmt, error) {
	db.Lock()
	defer db.Unlock()
	stmt, err := db.delegate.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return &sqlSafeStmt{delegate: stmt, db: db}, nil
}

func (db *sqlSafeDB) Exec(query string, args ...any) (sql.Result, error) {
	db.Lock()
	defer db.Unlock()
	return db.delegate.ExecContext(context.Background(), query, args...)
}

func (db *sqlSafeDB) Query(query string, args ...any) (*sql.Rows, error) {
	db.Lock()
	defer db.Unlock()
	return db.delegate.QueryContext(context.Background(), query, args...)
}

func (db *sqlSafeDB) QueryRow(query string, args ...any) *sql.Row {
	db.Lock()
	defer db.Unlock()
	return db.delegate.QueryRowContext(context.Background(), query, args...)
}

func (db *sqlSafeDB) Begin() (*sqlSafeTx, error) {
	db.Lock()
	defer db.Unlock()
	tx, err := db.delegate.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &sqlSafeTx{delegate: tx, db: db}, nil
}

func newSQLSafeDB(db *sql.DB) *sqlSafeDB {
	return &sqlSafeDB{delegate: db}
}

type sqlSafeStmt struct {
	delegate *sql.Stmt
	db       *sqlSafeDB
}

func (stmt *sqlSafeStmt) Exec(args ...any) (sql.Result, error) {
	stmt.db.Lock()
	defer stmt.db.Unlock()
	return stmt.delegate.ExecContext(context.Background(), args...)
}

func (stmt *sqlSafeStmt) Query(args ...any) (*sql.Rows, error) {
	stmt.db.Lock()
	defer stmt.db.Unlock()
	return stmt.delegate.QueryContext(context.Background(), args...)
}

func (stmt *sqlSafeStmt) QueryRow(args ...any) *sql.Row {
	stmt.db.Lock()
	defer stmt.db.Unlock()
	return stmt.delegate.QueryRowContext(context.Background(), args...)
}

type sqlSafeTx struct {
	delegate *sql.Tx
	db       *sqlSafeDB
}

func (tx *sqlSafeTx) Commit() error {
	tx.db.Lock()
	defer tx.db.Unlock()
	return tx.delegate.Commit()
}

func (tx *sqlSafeTx) Rollback() error {
	tx.db.Lock()
	defer tx.db.Unlock()
	return tx.delegate.Rollback()
}

func (tx *sqlSafeTx) Stmt(stmt *sqlSafeStmt) *sqlSafeStmt {
	tx.db.Lock()
	defer tx.db.Unlock()
	return &sqlSafeStmt{delegate: tx.delegate.Stmt(stmt.delegate), db: tx.db}
}

type sqlBase struct {
	db              *sqlSafeDB
	schema          *sqlSchema
	uri             *URI
	idAttributeName string
	idManager       EObjectIDManager
}

func (s *sqlBase) encodeProperties() error {
	// properties
	// propertiesQuery := `
	// PRAGMA synchronous = NORMAL;
	// PRAGMA journal_mode = WAL;
	// `
	// _, err := s.db.Exec(propertiesQuery)
	// return err
	return nil
}

func (s *sqlBase) encodeSchema() error {
	// tables
	for _, table := range []*sqlTable{
		s.schema.packagesTable,
		s.schema.classesTable,
		s.schema.objectsTable,
		s.schema.contentsTable,
		s.schema.enumsTable,
	} {
		if _, err := s.db.Exec(table.createQuery()); err != nil {
			return err
		}
	}
	return nil
}
