# Golang ORM

ORM library for Go Golang

## Environment

### Database

	DROP TABLE IF EXISTS test_user;
	CREATE TABLE test_user (
	  id int(11) NOT NULL AUTO_INCREMENT,
	  username varchar(16) CHARACTER SET ascii NOT NULL,
	  password varchar(32) CHARACTER SET ascii NOT NULL,
	  reg_time int(11) NOT NULL,
	  reg_ip int(11) NOT NULL,
	  update_time int(11) NOT NULL,
	  update_ip int(11) NOT NULL,
	  PRIMARY KEY (id),
	  UNIQUE KEY username_UNIQUE (username)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;

	DROP TABLE IF EXISTS test_blog;
	CREATE TABLE test_blog (
	  id int(11) NOT NULL AUTO_INCREMENT,
	  user_id int(11) NOT NULL,
	  title varchar(45) NOT NULL,
	  content text NOT NULL,
	  add_time int(11) NOT NULL,
	  update_time int(11) NOT NULL,
	  PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
}

## Initialization

### import

	import (
		"github.com/dotcoo/orm"
		"github.com/dotcoo/orm/ormweb"
	)

### struct

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

	type Blog struct {
		ID         int64 `orm:"pk"`
		UserID     int64
		Title      string
		Content    string
		AddTime    int64 `orm:"created"`
		UpdateTime int64 `orm:"updated"`
	}

### SetDB

	db, err := sql.Open("mysql", "root:123456@/mingo?charset=utf8")
	if err != nil {
		panic(err)
	}

	ormweb.SetDB(db)

### Prefix

	ormweb.SetPrefix("test_")

### Variables

	var user *User
	var users []User
	var users_map map[int64]User
	var result sql.Result
	var ok bool
	var sq *orm.SQL
	var n int

## CRUD

### Create

	user = new(User)
	user.Username = "dotcoo"
	user.Password = "123456"

	result = ormweb.Add(user)
	// result = ormweb.Add(user, "id, username")
	// result = ormweb.Add(user, []string{"id", "username"}...)

	log.Println(result.LastInsertId())

### Read/Retrieve

	user = new(User)
	user.ID = 1

	ok = ormweb.Get(user)
	// ok = ormweb.Get(user, "id, username")
	// ok = ormweb.Get(user, []string{"id", "username"}...)

	if ok {
		log.Println(user)
	} else {
		log.Println("user not find")
	}

### Update

	user = new(User)
	user.ID = 1
	user.Password = "654321"

	result = ormweb.Up(user, "password")
	// result = ormweb.Up(user, "id, username")
	// result = ormweb.Up(user, []string{"id", "username"}...)

	log.Println(result.RowsAffected())

### Delete

	user = new(User)
	user.ID = 1

	result = ormweb.Del(user)

	log.Println(result.RowsAffected())

### Save

	// insert
	user = new(User)
	user.ID = 0
	user.Username = "dotcoo2"
	user.Password = "123456"

	result = ormweb.Save(user)
	// result = ormweb.Save(user, "id, username")
	// result = ormweb.Save(user, []string{"id", "username"}...)

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())

	// update
	user = new(User)
	user.ID = 2
	user.Username = "dotcoo2"
	user.Password = "654321"

	result = ormweb.Save(user, "username, password")

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())

## SQL CRUD

### Insert

	user = new(User)
	user.ID = 1
	user.Username = "dotcoo"
	user.Password = "123456"

	result = ormweb.Insert(user, "id, username, password")
	// result = ormweb.Insert(user, "id, username")
	// result = ormweb.Insert(user, []string{"id", "username", "password"}...)

	log.Println(result.LastInsertId())

### Select Row

	user = new(User)

	sq = ormweb.NewSQL().Where("username = ?", "dotcoo")

	ok = ormweb.Select(user, sq)
	// ok = ormweb.Select(user, sq, "id, username, password")
	// ok = ormweb.Select(user, sq, []string{"id", "username", "password"}...)

	if ok {
		log.Println(user)
	} else {
		log.Println("user not find")
	}

### Select Rows

#### Slice

	users = make([]User, 0, 10)

	sq = ormweb.NewSQL().Where("username like ?", "dotcoo%")

	ok = ormweb.Select(&users, sq)
	// ok = ormweb.Select(&users, sq, "id, username, password")
	// ok = ormweb.Select(&users, sq, []string{"id", "username", "password"}...)

	log.Println(users)

#### Map

	users_map = make(map[int64]User)

	sq = ormweb.NewSQL().Where("username like ?", "dotcoo%")

	ok = ormweb.Select(&users_map, sq)
	// ok = ormweb.Select(&users_map, sq, "id, username, password")
	// ok = ormweb.Select(&users_map, sq, []string{"id", "username", "password"}...)

	log.Println(users_map)

### Count

	n = ormweb.Count(sq)

	log.Println(n)

