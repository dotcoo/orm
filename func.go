// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
)

// default ORM method

var DefaultORM *ORM = NewORM(nil)

func SetDB(db *sql.DB) {
	DefaultORM.SetDB(db)
}

func SetPrefix(prefix string) {
	DefaultORM.SetPrefix(prefix)
}

func NewSQL() *SQL {
	return DefaultORM.NewSQL()
}

func Exec(query string, args ...interface{}) sql.Result {
	return DefaultORM.Exec(query, args...)
}

func Query(query string, args ...interface{}) *sql.Rows {
	return DefaultORM.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	return DefaultORM.QueryRow(query, args...)
}

func Begin() *ORM {
	return DefaultORM.Begin()
}

func Commit() {
	DefaultORM.Commit()
}

func Rollback() {
	DefaultORM.Rollback()
}

func Select(s *SQL, model interface{}, columns ...string) bool {
	return DefaultORM.Select(s, model, columns...)
}

func SelectVal(s *SQL, vals ...interface{}) bool {
	return DefaultORM.SelectVal(s, vals...)
}

func Count(s *SQL) int {
	return DefaultORM.Count(s)
}

func Insert(model interface{}, columns ...string) sql.Result {
	return DefaultORM.Insert(model, columns...)
}

func Replace(model interface{}, columns ...string) sql.Result {
	return DefaultORM.Replace(model, columns...)
}

func Update(s *SQL, model interface{}, columns ...string) sql.Result {
	return DefaultORM.Update(s, model, columns...)
}

func Delete(s *SQL, model interface{}) sql.Result {
	return DefaultORM.Delete(s, model)
}

func BatchInsert(models interface{}, columns ...string) {
	DefaultORM.BatchInsert(models, columns...)
}

func BatchReplace(models interface{}, columns ...string) {
	DefaultORM.BatchReplace(models, columns...)
}

func Add(model interface{}, columns ...string) sql.Result {
	return DefaultORM.Add(model, columns...)
}

func Get(model interface{}, columns ...string) bool {
	return DefaultORM.Get(model, columns...)
}

func GetBy(model interface{}, columns ...string) bool {
	return DefaultORM.GetBy(model, columns...)
}

func Up(model interface{}, columns ...string) sql.Result {
	return DefaultORM.Up(model, columns...)
}

func Del(model interface{}) sql.Result {
	return DefaultORM.Del(model)
}

func Save(model interface{}, columns ...string) sql.Result {
	return DefaultORM.Save(model, columns...)
}

func ForeignKey(sources interface{}, fk_column string, models interface{}, pk_column string, columns ...string) {
	DefaultORM.ForeignKey(sources, fk_column, models, pk_column, columns...)
}
