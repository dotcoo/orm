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
	return strings.ToLower(strings.Trim(re.ReplaceAllString(column, "_$1"), "_"))
}

func fieldsFunc(d rune) func(rune) bool {
	return func(r rune) bool {
		return r == d
	}
}

var commaFieldsFunc = fieldsFunc(',')

type ModelField struct {
	Field   string
	Column  string
	PK      bool
	Kind    reflect.Kind
	Created bool
	Updated bool
}

type ModelInfo struct {
	Value reflect.Value
	Type  reflect.Type

	Map       bool
	Slice     bool
	KeyPtr    bool
	KeyType   reflect.Type
	ValPtr    bool
	ValType   reflect.Type
	ModelType reflect.Type

	Table         string
	PK            *ModelField
	Columns       []*ModelField
	Fields        []*ModelField
	Column2Field  map[string]*ModelField
	Field2Column  map[string]*ModelField
	ColumnNames   []string
	FieldNames    []string
	FieldsCreated []string
	FieldsUpdated []string
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
		panic("register model must be a pointer struct or pointer slice or pointer map")
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

	if table == "" {
		table = field2Column(mi.ModelType.Name())
	}

	mi.Table = prefix + table
	mi.PK = nil
	mi.Columns = make([]*ModelField, 0, mi.ModelType.NumField())
	mi.Fields = make([]*ModelField, 0, mi.ModelType.NumField())
	mi.Column2Field = make(map[string]*ModelField)
	mi.Field2Column = make(map[string]*ModelField)
	mi.ColumnNames = make([]string, 0, mi.ModelType.NumField())
	mi.FieldNames = make([]string, 0, mi.ModelType.NumField())
	mi.FieldsCreated = make([]string, 0, mi.ModelType.NumField())
	mi.FieldsUpdated = make([]string, 0, mi.ModelType.NumField())

CONTINUE_FIELD:
	for i := 0; i < mi.ModelType.NumField(); i++ {
		tf := mi.ModelType.Field(i)

		mf := new(ModelField)
		mf.Field = tf.Name
		mf.Column = field2Column(tf.Name)
		mf.Kind = tf.Type.Kind()

		ss := strings.FieldsFunc(tf.Tag.Get("orm"), commaFieldsFunc)
		for _, s := range ss {
			s = strings.ToLower(s)
			switch s {
			case "-":
				continue CONTINUE_FIELD
			case "pk":
				mf.PK = true
				mi.PK = mf
			case "unique":
			case "index":
			case "fk":
			case "created":
				mf.Created = true
				mi.FieldsCreated = append(mi.FieldsCreated, mf.Field)
			case "updated":
				mf.Updated = true
				mi.FieldsUpdated = append(mi.FieldsUpdated, mf.Field)
			default:
				mf.Column = s
			}
		}

		mi.Columns = append(mi.Columns, mf)
		mi.Fields = append(mi.Fields, mf)
		mi.Column2Field[mf.Column] = mf
		mi.Field2Column[mf.Field] = mf
		mi.ColumnNames = append(mi.ColumnNames, mf.Column)
		mi.FieldNames = append(mi.FieldNames, mf.Field)
	}

	return mi
}

func (mi *ModelInfo) Column(field string) *ModelField {
	mf, exist := mi.Field2Column[field]
	if exist {
		return mf
	}
	mf, exist = mi.Column2Field[field]
	if exist {
		return mf
	}
	panic("column " + field + " not found!")
}

func (mi *ModelInfo) Field(column string) *ModelField {
	mf, exist := mi.Column2Field[column]
	if exist {
		return mf
	}
	mf, exist = mi.Field2Column[column]
	if exist {
		return mf
	}
	panic("field " + column + " not found!")
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

func (m *ModelInfoManager) Get(t reflect.Type) *ModelInfo {
	m.mtx.RLock()
	mi := m.modelInfos[t]
	m.mtx.RUnlock()
	return mi
}

func (m *ModelInfoManager) ValueOf(model interface{}) (*ModelInfo, reflect.Value) {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		panic("register model must be a pointer!")
	}
	v = reflect.Indirect(v)
	t := v.Type()

	mi := m.Get(t)
	if mi == nil {
		mi = NewModelInfo(model, m.prefix, "")
		m.Set(mi)
	}
	return mi, v
}
