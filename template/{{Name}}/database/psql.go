package database

import (
	"strings"
	"sync"
	"fmt"
	
	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"

	"bitbucket.org/gismart/{{Name}}/app/models"
	"bitbucket.org/gismart/{{Name}}/config"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)

var (
	postgresOnce       sync.Once
	postgresConnection *dbx.DB
)

type Postgres interface {
	DBX() dbx.Builder
	Create(model dbx.TableModel, opts ...QueryOption) error
	GetByID(id string, model interface{}) error
	GetAll(m interface{}, query models.QParams) error
	Delete(model dbx.TableModel) error
	Update(model dbx.TableModel, opts ...QueryOption) error
	Commit() error
	Rollback() error
}

func PostgresConnect() error {
	var err error
	c := config.Config.Postgres

	f := func() {
		dsnFormat := "postgres://%s:%s@%s:%d/%s?sslmode=disable"

		dsn := fmt.Sprintf(
			dsnFormat,
			c.User,
			c.Password,
			c.Host,
			c.Port,
			c.DB)

		log.Infof(
			"postgres dsn: "+dsnFormat,
			"******", "******", c.Host, c.Port, c.DB)

		postgresConnection, err = dbx.MustOpen("postgres", dsn)
		if err != nil {
			return
		}

		// maxConn - 1 - for migrations package connection (healthcheck use go-pg connection)
		postgresConnection.DB().SetMaxOpenConns(c.MaxOpenConns - 1)

		if strings.ToLower(config.Config.LogLevel) == log.DebugLevel.String() {
			postgresConnection.LogFunc = log.Infof
		}
	}

	postgresOnce.Do(f)

	return err
}

func PostgresClose() error {
	return closeDB(postgresConnection)

}

func PostgresPing() error {
	return ping(postgresConnection)
}

func closeDB(connection *dbx.DB) error {
	if connection == nil {
		return errors.New("Can't close not opened connection")
	}

	return connection.Close()
}

func ping(connection *dbx.DB) error {
	if connection == nil {
		return errors.New("Can't ping not opened connection")
	}
	return connection.DB().Ping()
}
