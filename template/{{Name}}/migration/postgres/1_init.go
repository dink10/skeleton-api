package postgres

import (
	"fmt"
	"log"

	"github.com/go-pg/migrations"
)

func init() {
	err := Migrations.RegisterTx(
		func(db migrations.DB) error {
			return initChain(
				createUser,
			).Run(db)
		},
		func(db migrations.DB) error {
			return initChain(
				dropUser,
			).Run(db)
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}

func createUser(db migrations.DB) error {
	fmt.Println("creating user table...")

	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS "user"
      (
        id serial NOT NULL
          CONSTRAINT user_pkey
            PRIMARY KEY,
        email varchar(255) NOT NULL,
        name varchar(255)
      );
    CREATE UNIQUE INDEX IF NOT EXISTS user_unique ON "user"(email)
  `)

	return err
}

func dropUser(db migrations.DB) error {
	fmt.Println("dropping user table...")

	_, err := db.Exec(`DROP TABLE IF EXISTS "user"`)

	return err
}
