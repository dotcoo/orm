// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
)

// sql and orm

func (s *SQL) RawExec() (sql.Result, error) {
	query, args := s.SQL()
	return s.orm.RawExec(query, args...)
}

func (s *SQL) RawQuery() (*sql.Rows, error) {
	query, args := s.SQL()
	return s.orm.RawQuery(query, args...)
}

func (s *SQL) RawQueryRow() (*sql.Row, error) {
	query, args := s.SQL()
	return s.orm.RawQueryRow(query, args...)
}

func (s *SQL) RawSelect(model interface{}, columns ...string) (bool, error) {
	return s.orm.RawSelect(s, model, columns...)
}

func (s *SQL) RawSelectVal(vals ...interface{}) (bool, error) {
	return s.orm.RawSelectVal(s, vals...)
}

func (s *SQL) RawCount() (int, error) {
	return s.orm.RawCount(s)
}

func (s *SQL) RawUpdate(model interface{}, columns ...string) (sql.Result, error) {
	return s.orm.RawUpdate(s, model, columns...)
}

func (s *SQL) RawDelete(model interface{}) (sql.Result, error) {
	return s.orm.RawDelete(s, model)
}

func (s *SQL) Exec() sql.Result {
	query, args := s.SQL()
	return s.orm.Exec(query, args...)
}

func (s *SQL) Query() *sql.Rows {
	query, args := s.SQL()
	return s.orm.Query(query, args...)
}

func (s *SQL) QueryRow() *sql.Row {
	query, args := s.SQL()
	return s.orm.QueryRow(query, args...)
}

func (s *SQL) Select(model interface{}, columns ...string) bool {
	return s.orm.Select(s, model, columns...)
}

func (s *SQL) SelectVal(vals ...interface{}) bool {
	return s.orm.SelectVal(s, vals...)
}

func (s *SQL) Count() int {
	return s.orm.Count(s)
}

func (s *SQL) Update(model interface{}, columns ...string) sql.Result {
	return s.orm.Update(s, model, columns...)
}

func (s *SQL) Delete(model interface{}) sql.Result {
	return s.orm.Delete(s, model)
}
