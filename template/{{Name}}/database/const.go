package database

import (
	"context"

	"github.com/pkg/errors"
)

type ContextKey string

const (
	ContextStorage = ContextKey("storage")
)

func GetFromContext(ctx context.Context) (*Psql, error) {
	if db, ok := ctx.Value(ContextStorage).(*Psql); ok {
		return db, nil
	}

	return nil, errors.New("can not cast to *Psql")
}
