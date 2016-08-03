// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
)

// default ORM method

func SetDB(db *sql.DB) {
	DefaultORM.db = db
}

func SetPrefix(prefix string) {
	DefaultORM.Manager().SetPrefix(prefix)
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

func QueryOne(val interface{}, query string, args ...interface{}) bool {
	return DefaultORM.QueryOne(val, query, args...)
}

func Select(model interface{}, s *SQL, columns ...string) bool {
	return DefaultORM.Select(model, s, columns...)
}

func Count(s *SQL) int {
	return DefaultORM.Count(s)
}

func CountMySQL(s *SQL) int {
	return DefaultORM.CountMySQL(s)
}

func Insert(model interface{}, columns ...string) sql.Result {
	return DefaultORM.Insert(model, columns...)
}

func Replace(model interface{}, columns ...string) sql.Result {
	return DefaultORM.Replace(model, columns...)
}

func Update(model interface{}, s *SQL, columns ...string) sql.Result {
	return DefaultORM.Update(model, s, columns...)
}

func Delete(model interface{}, s *SQL) sql.Result {
	return DefaultORM.Delete(model, s)
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
