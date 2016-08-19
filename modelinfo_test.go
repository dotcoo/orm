// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"reflect"
	"testing"
)

type User struct {
	ID         int `orm:"pk"`
	Username   string
	Password   string
	RegTime    int `orm:"created"`
	RegIP      uint32
	UpdateTime int `orm:"updated"`
	UpdateIP   uint32
	OtherField string `orm:"_"`
}

type Category struct {
	ID   uint64 `orm:"pk"`
	Name string
}

type Blog struct {
	ID         uint64 `orm:"pk"`
	CategoryID uint64
	Title      string
	Content    string
	AddTime    int `orm:"created"`
	UpdateTime int `orm:"updated"`
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

var result *ModelInfo

func init() {
	user := new(User)

	id := &ModelField{Field: "ID", Column: "id", PK: true, Kind: reflect.Int, Created: false, Updated: false}
	username := &ModelField{Field: "Username", Column: "username", PK: false, Kind: reflect.String, Created: false, Updated: false}
	password := &ModelField{Field: "Password", Column: "password", PK: false, Kind: reflect.String, Created: false, Updated: false}
	reg_time := &ModelField{Field: "RegTime", Column: "reg_time", PK: false, Kind: reflect.Int, Created: true, Updated: false}
	reg_ip := &ModelField{Field: "RegIP", Column: "reg_ip", PK: false, Kind: reflect.Uint32, Created: false, Updated: false}
	update_time := &ModelField{Field: "UpdateTime", Column: "update_time", PK: false, Kind: reflect.Int, Created: false, Updated: true}
	update_ip := &ModelField{Field: "UpdateIP", Column: "update_ip", PK: false, Kind: reflect.Uint32, Created: false, Updated: false}
	result = &ModelInfo{
		Value: reflect.ValueOf(user).Elem(),
		Type:  reflect.ValueOf(user).Elem().Type(),

		Map:       false,
		Slice:     false,
		KeyPtr:    false,
		KeyType:   nil,
		ValPtr:    false,
		ValType:   nil,
		ModelType: reflect.ValueOf(user).Elem().Type(),

		Table:         "user",
		PK:            id,
		Columns:       []*ModelField{id, username, password, reg_time, reg_ip, update_time, update_ip},
		Fields:        []*ModelField{id, username, password, reg_time, reg_ip, update_time, update_ip},
		Column2Field:  map[string]*ModelField{"id": id, "username": username, "password": password, "reg_time": reg_time, "reg_ip": reg_ip, "update_time": update_time, "update_ip": update_ip},
		Field2Column:  map[string]*ModelField{"ID": id, "Username": username, "Password": password, "RegTime": reg_time, "RegIP": reg_ip, "UpdateTime": update_time, "UpdateIP": update_ip},
		ColumnNames:   []string{"id", "username", "password", "reg_time", "reg_ip", "update_time", "update_ip"},
		FieldNames:    []string{"ID", "Username", "Password", "RegTime", "RegIP", "UpdateTime", "UpdateIP"},
		FieldsCreated: []string{"RegTime"},
		FieldsUpdated: []string{"UpdateTime"},
	}
}

func TestNewModelInfo(t *testing.T) {
	user := new(User)

	result.Value = reflect.ValueOf(user).Elem()
	result.Type = reflect.ValueOf(user).Elem().Type()
	result.Map = false
	result.Slice = false
	result.KeyPtr = false
	result.KeyType = nil
	result.ValPtr = false
	result.ValType = nil
	result.ModelType = reflect.ValueOf(user).Elem().Type()

	mi := NewModelInfo(user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo error: \n%#v\n%#v", mi, result)
	}
}

func TestNewModelInfo1_Slice(t *testing.T) {
	user := make([]User, 0)

	result.Value = reflect.ValueOf(&user).Elem()
	result.Type = reflect.ValueOf(&user).Elem().Type()
	result.Map = false
	result.Slice = true
	result.KeyPtr = false
	result.KeyType = nil
	result.ValPtr = false
	result.ValType = reflect.ValueOf(&user).Elem().Type().Elem()
	result.ModelType = reflect.ValueOf(&user).Elem().Type().Elem()

	mi := NewModelInfo(&user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo1_Slice error: \n%#v\n%#v", mi, result)
	}
}

func TestNewModelInfo_SliceValPtr(t *testing.T) {
	user := make([]*User, 0)

	result.Value = reflect.ValueOf(&user).Elem()
	result.Type = reflect.ValueOf(&user).Elem().Type()
	result.Map = false
	result.Slice = true
	result.KeyPtr = false
	result.KeyType = nil
	result.ValPtr = true
	result.ValType = reflect.ValueOf(&user).Elem().Type().Elem().Elem()
	result.ModelType = reflect.ValueOf(&user).Elem().Type().Elem().Elem()

	mi := NewModelInfo(&user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo_SliceValPtr error: \n%#v\n%#v", mi, result)
	}
}

func TestNewModelInfo_Map(t *testing.T) {
	user := make(map[int]User)

	result.Value = reflect.ValueOf(&user).Elem()
	result.Type = reflect.ValueOf(&user).Elem().Type()
	result.Map = true
	result.Slice = false
	result.KeyPtr = false
	result.KeyType = reflect.ValueOf(&user).Elem().Type().Key()
	result.ValPtr = false
	result.ValType = reflect.ValueOf(&user).Elem().Type().Elem()
	result.ModelType = reflect.ValueOf(&user).Elem().Type().Elem()

	mi := NewModelInfo(&user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo_Map error: \n%#v\n%#v", mi, result)
	}
}

func TestNewModelInfo_MapKeyPtr(t *testing.T) {
	user := make(map[*int]User)

	result.Value = reflect.ValueOf(&user).Elem()
	result.Type = reflect.ValueOf(&user).Elem().Type()
	result.Map = true
	result.Slice = false
	result.KeyPtr = true
	result.KeyType = reflect.ValueOf(&user).Elem().Type().Key().Elem()
	result.ValPtr = false
	result.ValType = reflect.ValueOf(&user).Elem().Type().Elem()
	result.ModelType = reflect.ValueOf(&user).Elem().Type().Elem()

	mi := NewModelInfo(&user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo_MapKeyPtr error: \n%#v\n%#v", mi, result)
	}
}

func TestNewModelInfo_MapValPtr(t *testing.T) {
	user := make(map[int]*User)

	result.Value = reflect.ValueOf(&user).Elem()
	result.Type = reflect.ValueOf(&user).Elem().Type()
	result.Map = true
	result.Slice = false
	result.KeyPtr = false
	result.KeyType = reflect.ValueOf(&user).Elem().Type().Key()
	result.ValPtr = true
	result.ValType = reflect.ValueOf(&user).Elem().Type().Elem().Elem()
	result.ModelType = reflect.ValueOf(&user).Elem().Type().Elem().Elem()

	mi := NewModelInfo(&user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo_MapValPtr error: \n%#v\n%#v", mi, result)
	}
}

func TestNewModelInfo_MapKeyPtrValPtr(t *testing.T) {
	user := make(map[*int]*User)

	result.Value = reflect.ValueOf(&user).Elem()
	result.Type = reflect.ValueOf(&user).Elem().Type()
	result.Map = true
	result.Slice = false
	result.KeyPtr = true
	result.KeyType = reflect.ValueOf(&user).Elem().Type().Key().Elem()
	result.ValPtr = true
	result.ValType = reflect.ValueOf(&user).Elem().Type().Elem().Elem()
	result.ModelType = reflect.ValueOf(&user).Elem().Type().Elem().Elem()

	mi := NewModelInfo(&user, "", "")
	if !reflect.DeepEqual(mi, result) {
		t.Errorf("TestNewModelInfo_MapKeyPtrValPtr error: \n%#v\n%#v", mi, result)
	}
}

func TestvalueModelInfo(t *testing.T) {
	user1 := new(User)
	user2 := new(User)
	mi1, _ := DefaultModelInfoManager.ValueOf(user1)
	mi2, _ := DefaultModelInfoManager.ValueOf(user2)

	if mi1 != mi2 {
		t.Errorf("TestvalueModelInfo error: \n%#v\n%#v", mi1, mi2)
	}
}
