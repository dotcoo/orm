package main

import (
	"database/sql"
	"log"
	"reflect"

	// ## import
	"github.com/dotcoo/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

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
	  PRIMARY KEY (id),
	  UNIQUE KEY username_UNIQUE (username)
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
	  user_id int(11) NOT NULL COMMENT '分类ID',
	  title varchar(45) NOT NULL COMMENT '标题',
	  content text NOT NULL COMMENT '内容',
	  add_time int(11) NOT NULL COMMENT '添加时间',
	  update_time int(11) NOT NULL COMMENT '更新时间',
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='博客表';`,
}

type User struct {
	ID         int64 `orm:"pk"`
	Username   string
	Password   string
	RegTime    int64 `orm:"created"`
	RegIP      uint32
	UpdateTime int64 `orm:"updated"`
	UpdateIP   uint32
	OterField  string `orm:"-"`
}

type Category struct {
	ID   int64 `orm:"pk"`
	Name string
}

type Blog struct {
	ID         int64 `orm:"pk"`
	CategoryID int64
	UserID     int64
	Title      string
	Content    string
	AddTime    int64 `orm:"created"`
	UpdateTime int64 `orm:"updated"`
}

func main() {
	// init
	db, err := sql.Open("mysql", "root:123456@/mingo?charset=utf8")
	if err != nil {
		panic(err)
	}

	for _, init_sql := range init_sqls {
		db.Exec(init_sql)
	}

	orm.SetDB(db)
	orm.SetPrefix("test_")

	var user *User
	var users []User
	var users_map map[int64]User
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
	// result = orm.Add(user, "id, username")
	// result = orm.Add(user, []string{"id", "username"}...)

	log.Println(result.LastInsertId())

	// ### Retrieve

	user = new(User)
	user.ID = 1

	ok = orm.Get(user)
	// ok = orm.Get(user, "id, username")
	// ok = orm.Get(user, []string{"id", "username"}...)

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
	// result = orm.Up(user, "id, username")
	// result = orm.Up(user, []string{"id", "username"}...)

	log.Println(result.RowsAffected())

	// ### Delete

	user = new(User)
	user.ID = 1

	result = orm.Del(user)

	log.Println(result.RowsAffected())

	// ### Save

	user = new(User)
	user.ID = 0
	user.Username = "dotcoo2"
	user.Password = "123456"

	result = orm.Save(user)
	// result = orm.Save(user, "id, username")
	// result = orm.Save(user, []string{"id", "username"}...)

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())

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
	// result = orm.Insert(user, "id, username")
	// result = orm.Insert(user, []string{"id", "username", "password"}...)

	log.Println(result.LastInsertId())

	// ### Select Row

	user = new(User)

	sq = orm.NewSQL().Where("username = ?", "dotcoo")

	ok = orm.Select(user, sq)
	// ok = orm.Select(user, sq, "id, username, password")
	// ok = orm.Select(user, sq, []string{"id", "username", "password"}...)

	if ok {
		log.Println(user)
	} else {
		log.Println("user not find")
	}

	// ### Select Rows

	// #### Slice

	users = make([]User, 0, 10)

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	ok = orm.Select(&users, sq)
	// ok = orm.Select(&users, sq, "id, username, password")
	// ok = orm.Select(&users, sq, []string{"id", "username", "password"}...)

	log.Println(users)

	// #### Map

	users_map = make(map[int64]User)

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	ok = orm.Select(&users_map, sq)
	// ok = orm.Select(&users_map, sq, "id, username, password")
	// ok = orm.Select(&users_map, sq, []string{"id", "username", "password"}...)

	log.Println(users_map)

	// ### Count

	n = orm.Count(sq)

	log.Println(n)

	// ### Update

	user = new(User)
	user.Password = "123321"

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	result = orm.Update(user, sq, "password")
	// result = orm.Update(user, sq, "*") // Warning: Update All Columns
	// result = orm.Update(user, sq, "username, password")
	// result = orm.Update(user, sq, []string{"username", "password"}...)

	log.Println(result.RowsAffected())

	// ### Delete

	user = new(User)

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	result = orm.Delete(user, sq)

	log.Println(result.RowsAffected())

	// ## SQL

	// ### Where

	log.Println(orm.NewSQL("test_user").Where("username = ?", "dotcoo").ToSelect())
	// SELECT * FROM test_user WHERE username = ? [dotcoo]

	// ### Where OR

	log.Println(orm.NewSQL("test_user").Where("username = ? or username = ?", "dotcoo", "dotcoo2").ToSelect())
	// SELECT * FROM test_user WHERE username = ? or username = ? [dotcoo dotcoo2]

	// ### Columns and Table

	log.Println(orm.NewSQL().Columns("id", "username").From("test_user").Where("username = ?", "dotcoo").ToSelect())
	// SELECT id, username FROM test_user WHERE username = ? [dotcoo]

	// ### Group

	log.Println(orm.NewSQL("test_user").Group("username").Having("id > ?", 100).ToSelect())
	// SELECT * FROM test_user GROUP BY username HAVING id > ? [100]

	// ### Order

	log.Println(orm.NewSQL("test_user").Group("username desc, id asc").ToSelect())
	// SELECT * FROM test_user GROUP BY username desc, id asc []

	// ### Limit Offset

	log.Println(orm.NewSQL("test_user").Limit(10).Offset(30).ToSelect())
	// SELECT * FROM test_user LIMIT 10 OFFSET 30 []

	// ### Update

	log.Println(orm.NewSQL("test_user").Set("password", "123123").Set("age", 28).Where("id = ?", 1).ToUpdate())
	// UPDATE test_user SET password = ?, age = ? WHERE id = ? [123123 28 1]

	// ### Delete

	log.Println(orm.NewSQL("test_user").Where("id = ?", 1).ToDelete())
	// DELETE FROM test_user WHERE id = ? [1]

	// ### Plus

	log.Println(orm.NewSQL("test_user").Plus("age", 1).Where("id = ?", 1).ToUpdate())
	// UPDATE test_user SET age = age + ? WHERE id = ? [1 1]

	// ### Incr

	log.Println(orm.NewSQL("test_user").Incr("age", 1).Where("id = ?", 1).ToUpdate())
	// UPDATE test_user SET age = last_insert_id(age + ?) WHERE id = ? [1 1]

	// ## Custom SQL

	// ### Exec

	result = orm.Exec("delete from test_user where id < ?", 10)

	// ### Query

	rows := orm.Query("select * from test_user where id < ?", 10)

	log.Println(rows)

	// ### QueryRow

	row := orm.Query("select * from test_user where id = ?", 10)

	log.Println(row)

	// ### QueryOne

	count := 0
	ok = orm.QueryOne(&count, "select count(*) as c from test_user")

	// ## Other

	// ### BatchInsert

	users = []User{
		{Username: "dotcoo3", Password: "123456", RegTime: 100},
		{Username: "dotcoo4", Password: "123456", RegTime: 101},
	}

	orm.BatchInsert(&users, "username, password, reg_time")
	// orm.BatchInsert(&users, []string{"username", "password"}...)

	// ### BatchReplace

	users = []User{
		{ID: 3, Username: "dotcoo3", Password: "654321"},
		{ID: 4, Username: "dotcoo4", Password: "654321"},
	}

	orm.BatchReplace(&users, "id, username, password")
	// orm.BatchReplace(&users, []string{"id", "username", "password"}...)

	// ### ForeignKey

	blogs := []Blog{
		{ID: 1, Title: "blog title 1", UserID: 3},
		{ID: 2, Title: "blog title 2", UserID: 4},
		{ID: 3, Title: "blog title 3", UserID: 3},
	}

	users_map = make(map[int64]User)

	orm.ForeignKey(&blogs, "user_id", &users_map, "id")
	// orm.ForeignKey(&blogs, "user_id", &users_map, "id", "id, username, password")
	// orm.ForeignKey(&blogs, "user_id", &users_map, "id", []string{"id", "username", "password"}...)

	for _, b := range blogs {
		log.Println(b.ID, b.Title, users_map[b.UserID].Username)
	}

	// ## Transaction

	o := orm.DefaultORM

	otx, _ := o.Begin()

	user = new(User)
	sq = otx.NewSQL().Where("id = ?", 3).ForUpdate()
	ok = otx.Select(user, sq)

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

	mi, ok := m.Get(reflect.ValueOf(user).Elem().Type())

	if ok {
		// if exist, modify table name
		mi.Table = "prefix_users"
	} else {
		// if not exist, modify table name
		mi = orm.NewModelInfo(user, "prefix_", "users")
		m.Set(mi)
	}
}
