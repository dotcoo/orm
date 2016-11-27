// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var o *ORM

var init_sqls []string = []string{
	`DROP TABLE IF EXISTS test_user;`,
	`CREATE TABLE test_user (
	  id int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
	  username varchar(16) CHARACTER SET ascii NOT NULL COMMENT '用户名',
	  password varchar(32) CHARACTER SET ascii NOT NULL COMMENT '密码',
	  reg_time int(11) NOT NULL COMMENT '注册时间',
	  reg_ip int(11) NOT NULL COMMENT '注册IP',
	  update_time int(11) NOT NULL COMMENT '更新时间',
	  update_ip int(11) NOT NULL COMMENT '更新IP',
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';`,

	`DROP TABLE IF EXISTS test_category;`,
	`CREATE TABLE test_category (
	  id int(11) NOT NULL AUTO_INCREMENT COMMENT '分类编号',
	  name varchar(45) NOT NULL COMMENT '分类名称',
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='分类表';`,

	`DROP TABLE IF EXISTS test_blog;`,
	`CREATE TABLE test_blog (
	  id int(11) NOT NULL AUTO_INCREMENT COMMENT '域名ID',
	  category_id int(11) NOT NULL COMMENT '分类ID',
	  title varchar(45) NOT NULL COMMENT '标题',
	  content text NOT NULL COMMENT '内容',
	  add_time int(11) NOT NULL COMMENT '添加时间',
	  update_time int(11) NOT NULL COMMENT '更新时间',
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='博客表';`,
}

func init() {
	db, err := sql.Open("mysql", "root:123456@/mingo?charset=utf8")
	if err != nil {
		panic(err)
	}

	for _, init_sql := range init_sqls {
		db.Exec(init_sql)
	}

	o = NewORM(db)

	o.SetPrefix("test_")
}

func TestOrmInsert(t *testing.T) {
	u := new(User)
	u.Username = "dotcoo"
	u.Password = "dotcoopwd"
	result, err := o.RawInsert(u, "username, password")
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 || id != u.ID {
		t.Fatal("id != u.ID")
	}

	c := new(Category)
	c.Name = "Golang"
	result, err = o.RawInsert(c)
	if err != nil {
		t.Fatal(err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 || id != u.ID {
		t.Fatal("id != u.ID")
	}

	b := new(Blog)
	b.CategoryID = c.ID
	b.Title = "Golang ORM"
	b.Content = "Golang ORM Content"
	result, err = o.RawInsert(b, "category_id", "title", "content")
	if err != nil {
		t.Fatal(err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 || id != u.ID {
		t.Fatal("id != u.ID")
	}
}

func TestOrmUpdate(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoo123"
	result, err := o.RawUpdate(new(SQL).Where("id = ?", u.ID), u, "username, password")
	if err != nil {
		t.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("count != 1")
	}
}

func TestOrmSelect(t *testing.T) {
	u := new(User)
	u.ID = 1
	exist, err := o.RawSelect(new(SQL).Where("id = ?", u.ID), u)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatal("user not found")
	}

	users := make([]User, 0, 100)
	s := new(SQL)
	exist, err = o.RawSelect(s, &users, "id, username")
	if err != nil {
		t.Fatal(err)
	}
	if !exist || len(users) != 1 {
		t.Fatal("len(users) != 1")
	}

	// count
	count, err := o.RawCount(s.NewCount())
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("count != 1")
	}
}

func TestOrmSelectVal(t *testing.T) {
	s := o.NewSQL().From("user").Columns("count(*)", "sum(id)", "avg(id)")
	var count, sum int
	var avg float64
	exist, err := o.RawSelectVal(s, &count, &sum, &avg)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatal("row not found")
	}
	if count != 1 || sum != 1 || avg != 1 {
		t.Fatal("count != 1 || sum != 1 || avg != 1")
	}
}

func TestOrmReplace(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoo456"
	result, err := o.RawReplace(u)
	if err != nil {
		t.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal("count != 2")
	}
}

func TestOrmDelete(t *testing.T) {
	u := new(User)
	u.ID = 1
	result, err := o.RawDelete(new(SQL).Where("id = ?", u.ID), u)
	if err != nil {
		t.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("count != 1")
	}

	b := new(Blog)
	b.ID = 1
	result, err = o.RawDelete(new(SQL).Where("id = ?", b.ID), b)
	if err != nil {
		t.Fatal(err)
	}
	count, err = result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("count != 1")
	}
}

