// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
)

// No Error

func (o *ORM) Exec(query string, args ...interface{}) sql.Result {
	result, err := o.RawExec(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Query(query string, args ...interface{}) *sql.Rows {
	rows, err := o.RawQuery(query, args...)
	if err != nil {
		panic(err)
	}
	return rows
}

func (o *ORM) QueryRow(query string, args ...interface{}) *sql.Row {
	row, err := o.RawQueryRow(query, args...)
	if err != nil {
		panic(err)
	}
	return row
}

func (o *ORM) Begin() *ORM {
	otx, err := o.RawBegin()
	if err != nil {
		panic(err)
	}
	return otx
}

func (o *ORM) Commit() {
	err := o.RawCommit()
	if err != nil {
		panic(err)
	}
}

func (o *ORM) Rollback() {
	err := o.RawRollback()
	if err != nil {
		panic(err)
	}
}

func (o *ORM) Select(s *SQL, model interface{}, columns ...string) bool {
	exist, err := o.RawSelect(s, model, columns...)
	if err != nil {
		panic(err)
	}
	return exist
}

func (o *ORM) SelectVal(s *SQL, vals ...interface{}) bool {
	exist, err := o.RawSelectVal(s, vals...)
	if err != nil {
		panic(err)
	}
	return exist
}

func (o *ORM) Count(s *SQL) int {
	count, err := o.RawCount(s)
	if err != nil {
		panic(err)
	}
	return count
}

func (o *ORM) Insert(model interface{}, columns ...string) sql.Result {
	result, err := o.RawInsert(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Replace(model interface{}, columns ...string) sql.Result {
	result, err := o.RawReplace(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Update(s *SQL, model interface{}, columns ...string) sql.Result {
	result, err := o.RawUpdate(s, model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Delete(s *SQL, model interface{}) sql.Result {
	result, err := o.RawDelete(s, model)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) BatchInsert(models interface{}, columns ...string) {
	err := o.RawBatchInsert(models, columns...)
	if err != nil {
		panic(err)
	}
}

func (o *ORM) BatchReplace(models interface{}, columns ...string) {
	err := o.RawBatchReplace(models, columns...)
	if err != nil {
		panic(err)
	}
}

func (o *ORM) Add(model interface{}, columns ...string) sql.Result {
	result, err := o.RawAdd(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Get(model interface{}, columns ...string) bool {
	exist, err := o.RawGet(model, columns...)
	if err != nil {
		panic(err)
	}
	return exist
}

func (o *ORM) GetBy(model interface{}, columns ...string) bool {
	exist, err := o.RawGetBy(model, columns...)
	if err != nil {
		panic(err)
	}
	return exist
}

func (o *ORM) Up(model interface{}, columns ...string) sql.Result {
	result, err := o.RawUp(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Del(model interface{}) sql.Result {
	result, err := o.RawDel(model)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Save(model interface{}, columns ...string) sql.Result {
	result, err := o.RawSave(model, columns...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) ForeignKey(sources interface{}, fk_column string, models interface{}, pk_column string, columns ...string) {
	err := o.RawForeignKey(sources, fk_column, models, pk_column, columns...)
	if err != nil {
		panic(err)
	}
}
