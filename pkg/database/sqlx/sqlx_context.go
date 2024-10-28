// +build go1.8

package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"ascale/pkg/conf/env"
	"ascale/pkg/stat"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

var stats = stat.DB

// Config mysql config.
type Config struct {
	Addr         string        // for trace
	DSN          string        // write data source name.
	ReadDSN      []string      // read data source name.
	Active       int           // pool
	Idle         int           // pool
	IdleTimeout  time.Duration // connect max life time.
	QueryTimeout time.Duration // query sql timeout
	ExecTimeout  time.Duration // execute sql timeout
	TranTimeout  time.Duration // transaction sql timeout
	// Breaker      *breaker.Config // breaker
}

// // ConnectContext to a database and verify with a ping.
// func ConnectContext(ctx context.Context, driverName, dataSourceName string) (*DB, error) {
// 	db, err := Open(driverName, dataSourceName)
// 	if err != nil {
// 		return db, err
// 	}
// 	err = db.PingContext(ctx)
// 	return db, err
// }

// QueryerContext is an interface used by GetContext and SelectContext
type QueryerContext interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *Row
}

// PreparerContext is an interface used by PreparexContext.
type PreparerContext interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// SelectContext executes a query using the provided Queryer, and StructScans
// each row into dest, which must be a slice.  If the slice elements are
// scannable, then the result set must have only one column.  Otherwise,
// StructScan is used. The *sql.Rows are closed automatically.
// Any placeholder parameters are replaced with supplied args.
func SelectContext(ctx context.Context, q QueryerContext, dest interface{}, query string, args ...interface{}) error {
	rows, err := q.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()
	return scanAll(rows, dest, false)
}

// PreparexContext prepares a statement.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func PreparexContext(ctx context.Context, p PreparerContext, query string) (*Stmt, error) {
	s, err := p.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{Stmt: s, unsafe: isUnsafe(p), Mapper: mapperFor(p)}, err
}

// GetContext does a QueryRow using the provided Queryer, and scans the
// resulting row to dest.  If dest is scannable, the result must only have one
// column. Otherwise, StructScan is used.  Get will return sql.ErrNoRows like
// row.Scan would. Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func GetContext(ctx context.Context, q QueryerContext, dest interface{}, query string, args ...interface{}) error {
	r := q.QueryRowxContext(ctx, query, args...)
	return r.scanAny(dest, false)
}

// SelectContext using this DB.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	now := time.Now()
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		ctx, span = global.Tracer(env.AppID).Start(ctx, "select")
		span.SetAttribute("peer.address", db.Addr)
		span.SetAttribute("sql", query)
		defer span.End()
	}

	if err = db.breaker.Allow(); err != nil {
		stats.Incr("mysql:select", "breaker")
		return
	}

	_, c, cancel := db.QueryTimeout.Shrink(ctx)
	err = SelectContext(c, db, dest, query, args...)
	db.onBreaker(&err)
	stats.Timing("mysql:query", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.Wrapf(err, "query:%s, args:%+v", query, args)
		cancel()
		return
	}
	return
}

// GetContext using this DB.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return GetContext(ctx, db, dest, query, args...)
}

// PreparexContext returns an sqlx.Stmt instead of a sql.Stmt.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func (db *DB) PreparexContext(ctx context.Context, query string) (*Stmt, error) {
	return PreparexContext(ctx, db, query)
}

// QueryxContext queries the database and returns an *sqlx.Rows.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) QueryxContext(ctx context.Context, query string, args ...interface{}) (rows *Rows, err error) {
	now := time.Now()
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		ctx, span = global.Tracer(env.AppID).Start(ctx, "query")
		span.SetAttribute("peer.address", db.Addr)
		span.SetAttribute("sql", query)
		defer span.End()
	}

	if err = db.breaker.Allow(); err != nil {
		stats.Incr("mysql:query", "breaker")
		return
	}

	_, c, cancel := db.QueryTimeout.Shrink(ctx)

	r, err := db.DB.QueryContext(c, query, args...)
	db.onBreaker(&err)
	stats.Timing("mysql:query", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.Wrapf(err, "query:%s, args:%+v", query, args)
		cancel()
		return
	}
	return &Rows{Rows: r, unsafe: db.unsafe, Mapper: db.Mapper}, err
}

// QueryRowxContext queries the database and returns an *sqlx.Row.
// Any placeholder parameters are replaced with supplied args.
func (db *DB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *Row {
	now := time.Now()
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		ctx, span = global.Tracer(env.AppID).Start(ctx, "queryrow")
		span.SetAttribute("address", db.Addr)
		span.SetAttribute("sql", query)
		defer span.End()
	}

	if err := db.breaker.Allow(); err != nil {
		stats.Incr("mysql:queryrow", "breaker")
		return &Row{err: err, unsafe: db.unsafe, Mapper: db.Mapper}
	}

	_, c, cancel := db.QueryTimeout.Shrink(ctx)
	rows, err := db.DB.QueryContext(c, query, args...)
	stats.Timing("mysql:queryrow", int64(time.Since(now)/time.Millisecond))

	return &Row{rows: rows, err: err, unsafe: db.unsafe, Mapper: db.Mapper, cancel: cancel}
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	now := time.Now()
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		ctx, span = global.Tracer(env.AppID).Start(ctx, "exec")
		span.SetAttribute("address", db.Addr)
		span.SetAttribute("sql", query)
		defer span.End()
	}

	if err = db.breaker.Allow(); err != nil {
		stats.Incr("mysql:exec", "breaker")
		return
	}

	_, c, cancel := db.ExecTimeout.Shrink(ctx)

	result, err = db.DB.ExecContext(c, query, args...)
	cancel()
	db.onBreaker(&err)
	stats.Timing("mysql:exec", int64(time.Since(now)/time.Millisecond))
	if err != nil {
		err = errors.Wrapf(err, "exec:%s, args:%+v", query, args)
		cancel()
		return
	}
	return
}

