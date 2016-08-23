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
	o.Insert(u, "username, password")
	t.Log(u.ID)

	c := new(Category)
	c.Name = "Golang"
	o.Insert(c)
	t.Log(c.ID)

	b := new(Blog)
	b.CategoryID = c.ID
	b.Title = "Golang ORM"
	b.Content = "Golang ORM Content"
	o.Insert(b, "category_id", "title", "content")
	t.Log(b.ID)
}

func TestOrmUpdate(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoo123"
	count, err := o.Update(o.NewSQL().Where("id = ?", u.ID), u, "username, password").RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(count)
}

func TestOrmSelect(t *testing.T) {
	u := new(User)
	u.ID = 1
	if !o.Select(o.NewSQL().Where("id = ?", u.ID), u) {
		panic("user not found")
	}
	t.Log(u)

	users := make([]User, 0, 100)
	o.Select(o.NewSQL(), &users, "id, username")
	t.Log(users)

	testOrmCount(t)
}

func TestOrmSelectRow(t *testing.T) {
	s := o.NewSQL("user").Columns("count(*)", "sum(id)", "avg(id)")
	var count, sum int
	var avg float64
	o.SelectRow(s, &count, &sum, &avg)
	t.Log(count, sum, avg)
}

func testOrmCount(t *testing.T) {
	blogs := make([]Blog, 0, 100)
	s := o.NewSQL().Where("id > ?", 10).Order("id").Page(3, 10)
	o.Select(s, &blogs)
	t.Log(blogs)

	count := s.Count()
	t.Log(count)
}

func TestOrmReplace(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoo456"
	count, err := o.Replace(u).RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(count)
}

func TestOrmDelete(t *testing.T) {
	u := new(User)
	u.ID = 1
	count, err := o.Delete(o.NewSQL().Where("id = ?", u.ID), u).RowsAffected()
	if err != nil {
		panic(err)
	}
	t.Log(count)

	b := new(Blog)
	b.ID = 1
	count, err = o.Delete(o.NewSQL().Where("id = ?", b.ID), b).RowsAffected()
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
		b.ID = uint64(i)
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
	t.Log(o.Add(u).LastInsertId())
}

func TestOrmGet(t *testing.T) {
	u := new(User)
	u.ID = 1
	exist := o.Get(u)
	t.Log(u, exist)
}

func TestOrmGetBy(t *testing.T) {
	u := new(User)
	u.Username = "dotcoo"
	exist := o.GetBy(u, "username")
	t.Log(u, exist)
}

func TestOrmUp(t *testing.T) {
	u := new(User)
	u.ID = 1
	u.Username = "dotcoo"
	u.Password = "dotcoopwd2"
	n, err := o.Up(u, "username, password").RowsAffected()
	t.Log(u, n, err)
}

func TestOrmDel(t *testing.T) {
	u := new(User)
	u.ID = 1
	n, err := o.Del(u).RowsAffected()
	t.Log(u, n, err)
}

func TestOrmSave(t *testing.T) {
	u1 := new(User)
	u1.ID = 1
	u1.Username = "dotcoo"
	o.Add(u1)

	result := o.Del(u1)
	id, id_err := result.LastInsertId()
	row, row_err := result.RowsAffected()
	t.Log(u1, id, id_err, row, row_err)

	exist := o.Get(u1)
	t.Log(u1, exist)

	u := new(User)
	u.Username = "dotcoo"
	u.Password = "dotcoopwd3"
	result = o.Save(u, "*")
	id, id_err = result.LastInsertId()
	row, row_err = result.RowsAffected()
	t.Log(u, id, id_err, row, row_err)

	u2 := new(User)
	u2.ID = 2
	exist = o.Get(u2)
	t.Log(u2, exist)

	u.Password = "dotcoopwd4"
	result = o.Save(u, "*")
	id, id_err = result.LastInsertId()
	row, row_err = result.RowsAffected()
	t.Log(u, id, id_err, row, row_err)

	u3 := new(User)
	u3.ID = 2
	exist = o.Get(u3)
	t.Log(u3, exist)
}

func TestOrmForeignKey_Slice(t *testing.T) {
	blogs := make([]Blog, 0, 100)
	s := o.NewSQL().Where("id > ?", 10).Order("id").Page(3, 10)
	o.Select(s, &blogs)
	t.Log(blogs)

	categorys := make([]Category, 0, 20)
	err := o.ForeignKey(&blogs, "category_id", &categorys, "id")
	if err != nil {
		panic(err)
	}
	t.Log(categorys)
}

func TestOrmForeignKey_Map(t *testing.T) {
	blogs := make([]*Blog, 0, 100)
	s := o.NewSQL().Where("id > ?", 10).Order("id").Page(3, 10)
	o.Select(s, &blogs)
	t.Log(blogs)

	categorys := make(map[uint64]Category)
	err := o.ForeignKey(&blogs, "category_id", &categorys, "id")
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
	if !otx.Select(s, u) {
		panic("user 1 not found!")
	}

	u.Password = "haha"
	otx.Up(u, "password")

	err = otx.Commit()
	if err != nil {
		panic(err)
	}
}
