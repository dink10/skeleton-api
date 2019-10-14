package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/pg"

	"bitbucket.org/gismart/{{Name}}/config"
	"bitbucket.org/gismart/{{Name}}/migration/postgres"
)

const usageText = `This program runs command on the db. Supported commands are:
  - up - runs all available migrations.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.

Usage:
  go run *.go <command> [args]
`

func main() {
	flag.Usage = usage
	flag.Parse()
	config.Init(&config.Config)
	addr := fmt.Sprintf("%s:%d", config.Config.Postgres.Host, config.Config.Postgres.Port)
	fmt.Println(addr)

	db := pg.Connect(&pg.Options{
		Addr:     addr,
		User:     config.Config.Postgres.User,
		Password: config.Config.Postgres.Password,
		Database: config.Config.Postgres.DB,
		PoolSize: 1,
	})

	fmt.Println(flag.Args())

	oldVersion, newVersion, err := postgres.Migrations.Run(db, flag.Args()...)
	if err != nil {
		e, ok := err.(pg.Error)
		if !ok {
			panic(err)
		}
		panic(fmt.Errorf("%s \n Details: %s \n Position: %s \n Hint: %s \n Query: %s \n Line: %s", err.Error(), e.Field('D'), e.Field('P'), e.Field('H'), e.Field('q'), e.Field('L')))
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}
