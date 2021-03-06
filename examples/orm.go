// # Golang ORM

// ORM library for Go Golang

package main

import (
	"database/sql"
	"log"
	"reflect"

	// ## Import
	"github.com/dotcoo/orm"
	_ "github.com/go-sql-driver/mysql"
)

// ## Environment

// ### Database

var init_sqls []string = []string{
	`DROP TABLE IF EXISTS test_user;`,
	`CREATE TABLE test_user (
	  id int(11) NOT NULL AUTO_INCREMENT,
	  username varchar(16) CHARACTER SET ascii NOT NULL,
	  password varchar(32) CHARACTER SET ascii NOT NULL,
	  reg_time int(11) NOT NULL,
	  reg_ip int(11) NOT NULL,
	  update_time int(11) NOT NULL,
	  update_ip int(11) NOT NULL,
	  PRIMARY KEY (id),
	  UNIQUE KEY username_UNIQUE (username)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

	`DROP TABLE IF EXISTS test_blog;`,
	`CREATE TABLE test_blog (
	  id int(11) NOT NULL AUTO_INCREMENT,
	  user_id int(11) NOT NULL,
	  title varchar(45) NOT NULL,
	  content text NOT NULL,
	  add_time int(11) NOT NULL,
	  update_time int(11) NOT NULL,
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// ## Models

type User struct {
	ID         int `orm:"pk"`
	Username   string
	Password   string
	RegTime    int `orm:"created"`
	RegIP      uint32
	UpdateTime int `orm:"updated"`
	UpdateIP   uint32
	OtherField string `orm:"-"`
}

type Blog struct {
	ID         int `orm:"pk"`
	UserID     int
	Title      string
	Content    string
	AddTime    int `orm:"created"`
	UpdateTime int `orm:"updated"`
}

func main() {
	// ## Initialization

	// ### SetDB

	db, err := sql.Open("mysql", "root:123456@/mingo?charset=utf8")
	if err != nil {
		panic(err)
	}
	orm.SetDB(db)

	for _, init_sql := range init_sqls {
		db.Exec(init_sql)
	}

	// ### SetPrefix

	orm.SetPrefix("test_")

	// ### Variables

	var user *User
	var users []User
	var users_map map[int]User
	var result sql.Result
	var ok bool
	var sq *orm.SQL
	var n int

	// ## CRUD

	// ### Create

	user = new(User)
	user.Username = "dotcoo"
	user.Password = "123456"

	result = orm.Add(user)
	// result = orm.Add(user, "id", "username")

	log.Println(user.ID)
	log.Println(result.LastInsertId())

	// ### Read/Retrieve

	user = new(User)
	user.ID = 1

	ok = orm.Get(user)
	// ok = orm.Get(user, "id", "username")

	if ok {
		log.Println(user)
	} else {
		log.Println("user not find")
	}

	// #### GetBy

	user = new(User)
	user.Username = "dotcoo"
	ok = orm.GetBy(user, "username")
	// ok = orm.GetBy(user, "id", "username")

	if ok {
		log.Println(user)
	} else {
		log.Println("user not find")
	}

	// ### Update

	user = new(User)
	user.ID = 1
	user.Password = "654321"

	result = orm.Up(user, "password")
	// result = orm.Up(user, "id", "username")

	log.Println(result.RowsAffected())

	// ### Delete

	user = new(User)
	user.ID = 1

	result = orm.Del(user)

	log.Println(result.RowsAffected())

	// ### Save

	// insert
	user = new(User)
	user.ID = 0
	user.Username = "dotcoo2"
	user.Password = "123456"

	result = orm.Save(user)
	// result = orm.Save(user, "id", "username")

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())

	// update
	user = new(User)
	user.ID = 2
	user.Username = "dotcoo2"
	user.Password = "654321"

	result = orm.Save(user, "username, password")

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())

	// ## SQL CRUD

	// ### Insert

	user = new(User)
	user.ID = 1
	user.Username = "dotcoo"
	user.Password = "123456"

	result = orm.Insert(user, "id, username, password")
	// result = orm.Insert(user, "id", "username", "password")

	log.Println(result.LastInsertId())

	// ### Select Model

	sq = orm.NewSQL().Where("username = ?", "dotcoo")

	user = new(User)

	ok = orm.Select(sq, user)
	// ok = orm.Select(sq, user, "id", "username", "password")

	if ok {
		log.Println(user)
	} else {
		log.Println("user not find")
	}

	// ### Select Models

	// #### Slice

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	users = make([]User, 0, 10)

	ok = orm.Select(sq, &users)
	// ok = orm.Select(sq, &users, "id", "username", "password")

	log.Println(users)

	// #### Map

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	users_map = make(map[int]User)

	ok = orm.Select(sq, &users_map)
	// ok = orm.Select(sq, &users_map, "id", "username", "password")

	log.Println(users_map)

	// ### Count

	n = orm.Count(sq)

	log.Println(n)

	// ### Select Row

	sq = orm.NewSQL().From("user").Columns("count(*)", "sum(id)", "avg(id)")

	var count, sum int
	var avg float64

	ok = orm.SelectVal(sq, &count, &sum, &avg)

	log.Println(count, sum, avg)

	// ### Update

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	user = new(User)
	user.Password = "123321"

	result = orm.Update(sq, user, "password")
	// result = orm.Update(sq, user) // Error
	// result = orm.Update(sq, user, "*") // Correct

	log.Println(result.RowsAffected())

	// ### Delete

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	user = new(User)

	result = orm.Delete(sq, user)

	log.Println(result.RowsAffected())

	// ## SQL

	// ### Where

	log.Println(orm.NewSQL().From("user").Where("username = ?", "dotcoo").ToSelect())
	// SELECT * FROM `test_user` WHERE username = ? [dotcoo]

	// ### Where OR

	log.Println(orm.NewSQL().From("user").Where("(username = ?", "dotcoo").Where(" OR username = ?", "dotcoo2)").ToSelect())
	// SELECT * FROM `test_user` WHERE (username = ? OR username = ?) [dotcoo dotcoo2]

	// ### Columns and Table

	log.Println(orm.NewSQL().Columns("id", "username").From("user").Where("username = ?", "dotcoo").ToSelect())
	// SELECT id, username FROM `test_user` WHERE username = ? [dotcoo]

	// ### Group

	log.Println(orm.NewSQL().From("user").Group("username").Having("id > ?", 100).ToSelect())
	// SELECT * FROM `test_user` GROUP BY username HAVING id > ? [100]

	// ### Order

	log.Println(orm.NewSQL().From("user").Group("username desc, id asc").ToSelect())
	// SELECT * FROM `test_user` GROUP BY username desc, id asc []

	// ### Limit Offset

	log.Println(orm.NewSQL().From("user").Limit(10).Offset(30).ToSelect())
	// SELECT * FROM `test_user` LIMIT 10 OFFSET 30 []

	// ### Join

	log.Println(orm.NewSQL().From("blog AS b").Join("user AS u", "b.user_id = u.id").ToSelect())
	// SELECT * FROM `test_blog` AS b LEFT JOIN `test_user` AS u ON b.user_id = u.id []

	// ### Update

	log.Println(orm.NewSQL().From("user").Set("password", "123123").Set("age", 28).Where("id = ?", 1).ToUpdate())
	// UPDATE `test_user` SET `password` = ?, `age` = ? WHERE id = ? [123123 28 1]

	// ### Delete

	log.Println(orm.NewSQL().From("user").Where("id = ?", 1).ToDelete())
	// DELETE FROM `test_user` WHERE id = ? [1]

	// ### Plus

	log.Println(orm.NewSQL().From("user").Plus("age", 1).Where("id = ?", 1).ToUpdate())
	// UPDATE `test_user` SET `age` = `age` + ? WHERE id = ? [1 1]

	// ### Incr

	log.Println(orm.NewSQL().From("user").Incr("age", 1).Where("id = ?", 1).ToUpdate())
	// UPDATE `test_user` SET `age` = last_insert_id(`age` + ?) WHERE id = ? [1 1]

	// ## Custom SQL

	// ### Exec

	result = orm.Exec("delete from test_user where id < ?", 10)
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())

	// ### Query

	rows := orm.Query("select * from test_user where id < ?", 10)
	log.Println(rows)

	// ### QueryRow

	row := orm.Query("select * from test_user where id = ?", 10)
	log.Println(row)

	// ## Other Method

	// ### BatchInsert

	users = []User{
		{Username: "dotcoo3", Password: "123456", RegTime: 100},
		{Username: "dotcoo4", Password: "123456", RegTime: 101},
	}

	orm.BatchInsert(&users, "username, password, reg_time")
	// orm.BatchInsert(&users, "username", "password", "reg_time")

	// ### BatchReplace

	users = []User{
		{ID: 3, Username: "dotcoo3", Password: "654321"},
		{ID: 4, Username: "dotcoo4", Password: "654321"},
	}

	orm.BatchReplace(&users, "id, username, password")
	// orm.BatchReplace(&users, "id", "username", "password")

	// ### ForeignKey

	blogs := []Blog{
		{ID: 1, Title: "blog title 1", UserID: 3},
		{ID: 2, Title: "blog title 2", UserID: 4},
		{ID: 3, Title: "blog title 3", UserID: 3},
	}

	users_map = make(map[int]User)

	orm.ForeignKey(&blogs, "user_id", &users_map, "id")
	// orm.ForeignKey(&blogs, "user_id", &users_map, "id", "id", "username", "password")

	for _, b := range blogs {
		log.Println(b.ID, b.Title, users_map[b.UserID].Username)
	}

	// ## Transaction

	o := orm.DefaultORM

	otx := o.Begin()

	sq = otx.NewSQL().Where("id = ?", 3).ForUpdate()
	user = new(User)
	ok = otx.Select(sq, user)

	if !ok {
		otx.Rollback()
		log.Println("Rollback")
	} else {
		user.RegTime++
		otx.Up(user, "reg_time")

		otx.Commit()
		log.Println("Commit")
	}

	// ## ModelInfo

	m := orm.DefaultORM.Manager()

	user = new(User)

	mi := m.Get(reflect.ValueOf(user).Elem().Type())

	if mi == nil {
		// if not exist, modify table name
		mi = orm.NewModelInfo(user, "prefix_", "users")
	} else {
		// if exist, modify table name
		mi.Table = "prefix_users"
	}

	m.Set(mi)
}
