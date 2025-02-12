package ecore

import (
	"time"

	"go.uber.org/zap"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type executeQueryFn func(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error

type sqlBase struct {
	codecVersion          int64
	schema                *sqlSchema
	uri                   *URI
	objectIDName          string
	objectIDManager       EObjectIDManager
	isObjectID            bool
	isContainerID         bool
	logger                *zap.Logger
	executeQuery          executeQueryFn
	executeQueryTransient executeQueryFn
	executeQueryScript    executeQueryFn
}

var logQuery bool = true

func getExecuteQueryWithLoggerFn(fn executeQueryFn, logger *zap.Logger) executeQueryFn {
	if logQuery {
		return func(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) (err error) {
			start := time.Now()
			defer func() {
				args := []zap.Field{zap.String("query", query), zap.Duration("duration", time.Since(start))}
				if opts != nil {
					args = append(args, zap.Any("args", opts.Args))
				}
				logger.Debug("execute", args...)
			}()
			return fn(conn, query, opts)
		}
	}
	return fn
}

func (d *sqlBase) decodeProperties(conn *sqlite.Conn) (map[string]string, error) {
	// result
	properties := map[string]string{}

	// check if properties table exists
	tableExists := false
	if err := d.executeQuery(conn, `SELECT name FROM sqlite_master WHERE type='table' AND name='.properties';`, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			tableExists = true
			return nil
		},
	}); err != nil {
		return nil, err
	}
	if !tableExists {
		return properties, nil
	}

	// retrieve properties from table
	if err := d.executeQuery(conn, `SELECT key,value FROM ".properties" `, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			key := stmt.ColumnText(0)
			value := stmt.ColumnText(1)
			properties[key] = value
			return nil
		},
	}); err != nil {
		return nil, err
	}
	return properties, nil
}
