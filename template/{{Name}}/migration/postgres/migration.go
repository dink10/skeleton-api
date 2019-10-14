package postgres

import (
	"fmt"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"

	"bitbucket.org/gismart/{{Name}}/config"
)

var (
	Migrations   = migrations.NewCollection()
	pgConnection *pg.DB
)

func GetRegisteredMigrations() []*migrations.Migration {
	return Migrations.Migrations()
}

func GetLastVersion() (int64, error) {
	return Migrations.Version(GetPgConnection())
}

func GetPgConnection() *pg.DB {
	if pgConnection != nil {
		return pgConnection
	}

	addr := fmt.Sprintf("%s:%d", config.Config.Postgres.Host, config.Config.Postgres.Port)
	fmt.Println(addr)

	pgConnection = pg.Connect(&pg.Options{
		Addr:     addr,
		User:     config.Config.Postgres.User,
		Password: config.Config.Postgres.Password,
		Database: config.Config.Postgres.DB,
		PoolSize: 1,
	})

	return pgConnection
}

type iterator struct {
	fns           []func(db migrations.DB) error
	nextCallIndex int
	closed        bool
	err           error
}

func initChain(fns ...func(db migrations.DB) error) *iterator {
	return &iterator{fns: fns}
}

func (fns *iterator) Err() error {
	return fns.err
}

func (fns *iterator) RunNext(db migrations.DB) {
	if err := fns.fns[fns.nextCallIndex](db); err != nil {
		fns.err = err
		return
	}

	fns.nextCallIndex += 1

	if fns.nextCallIndex >= len(fns.fns) {
		fns.closed = true
	}
}

func (fns *iterator) Next() bool {
	return fns.err == nil && !fns.closed
}

func (fns *iterator) Run(db migrations.DB) error {
	if len(fns.fns) == 0 {
		return nil
	}

	for fns.Next() {
		fns.RunNext(db)
	}

	return fns.Err()
}
