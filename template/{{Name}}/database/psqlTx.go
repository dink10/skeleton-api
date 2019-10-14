package database

import (
	"context"
	"database/sql"

	"bitbucket.org/gismart/ddtracer"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

type Txer interface {
	Commit() error
	Rollback() error
}

type Tx struct {
	dbx.Builder
	tx *sql.Tx
}

func (t *Tx) Commit() error {
	err := t.tx.Commit()
	return err
}

func (t *Tx) Rollback() error {
	if err := t.tx.Rollback(); err != sql.ErrTxDone {
		return err
	}
	return nil
}

func PostgresTransaction(ctx context.Context, withTracing bool) (*Psql, error) {
	tx, err := postgresConnection.DB().Begin()
	if err != nil {
		return nil, err
	}

	var builder dbx.Builder

	if !withTracing {
		builder = postgresConnection.Builder
	} else {
		exc := &ddtracer.DBXExecutor{Tx: tx, Ctx: ctx}
		builderFn, ok := dbx.BuilderFuncMap[postgresConnection.DriverName()]
		if !ok {
			builderFn = dbx.NewStandardBuilder
		}
		builder = builderFn(postgresConnection, exc)
	}

	return &Psql{&Tx{
		Builder: builder,
		tx:      tx,
	}}, nil
}
