// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ormweb

import (
	"database/sql"

	"github.com/dotcoo/orm"
)

type ORM struct {
	*orm.ORM
}

var DefaultORM *ORM = NewORM(nil)

func NewORM(db *sql.DB) *ORM {
	o := new(ORM)
	o.ORM = orm.NewORM(db)
	return o
}

func (o *ORM) SetPrefix(prefix string) {
	o.ORM.SetPrefix(prefix)
}

// query

func (o *ORM) Exec(query string, args ...interface{}) sql.Result {
	result, err := o.ORM.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Query(query string, args ...interface{}) *sql.Rows {
	rows, err := o.ORM.Query(query, args...)
	if err != nil {
		panic(err)
	}
	return rows
}

func (o *ORM) QueryRow(query string, args ...interface{}) *sql.Row {
	row, err := o.ORM.QueryRow(query, args...)
	if err != nil {
		panic(err)
	}
	return row
}

func (o *ORM) QueryOne(val interface{}, query string, args ...interface{}) bool {
	ok, err := o.ORM.QueryOne(val, query, args...)
	if err != nil {
		panic(err)
	}
	return ok
}

// transaction

func (o *ORM) Begin() (otx *ORM, err error) {
	otx = new(ORM)
	otx.ORM, err = o.ORM.Begin()
	return otx, err
}

func (o *ORM) Commit() error {
	return o.ORM.Commit()
}

func (o *ORM) Rollback() error {
	return o.ORM.Rollback()
}

// select

func (o *ORM) Manager() *orm.ModelInfoManager {
	return o.ORM.Manager()
}

func (o *ORM) Select(model interface{}, s *orm.SQL, columns ...string) bool {
	exist, err := o.ORM.Select(model, s, columns...)
	if err != nil {
		panic(err)
	}
	return exist
}

func (o *ORM) Count(s *orm.SQL) (count int) {
	var err error
	count, err = o.ORM.Count(s)
	if err != nil {
		panic(err)
	}
	return count
}

func (o *ORM) CountMySQL(s *orm.SQL) (count int) {
	var err error
	count, err = o.ORM.CountMySQL(s)
	if err != nil {
		panic(err)
	}
	return count
}

func (o *ORM) Insert(model interface{}, columns ...string) sql.Result {
	result, err := o.ORM.Insert(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Replace(model interface{}, columns ...string) sql.Result {
	result, err := o.ORM.Replace(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Update(model interface{}, s *orm.SQL, columns ...string) sql.Result {
	result, err := o.ORM.Update(model, s, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Delete(model interface{}, s *orm.SQL) sql.Result {
	result, err := o.ORM.Delete(model, s)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) BatchInsert(models interface{}, columns ...string) {
	err := o.ORM.BatchInsert(models, columns...)
	if err != nil {
		panic(err)
	}
}

func (o *ORM) BatchReplace(models interface{}, columns ...string) {
	err := o.ORM.BatchReplace(models, columns...)
	if err != nil {
		panic(err)
	}
}

// quick method

func (o *ORM) Add(model interface{}, columns ...string) sql.Result {
	result, err := o.ORM.Add(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Get(model interface{}, columns ...string) bool {
	result, err := o.ORM.Get(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Up(model interface{}, columns ...string) sql.Result {
	result, err := o.ORM.Up(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Del(model interface{}) sql.Result {
	result, err := o.ORM.Del(model)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Save(model interface{}, columns ...string) sql.Result {
	result, err := o.ORM.Save(model, columns ...)
	if err != nil {
		panic(err)
	}
	return result
}

// foreign key

func (o *ORM) ForeignKey(sources interface{}, fk_column string, models interface{}, pk_column string, columns ...string) {
	err := o.ORM.ForeignKey(sources, fk_column, models, pk_column, columns...)
	if err != nil {
		panic(err)
	}
}

// SQL

func (o *ORM) NewSQL(table ...string) *orm.SQL {
	return o.ORM.NewSQL(table...)
}
