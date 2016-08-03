// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type dber interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type ORM struct {
	db               *sql.DB
	tx               *sql.Tx
	modelInfoManager *ModelInfoManager
	BatchRow         int
}

var DefaultORM *ORM = NewORM(nil)

func NewORM(db *sql.DB) *ORM {
	o := new(ORM)
	o.db = db
	o.tx = nil
	o.BatchRow = 100
	return o
}

func (o *ORM) SetPrefix(prefix string) {
	o.Manager().SetPrefix(prefix)
}

// query

func (o *ORM) getTxOrDB() dber {
	if o.tx != nil {
		return o.tx
	}
	if o.db != nil {
		return o.db
	}
	if DefaultORM.db != nil {
		return DefaultORM.db
	}
	panic("DB is nil!")
}

func (o *ORM) Exec(query string, args ...interface{}) sql.Result {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	result, err := o.getTxOrDB().Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (o *ORM) Query(query string, args ...interface{}) *sql.Rows {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	rows, err := o.getTxOrDB().Query(query, args...)
	if err != nil {
		panic(err)
	}
	return rows
}

func (o *ORM) QueryRow(query string, args ...interface{}) *sql.Row {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	return o.getTxOrDB().QueryRow(query, args...)
}

func (o *ORM) QueryOne(val interface{}, query string, args ...interface{}) bool {
	err := o.QueryRow(query, args...).Scan(val)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

// transaction

func (o *ORM) Begin() (otx *ORM, err error) {
	otx = NewORM(o.db)
	otx.tx, err = o.db.Begin()
	if err != nil {
		return nil, err
	}
	return otx, nil
}

func (o *ORM) Commit() error {
	tx := o.tx
	o.tx = nil
	return tx.Commit()
}

func (o *ORM) Rollback() error {
	tx := o.tx
	o.tx = nil
	return tx.Rollback()
}

// select

func (o *ORM) Manager() *ModelInfoManager {
	if o.modelInfoManager != nil {
		return o.modelInfoManager
	}

	return DefaultModelInfoManager
}

func fillModel(v reflect.Value, mi *ModelInfo, columns []string) []interface{} {
	vals := make([]interface{}, 0, len(columns))
	for _, column := range columns {
		field := mi.GetField(column)
		fp := v.FieldByName(field).Addr().Interface()
		vals = append(vals, fp)
	}
	return vals
}

func (o *ORM) Select(model interface{}, s *SQL, columns ...string) bool {
	mi, v := o.Manager().ValueOf(model)

	s.From(mi.Table).Columns(columns...)

	query, args := s.ToSelect()
	rows := o.Query(query, args...)
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	switch {
	case mi.Slice:
		for rows.Next() {
			ev := reflect.New(mi.ValType)
			vals := fillModel(ev.Elem(), mi, columns)
			err = rows.Scan(vals...)
			if err != nil {
				panic(err)
			}

			if mi.ValPtr {
				v.Set(reflect.Append(v, ev))
			} else {
				v.Set(reflect.Append(v, ev.Elem()))
			}
		}
	case mi.Map:
		for rows.Next() {
			ev := reflect.New(mi.ValType)
			vals := fillModel(ev.Elem(), mi, columns)
			err = rows.Scan(vals...)
			if err != nil {
				panic(err)
			}

			field := mi.GetField(columns[0])

			if mi.ValPtr {
				if mi.KeyPtr {
					v.SetMapIndex(ev.Elem().FieldByName(field).Addr(), ev)
				} else {
					v.SetMapIndex(ev.Elem().FieldByName(field), ev)
				}
			} else {
				if mi.KeyPtr {
					v.SetMapIndex(ev.Elem().FieldByName(field).Addr(), ev.Elem())
				} else {
					v.SetMapIndex(ev.Elem().FieldByName(field), ev.Elem())
				}
			}
		}
	default:
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

func (o *ORM) Count(s *SQL) (count int) {
	query, args := s.ToCount()
	o.QueryOne(&count, query, args...)
	return
}

func (o *ORM) CountMySQL(s *SQL) (count int) {
	query, args := s.ToCountMySQL()
	o.QueryOne(&count, query, args...)
	return
}

func columnsDefault(mi *ModelInfo, columns ...string) []string {
	switch len(columns) {
	case 0:
		columns = mi.Columns
	case 1:
		if columns[0] == "*" {
			columns = mi.Columns
		} else {
			columns = strings.Split(columns[0], ",")
			for i, column := range columns {
				columns[i] = strings.TrimSpace(column)
			}
		}
	}
	return columns
}

func setModel(s *SQL, v reflect.Value, mi *ModelInfo, skipPK bool, columns ...string) {
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
		s.Set(column, val.Interface())
	}
}

func (o *ORM) Insert(model interface{}, columns ...string) sql.Result {
	mi, v := o.Manager().ValueOf(model)

	for _, field := range mi.FieldsCreated {
		v.FieldByName(field).SetInt(time.Now().Unix())
	}

	s := o.NewSQL().From(mi.Table)
	setModel(s, v, mi, false, columns...)

	query, args := s.ToInsert()
	return o.Exec(query, args...)
}

func (o *ORM) Replace(model interface{}, columns ...string) sql.Result {
	mi, v := o.Manager().ValueOf(model)

	s := o.NewSQL().From(mi.Table)
	setModel(s, v, mi, false, columns...)

	query, args := s.ToReplace()
	return o.Exec(query, args...)
}

func (o *ORM) Update(model interface{}, s *SQL, columns ...string) sql.Result {
	if len(s.sets) == 0 && len(columns) == 0 {
		panic("Update columns cannot be empty! if update all columns, please input \"*\".")
	}

	mi, v := o.Manager().ValueOf(model)

	for _, field := range mi.FieldsUpdated {
		v.FieldByName(field).SetInt(time.Now().Unix())
	}

	s.From(mi.Table)
	setModel(s, v, mi, true, columns...)

	query, args := s.ToUpdate()
	return o.Exec(query, args...)
}

func (o *ORM) Delete(model interface{}, s *SQL) sql.Result {
	mi, _ := o.Manager().ValueOf(model)

	s.From(mi.Table)

	query, args := s.ToDelete()
	return o.Exec(query, args...)
}

func (o *ORM) batchInsertOrReplace(mode string, lineBatch int, models interface{}, columns ...string) {
	mi, vs := o.Manager().ValueOf(models)

	columns = columnsDefault(mi, columns...)

	fields := make([]string, 0, len(columns))
	for _, column := range columns {
		field := mi.GetField(column)
		fields = append(fields, field)
	}

	column := strings.Join(columns, "`,`")
	value := ",(?" + strings.Repeat(",?", len(columns)-1) + ")"

	args := make([]interface{}, 0, lineBatch)
	models_len := vs.Len()
	for i := 0; i < models_len; i++ {
		v := reflect.Indirect(vs.Index(i))
		for _, field := range fields {
			args = append(args, v.FieldByName(field).Interface())
		}
		if (i+1)%lineBatch == 0 {
			query := fmt.Sprintf("%s INTO `%s` (`%s`) VALUES %s%s", mode, mi.Table, column, value[1:], strings.Repeat(value, lineBatch-1))
			o.Exec(query, args...)
			args = args[0:0:lineBatch]
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

func whereById(o *ORM, model interface{}) *SQL {
	mi, v := o.Manager().ValueOf(model)
	return o.NewSQL().Where(fmt.Sprintf("`%s` = ?", mi.PK), v.FieldByName(mi.GetField(mi.PK)).Interface())
}

func (o *ORM) Add(model interface{}, columns ...string) sql.Result {
	return o.Insert(model, columns...)
}

func (o *ORM) Get(model interface{}, columns ...string) bool {
	return o.Select(model, whereById(o, model), columns...)
}

func (o *ORM) Up(model interface{}, columns ...string) sql.Result {
	return o.Update(model, whereById(o, model), columns...)
}

func (o *ORM) Del(model interface{}) sql.Result {
	return o.Delete(model, whereById(o, model))
}

func (o *ORM) Save(model interface{}, columns ...string) sql.Result {
	mi, v := o.Manager().ValueOf(model)
	if v.FieldByName(mi.GetField(mi.PK)).Int() > 0 {
		return o.Up(model, columns...)
	} else {
		return o.Add(model, columns...)
	}
}

// foreign key

func (o *ORM) ForeignKey(sources interface{}, fk_column string, models interface{}, pk_column string, columns ...string) {
	mi, vs := o.Manager().ValueOf(sources)

	if vs.Len() == 0 {
		return
	}

	field := mi.GetField(fk_column)
	sf, exist := mi.ValType.FieldByName(field)
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

	s := o.NewSQL().WhereIn(fmt.Sprintf("`%s` in (?)", pk_column), ids...)
	o.Select(models, s, columns...)
}

// SQL

func (o *ORM) NewSQL(table ...string) *SQL {
	s := NewSQL(table...)
	s.SetORM(o)
	return s
}
