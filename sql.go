// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
	"fmt"
	"strings"
)

type SQL struct {
	keywords          []string      // keywords
	columns           []string      // columns
	table             string        // table
	sets              []string      // sets
	setsParams        []interface{} // sets params
	joins             []string      // joins
	wheres            []string      // where
	wheresParams      []interface{} // where params
	countWheres       []string      // count where
	countWheresParams []interface{} // count where params
	groups            []string      // group
	havings           []string      // having
	havingsParams     []interface{} // having params
	orders            []string      // order
	limit             int           // limit
	offset            int           // offset
	forUpdate         string        // read lock
	lockInShareMode   string        // write lock
	orm               *ORM          // ORM
}

func NewSQL(table ...string) *SQL {
	sql := new(SQL)
	sql.keywords = make([]string, 0, 0)
	sql.columns = make([]string, 0, 0)
	sql.table = ""
	if len(table) > 0 {
		sql.table = table[0]
	}
	sql.sets = make([]string, 0, 20)
	sql.setsParams = make([]interface{}, 0, 20)
	sql.joins = make([]string, 0, 0)
	sql.wheres = make([]string, 0, 5)
	sql.wheresParams = make([]interface{}, 0, 20)
	sql.groups = make([]string, 0)
	sql.havings = make([]string, 0)
	sql.havingsParams = make([]interface{}, 0)
	sql.orders = make([]string, 0, 1)
	sql.limit = -1
	sql.offset = -1
	sql.forUpdate = ""
	sql.lockInShareMode = ""
	sql.orm = DefaultORM
	return sql
}

func (s *SQL) Reset() *SQL {
	s.keywords = s.keywords[0:0]
	s.columns = s.columns[0:0]
	s.sets = s.sets[0:0]
	s.setsParams = s.setsParams[0:0]
	s.joins = s.joins[0:0]
	s.wheres = s.wheres[0:0]
	s.wheresParams = s.wheresParams[0:0]
	s.groups = s.groups[0:0]
	s.havings = s.havings[0:0]
	s.havingsParams = s.havingsParams[0:0]
	s.orders = s.orders[0:0]
	s.limit = -1
	s.offset = -1
	s.forUpdate = ""
	s.lockInShareMode = ""
	return s
}

// sql syntax

func (s *SQL) Keywords(keywords ...string) *SQL {
	s.keywords = append(s.keywords, keywords...)
	return s
}

func (s *SQL) CalcFoundRows() *SQL {
	return s.Keywords("SQL_CALC_FOUND_ROWS")
}

func (s *SQL) Columns(columns ...string) *SQL {
	s.columns = append(s.columns, columns...)
	return s
}

func (s *SQL) Table(table string, alias ...string) *SQL {
	return s.From(table, alias...)
}

func (s *SQL) From(table string, alias ...string) *SQL {
	if len(alias) == 0 {
		s.table = fmt.Sprintf("`%s`", table)
	} else {
		s.table = fmt.Sprintf("`%s` AS `%s`", table, alias[0])
	}
	return s
}

func (s *SQL) Set(col string, val interface{}) *SQL {
	s.sets = append(s.sets, fmt.Sprintf("`%s` = ?", col))
	s.setsParams = append(s.setsParams, val)
	return s
}

func (s *SQL) Join(table, alias, cond string) *SQL {
	if alias == "" {
		s.joins = append(s.joins, fmt.Sprintf("`%s` ON %s", table, cond))
	} else {
		s.joins = append(s.joins, fmt.Sprintf("`%s` AS `%s` ON %s", table, alias, cond))
	}
	return s
}

func (s *SQL) Where(where string, params ...interface{}) *SQL {
	s.wheres = append(s.wheres, where)
	s.wheresParams = append(s.wheresParams, params...)
	return s
}

func (s *SQL) WhereIn(where string, params ...interface{}) *SQL {
	where = strings.Replace(where, "?", "?"+strings.Repeat(", ?", len(params)-1), 1)
	s.wheres = append(s.wheres, where)
	s.wheresParams = append(s.wheresParams, params...)
	return s
}

func (s *SQL) Group(groups ...string) *SQL {
	s.groups = append(s.groups, groups...)
	return s
}

func (s *SQL) Having(having string, params ...interface{}) *SQL {
	s.havings = append(s.havings, having)
	s.havingsParams = append(s.havingsParams, params...)
	return s
}

func (s *SQL) Order(orders ...string) *SQL {
	s.orders = append(s.orders, orders...)
	return s
}

func (s *SQL) Limit(limit int) *SQL {
	s.limit = limit
	return s
}

func (s *SQL) Offset(offset int) *SQL {
	s.offset = offset
	return s
}

func (s *SQL) ForUpdate() *SQL {
	s.forUpdate = " FOR UPDATE"
	return s
}

func (s *SQL) LockInShareMode() *SQL {
	s.lockInShareMode = " LOCK IN SHARE MODE"
	return s
}

// sql tool

func (s *SQL) SetMap(data map[string]interface{}) *SQL {
	for col, val := range data {
		s.sets = append(s.sets, fmt.Sprintf("`%s` = ?", col))
		s.setsParams = append(s.setsParams, val)
	}
	return s
}

func (s *SQL) Page(page, pagesize int) *SQL {
	s.limit = pagesize
	s.offset = (page - 1) * pagesize
	return s
}

func (s *SQL) Plus(col string, val int) *SQL {
	s.sets = append(s.sets, fmt.Sprintf("`%s` = `%s` + ?", col, col))
	s.setsParams = append(s.setsParams, val)
	return s
}

