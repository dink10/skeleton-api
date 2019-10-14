package database

import (
	"context"
	
	"bitbucket.org/gismart/{{Name}}/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)

type QueryOption func(*dbx.ModelQuery) *dbx.ModelQuery

func WithExclude(excludedFields ...string) QueryOption {
	return func(mq *dbx.ModelQuery) *dbx.ModelQuery {
		return mq.Exclude(excludedFields...)
	}
}

type Builder interface {
	dbx.Builder
	Txer
}

type Psql struct {
	Builder
}

func (p *Psql) DBX() dbx.Builder {
	return p.Builder
}

func (p *Psql) Create(m dbx.TableModel, opts ...QueryOption) error {
	q := p.Model(m)

	for _, opt := range opts {
		q = opt(q)
	}

	return q.Insert()
}

func (p *Psql) GetByID(id string, m interface{}) error {
	return p.Select().Model(id, m)
}

func (p *Psql) GetBy(cond map[string]interface{}, m interface{}) error {
	return p.Select().Where(dbx.HashExp(cond)).One(m)
}

func (p *Psql) GetAll(m interface{}, query models.QParams) error {
	return p.Select().Where(dbx.HashExp(query.GetFilter())).OrderBy(query.GetOrder()).All(m)
}

func (p *Psql) Delete(m dbx.TableModel, ctx context.Context) error {

	err := p.Model(m).Delete()

	return err
}

func (p *Psql) Update(m dbx.TableModel, opts ...QueryOption) error {

	q := p.Model(m)

	for _, opt := range opts {
		q = opt(q)
	}

	err := q.Update()

	return err
}
