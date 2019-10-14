package auth

import (
	dbx "github.com/go-ozzo/ozzo-dbx"

	"bitbucket.org/gismart/{{Name}}/app/models"
	"bitbucket.org/gismart/{{Name}}/database"
)

type repository interface {
	database.Postgres
	Upsert(m *models.User) error
}

type postgres struct {
	*database.Psql
}

func (r *postgres) Upsert(m *models.User) error {
	return r.DBX().NewQuery(`
    INSERT INTO "user"(email, name) VALUES ({:email}, {:name})
      ON CONFLICT(email) DO UPDATE SET name = {:name}
      RETURNING *
  `).Bind(dbx.Params{
		"email": m.Email,
		"name":  m.Name,
	}).One(m)
}