### Update

	user = new(User)
	user.Password = "123321"

	sq = ormweb.NewSQL().Where("username like ?", "dotcoo%")

	result = ormweb.Update(user, sq, "password")
	// result = ormweb.Update(user, sq, "*") // Warning: Update All Columns
	// result = ormweb.Update(user, sq, "username, password")
	// result = ormweb.Update(user, sq, []string{"username", "password"}...)

	log.Println(result.RowsAffected())

### Delete

	user = new(User)

	sq = ormweb.NewSQL().Where("username like ?", "dotcoo%")

	result = ormweb.Delete(user, sq)

	log.Println(result.RowsAffected())

## SQL

### Where

	ormweb.NewSQL("test_user").Where("username = ?", "dotcoo").ToSelect()
	// SELECT * FROM test_user WHERE username = ? [dotcoo]

### Where OR

	ormweb.NewSQL("test_user").Where("username = ? or username = ?", "dotcoo", "dotcoo2").ToSelect()
	// SELECT * FROM test_user WHERE username = ? or username = ? [dotcoo dotcoo2]

### Columns and Table

	ormweb.NewSQL().Columns("id", "username").From("test_user").Where("username = ?", "dotcoo").ToSelect()
	// SELECT id, username FROM test_user WHERE username = ? [dotcoo]

### Group

	ormweb.NewSQL("test_user").Group("username").Having("id > ?", 100).ToSelect()
	// SELECT * FROM test_user GROUP BY username HAVING id > ? [100]

### Order

	ormweb.NewSQL("test_user").Group("username desc, id asc").ToSelect()
	// SELECT * FROM test_user GROUP BY username desc, id asc []

### Limit Offset

	ormweb.NewSQL("test_user").Limit(10).Offset(30).ToSelect()
	// SELECT * FROM test_user LIMIT 10 OFFSET 30 []

### Update

	ormweb.NewSQL("test_user").Set("password", "123123").Set("age", 28).Where("id = ?", 1).ToUpdate()
	// UPDATE test_user SET password = ?, age = ? WHERE id = ? [123123 28 1]

### Delete

	ormweb.NewSQL("test_user").Where("id = ?", 1).ToDelete()
	// DELETE FROM test_user WHERE id = ? [1]

### Plus

	ormweb.NewSQL("test_user").Plus("age", 1).Where("id = ?", 1).ToUpdate()
	// UPDATE test_user SET age = age + ? WHERE id = ? [1 1]

### Incr

	ormweb.NewSQL("test_user").Incr("age", 1).Where("id = ?", 1).ToUpdate()
	// UPDATE test_user SET age = last_insert_id(age + ?) WHERE id = ? [1 1]

## Custom SQL

### Exec

	result = ormweb.Exec("delete from test_user where id < ?", 10)

### Query

	rows := ormweb.Query("select * from test_user where id < ?", 10)

### QueryRow

	row := ormweb.Query("select * from test_user where id = ?", 10)

### QueryOne

	count := 0
	ok = ormweb.QueryOne(&count, "select count(*) as c from test_user")

## Other Method

### BatchInsert

	users = []User{
		{Username: "dotcoo3", Password: "123456", RegTime: 100},
		{Username: "dotcoo4", Password: "123456", RegTime: 101},
	}

	ormweb.BatchInsert(&users, "username, password, reg_time")
	// ormweb.BatchInsert(&users, []string{"username", "password"}...)

### BatchReplace

	users = []User{
		{ID: 3, Username: "dotcoo3", Password: "654321"},
		{ID: 4, Username: "dotcoo4", Password: "654321"},
	}

	ormweb.BatchReplace(&users, "id, username, password")
	// ormweb.BatchReplace(&users, []string{"id", "username", "password"}...)

### ForeignKey

	blogs := []Blog{
		{ID: 1, Title: "blog title 1", UserID: 3},
		{ID: 2, Title: "blog title 2", UserID: 4},
		{ID: 3, Title: "blog title 3", UserID: 3},
	}

	users_map = make(map[int64]User)

	ormweb.ForeignKey(&blogs, "user_id", &users_map, "id")
	// ormweb.ForeignKey(&blogs, "user_id", &users_map, "id", "id, username, password")
	// ormweb.ForeignKey(&blogs, "user_id", &users_map, "id", []string{"id", "username", "password"}...)

	for _, b := range blogs {
		log.Println(b.ID, b.Title, users_map[b.UserID].Username)
	}

## Transaction

	o := ormweb.DefaultORM

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

## ModelInfo

	m := ormweb.DefaultORM.Manager()

	user = new(User)

	mi, ok := m.Get(reflect.ValueOf(user).Elem().Type())

	if ok {
		// if exist, modify table name
		mi.Table = "prefix_users"
	} else {
		// if not exist, modify table name
		mi, err = orm.NewModelInfo(user, "prefix_", "users")
		if err != nil {
			panic(err)
		}
		m.Set(mi)
	}

