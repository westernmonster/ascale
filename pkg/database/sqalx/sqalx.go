package sqalx

import (
	"ascale/pkg/database/sqlx"
	"ascale/pkg/def"
	"ascale/pkg/ecode"
	"ascale/pkg/log"
	"ascale/pkg/net/netutil/breaker"
	"ascale/pkg/xtime"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type masterKey struct{}

func NewContext(c context.Context, useMaster bool) context.Context {
	return context.WithValue(c, masterKey{}, useMaster)
}

func UseMaster(c context.Context) bool {
	useMaster, ok := c.Value(masterKey{}).(bool)
	if !ok {
		return false
	}

	return useMaster
}

type Config struct {
	Addr         string          // for trace
	DSN          string          // write data source name.
	ReadDSN      []string        // read data source name.
	Active       int             // pool
	Idle         int             // pool
	IdleTimeout  xtime.Duration  // connect max life time.
	QueryTimeout xtime.Duration  // query sql timeout
	ExecTimeout  xtime.Duration  // execute sql timeout
	TranTimeout  xtime.Duration  // transaction sql timeout
	Breaker      *breaker.Config // breaker
}

var (
	// ErrNotInTransaction is returned when using Commit
	// outside of a transaction.
	ErrNotInTransaction = errors.New("not in transaction")

	// ErrIncompatibleOption is returned when using an option incompatible
	// with the selected driver.
	ErrIncompatibleOption = errors.New("incompatible option")
)

func NewMySQL(c *Config) Node {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("mysql must be set query/execute/transction timeout")
	}

	brkGroup := breaker.NewGroup(c.Breaker)
	brk := brkGroup.Get(c.Addr)

	w, err := connect(c, c.DSN, brk)
	if err != nil {
		log.Error(fmt.Sprintf("open mysql error(%v)", err))
		panic(err)
	}

	rs := make([]*sqlx.DB, 0, len(c.ReadDSN))
	for _, rd := range c.ReadDSN {
		brk := brkGroup.Get(parseDSNAddr(rd))
		d, err := connect(c, rd, brk)
		if err != nil {
			log.Error(fmt.Sprintf("open mysql error(%v)", err))
			panic(err)
		}
		rs = append(rs, d)
	}

	n := node{
		write:  w,
		read:   rs,
		master: w,
		Driver: w,
		Config: c,
	}

	return &n
}

func connect(c *Config, dataSourceName string, breaker breaker.Breaker) (db *sqlx.DB, err error) {
	addr := parseDSNAddr(dataSourceName)
	db, err = sqlx.Open(
		"mysql",
		dataSourceName,
		breaker,
		addr,
		c.QueryTimeout,
		c.ExecTimeout,
		c.TranTimeout,
	)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	db.DB.SetMaxOpenConns(c.Active)
	db.DB.SetMaxIdleConns(c.Idle)
	db.DB.SetConnMaxLifetime(time.Duration(c.IdleTimeout))
	return db, nil
}

// A Node is a database driver that can manage nested transactions.
type Node interface {
	// Close the underlying sqlx connection.
	Close() error
	// Begin a new transaction.
	Beginx(c context.Context) (Node, error)
	// Rollback the associated transaction.
	Rollback() error
	// Commit the assiociated transaction.
	Commit() error
	// Tx returns the underlying transaction.
	Tx() *sqlx.Tx

	SelectContext(
		ctx context.Context,
		dest interface{},
		query string,
		args ...interface{},
	) (err error)

	ExecContext(
		ctx context.Context,
		query string,
		args ...interface{},
	) (result sql.Result, err error)

	// execute sql with retry only deadlock
	ExecContextWithRetry(
		ctx context.Context,
		query string,
		args ...interface{},
	) (result sql.Result, err error)

	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)

	GetContextByMaster(
		ctx context.Context,
		dest interface{},
		query string,
		args ...interface{},
	) (err error)

	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	Ping(c context.Context) (err error)

	Rebind(query string) string
}

// A Driver can query the database. It can either be a *sqlx.DB or a *sqlx.Tx
// and therefore is limited to the methods they have in common.
type Driver interface {
	sqlx.QueryerContext
	sqlx.PreparerContext
	// BindNamed(query string, arg interface{}) (string, []interface{}, error)
	// DriverName() string
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	Rebind(query string) string
}

type node struct {
	Driver Driver
	Config *Config
	write  *sqlx.DB
	read   []*sqlx.DB
	idx    int64
	master *sqlx.DB
	tx     *sqlx.Tx
	nested bool
}

