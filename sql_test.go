// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"reflect"
	"testing"
)

func TestSQLSelect(t *testing.T) {
	// all method
	s := new(SQL).
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

	// count
	sq, params = s.NewCount().ToSelect()
	sq_count := "SELECT count(*) AS count FROM `user` WHERE username = ? AND age BETWEEN ? AND ? AND no IN (?, ?, ?, ?, ?) GROUP BY age HAVING count > ? AND count < ?"
	params_count := []interface{}{"dotcoo", 18, 25, 1, 2, 3, 4, 5, 3, 10}
	if sq != sq_count || !reflect.DeepEqual(params, params_count) {
		t.Errorf("sq_count error: %s, %v", sq, params)
	}

	// where or
	sq, params = new(SQL).From("user").Where("username = ?", "dotcoo").Where("(uid = ?", 1).Where(" OR uid = ?", 2).Where(" OR uid BETWEEN ? AND ?", 5, 9).WhereIn(" OR uid IN (?)", 11, 12, 13, 14, 15).Where(" OR uid >= ?)", 20).ToSelect()
	sq_where := "SELECT * FROM `user` WHERE username = ? AND (uid = ? OR uid = ? OR uid BETWEEN ? AND ? OR uid IN (?, ?, ?, ?, ?) OR uid >= ?)"
	params_where := []interface{}{"dotcoo", 1, 2, 5, 9, 11, 12, 13, 14, 15, 20}
	if sq != sq_where || !reflect.DeepEqual(params, params_where) {
		t.Errorf("sq_where error: %s, %v", sq, params)
	}

	// select
	sq, params = new(SQL).From("blog AS b").Join("user AS u", "b.user_id = u.id").Where("b.start > ?", 200).Page(3, 10).ToSelect()
	sq_join := "SELECT * FROM `blog` AS `b` LEFT JOIN `user` AS `u` ON b.user_id = u.id WHERE b.start > ? LIMIT 10 OFFSET 20"
	params_join := []interface{}{200}
	if sq != sq_join || !reflect.DeepEqual(params, params_join) {
		t.Errorf("sq_join error: %s, %v", sq, params)
	}

	// insert
	sq, params = new(SQL).From("user").Set("username", "dotcoo").Set("password", "dotcoopwd").Set("age", 1).ToInsert()
	sq_insert := "INSERT INTO `user` (`username`, `password`, `age`) VALUES (?, ?, ?)"
	params_insert := []interface{}{"dotcoo", "dotcoopwd", 1}
	if sq != sq_insert || !reflect.DeepEqual(params, params_insert) {
		t.Errorf("sq_insert error: %s, %v", sq, params)
	}

	// replace
	sq, params = new(SQL).From("user").Set("username", "dotcoo").Set("password", "dotcoopwd").Set("age", 1).ToReplace()
	sq_replace := "REPLACE INTO `user` (`username`, `password`, `age`) VALUES (?, ?, ?)"
	params_replace := []interface{}{"dotcoo", "dotcoopwd", 1}
	if sq != sq_replace || !reflect.DeepEqual(params, params_replace) {
		t.Errorf("sq_replace error: %s, %v", sq, params)
	}

	// update
	sq, params = new(SQL).From("user").Set("username", "dotcoo").Set("password", "dotcoopwd").Set("age", 1).Where("id = ?", 1).ToUpdate()
	sq_update := "UPDATE `user` SET `username` = ?, `password` = ?, `age` = ? WHERE id = ?"
	params_update := []interface{}{"dotcoo", "dotcoopwd", 1, 1}
	if sq != sq_update || !reflect.DeepEqual(params, params_update) {
		t.Errorf("sq_update error: %s, %v", sq, params)
	}

	// delete
	sq, params = new(SQL).From("user").Where("id = ?", 1).ToDelete()
	sq_delete := "DELETE FROM `user` WHERE id = ?"
	params_delete := []interface{}{1}
	if sq != sq_delete || !reflect.DeepEqual(params, params_delete) {
		t.Errorf("sq_delete error: %s, %v", sq, params)
	}

	// page
	sq, params = new(SQL).From("user").Page(3, 10).ToSelect()
	sq_page := "SELECT * FROM `user` LIMIT 10 OFFSET 20"
	params_page := []interface{}{}
	if sq != sq_page || !reflect.DeepEqual(params, params_page) {
		t.Errorf("sq_page error: %s, %v", sq, params)
	}

	// plus
	sq, params = new(SQL).From("user").Plus("age", 1).Where("id = ?", 1).ToUpdate()
	sq_plus := "UPDATE `user` SET `age` = `age` + ? WHERE id = ?"
	params_plus := []interface{}{1, 1}
	if sq != sq_plus || !reflect.DeepEqual(params, params_plus) {
		t.Errorf("sq_plus error: %s, %v", sq, params)
	}

	// incr
	sq, params = new(SQL).From("user").Incr("age", 1).Where("id = ?", 1).ToUpdate()
	sq_incr := "UPDATE `user` SET `age` = last_insert_id(`age` + ?) WHERE id = ?"
	params_incr := []interface{}{1, 1}
	if sq != sq_incr || !reflect.DeepEqual(params, params_incr) {
		t.Errorf("sq_incr error: %s, %v", sq, params)
	}
}

func BenchmarkSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		new(SQL).
			Keywords("SQL_NO_CACHE").
			CalcFoundRows().
			Columns("username, password", "email", "count(*) as count").
			From("user").
			Where("username = ?", "dotcoo").
			Where("age BETWEEN ? AND ?", 18, 25).
			WhereIn("no IN (?)", 1, 2, 3, 4, 5).
			// WhereIn("no IN (?)", []int{1, 2, 3, 4, 5}...).
			Group("age").
			Having("count > ?", 3).
			Having("count < ?", 10).
			Order("id DESC, username", "password DESC").
			Limit(10).
			Offset(20).
			ForUpdate().
			LockInShareMode().
			ToSelect()
		new(SQL).
			From("user").
			Set("username", "dotcoo").
			Set("password", "dotcoopwd").
			Set("col1", "col1").
			Set("col2", "col2").
			Set("col3", "col3").
			Set("col4", "col4").
			Set("col5", "col5").
			Set("col6", "col6").
			Set("col7", "col7").
			Set("col8", "col8").
			Set("col9", "col9").
			Set("col0", "col0").
			ToInsert()
		new(SQL).
			From("user").
			Set("id", 1).
			Set("username", "dotcoo").
			Set("password", "dotcoopwd").
			Set("col1", "col1").
			Set("col2", "col2").
			Set("col3", "col3").
			Set("col4", "col4").
			Set("col5", "col5").
			Set("col6", "col6").
			Set("col7", "col7").
			Set("col8", "col8").
			Set("col9", "col9").
			Set("col0", "col0").
			ToReplace()
		new(SQL).
			From("user").
			Set("password", "dotcoopwd").
			Set("col1", "col1").
			Set("col2", "col2").
			Set("col3", "col3").
			Where("username = ?", "dotcoo").
			ToUpdate()
		new(SQL).
			From("user").
			Where("username = ?", "dotcoo").
			ToDelete()
	}
}
