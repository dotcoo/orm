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
	sql, params := s.ToSelect()
	sql_select := "SELECT SQL_NO_CACHE SQL_CALC_FOUND_ROWS username, password, email, count(*) as count FROM `user` WHERE username = ? AND age BETWEEN ? AND ? AND no IN (?, ?, ?, ?, ?) GROUP BY age HAVING count > ? AND count < ? ORDER BY id DESC, username, password DESC LIMIT 10 OFFSET 20 FOR UPDATE LOCK IN SHARE MODE"
	params_select := []interface{}{"dotcoo", 18, 25, 1, 2, 3, 4, 5, 3, 10}
	if sql != sql_select || !reflect.DeepEqual(params, params_select) {
		t.Errorf("sql_select error: %s, %v", sql, params)
	}

	sql, params = s.ToCount()
	sql_count := "SELECT count(*) AS count FROM `user` WHERE username = ? AND age BETWEEN ? AND ? AND no IN (?, ?, ?, ?, ?)"
	params_count := []interface{}{"dotcoo", 18, 25, 1, 2, 3, 4, 5}
	if sql != sql_count || !reflect.DeepEqual(params, params_count) {
		t.Errorf("sql_count error: %s, %v", sql, params)
	}

	sql, params = s.ToCountMySQL()
	sql_count_mysql := "SELECT FOUND_ROWS()"
	params_count_mysql := []interface{}{}
	if sql != sql_count_mysql || !reflect.DeepEqual(params, params_count_mysql) {
		t.Errorf("sql_count_mysql error: %s, %v", sql, params)
	}

	sql, params = s.From("blog", "b").Join("user", "u", "b.user_id = u.id").Where("b.start > ?", 200).Page(3, 10).ToSelect()
	sql_join := "SELECT * FROM `blog` AS `b` LEFT JOIN `user` AS `u` ON b.user_id = u.id WHERE b.start > ? LIMIT 10 OFFSET 20"
	params_join := []interface{}{200}
	if sql != sql_join || !reflect.DeepEqual(params, params_join) {
		t.Errorf("sql_join error: %s, %v", sql, params)
	}

	sql, params = s.From("user").Set("username", "dotcoo").Set("password", "dotcoopwd").Set("age", 1).Where("id = ?", 1).ToUpdate()
	sql_update := "UPDATE `user` SET `username` = ?, `password` = ?, `age` = ? WHERE id = ?"
	params_update := []interface{}{"dotcoo", "dotcoopwd", 1, 1}
	if sql != sql_update || !reflect.DeepEqual(params, params_update) {
		t.Errorf("sql_update error: %s, %v", sql, params)
	}

	sql, params = s.Set("username", "dotcoo").Set("password", "dotcoopwd").Set("age", 1).ToReplace()
	sql_replace := "REPLACE `user` SET `username` = ?, `password` = ?, `age` = ?"
	params_replace := []interface{}{"dotcoo", "dotcoopwd", 1}
	if sql != sql_replace || !reflect.DeepEqual(params, params_replace) {
		t.Errorf("sql_replace error: %s, %v", sql, params)
	}

	sql, params = s.Where("id = ?", 1).ToDelete()
	sql_delete := "DELETE FROM `user` WHERE id = ?"
	params_delete := []interface{}{1}
	if sql != sql_delete || !reflect.DeepEqual(params, params_delete) {
		t.Errorf("sql_delete error: %s, %v", sql, params)
	}

	sql, params = s.Plus("age", 1).Where("id = ?", 1).ToUpdate()
	sql_plus := "UPDATE `user` SET `age` = `age` + ? WHERE id = ?"
	params_plus := []interface{}{1, 1}
	if sql != sql_plus || !reflect.DeepEqual(params, params_plus) {
		t.Errorf("sql_plus error: %s, %v", sql, params)
	}

	sql, params = s.Incr("age", 1).Where("id = ?", 1).ToUpdate()
	sql_incr := "UPDATE `user` SET `age` = last_insert_id(`age` + ?) WHERE id = ?"
	params_incr := []interface{}{1, 1}
	if sql != sql_incr || !reflect.DeepEqual(params, params_incr) {
		t.Errorf("sql_incr error: %s, %v", sql, params)
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
