package ecore

import (
	"context"
	"errors"
	"io"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/chebyrash/promise"
	"github.com/panjf2000/ants/v2"
	"github.com/petermattis/goid"
	"github.com/rqlite/sql"
	"go.uber.org/zap"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type executeQueryFn func(conn *sqlite.Conn, query string, opts *sqlitex.ExecOptions) error

var queryID atomic.Int64

type queryType uint8

const (
	queryRead  queryType = 1 << 0
	queryWrite queryType = 1 << 1
)

type query struct {
	id      int64
	cmd     string
	type_   queryType
	tables  []string
	promise *promise.Promise[any]
}

func getTableIdentifiers(sources []sql.Source) (result []string) {
	for _, source := range sources {
		switch s := source.(type) {
		case *sql.JoinClause:
			result = append(result, getTableIdentifiers([]sql.Source{s.X, s.Y})...)
		case *sql.ParenSource:
			result = append(result, getTableIdentifiers([]sql.Source{s.X})...)
		case *sql.QualifiedTableName:
			result = append(result, s.Name.Name)
		}
	}
	return
}

func newQuery(cmd string) *query {
	var queryType queryType
	var queryTables []string
	parser := sql.NewParser(strings.NewReader(cmd))
loop:
	for {
		stmt, err := parser.ParseStatement()
		if err != nil {
			if err != io.EOF {
				// parser error - all tables - write
				queryType = queryWrite
				queryTables = nil
			}
			break loop
		}

		// compute query type and query tables
		switch s := stmt.(type) {
		case *sql.SelectStatement:
			queryType = queryRead
			queryTables = append(queryTables, getTableIdentifiers([]sql.Source{s.Source})...)
		case *sql.CreateTableStatement:
			queryType = queryWrite
			queryTables = append(queryTables, s.Name.Name)
		case *sql.InsertStatement:
			queryType = queryWrite
			queryTables = append(queryTables, s.Table.Name)
		case *sql.DeleteStatement:
			queryType = queryWrite
			queryTables = append(queryTables, s.Table.Name.Name)
		case *sql.UpdateStatement:
			queryType = queryWrite
			queryTables = append(queryTables, s.Table.Name.Name)
		default:
			// unknown statement - all tables - write
			queryType = queryWrite
			queryTables = nil
			break loop
		}
	}

	// add empty table if in write mode
	// to ensure query will be global locked
	if queryType == queryWrite {
		queryTables = append(queryTables, "")
	}

	return &query{
		id:     queryID.Add(1),
		cmd:    cmd,
		type_:  queryType,
		tables: queryTables,
	}
}

type sqlBase struct {
	codecVersion     int64
	schema           *sqlSchema
	uri              *URI
	objectIDName     string
	objectIDManager  EObjectIDManager
	isObjectID       bool
	isContainerID    bool
	sqliteMutex      sync.Mutex
	sqliteQueries    map[string][]*query
	logger           *zap.Logger
	antsPool         *ants.Pool
	promisePool      promise.Pool
	connPool         *sqlitex.Pool
	connPoolProvider func() (*sqlitex.Pool, error)
	connPoolClose    func(conn *sqlitex.Pool) error
}

// execute sqlite cmd
func (s *sqlBase) executeSqlite(fn executeQueryFn, cmd string, opts *sqlitex.ExecOptions) error {
	s.sqliteMutex.Lock()

	// create query
	q := newQuery(cmd)

	// log
	if e := s.logger.Named("sqlite").Check(zap.DebugLevel, "schedule"); e != nil {
		args := []zap.Field{zap.Int64("id", q.id), zap.Int64("goid", goid.Get()), zap.String("query", cmd)}
		if opts != nil {
			args = append(args, zap.Any("args", opts.Args))
		}
		s.logger.Named("sqlite").Debug("schedule", args...)
	}

	// compute previous query
	// only one write to db is active
	// multiple read to db is allowed
	previous := map[*query]struct{}{}
	for _, table := range q.tables {
		// previous
		sqliteQueries := s.sqliteQueries[table]
		switch q.type_ {
		case queryRead:
			for i := len(sqliteQueries) - 1; i >= 0; i-- {
				if query := sqliteQueries[i]; query.type_ == queryWrite {
					previous[query] = struct{}{}
					break
				}
			}
		case queryWrite:
			if len(sqliteQueries) > 0 {
				query := sqliteQueries[len(sqliteQueries)-1]
				previous[query] = struct{}{}
			}
		}
		// register query for thid table
		s.sqliteQueries[table] = append(sqliteQueries, q)
	}

	// create query promise
	q.promise = promise.NewWithPool(func(resolve func(any), reject func(error)) {
		logger := s.logger.Named("sqlite").With(zap.Int64("id", q.id), zap.Int64("goid", goid.Get()))

		if len(previous) > 0 {
			if e := logger.Check(zap.DebugLevel, "waiting previous queries"); e != nil {
				e.Write(zap.Int64s("previous", mapSet(previous, func(query *query) int64 { return query.id })))
			}
			promises := mapSet(previous, func(query *query) *promise.Promise[any] { return query.promise })
			if _, err := promise.All(context.Background(), promises...).Await(context.Background()); err != nil {
				logger.Debug("error in previous query", zap.Error(err))
				reject(err)
				return
			}
		}

		// execute query
		conn, err := s.connPool.Take(context.Background())
		if err != nil {
			reject(err)
			return
		}
		defer s.connPool.Put(conn)

		args := []zap.Field{zap.String("query", cmd)}
		if opts != nil {
			args = append(args, zap.Any("args", opts.Args))
		}
		executeLogger := logger.With(args...)
		executeLogger.Debug("executing")
		if err := fn(conn, cmd, opts); err != nil {
			executeLogger.Error("executed", zap.Error(err))
			reject(err)
			return
		} else {
			executeLogger.Debug("executed")
		}

		// clean query
		logger.Debug("cleaning")
		s.sqliteMutex.Lock()
		defer s.sqliteMutex.Unlock()

		// unregister query from all tables
		for _, table := range q.tables {
			queries := s.sqliteQueries[table]
			index := slices.Index(queries, q)
			if index == -1 {
				reject(errors.New("unable to find query index"))
				return
			}
			// remove query from queries
			copy(queries[index:], queries[index+1:])
			queries[len(queries)-1] = nil
			queries = queries[:len(queries)-1]
			if len(queries) > 0 {
				s.sqliteQueries[table] = queries
			} else {
				delete(s.sqliteQueries, table)
			}
			logger.Debug("cleaned")
		}
		if len(s.sqliteQueries) == 0 {
			logger.Debug("no pending")
		}
		resolve(nil)
	}, s.promisePool)

	s.sqliteMutex.Unlock()

	// wait for query to be finished
	_, err := q.promise.Await(context.Background())
	return err
}

func (d *sqlBase) executeQuery(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.Execute, query, opts)
}

func (d *sqlBase) executeQueryTransient(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.ExecuteTransient, query, opts)
}

func (d *sqlBase) executeQueryScript(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.ExecuteScript, query, opts)
}

func (d *sqlBase) decodeProperties() (map[string]string, error) {
	// result
	properties := map[string]string{}

	// check if properties table exists
	tableExists := false
	if err := d.executeQuery(`SELECT name FROM sqlite_master WHERE type='table' AND name='.properties';`, &sqlitex.ExecOptions{
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
	if err := d.executeQuery(`SELECT "key",value FROM ".properties";`, &sqlitex.ExecOptions{
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
