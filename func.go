// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	stdsql "database/sql"
)

// Default ORM

var DefaultORM *ORM = NewORM(nil)

func SetDB(db *stdsql.DB) {
	DefaultORM.DB = db
}

func GetDB() *stdsql.DB {
	return DefaultORM.DB
}

func Exec(query string, args ...interface{}) stdsql.Result {
	return DefaultORM.Exec(query, args...)
}

func Query(query string, args ...interface{}) *stdsql.Rows {
	return DefaultORM.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *stdsql.Row {
	return DefaultORM.QueryRow(query, args...)
}

func QueryOne(val interface{}, query string, args ...interface{}) bool {
	return DefaultORM.QueryOne(val, query, args...)
}

func Begin() (*ORM, error) {
	return DefaultORM.Begin()
}

func Commit() error {
	return DefaultORM.Commit()
}

func Rollback() error {
	return DefaultORM.Rollback()
}

func Select(model interface{}, sql *SQL, columns ...string) bool {
	return DefaultORM.Select(model, sql, columns...)
}

func Count(sql *SQL) int {
	return DefaultORM.Count(sql)
}

func CountMySQL(sql *SQL) int {
	return DefaultORM.CountMySQL(sql)
}

func Insert(model interface{}, columns ...string) stdsql.Result {
	return DefaultORM.Insert(model, columns...)
}

func Replace(model interface{}, columns ...string) stdsql.Result {
	return DefaultORM.Replace(model, columns...)
}

func Update(model interface{}, sql *SQL, columns ...string) stdsql.Result {
	return DefaultORM.Update(model, sql, columns...)
}

func Delete(model interface{}, sql *SQL) stdsql.Result {
	return DefaultORM.Delete(model, sql)
}

func BatchInsert(models interface{}, columns ...string) {
	DefaultORM.BatchInsert(models, columns...)
}

func BatchReplace(models interface{}, columns ...string) {
	DefaultORM.BatchReplace(models, columns...)
}

func Add(model interface{}, columns ...string) stdsql.Result {
	return DefaultORM.Add(model, columns...)
}

func Get(model interface{}, columns ...string) bool {
	return DefaultORM.Get(model, columns...)
}

func Up(model interface{}, columns ...string) stdsql.Result {
	return DefaultORM.Up(model, columns...)
}

func Del(model interface{}) stdsql.Result {
	return DefaultORM.Del(model)
}

func Save(model interface{}, columns ...string) stdsql.Result {
	return DefaultORM.Save(model, columns...)
}

func ForeignKey(models interface{}, foreign_key_column string, foreign_models interface{}, key_column string, columns ...string) {
	DefaultORM.ForeignKey(models, foreign_key_column, foreign_models, key_column, columns...)
}