func (s *SQL) Incr(col string, val int) *SQL {
	s.sets = append(s.sets, fmt.Sprintf("`%s` = last_insert_id(`%s` + ?)", col, col))
	s.setsParams = append(s.setsParams, val)
	return s
}

// build sql

func (s *SQL) ToSelect(columns ...string) (string, []interface{}) {
	s.columns = append(s.columns, columns...)

	defer s.Reset()

	keyword := ""
	if len(s.keywords) > 0 {
		keyword = " " + strings.Join(s.keywords, " ")
	}
	column := "*"
	if len(s.columns) > 0 {
		column = strings.Join(s.columns, ", ")
	}
	join := ""
	if len(s.joins) > 0 {
		join = " LEFT JOIN " + strings.Join(s.joins, " LEFT JOIN ")
	}
	where := ""
	if len(s.wheres) > 0 {
		where = " WHERE " + strings.Join(s.wheres, " AND ")
	}
	group := ""
	if len(s.groups) > 0 {
		group = " GROUP BY " + strings.Join(s.groups, ", ")
	}
	having := ""
	if len(s.havings) > 0 {
		having = " HAVING " + strings.Join(s.havings, " AND ")
	}
	order := ""
	if len(s.orders) > 0 {
		order = " ORDER BY " + strings.Join(s.orders, ", ")
	}
	limit := ""
	if s.limit > -1 {
		limit = fmt.Sprintf(" LIMIT %d", s.limit)
	}
	offset := ""
	if s.offset > -1 {
		offset = fmt.Sprintf(" OFFSET %d", s.offset)
	}
	forUpdate := s.forUpdate
	lockInShareMode := s.lockInShareMode

	sql := fmt.Sprintf("SELECT%s %s FROM %s%s%s%s%s%s%s%s%s%s", keyword, column, s.table, join, where, group, having, order, limit, offset, forUpdate, lockInShareMode)

	params := make([]interface{}, 0, 20)
	params = append(params, s.wheresParams...)
	params = append(params, s.havingsParams...)

	s.countWheres = s.countWheres[0:0]
	s.countWheresParams = s.countWheresParams[0:0]
	s.countWheres = append(s.countWheres, s.wheres...)
	s.countWheresParams = append(s.countWheresParams, s.wheresParams...)

	return sql, params
}

func (s *SQL) ToCount() (string, []interface{}) {
	where := ""
	if len(s.countWheres) > 0 {
		where = " WHERE " + strings.Join(s.countWheres, " AND ")
	}

	return fmt.Sprintf("SELECT count(*) AS count FROM %s%s", s.table, where), s.countWheresParams
}

func (s *SQL) ToCountMySQL() (string, []interface{}) {
	return "SELECT FOUND_ROWS()", make([]interface{}, 0, 0)
}

func (s *SQL) ToInsert() (string, []interface{}) {
	if len(s.sets) == 0 {
		panic("Insert sets is empty!")
	}

	defer s.Reset()

	return fmt.Sprintf("INSERT %s SET %s", s.table, strings.Join(s.sets, ", ")), s.setsParams
}

func (s *SQL) ToReplace() (string, []interface{}) {
	if len(s.sets) == 0 {
		panic("Replace sets is empty!")
	}

	defer s.Reset()

	return fmt.Sprintf("REPLACE %s SET %s", s.table, strings.Join(s.sets, ", ")), s.setsParams
}

func (s *SQL) ToUpdate() (string, []interface{}) {
	if len(s.sets) == 0 {
		panic("Update sets is empty!")
	}
	if len(s.wheres) == 0 {
		panic("Update where is empty!")
	}

	defer s.Reset()

	set := strings.Join(s.sets, ", ")
	where := strings.Join(s.wheres, " AND ")
	order := ""
	if len(s.orders) > 0 {
		order = " ORDER BY " + strings.Join(s.orders, ", ")
	}
	limit := ""
	if s.limit > -1 {
		limit = fmt.Sprintf(" LIMIT %d", s.limit)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s%s%s", s.table, set, where, order, limit)

	params := make([]interface{}, 0, 20)
	params = append(params, s.setsParams...)
	params = append(params, s.wheresParams...)

	return sql, params
}

func (s *SQL) ToDelete() (string, []interface{}) {
	if len(s.wheres) == 0 {
		panic("Delete wheres is empty!")
	}

	defer s.Reset()

	where := strings.Join(s.wheres, " AND ")
	order := ""
	if len(s.orders) > 0 {
		order = " ORDER BY " + strings.Join(s.orders, ", ")
	}
	limit := ""
	if s.limit > -1 {
		limit = fmt.Sprintf(" LIMIT %d", s.limit)
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s%s%s", s.table, where, order, limit), s.wheresParams
}

func (s *SQL) String() string {
	sql, params := s.ToSelect()
	return fmt.Sprintf("%s %v", sql, params)
}

// orm

func (s *SQL) SetORM(orm *ORM) *SQL {
	s.orm = orm
	return s
}

func (s *SQL) Select(model interface{}, columns ...string) bool {
	return s.orm.Select(model, s, columns...)
}

func (s *SQL) Count() int {
	return s.orm.Count(s)
}

func (s *SQL) CountMySQL() int {
	return s.orm.CountMySQL(s)
}

func (s *SQL) Update(model interface{}, columns ...string) sql.Result {
	return s.orm.Update(model, s, columns...)
}

func (s *SQL) Delete(model interface{}) sql.Result {
	return s.orm.Delete(model, s)
}
