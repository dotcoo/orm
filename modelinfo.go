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

type modelInfo struct {
	Value    reflect.Value
	Type     reflect.Type
	Slice    bool
	ElemPtr  bool
	ElemType reflect.Type

	ModelType reflect.Type

	Fields        []string
	Field2Column  map[string]string
	FieldsCreated map[string]bool
	FieldsUpdated map[string]bool

	Table string
	PK    string

	Columns      []string
	Column2Field map[string]string
}

func newModelInfo(model interface{}, prefix ...string) *modelInfo {
	mi := new(modelInfo)

	mi.Value = reflect.ValueOf(model)
	if mi.Value.Kind() != reflect.Ptr {
		panic("register model must be a pointer!")
	}
	mi.Value = reflect.Indirect(mi.Value)

	mi.Type = mi.Value.Type()
	if mi.Type.Kind() != reflect.Struct && mi.Type.Kind() != reflect.Slice {
		panic("register model must be a pointer struct or pointer slice!")
	}

	mi.ModelType = mi.Type
	if mi.Value.Kind() == reflect.Slice {
		mi.Slice = true
		mi.ElemType = mi.Type.Elem()
		if mi.ElemType.Kind() == reflect.Ptr {
			mi.ElemPtr = true
			mi.ElemType = mi.ElemType.Elem()
		}
		if mi.ElemType.Kind() != reflect.Struct {
			panic("register model slice element must be a struct or pointer struct!")
		}
		mi.ModelType = mi.ElemType
	}

	mi.Fields = make([]string, 0, mi.ModelType.NumField())
	mi.Field2Column = make(map[string]string)
	mi.FieldsCreated = make(map[string]bool)
	mi.FieldsUpdated = make(map[string]bool)

	n := len(prefix)
	switch {
	case n >= 2:
		mi.Table = prefix[0] + prefix[1]
	case n == 1:
		mi.Table = prefix[0] + field2Column(mi.ModelType.Name())
	default:
		mi.Table = field2Column(mi.ModelType.Name())
	}
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
				mi.FieldsCreated[field] = true
			case "updated":
				mi.FieldsUpdated[field] = true
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

func (mi *modelInfo) GetColumn(field string) string {
	column, exist := mi.Field2Column[field]
	if !exist {
		panic("field " + field + " not found!")
	}
	return column
}

func (mi *modelInfo) GetField(column string) string {
	field, exist := mi.Column2Field[column]
	if !exist {
		panic("column " + column + " not found!")
	}
	return field
}

var modelInfos map[reflect.Type]*modelInfo = make(map[reflect.Type]*modelInfo)

var modelInfosMtx sync.RWMutex

func setModelInfo(mi *modelInfo) {
	modelInfosMtx.Lock()
	modelInfos[mi.Type] = mi
	modelInfosMtx.Unlock()
}

func getModelInfo(t reflect.Type) (*modelInfo, bool) {
	modelInfosMtx.RLock()
	mi, exist := modelInfos[t]
	modelInfosMtx.RUnlock()
	return mi, exist
}

func valueModelInfo(model interface{}) (*modelInfo, reflect.Value) {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		panic("register model must be a pointer!")
	}
	v = reflect.Indirect(v)
	t := v.Type()

	mi, exist := getModelInfo(t)
	if !exist {
		mi = newModelInfo(model)
		setModelInfo(mi)
	}
	return mi, v
}