// BeginTxx begins a transaction and returns an *sqlx.Tx instead of an
// *sql.Tx.
//
// The provided context is used until the transaction is committed or rolled
// back. If the context is canceled, the sql package will roll back the
// transaction. Tx.Commit will return an error if the context provided to
// BeginxContext is canceled.
func (db *DB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx, driverName: db.driverName, unsafe: db.unsafe, Mapper: db.Mapper}, err
}

// StmtxContext returns a version of the prepared statement which runs within a
// transaction. Provided stmt can be either *sql.Stmt or *sqlx.Stmt.
func (tx *Tx) StmtxContext(ctx context.Context, stmt interface{}) *Stmt {
	var s *sql.Stmt
	switch v := stmt.(type) {
	case Stmt:
		s = v.Stmt
	case *Stmt:
		s = v.Stmt
	case *sql.Stmt:
		s = v
	default:
		panic(fmt.Sprintf("non-statement type %v passed to Stmtx", reflect.ValueOf(stmt).Type()))
	}
	return &Stmt{Stmt: tx.StmtContext(ctx, s), Mapper: tx.Mapper}
}

// PreparexContext returns an sqlx.Stmt instead of a sql.Stmt.
//
// The provided context is used for the preparation of the statement, not for
// the execution of the statement.
func (tx *Tx) PreparexContext(ctx context.Context, query string) (*Stmt, error) {
	if tx.span != nil {
		tx.span.AddEvent(ctx, "prepare", label.String("sql", query))
	}

	return PreparexContext(ctx, tx, query)
}

// QueryxContext within a transaction and context.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryxContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	if tx.span != nil {
		tx.span.AddEvent(ctx, "query", label.String("sql", query))
	}

	now := time.Now()
	defer func() {
		stats.Timing("mysql:tx:query", int64(time.Since(now)/time.Millisecond))
	}()

	r, err := tx.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "query:%s, args:%+v", query, args)
	}
	return &Rows{Rows: r, unsafe: tx.unsafe, Mapper: tx.Mapper}, err
}

// SelectContext within a transaction and context.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if tx.span != nil {
		tx.span.AddEvent(ctx, "select", label.String("sql", query))
	}

	now := time.Now()
	defer func() {
		stats.Timing("mysql:tx:query", int64(time.Since(now)/time.Millisecond))
	}()

	return SelectContext(ctx, tx, dest, query, args...)
}

// GetContext within a transaction and context.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (tx *Tx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if tx.span != nil {
		tx.span.AddEvent(ctx, "get", label.String("sql", query))
	}
	return GetContext(ctx, tx, dest, query, args...)
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	if tx.span != nil {
		tx.span.AddEvent(ctx, "exec", label.String("sql", query))
	}
	return tx.Tx.ExecContext(ctx, query, args...)
}

// QueryRowxContext within a transaction and context.
// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *Row {
	if tx.span != nil {
		tx.span.AddEvent(ctx, "queryrow", label.String("sql", query))
	}

	now := time.Now()
	defer func() {
		stats.Timing("mysql:tx:queryrow", int64(time.Since(now)/time.Millisecond))
	}()
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	return &Row{rows: rows, err: err, unsafe: tx.unsafe, Mapper: tx.Mapper}
}

// SelectContext using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	return SelectContext(ctx, &qStmt{s}, dest, "", args...)
}

// GetContext using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func (s *Stmt) GetContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	return GetContext(ctx, &qStmt{s}, dest, "", args...)
}

// QueryRowxContext using this statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) QueryRowxContext(ctx context.Context, args ...interface{}) *Row {
	qs := &qStmt{s}
	return qs.QueryRowxContext(ctx, "", args...)
}

// QueryxContext using this statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) QueryxContext(ctx context.Context, args ...interface{}) (*Rows, error) {
	qs := &qStmt{s}
	return qs.QueryxContext(ctx, "", args...)
}

func (q *qStmt) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return q.Stmt.QueryContext(ctx, args...)
}

func (q *qStmt) QueryxContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	r, err := q.Stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, unsafe: q.Stmt.unsafe, Mapper: q.Stmt.Mapper}, err
}

func (q *qStmt) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *Row {
	rows, err := q.Stmt.QueryContext(ctx, args...)
	return &Row{rows: rows, err: err, unsafe: q.Stmt.unsafe, Mapper: q.Stmt.Mapper}
}

func (q *qStmt) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return q.Stmt.ExecContext(ctx, args...)
}
