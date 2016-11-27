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
	prefix           string
	BatchRow         int
}

func NewORM(db *sql.DB) *ORM {
	o := new(ORM)
	o.db = db
	o.tx = nil
	o.BatchRow = 100
	return o
}

func (o *ORM) SetDB(db *sql.DB) {
	o.db = db
}

func (o *ORM) SetPrefix(prefix string) {
	o.prefix = prefix
	o.Manager().SetPrefix(prefix)
}

func (o *ORM) Manager() *ModelInfoManager {
	if o.modelInfoManager != nil {
		return o.modelInfoManager
	}

	return DefaultModelInfoManager
}

func (o *ORM) NewManager() {
	o.modelInfoManager = NewModelInfoManager()
}

// query

func (o *ORM) getTxOrDB() dber {
	if o.tx != nil {
		return o.tx
	}
	if o.db != nil {
		return o.db
	}
	panic("DB is nil!")
}

func (o *ORM) RawExec(query string, args ...interface{}) (sql.Result, error) {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	result, err := o.getTxOrDB().Exec(query, args...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (o *ORM) RawQuery(query string, args ...interface{}) (*sql.Rows, error) {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	rows, err := o.getTxOrDB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (o *ORM) RawQueryRow(query string, args ...interface{}) (*sql.Row, error) {
	if strings.IndexByte(query, '\'') >= 0 {
		panic("SQL statement cannot contain single quotes!")
	}
	return o.getTxOrDB().QueryRow(query, args...), nil
}

// transaction

func (o *ORM) RawBegin() (*ORM, error) {
	var err error
	otx := NewORM(o.db)
	otx.tx, err = o.db.Begin()
	if err != nil {
		return nil, err
	}
	return otx, nil
}

func (o *ORM) RawCommit() error {
	err := o.tx.Commit()
	o.tx = nil
	return err
}

func (o *ORM) RawRollback() error {
	err := o.tx.Rollback()
	o.tx = nil
	return err
}

// select

func fillModel(v reflect.Value, mi *ModelInfo, columns []string) ([]interface{}, error) {
	vals := make([]interface{}, 0, len(columns))
	for _, column := range columns {
		vals = append(vals, v.FieldByName(mi.Field(column).Field).Addr().Interface())
	}
	return vals, nil
}

func (o *ORM) RawSelect(s *SQL, model interface{}, columns ...string) (bool, error) {
	mi, v := o.Manager().ValueOf(model)

	s.From(mi.Table).Columns(columns...)

	query, args := s.ToSelect()
	rows, err := o.RawQuery(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	columns, err = rows.Columns()
	if err != nil {
		return false, err
	}

	switch {
	case mi.Slice:
		for rows.Next() {
			ev := reflect.New(mi.ValType)
			vals, err := fillModel(ev.Elem(), mi, columns)
			if err != nil {
				return false, err
			}
			err = rows.Scan(vals...)
			if err != nil {
				return false, err
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
			vals, err := fillModel(ev.Elem(), mi, columns)
			if err != nil {
				return false, err
			}
			err = rows.Scan(vals...)
			if err != nil {
				return false, err
			}

			field := mi.Field(columns[0]).Field

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
			return false, nil
		}
		vals, err := fillModel(v, mi, columns)
		if err != nil {
			return false, err
		}
		err = rows.Scan(vals...)
		if err != nil {
			return false, err
		}
	}

	err = rows.Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (o *ORM) RawSelectVal(s *SQL, vals ...interface{}) (bool, error) {
	query, args := s.ToSelect()
	row, err := o.RawQueryRow(query, args...)
	if err != nil {
		return false, err
	}
	err = row.Scan(vals...)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *ORM) RawCount(s *SQL) (count int, err error) {
	_, err = o.RawSelectVal(s.NewCount(), &count)
	return count, err
}

func columnsDefault(mi *ModelInfo, columns ...string) []string {
	if len(columns) == 0 || columns[0] == "*" {
		columns = mi.ColumnNames
	} else {
		for _, column := range strings.Split(columns[0], ",") {
			columns = append(columns, strings.TrimSpace(column))
		}
		columns = columns[1:]
	}
	return columns
}

func valInt(v reflect.Value) (int64, uint64) {
	kind := v.Kind()
	if kind >= reflect.Int && kind <= reflect.Int64 {
		return v.Int(), 0
	}
	if kind >= reflect.Uint && kind <= reflect.Uint64 {
		return 0, v.Uint()
	}
	return 0, 0
}

func valSetInt(v reflect.Value, i64 int64, u64 uint64) {
	kind := v.Kind()
	if kind >= reflect.Int && kind <= reflect.Int64 {
		v.SetInt(i64)
	}
	if kind >= reflect.Uint && kind <= reflect.Uint64 {
		v.SetUint(u64)
	}
}

func setModel(s *SQL, v reflect.Value, mi *ModelInfo, skipPK bool, columns ...string) {
	columns = columnsDefault(mi, columns...)
	for _, column := range columns {
		if skipPK && column == mi.PK.Column {
			continue
		}
		mf := mi.Field(column)
		if mf.Field == mi.PK.Field {
			i64, u64 := valInt(v.FieldByName(mf.Field))
			if i64 <= 0 && u64 <= 0 {
				continue
			}
		}
		s.Set(column, v.FieldByName(mf.Field).Interface())
	}
}

func (o *ORM) RawInsert(model interface{}, columns ...string) (sql.Result, error) {
	mi, v := o.Manager().ValueOf(model)

	u := time.Now().Unix()
	for _, field := range mi.FieldsCreated {
		valSetInt(v.FieldByName(field), u, uint64(u))
		if len(columns) > 0 && columns[0] != "*" {
			columns = append(columns, mi.Column(field).Column)
		}
	}

	s := o.NewSQL().From(mi.Table)
	setModel(s, v, mi, false, columns...)

	query, args := s.ToInsert()
	result, err := o.RawExec(query, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	valSetInt(v.FieldByName(mi.PK.Field), id, uint64(id))

	return result, err
}

func (o *ORM) RawReplace(model interface{}, columns ...string) (sql.Result, error) {
	mi, v := o.Manager().ValueOf(model)

	s := o.NewSQL().From(mi.Table)
	setModel(s, v, mi, false, columns...)

	query, args := s.ToReplace()
	return o.RawExec(query, args...)
}

func (o *ORM) RawUpdate(s *SQL, model interface{}, columns ...string) (sql.Result, error) {
	if len(s.sets) == 0 && len(columns) == 0 {
		panic("Update columns cannot be empty! if update all columns, please input \"*\".")
	}

	mi, v := o.Manager().ValueOf(model)

	u := time.Now().Unix()
	for _, field := range mi.FieldsUpdated {
		valSetInt(v.FieldByName(field), u, uint64(u))
		if len(columns) > 0 && columns[0] != "*" {
			columns = append(columns, mi.Column(field).Column)
		}
	}

	s.From(mi.Table)
	setModel(s, v, mi, true, columns...)

	query, args := s.ToUpdate()
	return o.RawExec(query, args...)
}

func (o *ORM) RawDelete(s *SQL, model interface{}) (sql.Result, error) {
	mi, _ := o.Manager().ValueOf(model)

	s.From(mi.Table)

	query, args := s.ToDelete()
	return o.RawExec(query, args...)
}

func (o *ORM) batchInsertOrReplace(mode string, lineBatch int, models interface{}, columns ...string) error {
	mi, vs := o.Manager().ValueOf(models)

	columns = columnsDefault(mi, columns...)

	fields := make([]string, 0, len(columns))
	for _, column := range columns {
		fields = append(fields, mi.Field(column).Field)
	}

	column := strings.Join(columns, "`,`")
	value := ",(" + strings.Repeat(",?", len(columns))[1:] + ")"

	args := make([]interface{}, 0, lineBatch)
	models_len := vs.Len()
	for i := 0; i < models_len; i++ {
		v := reflect.Indirect(vs.Index(i))
		for _, field := range fields {
			args = append(args, v.FieldByName(field).Interface())
		}
		if (i+1)%lineBatch == 0 {
			query := fmt.Sprintf("%s INTO `%s` (`%s`) VALUES %s", mode, mi.Table, column, strings.Repeat(value, lineBatch)[1:])
			_, err := o.RawExec(query, args...)
			if err != nil {
				return err
			}
			args = args[0:0:lineBatch]
		}
	}
	if models_len%lineBatch > 0 {
		query := fmt.Sprintf("%s INTO `%s` (`%s`) VALUES %s", mode, mi.Table, column, strings.Repeat(value, models_len%lineBatch)[1:])
		_, err := o.RawExec(query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ORM) RawBatchInsert(models interface{}, columns ...string) error {
	return o.batchInsertOrReplace("INSERT", o.BatchRow, models, columns...)
}

func (o *ORM) RawBatchReplace(models interface{}, columns ...string) error {
	return o.batchInsertOrReplace("REPLACE", o.BatchRow, models, columns...)
}

// quick method

func whereById(s *SQL, o *ORM, model interface{}) *SQL {
	mi, v := o.Manager().ValueOf(model)
	return s.Where(fmt.Sprintf("`%s` = ?", mi.PK.Column), v.FieldByName(mi.PK.Field).Interface())
}

func (o *ORM) RawAdd(model interface{}, columns ...string) (sql.Result, error) {
	return o.RawInsert(model, columns...)
}

func (o *ORM) RawGet(model interface{}, columns ...string) (bool, error) {
	return o.RawSelect(whereById(o.NewSQL(), o, model), model, columns...)
}

func stringsIndex(ss []string, val string) int {
	for i, v := range ss {
		if v == val {
			return i
		}
	}
	return -1
}

func (o *ORM) RawGetBy(model interface{}, cols_nil_columns ...string) (bool, error) {
	if len(cols_nil_columns) == 0 {
		panic("cols_nil_columns not can null!")
	}
	var cols, columns []string
	if i := stringsIndex(cols_nil_columns, ""); i == -1 {
		cols, columns = cols_nil_columns, columns
	} else {
		cols, columns = cols_nil_columns[:i], cols_nil_columns[i+1:]
	}
	mi, v := o.Manager().ValueOf(model)
	sq := o.NewSQL()
	for _, cols := range cols {
		sq.Where(fmt.Sprintf("`%s` = ?", cols), v.FieldByName(mi.Field(cols).Field).Interface())
	}
	return o.RawSelect(sq, model, columns...)
}

func (o *ORM) RawUp(model interface{}, columns ...string) (sql.Result, error) {
	return o.RawUpdate(whereById(o.NewSQL(), o, model), model, columns...)
}

func (o *ORM) RawDel(model interface{}) (sql.Result, error) {
	return o.RawDelete(whereById(o.NewSQL(), o, model), model)
}

func (o *ORM) RawSave(model interface{}, columns ...string) (sql.Result, error) {
	mi, v := o.Manager().ValueOf(model)
	i64, u64 := valInt(v.FieldByName(mi.PK.Field))
	if i64 > 0 || u64 > 0 {
		return o.RawUp(model, columns...)
	} else {
		return o.RawAdd(model, columns...)
	}
}

// foreign key

func (o *ORM) RawForeignKey(sources interface{}, fk_column string, models interface{}, pk_column string, columns ...string) error {
	mi, vs := o.Manager().ValueOf(sources)

	if vs.Len() == 0 {
		return nil
	}

	fk := mi.Field(fk_column)
	if fk.Kind < reflect.Int && fk.Kind > reflect.Uint64 {
		panic("field " + fk.Field + " not int type!")
	}

	isInt64 := fk.Kind >= reflect.Int && fk.Kind <= reflect.Int64

	ids_map_int64 := make(map[int64]bool)
	ids_map_uint64 := make(map[uint64]bool)
	models_len := vs.Len()
	for i := 0; i < models_len; i++ {
		i64, u64 := valInt(reflect.Indirect(vs.Index(i)).FieldByName(fk.Field))
		if isInt64 {
			ids_map_int64[i64] = true
		} else {
			ids_map_uint64[u64] = true
		}
	}

	ids := make([]interface{}, 0, 20)
	for id, _ := range ids_map_int64 {
		ids = append(ids, id)
	}
	for id, _ := range ids_map_uint64 {
		ids = append(ids, id)
	}

	s := o.NewSQL().WhereIn(fmt.Sprintf("`%s` in (?)", pk_column), ids...)
	_, err := o.RawSelect(s, models, columns...)
	return err
}

// SQL

func (o *ORM) NewSQL() *SQL {
	s := new(SQL)
	s.orm = o
	return s
}

func (o *ORM) sqlFrom(s *SQL, table string) string {
	if !strings.HasPrefix(table, o.prefix) {
		table = o.prefix + table
	}
	return table
}

func (o *ORM) sqlJoin(s *SQL, table, cond string) (string, string) {
	if !strings.HasPrefix(table, o.prefix) {
		table = o.prefix + table
	}
	return table, cond
}