func (n node) Beginx(c context.Context) (Node, error) {
	var err error

	switch {
	case n.tx == nil:
		// new actual transaction
		n.tx, err = n.write.Beginx(c)
		n.Driver = n.tx
	default:
		// already in a transaction: reusing current transaction
		n.nested = true
	}

	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (n *node) Rollback() error {
	if n.tx == nil {
		return nil
	}

	var err error

	if !n.nested {
		err = n.tx.Rollback()
	}

	if err != nil {
		return err
	}

	n.tx = nil
	n.Driver = nil

	return nil
}

func (n *node) Commit() error {
	if n.tx == nil {
		return ErrNotInTransaction
	}

	var err error

	if !n.nested {
		err = n.tx.Commit()
	}

	if err != nil {
		return err
	}

	n.tx = nil
	n.Driver = nil

	return nil
}

// Tx returns the underlying transaction.
func (n *node) Tx() *sqlx.Tx {
	return n.tx
}

// Ping verifies a connection to the database is still alive, establishing a
// connection if necessary.
func (n *node) Ping(c context.Context) (err error) {
	if err = n.write.DBPing(c); err != nil {
		return
	}
	for _, rd := range n.read {
		if err = rd.DBPing(c); err != nil {
			return
		}
	}
	return
}

func (n *node) Close() (err error) {
	if err = n.write.Close(); err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, rd := range n.read {
		if e := rd.Close(); e != nil {
			err = errors.WithStack(e)
			return
		}
	}
	return
}

func (n *node) SelectContext(
	ctx context.Context,
	dest interface{},
	query string,
	args ...interface{},
) (err error) {
	if n.tx == nil && UseMaster(ctx) {
		return n.write.SelectContext(ctx, dest, query, args...)
	}

	// By default, to write database is used for execution. If a transaction exists, the driver will use to write database.
	if n.tx == nil {
		idx := n.readIndex()
		for i := range n.read {
			if err = n.read[(idx+i)%len(n.read)].SelectContext(ctx, dest, query, args...); !ecode.ServiceUnavailable.Equal(
				err,
			) {
				return
			}
		}
	}
	return n.Driver.SelectContext(ctx, dest, query, args...)
}

func (n *node) ExecContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (result sql.Result, err error) {
	if n.tx == nil && UseMaster(ctx) {
		return n.write.ExecContext(ctx, query, args...)
	}

	// By default, to write database is used for execution. If a transaction exists, the driver will use to write database.
	if n.tx != nil {
		return n.tx.ExecContext(n.tx.Context, query, args...)
	}

	return n.write.ExecContext(ctx, query, args...)
}

func (n *node) ExecContextWithRetry(
	ctx context.Context,
	query string,
	args ...interface{},
) (result sql.Result, err error) {
	maxRetry := 3
	for try := 1; try <= maxRetry; try++ {
		if n.tx == nil && UseMaster(ctx) {
			result, err = n.write.ExecContext(ctx, query, args...)
		} else if n.tx != nil {
			result, err = n.tx.ExecContext(n.tx.Context, query, args...)
		} else {
			result, err = n.write.ExecContext(ctx, query, args...)
		}

		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok &&
				mysqlErr.Number == def.MysqlErrorCode.DeadLock {
				log.For(ctx).Errorf("ExecContextWithRetry deadlock, retry %d", try)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			log.For(ctx).Errorf("ExecContextWithRetry error, retry %d", try)
			return result, err
		}

		return result, nil
	}

	return result, err
}

func (n *node) GetContext(
	ctx context.Context,
	dest interface{},
	query string,
	args ...interface{},
) (err error) {
	if n.tx == nil && UseMaster(ctx) {
		return n.write.GetContext(ctx, dest, query, args...)
	}
	// 事务默认write库执行，如果没有事务，则从随机从只读库中读取
	if n.tx == nil {
		idx := n.readIndex()
		for i := range n.read {
			if err = n.read[(idx+i)%len(n.read)].GetContext(ctx, dest, query, args...); !ecode.ServiceUnavailable.Equal(
				err,
			) {
				return
			}
		}
	}
	return n.Driver.GetContext(ctx, dest, query, args...)
}

func (n *node) GetContextByMaster(
	ctx context.Context,
	dest interface{},
	query string,
	args ...interface{},
) (err error) {
	if n.tx == nil {
		return n.write.GetContext(ctx, dest, query, args...)
	}
	return n.Driver.GetContext(ctx, dest, query, args...)
}

func (n *node) Rebind(query string) string {
	return n.Driver.Rebind(query)
}

func (n *node) QueryxContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (rows *sqlx.Rows, err error) {
	if n.tx == nil && UseMaster(ctx) {
		return n.write.QueryxContext(ctx, query, args...)
	}

	if n.tx == nil {
		idx := n.readIndex()
		for i := range n.read {
			if rows, err = n.read[(idx+i)%len(n.read)].QueryxContext(ctx, query, args...); !ecode.ServiceUnavailable.Equal(
				err,
			) {
				return
			}
		}
	}

	return n.write.QueryxContext(ctx, query, args...)
}

func (n *node) readIndex() int {
	if len(n.read) == 0 {
		return 0
	}
	v := atomic.AddInt64(&n.idx, 1)
	return int(v) % len(n.read)
}

// parseDSNAddr parse dsn name and return addr.
func parseDSNAddr(dsn string) (addr string) {
	if dsn == "" {
		return
	}
	part0 := strings.Split(dsn, "@@")
	if len(part0) > 1 {
		part1 := strings.Split(part0[1], "?")
		if len(part1) > 0 {
			addr = part1[0]
		}
	}
	return
}
