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

var err error

var init_sqls []string = []string{
	`DROP TABLE IF EXISTS user;`,
	`CREATE TABLE user (
	  id int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
	  username varchar(16) CHARACTER SET ascii NOT NULL COMMENT '用户名',
	  password varchar(32) CHARACTER SET ascii NOT NULL COMMENT '密码',
	  reg_time int(11) NOT NULL COMMENT '注册时间',
	  reg_ip int(11) NOT NULL COMMENT '注册IP',
	  update_time int(11) NOT NULL COMMENT '更新时间',
	  update_ip int(11) NOT NULL COMMENT '更新IP',
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';`,

	`DROP TABLE IF EXISTS category;`,
	`CREATE TABLE category (
	  id int(11) NOT NULL AUTO_INCREMENT COMMENT '分类编号',
	  name varchar(45) NOT NULL COMMENT '分类名称',
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='分类表';`,

	`DROP TABLE IF EXISTS blog;`,
	`CREATE TABLE blog (
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
}

func TestOrmInsert(t *testing.T) {
	u := new(User)
	u.Username = "dotcoo"
	u.Password = "dotcoopwd"
	result, err := o.Insert(u, "username, password")
	if err != nil {
		panic(err)
	}
	u.ID, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}
	t.Log(u.ID)

	c := new(Category)
	c.Name = "Golang"
	result, err = o.Insert(c)
	if err != nil {
		panic(err)
	}
	c.ID, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}
	t.Log(c.ID)

	b := new(Blog)
	b.CategoryID = c.ID
	b.Title = "Golang ORM"
	b.Content = "Golang ORM Content"
	result, err = o.Insert(b, "category_id", "title", "content")
	b.ID, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}
	t.Log(b.ID)
}

func TestOrmUpdate(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoo123"
	result, err := o.Update(u, o.NewSQL().Where("id = ?", u.ID), "username, password")
	if err != nil {
		panic(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(count)
}

func TestOrmSelect(t *testing.T) {
	u := new(User)
	u.ID = 1
	ok, err := o.Select(u, o.NewSQL().Where("id = ?", u.ID))
	if err != nil {
		panic(err)
	}
	if !ok {
		panic("user not found")
	}
	t.Log(u)

	users := make([]User, 0, 100)
	_, err = o.Select(&users, o.NewSQL(), "id, username")
	if err != nil {
		panic(err)
	}
	t.Log(users)

	testOrmCount(t)

	testOrmCountMySQL(t)
}

func testOrmCount(t *testing.T) {
	blogs := make([]Blog, 0, 100)
	s := o.NewSQL().Where("id > ?", 10).Order("id").Page(3, 10)
	_, err := o.Select(&blogs, s)
	if err != nil {
		panic(err)
	}
	t.Log(blogs)

	count, err := s.Count()
	if err != nil {
		panic(err)
	}
	t.Log(count)
}

func testOrmCountMySQL(t *testing.T) {
	blogs := make([]Blog, 0, 100)
	s := o.NewSQL().CalcFoundRows().Where("id > ?", 0).Order("id").Page(3, 10)
	_, err := o.Select(&blogs, s)
	if err != nil {
		panic(err)
	}
	t.Log(blogs)

	count, err := s.CountMySQL()
	if err != nil {
		panic(err)
	}
	t.Log(count)
}

func TestOrmReplace(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoo456"
	result, err := o.Replace(u)
	if err != nil {
		panic(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(count)
}

func TestOrmDelete(t *testing.T) {
	u := new(User)
	u.ID = 1
	result, err := o.Delete(u, o.NewSQL().Where("id = ?", u.ID))
	if err != nil {
		panic(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(count)

	b := new(Blog)
	b.ID = 1
	result, err = o.Delete(b, o.NewSQL().Where("id = ?", b.ID))
	if err != nil {
		panic(err)
	}
	count, err = result.RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(count)
}

func TestOrmBatchInsert(t *testing.T) {
	blogs := make([]*Blog, 0, 100)
	for i := 1; i <= 100; i++ {
		b := new(Blog)
		b.Title = fmt.Sprintf("Golang ORM %d", i)
		b.Content = fmt.Sprintf("Golang ORM %d Content", i)
		blogs = append(blogs, b)
	}
	err := o.BatchInsert(&blogs)
	if err != nil {
		panic(err)
	}
}

func TestOrmBatchReplace(t *testing.T) {
	blogs := make([]*Blog, 0, 100)
	for i := 1; i <= 100; i++ {
		b := new(Blog)
		b.ID = int64(i)
		b.CategoryID = 1
		b.Title = fmt.Sprintf("Golang ORM %d", i)
		b.Content = fmt.Sprintf("Golang ORM %d Content", i)
		blogs = append(blogs, b)
	}
	err := o.BatchReplace(&blogs)
	if err != nil {
		panic(err)
	}
}

func TestOrmAdd(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoopwd"
	result, err := o.Add(u)
	if err != nil {
		panic(err)
	}
	t.Log(result.LastInsertId())
}

func TestOrmGet(t *testing.T) {
	u := new(User)
	u.ID = 1
	exist, err := o.Get(u)
	if err != nil {
		panic(err)
	}
	t.Log(u, exist)
}

func TestOrmUp(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoopwd"
	result, err := o.Up(u, "username, password")
	if err != nil {
		panic(err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(u, n, err)
}

func TestOrmDel(t *testing.T) {
	u := new(User)
	u.ID = 1
	result, err := o.Del(u)
	if err != nil {
		panic(err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(u, n, err)
}

func TestOrmSave(t *testing.T) {
	u1 := new(User)
	u1.ID = 1

	result, err := o.Del(u1)
	if err != nil {
		panic(err)
	}
	id, id_err := result.LastInsertId()
	row, row_err := result.RowsAffected()
	t.Log(u1, id, id_err, row, row_err)

	exist, err := o.Get(u1)
	if err != nil {
		panic(err)
	}
	t.Log(u1, exist)

	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoopwd2"
	result, err = o.Save(u, "*")
	if err != nil {
		panic(err)
	}
	id, id_err = result.LastInsertId()
	row, row_err = result.RowsAffected()
	t.Log(u, id, id_err, row, row_err)

	u2 := new(User)
	u2.ID = 1
	exist, err = o.Get(u2)
	if err != nil {
		panic(err)
	}
	t.Log(u2, exist)

	u.Password = "dotcoopwd3"
	result, err = o.Save(u, "*")
	if err != nil {
		panic(err)
	}
	id, id_err = result.LastInsertId()
	row, row_err = result.RowsAffected()
	t.Log(u, id, id_err, row, row_err)

	u3 := new(User)
	u3.ID = 1
	exist, err = o.Get(u3)
	if err != nil {
		panic(err)
	}
	t.Log(u3, exist)
}

func TestOrmForeignKey(t *testing.T) {
	blogs := make([]Blog, 0, 100)
	s := o.NewSQL().Where("id > ?", 10).Order("id").Page(3, 10)
	_, err := o.Select(&blogs, s)
	if err != nil {
		panic(err)
	}
	t.Log(blogs)

	categorys := make([]Category, 0, 20)
	err = o.ForeignKey(&blogs, "category_id", &categorys, "id")
	if err != nil {
		panic(err)
	}
	t.Log(categorys)
}

func TestOrmTransaction(t *testing.T) {
	TestOrmAdd(t)

	otx, err := o.Begin()
	if err != nil {
		panic(err)
	}

	u := new(User)
	u.ID = 1
	s := otx.NewSQL().Where("id = ?", u.ID).ForUpdate()

	ok, err := o.Select(u, s)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic("user 1 not found!")
	}

	u.Password = "haha"
	result, err := otx.Up(u, "password")
	if err != nil {
		panic(err)
	}
	t.Log(result.RowsAffected())

	err = otx.Commit()
	if err != nil {
		panic(err)
	}
}
