// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	stdsql "database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type DB interface {
	Exec(query string, args ...interface{}) (stdsql.Result, error)
	Prepare(query string) (*stdsql.Stmt, error)
	Query(query string, args ...interface{}) (*stdsql.Rows, error)
	QueryRow(query string, args ...interface{}) *stdsql.Row
}

type ORM struct {
	DB       *stdsql.DB
	Tx       *stdsql.Tx
	BatchRow int
}

func NewORM(db *stdsql.DB) *ORM {
	o := new(ORM)
	o.DB = db
	o.Tx = nil
	o.BatchRow = 100
	return o
}

// query

func (o *ORM) getTxOrDB() DB {
	if o.Tx != nil {
		return o.Tx
	}
	if o.DB != nil {
		return o.DB
	}
	panic("ORM.DB is nil!")
}

func (o *ORM) Exec(query string, args ...interface{}) stdsql.Result {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	result, err := o.getTxOrDB().Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Query(query string, args ...interface{}) *stdsql.Rows {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	rows, err := o.getTxOrDB().Query(query, args...)
	if err != nil {
		panic(err)
	}
	return rows
}

func (o *ORM) QueryRow(query string, args ...interface{}) *stdsql.Row {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	return o.getTxOrDB().QueryRow(query, args...)
}

func (o *ORM) QueryOne(val interface{}, query string, args ...interface{}) bool {
	err := o.QueryRow(query, args...).Scan(val)
	if err == stdsql.ErrNoRows {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

// transaction

func (o *ORM) Begin() (otx *ORM, err error) {
	otx = NewORM(o.DB)
	otx.Tx, err = o.DB.Begin()
	if err != nil {
		return nil, err
	}
	return otx, nil
}

func (o *ORM) Commit() error {
	tx := o.Tx
	o.Tx = nil
	return tx.Commit()
}

func (o *ORM) Rollback() error {
	tx := o.Tx
	o.Tx = nil
	return tx.Rollback()
}

// select

func fillModel(v reflect.Value, mi *ModelInfo, columns []string) []interface{} {
	vals := make([]interface{}, 0, len(columns))
	for _, column := range columns {
		field := mi.GetField(column)
		fp := v.FieldByName(field).Addr().Interface()
		vals = append(vals, fp)
	}
	return vals
}

func (o *ORM) Select(model interface{}, sql *SQL, columns ...string) bool {
	v, mi := ValueModelInfo(model)

	sql.From(mi.Table).Columns(columns...)

	query, args := sql.ToSelect()
	rows := o.Query(query, args...)
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	if mi.Slice {
		for rows.Next() {
			ev := reflect.New(mi.ElemType)
			vals := fillModel(ev.Elem(), mi, columns)
			err = rows.Scan(vals...)
			if err != nil {
				panic(err)
			}

			if mi.ElemPtr {
				v.Set(reflect.Append(v, ev))
			} else {
				v.Set(reflect.Append(v, ev.Elem()))
			}
		}
	} else {
		if !rows.Next() {
			return false
		}
		vals := fillModel(v, mi, columns)
		err = rows.Scan(vals...)
		if err != nil {
			panic(err)
		}
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return true
}

func (o *ORM) Count(sql *SQL) (count int) {
	query, args := sql.ToCount()
	o.QueryOne(&count, query, args...)
	return
}

func (o *ORM) CountMySQL(sql *SQL) (count int) {
	query, args := sql.ToCountMySQL()
	o.QueryOne(&count, query, args...)
	return
}

func columnsDefault(mi *ModelInfo, columns ...string) []string {
	switch len(columns) {
	case 0:
		columns = mi.Columns
	case 1:
		columns = strings.Split(columns[0], ",")
		for i, column := range columns {
			columns[i] = strings.TrimSpace(column)
		}
	}
	return columns
}

func setModel(sql *SQL, v reflect.Value, mi *ModelInfo, skipPK bool, columns ...string) {
	columns = columnsDefault(mi, columns...)
	for _, column := range columns {
		if skipPK && column == mi.PK {
			continue
		}
		field := mi.GetField(column)
		val := v.FieldByName(field)
		if field == mi.PK && val.Int() <= 0 {
			continue
		}
		sql.Set(column, val.Interface())
	}
}

func (o *ORM) Insert(model interface{}, columns ...string) stdsql.Result {
	v, mi := ValueModelInfo(model)

	for field, _ := range mi.FieldsCreated {
		v.FieldByName(field).SetInt(time.Now().Unix())
	}

	sql := o.NewSQL().From(mi.Table)
	setModel(sql, v, mi, false, columns...)

	query, args := sql.ToInsert()
	return o.Exec(query, args...)
}

func (o *ORM) Replace(model interface{}, columns ...string) stdsql.Result {
	v, mi := ValueModelInfo(model)

	sql := o.NewSQL().From(mi.Table)
	setModel(sql, v, mi, false, columns...)

	query, args := sql.ToReplace()
	return o.Exec(query, args...)
}

func (o *ORM) Update(model interface{}, sql *SQL, columns ...string) stdsql.Result {
	v, mi := ValueModelInfo(model)

	for field, _ := range mi.FieldsUpdated {
		v.FieldByName(field).SetInt(time.Now().Unix())
	}

	sql.From(mi.Table)
	setModel(sql, v, mi, true, columns...)

	query, args := sql.ToUpdate()
	return o.Exec(query, args...)
}

func (o *ORM) Delete(model interface{}, sql *SQL) stdsql.Result {
	_, mi := ValueModelInfo(model)

	sql.From(mi.Table)

	query, args := sql.ToDelete()
	return o.Exec(query, args...)
}

func (o *ORM) batchInsertOrReplace(mode string, lineBatch int, models interface{}, columns ...string) {
	vs, mi := ValueModelInfo(models)

	columns = columnsDefault(mi, columns...)

	fields := make([]string, 0, len(columns))
	for _, column := range columns {
		field := mi.GetField(column)
		fields = append(fields, field)
	}

	column := strings.Join(columns, "`,`")
	value := ",(?" + strings.Repeat(",?", len(columns)-1) + ")"

	args := make([]interface{}, 0, 100)
	models_len := vs.Len()
	for i := 0; i < models_len; i++ {
		v := reflect.Indirect(vs.Index(i))
		for _, field := range fields {
			args = append(args, v.FieldByName(field).Interface())
		}
		if (i+1)%lineBatch == 0 {
			query := fmt.Sprintf("%s INTO `%s` (`%s`) VALUES %s%s", mode, mi.Table, column, value[1:], strings.Repeat(value, lineBatch-1))
			o.Exec(query, args...)
			args = args[0:0:100]
		}
	}
	if models_len%lineBatch > 0 {
		query := fmt.Sprintf("%s INTO `%s` (`%s`) VALUES %s%s", mode, mi.Table, column, value[1:], strings.Repeat(value, models_len%lineBatch-1))
		o.Exec(query, args...)
	}
}

func (o *ORM) BatchInsert(models interface{}, columns ...string) {
	o.batchInsertOrReplace("INSERT", o.BatchRow, models, columns...)
}

func (o *ORM) BatchReplace(models interface{}, columns ...string) {
	o.batchInsertOrReplace("REPLACE", o.BatchRow, models, columns...)
}

// quick method

func whereById(model interface{}) *SQL {
	v, mi := ValueModelInfo(model)
	return o.NewSQL().Where(fmt.Sprintf("`%s` = ?", mi.PK), v.FieldByName(mi.GetField(mi.PK)).Interface())
}

func (o *ORM) Add(model interface{}, columns ...string) stdsql.Result {
	return o.Insert(model, columns...)
}

func (o *ORM) Get(model interface{}, columns ...string) bool {
	return o.Select(model, whereById(model), columns...)
}

func (o *ORM) Up(model interface{}, columns ...string) stdsql.Result {
	return o.Update(model, whereById(model), columns...)
}

func (o *ORM) Del(model interface{}) stdsql.Result {
	return o.Delete(model, whereById(model))
}

func (o *ORM) Save(model interface{}, columns ...string) stdsql.Result {
	v, mi := ValueModelInfo(model)
	if v.FieldByName(mi.GetField(mi.PK)).Int() > 0 {
		return o.Up(model, columns...)
	} else {
		return o.Add(model, columns...)
	}
}

// foreign key

func (o *ORM) ForeignKey(models interface{}, foreign_key_column string, foreign_models interface{}, key_column string, columns ...string) {
	vs, mi := ValueModelInfo(models)

	if vs.Len() == 0 {
		return
	}

	field := mi.GetField(foreign_key_column)
	sf, exist := mi.ElemType.FieldByName(field)
	if !exist {
		panic("field " + field + " not found!")
	}
	kind := sf.Type.Kind()
	if kind != reflect.Int && kind != reflect.Int32 && kind != reflect.Int64 && kind != reflect.Uint && kind != reflect.Uint32 && kind != reflect.Uint64 {
		panic("field " + field + " not int type!")
	}

	ids_map := make(map[int64]bool)
	models_len := vs.Len()
	for i := 0; i < models_len; i++ {
		v := reflect.Indirect(vs.Index(i))
		ids_map[v.FieldByName(field).Int()] = true
	}

	ids := make([]interface{}, 0, 20)
	for id, _ := range ids_map {
		ids = append(ids, id)
	}

	sql := o.NewSQL().WhereIn(fmt.Sprintf("`%s` in (?)", key_column), ids...)
	o.Select(foreign_models, sql, columns...)
}

// SQL

func (o *ORM) NewSQL(table ...string) *SQL {
	sql := NewSQL(table...)
	sql.SetORM(o)
	return sql
}
