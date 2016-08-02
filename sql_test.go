// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"reflect"
	"testing"
)

func TestSQLSelect(t *testing.T) {
	s := NewSQL().
		Keywords("SQL_NO_CACHE").
		CalcFoundRows().
		Columns("username, password", "email", "count(*) as count").
		From("user").
		Where("username = ?", "dotcoo").
		Where("age BETWEEN ? AND ?", 18, 25).
		WhereIn("no IN (?)", 1, 2, 3, 4, 5).
		Group("age").
		Having("count > ?", 3).
		Having("count < ?", 10).
		Order("id DESC, username", "password DESC").
		Limit(10).
		Offset(20).
		ForUpdate().
		LockInShareMode()
	sq, params := s.ToSelect()
	sq_select := "SELECT SQL_NO_CACHE SQL_CALC_FOUND_ROWS username, password, email, count(*) as count FROM `user` WHERE username = ? AND age BETWEEN ? AND ? AND no IN (?, ?, ?, ?, ?) GROUP BY age HAVING count > ? AND count < ? ORDER BY id DESC, username, password DESC LIMIT 10 OFFSET 20 FOR UPDATE LOCK IN SHARE MODE"
	params_select := []interface{}{"dotcoo", 18, 25, 1, 2, 3, 4, 5, 3, 10}
	if sq != sq_select || !reflect.DeepEqual(params, params_select) {
		t.Errorf("sq_select error: %s, %v", sq, params)
	}

	sq, params = s.ToCount()
	sq_count := "SELECT count(*) AS count FROM `user` WHERE username = ? AND age BETWEEN ? AND ? AND no IN (?, ?, ?, ?, ?)"
	params_count := []interface{}{"dotcoo", 18, 25, 1, 2, 3, 4, 5}
	if sq != sq_count || !reflect.DeepEqual(params, params_count) {
		t.Errorf("sq_count error: %s, %v", sq, params)
	}

	sq, params = s.ToCountMySQL()
	sq_count_mysql := "SELECT FOUND_ROWS()"
	params_count_mysql := []interface{}{}
	if sq != sq_count_mysql || !reflect.DeepEqual(params, params_count_mysql) {
		t.Errorf("sq_count_mysql error: %s, %v", sq, params)
	}

	sq, params = s.From("blog", "b").Join("user", "u", "b.user_id = u.id").Where("b.start > ?", 200).Page(3, 10).ToSelect()
	sq_join := "SELECT * FROM `blog` AS `b` LEFT JOIN `user` AS `u` ON b.user_id = u.id WHERE b.start > ? LIMIT 10 OFFSET 20"
	params_join := []interface{}{200}
	if sq != sq_join || !reflect.DeepEqual(params, params_join) {
		t.Errorf("sq_join error: %s, %v", sq, params)
	}

	sq, params = s.From("user").Set("username", "dotcoo").Set("password", "dotcoopwd").Set("age", 1).Where("id = ?", 1).ToUpdate()
	sq_update := "UPDATE `user` SET `username` = ?, `password` = ?, `age` = ? WHERE id = ?"
	params_update := []interface{}{"dotcoo", "dotcoopwd", 1, 1}
	if sq != sq_update || !reflect.DeepEqual(params, params_update) {
		t.Errorf("sq_update error: %s, %v", sq, params)
	}

	sq, params = s.Set("username", "dotcoo").Set("password", "dotcoopwd").Set("age", 1).ToReplace()
	sq_replace := "REPLACE `user` SET `username` = ?, `password` = ?, `age` = ?"
	params_replace := []interface{}{"dotcoo", "dotcoopwd", 1}
	if sq != sq_replace || !reflect.DeepEqual(params, params_replace) {
		t.Errorf("sq_replace error: %s, %v", sq, params)
	}

	sq, params = s.Where("id = ?", 1).ToDelete()
	sq_delete := "DELETE FROM `user` WHERE id = ?"
	params_delete := []interface{}{1}
	if sq != sq_delete || !reflect.DeepEqual(params, params_delete) {
		t.Errorf("sq_delete error: %s, %v", sq, params)
	}

	sq, params = s.Plus("age", 1).Where("id = ?", 1).ToUpdate()
	sq_plus := "UPDATE `user` SET `age` = `age` + ? WHERE id = ?"
	params_plus := []interface{}{1, 1}
	if sq != sq_plus || !reflect.DeepEqual(params, params_plus) {
		t.Errorf("sq_plus error: %s, %v", sq, params)
	}

	sq, params = s.Incr("age", 1).Where("id = ?", 1).ToUpdate()
	sq_incr := "UPDATE `user` SET `age` = last_insert_id(`age` + ?) WHERE id = ?"
	params_incr := []interface{}{1, 1}
	if sq != sq_incr || !reflect.DeepEqual(params, params_incr) {
		t.Errorf("sq_incr error: %s, %v", sq, params)
	}
}

func BenchmarkSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSQL()
		s.Keywords("SQL_NO_CACHE").
			CalcFoundRows().
			Columns("username, password", "email", "count(*) as count").
			From("user").
			Where("username = ?", "dotcoo").
			Where("age BETWEEN ? AND ?", 18, 25).
			WhereIn("no IN (?)", 1, 2, 3, 4, 5).
			// Where("no IN (?)", []int{1, 2, 3, 4, 5}).
			Group("age").
			Having("count > ?", 3).
			Having("count < ?", 10).
			Order("id DESC, username", "password DESC").
			Limit(10).
			Offset(20).
			ForUpdate().
			LockInShareMode()
		s.ToSelect()
	}
}
