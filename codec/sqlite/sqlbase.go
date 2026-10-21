package sqlite

import (
	"context"
	"errors"
	"io"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/chebyrash/promise"
	"github.com/masagroup/soft.go/ecore"
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

// boundConn is a connection bound to a goroutine : the queries of this goroutine
// are executed directly on this connection, in the order they are made, without
// the query scheduler
type boundConn struct {
	conn  *sqlite.Conn
	depth int // number of queries of this goroutine being executed
}

type sqlBase struct {
	codecVersion    int64
	schema          *sqlSchema
	uri             *ecore.URI
	objectIDName    string
	objectIDManager ecore.EObjectIDManager
	isObjectID      bool
	isContainerID   bool
	// isQualifiedTables names the class tables with their package: set for a new database, and
	// read from the "qualifiedTables" property of an existing one
	isQualifiedTables bool
	sqliteMutex       sync.Mutex
	sqliteQueries     map[string][]*query
	bindingsMutex     sync.Mutex
	bindings          map[int64]*boundConn
	logger            *zap.Logger
	sqliteLogger      *zap.Logger
	antsPool          *ants.Pool
	promisePool       promise.Pool
	connPool          *sqlitex.Pool
	connPoolProvider  func() (*sqlitex.Pool, error)
	connPoolClose     func(conn *sqlitex.Pool) error
}

// relativizeURI returns uri relative to the uri of the resource, uri itself for a resource without uri
func (s *sqlBase) relativizeURI(uri *ecore.URI) *ecore.URI {
	if s.uri != nil {
		return s.uri.Relativize(uri)
	}
	return uri
}

func (s *sqlBase) setLogger(logger *zap.Logger) {
	s.logger = logger
	s.sqliteLogger = logger.Named("sqlite")
}

// bindConn binds conn to the current goroutine until the returned function is called,
// which restores the connection bound before, if any
func (s *sqlBase) bindConn(conn *sqlite.Conn) (unbind func()) {
	id := goid.Get()
	s.bindingsMutex.Lock()
	if s.bindings == nil {
		s.bindings = map[int64]*boundConn{}
	}
	previous := s.bindings[id]
	s.bindings[id] = &boundConn{conn: conn}
	s.bindingsMutex.Unlock()
	return func() {
		s.bindingsMutex.Lock()
		if previous != nil {
			s.bindings[id] = previous
		} else {
			delete(s.bindings, id)
		}
		s.bindingsMutex.Unlock()
	}
}

// boundConn returns the connection bound to the current goroutine, or nil
func (s *sqlBase) boundConn() *boundConn {
	s.bindingsMutex.Lock()
	defer s.bindingsMutex.Unlock()
	return s.bindings[goid.Get()]
}

// execute sqlite cmd
// fn executes the query with the prepared statements cached by the connection,
// transientFn without : a query executed while another query of the same
// goroutine is being executed on the same connection (from its ResultFunc)
// must not reuse the cached statement, which may be the one being stepped
func (s *sqlBase) executeSqlite(fn, transientFn executeQueryFn, cmd string, opts *sqlitex.ExecOptions) error {

	// bound connection : execute directly
	if bound := s.boundConn(); bound != nil {
		loggerArgs := func() []zap.Field {
			args := []zap.Field{zap.Int64("goid", goid.Get()), zap.String("query", cmd)}
			if opts != nil {
				args = append(args, zap.Any("args", opts.Args))
			}
			return args
		}
		if e := s.logger.Check(zap.DebugLevel, "executing"); e != nil {
			e.Write(loggerArgs()...)
		}
		if bound.depth > 0 {
			fn = transientFn
		}
		bound.depth++
		err := fn(bound.conn, cmd, opts)
		bound.depth--
		if err != nil {
			s.logger.Error("executed", append(loggerArgs(), zap.Error(err))...)
			return err
		}
		if e := s.logger.Check(zap.DebugLevel, "executed"); e != nil {
			e.Write(loggerArgs()...)
		}
		return nil
	}

	// create query
	q := newQuery(cmd)

	// log
	if e := s.sqliteLogger.Check(zap.DebugLevel, "schedule"); e != nil {
		args := []zap.Field{zap.Int64("id", q.id), zap.Int64("goid", goid.Get()), zap.String("query", cmd)}
		if opts != nil {
			args = append(args, zap.Any("args", opts.Args))
		}
		e.Write(args...)
	}

	// lock sqlite
	s.sqliteMutex.Lock()

	// compute previous query
	// only one write to db is active
	// multiple read to db is allowed
	previous := map[*query]struct{}{}
	for _, table := range q.tables {
		// previous
		sqliteQueries := s.sqliteQueries[table]
		switch q.type_ {
		case queryRead:
			for _, query := range slices.Backward(sqliteQueries) {
				if query.type_ == queryWrite {
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
		goid := goid.Get()
		err := s.runQuery(q, goid, previous, fn, opts)
		// executed or failed : the next queries of its tables don't wait for it anymore
		if errUnregister := s.unregisterQuery(q, goid); err == nil {
			err = errUnregister
		}
		if err != nil {
			reject(err)
		} else {
			resolve(nil)
		}
	}, s.promisePool)
	s.sqliteMutex.Unlock()

	// wait for query to be finished
	_, err := q.promise.Await(context.Background())
	return err
}

// runQuery waits for the previous queries of a scheduled query, then executes it on a pool connection
func (s *sqlBase) runQuery(q *query, goid int64, previous map[*query]struct{}, fn executeQueryFn, opts *sqlitex.ExecOptions) error {
	if len(previous) > 0 {
		if e := s.sqliteLogger.Check(zap.DebugLevel, "waiting previous queries"); e != nil {
			e.Write(
				zap.Int64("id", q.id),
				zap.Int64("goid", goid),
				zap.Int64s("previous", mapSet(previous, func(query *query) int64 { return query.id })),
			)
		}
		promises := mapSet(previous, func(query *query) *promise.Promise[any] { return query.promise })
		if _, err := promise.AllWithPool(context.Background(), s.promisePool, promises...).Await(context.Background()); err != nil {
			if e := s.sqliteLogger.Check(zap.DebugLevel, "error in previous query"); e != nil {
				e.Write(
					zap.Int64("id", q.id),
					zap.Int64("goid", goid),
					zap.Error(err),
				)
			}
			return err
		}
	}

	// execute query
	conn, err := s.connPool.Take(context.Background())
	if err != nil {
		return err
	}
	defer s.connPool.Put(conn)

	loggerArgs := []zap.Field{zap.Int64("id", q.id), zap.Int64("goid", goid), zap.String("query", q.cmd)}
	if opts != nil {
		loggerArgs = append(loggerArgs, zap.Any("args", opts.Args))
	}
	s.logger.Log(zap.DebugLevel, "executing", loggerArgs...)
	if err := fn(conn, q.cmd, opts); err != nil {
		s.logger.Error("executed", append(loggerArgs, zap.Error(err))...)
		return err
	}
	s.logger.Log(zap.DebugLevel, "executed", loggerArgs...)
	return nil
}

// unregisterQuery removes a scheduled query from the queries of its tables
func (s *sqlBase) unregisterQuery(q *query, goid int64) error {
	if e := s.sqliteLogger.Check(zap.DebugLevel, "cleaning"); e != nil {
		e.Write(
			zap.Int64("id", q.id),
			zap.Int64("goid", goid),
		)
	}
	s.sqliteMutex.Lock()
	defer s.sqliteMutex.Unlock()

	for _, table := range q.tables {
		queries := s.sqliteQueries[table]
		index := slices.Index(queries, q)
		if index == -1 {
			return errors.New("unable to find query index")
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
	}
	if e := s.sqliteLogger.Check(zap.DebugLevel, "cleaned"); e != nil {
		e.Write(
			zap.Int64("id", q.id),
			zap.Int64("goid", goid),
		)
	}
	if len(s.sqliteQueries) == 0 {
		s.sqliteLogger.Debug("no pending")
	}
	return nil
}

func (d *sqlBase) executeQuery(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.Execute, sqlitex.ExecuteTransient, query, opts)
}

func (d *sqlBase) executeQueryTransient(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.ExecuteTransient, sqlitex.ExecuteTransient, query, opts)
}

func (d *sqlBase) executeQueryScript(query string, opts *sqlitex.ExecOptions) error {
	return d.executeSqlite(sqlitex.ExecuteScript, sqlitex.ExecuteScript, query, opts)
}

// hasTable returns true if the database has the table name
func (d *sqlBase) hasTable(name string) (bool, error) {
	exists := false
	err := d.executeQuery(`SELECT name FROM sqlite_master WHERE type='table' AND name=?;`, &sqlitex.ExecOptions{
		Args: []any{name},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			exists = true
			return nil
		},
	})
	return exists, err
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

// antsPromisePool runs the promises on an ants pool, or in a goroutine of their own when the
// pool refuses them (released by a close) : a promise always settles
type antsPromisePool struct {
	pool *ants.Pool
}

func newAntsPromisePool(pool *ants.Pool) promise.Pool {
	return &antsPromisePool{pool: pool}
}

func (p *antsPromisePool) Go(f func()) {
	if err := p.pool.Submit(f); err != nil {
		go f()
	}
}
