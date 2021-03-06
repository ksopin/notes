package db

import (
	"database/sql"
	"github.com/pkg/errors"
)

type Row interface {
	EmptyCopy() Row
	GetId() uint
	SetId(id uint)
	InsertArgs() []interface{}
	UpdateArgs() []interface{}
}


type Crud struct {
	db *sql.DB
	prototype Row
	sqlInsert string
	sqlUpdate string
	sqlDelete string
}
func NewCrud(db *sql.DB, prototype Row, sqlInsert, sqlUpdate, sqlDelete string) *Crud {
	return &Crud{
		db: db,
		prototype: prototype,
		sqlInsert: sqlInsert,
		sqlUpdate: sqlUpdate,
		sqlDelete: sqlUpdate,
	}
}

func (t *Crud) Save(r Row) error {
	return t.SaveTx(t.db, r)
}

func (t *Crud) SaveTx(e Execer, r Row) error {
	if r.GetId() > 0 {
		return t.UpdateTx(e, r)
	} else {
		return t.InsertTx(e, r)
	}
}

func (t *Crud) Insert(r Row) error {
	return t.InsertTx(t.db, r)
}

func (t *Crud) InsertTx(e Execer, r Row) error {
	res, err := e.Exec(t.sqlInsert, r.InsertArgs()...)
	if err != nil {
		return errors.WithStack(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return errors.WithStack(err)
	}
	r.SetId(uint(id))

	return nil
}

func (t *Crud) Update(r Row) error {
	return t.UpdateTx(t.db, r)
}

func (t *Crud) UpdateTx(e Execer, r Row) error {
	res, err := e.Exec(t.sqlUpdate, r.UpdateArgs()...)
	if err == nil {
		_, err = res.RowsAffected()
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (t *Crud) Delete(id uint) error {
	return t.DeleteTx(t.db, id)
}

func (t *Crud) DeleteTx(e Execer, id uint) error {
	res, err := e.Exec(t.sqlDelete, id)
	if err == nil {
		_, err = res.RowsAffected()
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
