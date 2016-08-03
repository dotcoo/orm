// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"reflect"
	"testing"
)

type User struct {
	ID         int64 `orm:"pk"`
	Username   string
	Password   string
	RegTime    int64 `orm:"created"`
	RegIP      uint32
	UpdateTime int64 `orm:"updated"`
	UpdateIP   uint32
	OtherField string `orm:"-"`
}

type Category struct {
	ID   int64 `orm:"pk"`
	Name string
}

type Blog struct {
	ID         int64 `orm:"pk"`
	CategoryID int64
	Title      string
	Content    string
	AddTime    int64 `orm:"created"`
	UpdateTime int64 `orm:"updated"`
}

func TestColumn2Field(t *testing.T) {
	columns := []string{"id", "category_id", "title", "content", "add_time", "add_ip"}
	fields := []string{}
	result := []string{"ID", "CategoryID", "Title", "Content", "AddTime", "AddIP"}

	for _, column := range columns {
		fields = append(fields, column2Field(column))
	}

	if !reflect.DeepEqual(fields, result) {
		t.Errorf("TestColumn2Field error: %v, %v", fields, result)
	}
}

func TestField2Column(t *testing.T) {
	columns := []string{}
	fields := []string{"ID", "CategoryID", "Title", "Content", "AddTime", "AddIP"}
	result := []string{"id", "category_id", "title", "content", "add_time", "add_ip"}

	for _, field := range fields {
		columns = append(columns, field2Column(field))
	}

	if !reflect.DeepEqual(columns, result) {
		t.Errorf("TestField2Column error: %v, %v", columns, result)
	}
}

func TestNewModelInfo(t *testing.T) {
	user := new(User)
	result := &ModelInfo{
		Value:   reflect.ValueOf(user).Elem(),
		Type:    reflect.ValueOf(user).Elem().Type(),
		Slice:   false,
		ValPtr:  false,
		ValType: nil,

		ModelType: reflect.ValueOf(user).Elem().Type(),

		Fields:        []string{"ID", "Username", "Password", "RegTime", "RegIP", "UpdateTime", "UpdateIP"},
		Field2Column:  map[string]string{"ID": "id", "Username": "username", "Password": "password", "RegTime": "reg_time", "RegIP": "reg_ip", "UpdateTime": "update_time", "UpdateIP": "update_ip"},
		FieldsCreated: []string{"RegTime"},
		FieldsUpdated: []string{"UpdateTime"},

		Table: "user",
		PK:    "id",

		Columns:      []string{"id", "username", "password", "reg_time", "reg_ip", "update_time", "update_ip"},
		Column2Field: map[string]string{"id": "ID", "username": "Username", "password": "Password", "reg_time": "RegTime", "reg_ip": "RegIP", "update_time": "UpdateTime", "update_ip": "UpdateIP"},
	}

	mi := NewModelInfo(user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestField2Column error: \n%#v\n%#v", mi, result)
	}
}

func TestNewModelInfo1_Slice(t *testing.T) {
	user := &[]User{}
	result := &ModelInfo{
		Value:   reflect.ValueOf(user).Elem(),
		Type:    reflect.ValueOf(user).Elem().Type(),
		Slice:   true,
		ValPtr:  false,
		ValType: reflect.ValueOf(user).Elem().Type().Elem(),

		ModelType: reflect.ValueOf(user).Elem().Type().Elem(),

		Fields:        []string{"ID", "Username", "Password", "RegTime", "RegIP", "UpdateTime", "UpdateIP"},
		Field2Column:  map[string]string{"ID": "id", "Username": "username", "Password": "password", "RegTime": "reg_time", "RegIP": "reg_ip", "UpdateTime": "update_time", "UpdateIP": "update_ip"},
		FieldsCreated: []string{"RegTime"},
		FieldsUpdated: []string{"UpdateTime"},

		Table: "user",
		PK:    "id",

		Columns:      []string{"id", "username", "password", "reg_time", "reg_ip", "update_time", "update_ip"},
		Column2Field: map[string]string{"id": "ID", "username": "Username", "password": "Password", "reg_time": "RegTime", "reg_ip": "RegIP", "update_time": "UpdateTime", "update_ip": "UpdateIP"},
	}

	mi := NewModelInfo(user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo1_Slice error: \n%#v\n%#v", mi, result)
	}
}

func TestNewModelInfo_SliceValPtr(t *testing.T) {
	user := &[]*User{}
	result := &ModelInfo{
		Value:   reflect.ValueOf(user).Elem(),
		Type:    reflect.ValueOf(user).Elem().Type(),
		Slice:   true,
		ValPtr:  true,
		ValType: reflect.ValueOf(user).Elem().Type().Elem().Elem(),

		ModelType: reflect.ValueOf(user).Elem().Type().Elem().Elem(),

		Fields:        []string{"ID", "Username", "Password", "RegTime", "RegIP", "UpdateTime", "UpdateIP"},
		Field2Column:  map[string]string{"ID": "id", "Username": "username", "Password": "password", "RegTime": "reg_time", "RegIP": "reg_ip", "UpdateTime": "update_time", "UpdateIP": "update_ip"},
		FieldsCreated: []string{"RegTime"},
		FieldsUpdated: []string{"UpdateTime"},

		Table: "user",
		PK:    "id",

		Columns:      []string{"id", "username", "password", "reg_time", "reg_ip", "update_time", "update_ip"},
		Column2Field: map[string]string{"id": "ID", "username": "Username", "password": "Password", "reg_time": "RegTime", "reg_ip": "RegIP", "update_time": "UpdateTime", "update_ip": "UpdateIP"},
	}

	mi := NewModelInfo(user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo_SliceValPtr error: \n%#v\n%#v", mi, result)
	}
}

func TestvalueModelInfo(t *testing.T) {
	user := &User{}
	mi1, _ := DefaultModelInfoManager.ValueOf(user)
	mi2, _ := DefaultModelInfoManager.ValueOf(user)

	if mi1 != mi2 {
		t.Errorf("TestvalueModelInfo error: \n%#v\n%#v", mi1, mi2)
	}
}
