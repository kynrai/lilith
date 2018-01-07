package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx"
)

var _ Repo = (*postgres)(nil)

type Repo interface {
	Conner
	Execer
	Transacter
}

type (
	// Conner is a datastore that can retrieve objects based on a query
	Conner interface {
		Conn(ctx context.Context, f func(*pgx.ConnPool) error) error
	}
	// ConnerFunc is a function adaptor for Conner
	ConnerFunc func(ctx context.Context, f func(*pgx.ConnPool) error) error
)

// Conner implements Queryer for ConnerFunc
func (c ConnerFunc) Query(ctx context.Context, f func(*pgx.ConnPool) error) error {
	return c(ctx, f)
}

type (
	// Execer runs the given query
	Execer interface {
		Exec(ctx context.Context, query string, args ...interface{}) error
	}
	// ListerFunc is a function adaptor for Execer
	ExecerFunc func(ctx context.Context, query string, args ...interface{}) error
)

// Exec implements Execer for ExecerFunc
func (f ExecerFunc) Exec(ctx context.Context, s string, a ...interface{}) error {
	return f(ctx, s, a...)
}

type (
	// Transacter is a datastore that can execute arbitary transactions
	Transacter interface {
		Transact(ctx context.Context, f func(*pgx.Tx) error) error
	}
	// TransacterFunc is a function adaptor for Transacter
	TransacterFunc func(ctx context.Context, f func(*pgx.Tx) error) error
)

// Transact implements Transacter for TransacterFunc
func (t TransacterFunc) Transact(ctx context.Context, f func(*pgx.Tx) error) (err error) {
	return t(ctx, f)
}

type postgres struct {
	psqlURI string
	psqlDSN string
	db      *pgx.ConnPool
}

type postgresOption func(*postgres)

func PsqlURI(uri string) postgresOption {
	return func(d *postgres) {
		d.psqlURI = uri
	}
}

func PsqlDSN(dsn string) postgresOption {
	return func(d *postgres) {
		d.psqlDSN = dsn
	}
}

func New(opts ...postgresOption) *postgres {
	p := new(postgres)
	for _, opt := range opts {
		opt(p)
	}

	// Connect to the database
	var config pgx.ConnConfig
	var err error
	switch {
	case p.psqlURI != "":
		config, err = pgx.ParseURI(p.psqlURI)
	case p.psqlDSN != "":
		config, err = pgx.ParseDSN(p.psqlDSN)
	}
	if err != nil {
		log.Fatal(err)
	}

	p.db, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 5,
		AcquireTimeout: time.Second * 5,
	})
	return p
}

func (p *postgres) Conn(ctx context.Context, f func(db *pgx.ConnPool) error) (err error) {
	done := make(chan struct{})
	go func() {
		defer func() {
			if p := recover(); p != nil {
				switch p := p.(type) {
				case error:
					err = p
				default:
					err = fmt.Errorf("%s", p)
				}
			}
			done <- struct{}{}
		}()
		err = f(p.db)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return
	}
}

func (p *postgres) Exec(ctx context.Context, query string, args ...interface{}) (err error) {
	done := make(chan struct{})
	go func() {
		defer func() {
			if p := recover(); p != nil {
				switch p := p.(type) {
				case error:
					err = p
				default:
					err = fmt.Errorf("%s", p)
				}
			}
			done <- struct{}{}
		}()
		_, err = p.db.Exec(query, args...)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return
	}
}

func (p *postgres) Transact(ctx context.Context, f func(*pgx.Tx) error) (err error) {
	done := make(chan struct{})
	go func() {
		defer func() {
			if p := recover(); p != nil {
				switch p := p.(type) {
				case error:
					err = p
				default:
					err = fmt.Errorf("%s", p)
				}
			}
			done <- struct{}{}
		}()
		err = transact(p.db, f)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return
	}
}

// transact takes a function and executes it within a DB transaction
func transact(db *pgx.ConnPool, f func(*pgx.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	return f(tx)
}

// InClause takes multiple strings, escapes and joins them so that they can be inserted in a IN SQL clause
func InClause(f func(string) string, strs ...string) string {
	a := make([]string, len(strs))
	for i, s := range strs {
		a[i] = f(s)
	}
	return strings.Join(a, ",")
}

var QuoteString = func(s string) string {
	return fmt.Sprintf("'%s'", s)
}