func TestOrmBatchInsert(t *testing.T) {
	blogs := make([]*Blog, 0, 100)
	for i := 1; i <= 100; i++ {
		b := new(Blog)
		b.Title = fmt.Sprintf("Golang ORM %d", i)
		b.Content = fmt.Sprintf("Golang ORM %d Content", i)
		blogs = append(blogs, b)
	}
	err := o.RawBatchInsert(&blogs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrmBatchReplace(t *testing.T) {
	blogs := make([]*Blog, 0, 100)
	for i := 1; i <= 100; i++ {
		b := new(Blog)
		b.ID = uint64(i)
		b.CategoryID = 1
		b.Title = fmt.Sprintf("Golang ORM %d Replace", i)
		b.Content = fmt.Sprintf("Golang ORM %d Content Replace", i)
		blogs = append(blogs, b)
	}
	err := o.RawBatchReplace(&blogs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrmAdd(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoopwd"
	result, err := o.RawAdd(u)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatal("id != 1")
	}
}

func TestOrmGet(t *testing.T) {
	u := new(User)
	u.ID = 1
	exist, err := o.RawGet(u)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatal("row not found")
	}
	if u.Username != "dotcoo" || u.Password != "dotcoopwd" {
		t.Fatal(`u.Username != "dotcoo" || u.Password != "dotcoopwd"`)
	}
}

func TestOrmGetBy(t *testing.T) {
	u := new(User)
	u.Username = "dotcoo"
	exist, err := o.RawGetBy(u, "username")
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatal("row not found")
	}
	if u.Username != "dotcoo" || u.Password != "dotcoopwd" {
		t.Fatal(`u.Username != "dotcoo" || u.Password != "dotcoopwd"`)
	}
}

func TestOrmUp(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoopwd2"
	result, err := o.RawUp(u, "username, password")
	if err != nil {
		t.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("count != 1")
	}
}

func TestOrmDel(t *testing.T) {
	u := new(User)
	u.ID = 1
	result, err := o.RawDel(u)
	if err != nil {
		t.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("count != 1")
	}
}

func TestOrmSave(t *testing.T) {
	_, err := o.RawExec("truncate test_user")
	if err != nil {
		t.Fatal(err)
	}

	u1 := new(User)
	u1.Username = "dotcoo"
	u1.Password = "dotcoopwd1"
	result, err := o.RawSave(u1, "*")
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatal("id != 1")
	}

	u2 := new(User)
	u2.ID = 1
	exist, err := o.RawGet(u2)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatal("row not found")
	}
	if u2.Username != "dotcoo" || u2.Password != "dotcoopwd1" {
		t.Fatal(`u2.Username != "dotcoo" || u2.Password != "dotcoopwd1"`)
	}

	u3 := new(User)
	u3.ID = 1
	u3.Username = "dotcoo"
	u3.Password = "dotcoopwd3"
	result, err = o.RawSave(u3, "*")
	if err != nil {
		t.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("count != 1")
	}

	u4 := new(User)
	u4.ID = 1
	exist, err = o.RawGet(u4)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatal("row not found")
	}
	if u4.Username != "dotcoo" || u4.Password != "dotcoopwd3" {
		t.Fatal(`u4.Username != "dotcoo" || u4.Password != "dotcoopwd3"`)
	}
}

func TestOrmForeignKey_Slice(t *testing.T) {
	blogs := make([]Blog, 0, 100)
	s := new(SQL).Where("id > ?", 10).Order("id").Page(3, 10)
	exist, err := o.RawSelect(s, &blogs)
	if err != nil {
		t.Fatal(err)
	}
	if !exist || len(blogs) != 10 {
		t.Fatal("len(blogs) != 10")
	}

	categorys := make([]Category, 0, len(blogs))
	err = o.RawForeignKey(&blogs, "category_id", &categorys, "id")
	if err != nil {
		t.Fatal(err)
	}
	if len(categorys) != 1 {
		t.Fatal("len(categorys) != 1")
	}
}

func TestOrmForeignKey_Map(t *testing.T) {
	blogs := make([]Blog, 0, 100)
	s := new(SQL).Where("id > ?", 10).Order("id").Page(3, 10)
	exist, err := o.RawSelect(s, &blogs)
	if err != nil {
		t.Fatal(err)
	}
	if !exist || len(blogs) != 10 {
		t.Fatal("len(blogs) != 10")
	}

	categorys := make(map[uint64]Category)
	err = o.RawForeignKey(&blogs, "category_id", &categorys, "id")
	if err != nil {
		t.Fatal(err)
	}
	if len(categorys) != 1 {
		t.Fatal("len(categorys) != 1")
	}
}

func TestOrmTransaction(t *testing.T) {
	u := new(User)
	u.ID = 1
	exist, err := o.RawGet(u)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		TestOrmAdd(t)
	}

	otx, err := o.RawBegin()
	if err != nil {
		panic(err)
	}

	u = new(User)
	u.ID = 1

	s := new(SQL).Where("id = ?", u.ID).ForUpdate()
	exist, err = otx.RawSelect(s, u)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		err = otx.RawRollback()
		if err != nil {
			t.Fatal(err)
		}
		t.Fatal("user 1 not found!")
	}

	u.Password = "haha"
	result, err := otx.RawUp(u, "password")
	if err != nil {
		t.Fatal(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("count != 1")
	}

	err = otx.RawCommit()
	if err != nil {
		panic(err)
	}
}
