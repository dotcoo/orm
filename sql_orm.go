// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
)

// sql and orm

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

func (s *SQL) Select(model interface{}, columns ...string) bool {
	return s.orm.Select(s, model, columns...)
}

func (s *SQL) SelectVal(vals ...interface{}) bool {
	return s.orm.SelectVal(s, vals...)
}

func (s *SQL) SelectCount(model interface{}, columns ...string) (bool, int) {
	return s.orm.Select(s, model, columns...), s.orm.Count(s)
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
