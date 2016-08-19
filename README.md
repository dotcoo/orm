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

	import "github.com/dotcoo/orm"

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

	orm.SetDB(db)

### Prefix

	orm.SetPrefix("test_")

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

	result = orm.Add(user)
	// result = orm.Add(user, "id, username")
	// result = orm.Add(user, []string{"id", "username"}...)

	log.Println(result.LastInsertId())

### Read/Retrieve

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

#### GetBy

	user = new(User)
	user.Username = "dotcoo"
	ok = orm.GetBy(user, "username")
	// ok = orm.GetBy(user, "id, username")
	// ok = orm.GetBy(user, []string{"id", "username"}...)

	if ok {
		log.Println(user)
	} else {
		log.Println("user not find")
	}

### Update

	user = new(User)
	user.ID = 1
	user.Password = "654321"

	result = orm.Up(user, "password")
	// result = orm.Up(user, "id, username")
	// result = orm.Up(user, []string{"id", "username"}...)

	log.Println(result.RowsAffected())

### Delete

	user = new(User)
	user.ID = 1

	result = orm.Del(user)

	log.Println(result.RowsAffected())

### Save

	// insert
	user = new(User)
	user.ID = 0
	user.Username = "dotcoo2"
	user.Password = "123456"

	result = orm.Save(user)
	// result = orm.Save(user, "id, username")
	// result = orm.Save(user, []string{"id", "username"}...)

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

## SQL CRUD

### Insert

	user = new(User)
	user.ID = 1
	user.Username = "dotcoo"
	user.Password = "123456"

	result = orm.Insert(user, "id, username, password")
	// result = orm.Insert(user, "id, username")
	// result = orm.Insert(user, []string{"id", "username", "password"}...)

	log.Println(result.LastInsertId())

### Select Row

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

### Select Rows

#### Slice

	users = make([]User, 0, 10)

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	ok = orm.Select(&users, sq)
	// ok = orm.Select(&users, sq, "id, username, password")
	// ok = orm.Select(&users, sq, []string{"id", "username", "password"}...)

	log.Println(users)

#### Map

	users_map = make(map[int64]User)

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	ok = orm.Select(&users_map, sq)
	// ok = orm.Select(&users_map, sq, "id, username, password")
	// ok = orm.Select(&users_map, sq, []string{"id", "username", "password"}...)

	log.Println(users_map)

### Count

	n = orm.Count(sq)

	log.Println(n)

### Update

	user = new(User)
	user.Password = "123321"

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	result = orm.Update(user, sq, "password")
	// result = orm.Update(user, sq, "*") // Warning: Update All Columns
	// result = orm.Update(user, sq, "username, password")
	// result = orm.Update(user, sq, []string{"username", "password"}...)

	log.Println(result.RowsAffected())

### Delete

	user = new(User)

	sq = orm.NewSQL().Where("username like ?", "dotcoo%")

	result = orm.Delete(user, sq)

	log.Println(result.RowsAffected())

## SQL

### Where

	orm.NewSQL("test_user").Where("username = ?", "dotcoo").ToSelect()
	// SELECT * FROM test_user WHERE username = ? [dotcoo]

### Where OR

	orm.NewSQL("test_user").Where("username = ? or username = ?", "dotcoo", "dotcoo2").ToSelect()
	// SELECT * FROM test_user WHERE username = ? or username = ? [dotcoo dotcoo2]

### Columns and Table

	orm.NewSQL().Columns("id", "username").From("test_user").Where("username = ?", "dotcoo").ToSelect()
	// SELECT id, username FROM test_user WHERE username = ? [dotcoo]

### Group

	orm.NewSQL("test_user").Group("username").Having("id > ?", 100).ToSelect()
	// SELECT * FROM test_user GROUP BY username HAVING id > ? [100]

### Order

	orm.NewSQL("test_user").Group("username desc, id asc").ToSelect()
	// SELECT * FROM test_user GROUP BY username desc, id asc []

### Limit Offset

	orm.NewSQL("test_user").Limit(10).Offset(30).ToSelect()
	// SELECT * FROM test_user LIMIT 10 OFFSET 30 []

### Update

	orm.NewSQL("test_user").Set("password", "123123").Set("age", 28).Where("id = ?", 1).ToUpdate()
	// UPDATE test_user SET password = ?, age = ? WHERE id = ? [123123 28 1]

### Delete

	orm.NewSQL("test_user").Where("id = ?", 1).ToDelete()
	// DELETE FROM test_user WHERE id = ? [1]

### Plus

	orm.NewSQL("test_user").Plus("age", 1).Where("id = ?", 1).ToUpdate()
	// UPDATE test_user SET age = age + ? WHERE id = ? [1 1]

### Incr

	orm.NewSQL("test_user").Incr("age", 1).Where("id = ?", 1).ToUpdate()
	// UPDATE test_user SET age = last_insert_id(age + ?) WHERE id = ? [1 1]

## Custom SQL

### Exec

	result, err = orm.Exec("delete from test_user where id < ?", 10)

### Query

	rows, err := orm.Query("select * from test_user where id < ?", 10)

### QueryRow

	row, err := orm.Query("select * from test_user where id = ?", 10)

### QueryOne

	count := 0
	ok, err = orm.QueryOne(&count, "select count(*) as c from test_user")

## Other Method

### BatchInsert

	users = []User{
		{Username: "dotcoo3", Password: "123456", RegTime: 100},
		{Username: "dotcoo4", Password: "123456", RegTime: 101},
	}

	orm.BatchInsert(&users, "username, password, reg_time")
	// orm.BatchInsert(&users, []string{"username", "password"}...)

### BatchReplace

	users = []User{
		{ID: 3, Username: "dotcoo3", Password: "654321"},
		{ID: 4, Username: "dotcoo4", Password: "654321"},
	}

	orm.BatchReplace(&users, "id, username, password")
	// orm.BatchReplace(&users, []string{"id", "username", "password"}...)

### ForeignKey

	blogs := []Blog{
		{ID: 1, Title: "blog title 1", UserID: 3},
		{ID: 2, Title: "blog title 2", UserID: 4},
		{ID: 3, Title: "blog title 3", UserID: 3},
	}

	users_map = make(map[int64]User)

	err = orm.ForeignKey(&blogs, "user_id", &users_map, "id")
	// err = orm.ForeignKey(&blogs, "user_id", &users_map, "id", "id, username, password")
	// err = orm.ForeignKey(&blogs, "user_id", &users_map, "id", []string{"id", "username", "password"}...)
	if err != nil {
		panic(err)
	}

	for _, b := range blogs {
		log.Println(b.ID, b.Title, users_map[b.UserID].Username)
	}

## Transaction

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

## ModelInfo

	m := orm.DefaultORM.Manager()

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

