// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var Column2Field map[string]string = map[string]string{"id": "ID", "ip": "IP"}

func column2Field(column string) (field string) {
	cs := strings.Split(strings.ToLower(column), "_")
	for _, c := range cs {
		f, exist := Column2Field[c]
		if exist {
			field += f
		} else {
			field += strings.Title(c)
		}
	}
	return field
}

func field2Column(field string) (column string) {
	column = field
	for c, f := range Column2Field {
		column = strings.Replace(column, f, "_"+c, -1)
	}
	re := regexp.MustCompile("([A-Z])")
	column = strings.ToLower(strings.Trim(re.ReplaceAllString(column, "_$1"), "_"))
	return
}

type ModelInfo struct {
	Value reflect.Value
	Type  reflect.Type

	Map     bool
	Slice   bool
	KeyPtr  bool
	KeyType reflect.Type
	ValPtr  bool
	ValType reflect.Type

	ModelType reflect.Type

	Fields        []string
	Field2Column  map[string]string
	FieldsCreated []string
	FieldsUpdated []string

	Table string
	PK    string

	Columns      []string
	Column2Field map[string]string
}

func NewModelInfo(model interface{}, prefix, table string) *ModelInfo {
	mi := new(ModelInfo)

	mi.Value = reflect.ValueOf(model)
	if mi.Value.Kind() != reflect.Ptr {
		panic("register model must be a pointer!")
	}
	mi.Value = reflect.Indirect(mi.Value)

	mi.Type = mi.Value.Type()
	if mi.Type.Kind() != reflect.Struct && mi.Type.Kind() != reflect.Slice && mi.Type.Kind() != reflect.Map {
		panic("register model must be a pointer struct or pointer slice! or pointer map")
	}

	mi.ModelType = mi.Type

	if mi.Value.Kind() == reflect.Slice {
		mi.Slice = true

		mi.ValType = mi.Type.Elem()
		if mi.ValType.Kind() == reflect.Ptr {
			mi.ValPtr = true
			mi.ValType = mi.ValType.Elem()
		}
		if mi.ValType.Kind() != reflect.Struct {
			panic("register model slice element must be a struct or pointer struct!")
		}

		mi.ModelType = mi.ValType
	}

	if mi.Value.Kind() == reflect.Map {
		mi.Map = true

		mi.KeyType = mi.Type.Key()
		if mi.KeyType.Kind() == reflect.Ptr {
			mi.KeyPtr = true
			mi.KeyType = mi.KeyType.Elem()
		}
		if mi.KeyType.Kind() > reflect.Uint64 && mi.KeyType.Kind() != reflect.String {
			panic("register model map key must be a int or string!")
		}

		mi.ValType = mi.Type.Elem()
		if mi.ValType.Kind() == reflect.Ptr {
			mi.ValPtr = true
			mi.ValType = mi.ValType.Elem()
		}
		if mi.ValType.Kind() != reflect.Struct {
			panic("register model map element must be a struct or pointer struct!")
		}

		mi.ModelType = mi.ValType
	}

	mi.Fields = make([]string, 0, mi.ModelType.NumField())
	mi.Field2Column = make(map[string]string)
	mi.FieldsCreated = make([]string, 0, 10)
	mi.FieldsUpdated = make([]string, 0, 10)

	if table == "" {
		table = field2Column(mi.ModelType.Name())
	}

	mi.Table = prefix + table
	mi.PK = "id"

	mi.Columns = make([]string, 0, mi.ModelType.NumField())
	mi.Column2Field = make(map[string]string)

CONTINUE_FIELD:
	for i := 0; i < mi.ModelType.NumField(); i++ {
		tf := mi.ModelType.Field(i)
		field := tf.Name
		column := field2Column(field)

		ss := strings.Fields(tf.Tag.Get("orm"))
		for _, s := range ss {
			s = strings.ToLower(s)
			switch s {
			case "-":
				continue CONTINUE_FIELD
			case "pk":
				mi.PK = column
			case "created":
				mi.FieldsCreated = append(mi.FieldsCreated, field)
			case "updated":
				mi.FieldsUpdated = append(mi.FieldsUpdated, field)
			default:
				column = s
			}
		}

		mi.Fields = append(mi.Fields, field)
		mi.Field2Column[field] = column
		mi.Columns = append(mi.Columns, column)
		mi.Column2Field[column] = field
	}

	return mi
}

func (mi *ModelInfo) GetColumn(field string) string {
	if field == "" {
		panic("field cannot be null string")
	}
	if field[0] >= 'A' && field[0] <= 'Z' {
		column, exist := mi.Field2Column[field]
		if !exist {
			panic("field " + field + " not found!")
		}
		return column
	} else {
		_, exist := mi.Column2Field[field]
		if !exist {
			panic("field " + field + " not found!")
		}
		return field
	}
}

func (mi *ModelInfo) GetField(column string) string {
	if column == "" {
		panic("column cannot be null string")
	}
	if column[0] >= 'a' && column[0] <= 'z' {
		field, exist := mi.Column2Field[column]
		if !exist {
			panic("column " + column + " not found!")
		}
		return field
	} else {
		_, exist := mi.Field2Column[column]
		if !exist {
			panic("column " + column + " not found!")
		}
		return column
	}
}

type ModelInfoManager struct {
	modelInfos map[reflect.Type]*ModelInfo
	mtx        sync.RWMutex
	prefix     string
}

var DefaultModelInfoManager *ModelInfoManager = NewModelInfoManager()

func NewModelInfoManager() *ModelInfoManager {
	m := new(ModelInfoManager)
	m.modelInfos = make(map[reflect.Type]*ModelInfo)
	return m
}

func (m *ModelInfoManager) SetPrefix(prefix string) {
	m.prefix = prefix
}

func (m *ModelInfoManager) Set(mi *ModelInfo) {
	m.mtx.Lock()
	for _, i := range m.modelInfos {
		if i.ModelType == mi.ModelType {
			i.Table = mi.Table
		}
	}
	m.modelInfos[mi.Type] = mi
	m.mtx.Unlock()
}

func (m *ModelInfoManager) Get(t reflect.Type) (*ModelInfo, bool) {
	m.mtx.RLock()
	mi, exist := m.modelInfos[t]
	m.mtx.RUnlock()
	return mi, exist
}

func (m *ModelInfoManager) ValueOf(model interface{}) (*ModelInfo, reflect.Value) {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		panic("register model must be a pointer!")
	}
	v = reflect.Indirect(v)
	t := v.Type()

	mi, exist := m.Get(t)
	if !exist {
		mi = NewModelInfo(model, m.prefix, "")
		m.Set(mi)
	}
	return mi, v
}
