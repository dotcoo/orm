// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"fmt"
	"strings"
)

const (
	sqlAnd = " AND "
	sqlOr  = " OR "
)

const (
	sqlSelect = iota
	sqlCount
	sqlInsert
	sqlReplace
	sqlUpdate
	sqlDelete
)

type SQL struct {
	mode            int           // sql mode
	table           string        // table
	alias           string        // table alias
	keywords        []string      // keywords
	columns         []string      // columns
	from            string        // from
	cols            []string      // cols
	sets            []string      // sets
	setsArgs        []interface{} // sets args
	joins           []string      // joins
	wheres          []string      // where
	wheresArgs      []interface{} // where args
	groups          []string      // group
	havings         []string      // having
	havingsArgs     []interface{} // having args
	orders          []string      // order
	limit           int           // limit
	offset          int           // offset
	forUpdate       string        // read lock
	lockInShareMode string        // write lock
	orm             *ORM          // ORM
}

func newSQLMode(mode int, table ...string) *SQL {
	s := new(SQL)
	s.mode = mode
	s.From(table...)
	s.cols = make([]string, 0, 20)
	s.sets = make([]string, 0, 20)
	s.setsArgs = make([]interface{}, 0, 20)
	s.wheres = make([]string, 0, 5)
	s.wheresArgs = make([]interface{}, 0, 20)
	s.orders = make([]string, 0, 3)
	s.limit = -1
	s.offset = -1
	s.orm = DefaultORM
	return s
}

func NewSQL(table ...string) *SQL {
	return newSQLMode(sqlSelect, table...)
}

func NewSelect(table ...string) *SQL {
	return newSQLMode(sqlSelect, table...)
}

func NewInsert(table ...string) *SQL {
	return newSQLMode(sqlInsert, table...)
}

func NewReplace(table ...string) *SQL {
	return newSQLMode(sqlReplace, table...)
}

func NewUpdate(table ...string) *SQL {
	return newSQLMode(sqlUpdate, table...)
}

func NewDelete(table ...string) *SQL {
	return newSQLMode(sqlDelete, table...)
}

func (s *SQL) Reset() *SQL {
	s.keywords = s.keywords[0:0]
	s.columns = s.columns[0:0]
	s.cols = s.cols[0:0]
	s.sets = s.sets[0:0]
	s.setsArgs = s.setsArgs[0:0]
	s.joins = s.joins[0:0]
	s.wheres = s.wheres[0:0]
	s.wheresArgs = s.wheresArgs[0:0]
	s.groups = s.groups[0:0]
	s.havings = s.havings[0:0]
	s.havingsArgs = s.havingsArgs[0:0]
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

func (s *SQL) Table(table ...string) *SQL {
	return s.From(table...)
}

func (s *SQL) From(table ...string) *SQL {
	if len(table) == 1 {
		s.table = table[0]
		s.from = fmt.Sprintf("`%s`", s.table)
	}
	if len(table) == 2 {
		s.table = table[0]
		s.alias = table[1]
		s.from = fmt.Sprintf("`%s` AS `%s`", s.table, s.alias)
	}
	return s
}

func (s *SQL) Set(col string, val interface{}) *SQL {
	s.cols = append(s.cols, fmt.Sprintf("`%s`", col))
	s.sets = append(s.sets, fmt.Sprintf("`%s` = ?", col))
	s.setsArgs = append(s.setsArgs, val)
	return s
}

func (s *SQL) Join(table, alias, cond string) *SQL {
	if alias == "" {
		s.joins = append(s.joins, fmt.Sprintf("`%s` ON `%s`", table, cond))
	} else {
		s.joins = append(s.joins, fmt.Sprintf("`%s` AS `%s` ON %s", table, alias, cond))
	}
	return s
}

func (s *SQL) where(wheres *[]string, wheresArgs *[]interface{}, and string, where string, args ...interface{}) *SQL {
	switch {
	case where == "(":
		*wheres = append(*wheres, and, where)
	case where == ")":
		*wheres = append(*wheres, where)
	case and == sqlAnd, and == sqlOr:
		if len(*wheres) == 0 || (*wheres)[len(*wheres)-1] == "(" {
			*wheres = append(*wheres, where)
		} else {
			*wheres = append(*wheres, and, where)
		}
		*wheresArgs = append(*wheresArgs, args...)
	default:
		panic("not reached")
	}
	return s
}

func (s *SQL) Where(where string, args ...interface{}) *SQL {
	return s.where(&s.wheres, &s.wheresArgs, sqlAnd, where, args...)
}

func (s *SQL) WhereOr(where string, args ...interface{}) *SQL {
	return s.where(&s.wheres, &s.wheresArgs, sqlOr, where, args...)
}

func (s *SQL) whereIn(and string, where string, args ...interface{}) *SQL {
	if len(args) == 0 {
		panic("args is null!")
	}
	where = strings.Replace(where, "?", strings.Repeat(", ?", len(args))[2:], 1)
	return s.where(&s.wheres, &s.wheresArgs, and, where, args...)
}

func (s *SQL) WhereIn(where string, args ...interface{}) *SQL {
	return s.whereIn(sqlAnd, where, args...)
}

func (s *SQL) WhereOrIn(where string, args ...interface{}) *SQL {
	return s.whereIn(sqlOr, where, args...)
}

func (s *SQL) Group(groups ...string) *SQL {
	s.groups = append(s.groups, groups...)
	return s
}

func (s *SQL) Having(having string, args ...interface{}) *SQL {
	return s.where(&s.havings, &s.havingsArgs, sqlAnd, having, args...)
}

func (s *SQL) HavingOr(having string, args ...interface{}) *SQL {
	return s.where(&s.havings, &s.havingsArgs, sqlOr, having, args...)
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
		s.Set(col, val)
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
	s.setsArgs = append(s.setsArgs, val)
	return s
}

func (s *SQL) Incr(col string, val int) *SQL {
	s.sets = append(s.sets, fmt.Sprintf("`%s` = last_insert_id(`%s` + ?)", col, col))
	s.setsArgs = append(s.setsArgs, val)
	return s
}

// build sql

func (s *SQL) toSelect(columns ...string) (string, []interface{}) {
	s.columns = append(s.columns, columns...)

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
		where = " WHERE " + strings.Join(s.wheres, "")
	}
	group := ""
	if len(s.groups) > 0 {
		group = " GROUP BY " + strings.Join(s.groups, ", ")
	}
	having := ""
	if len(s.havings) > 0 {
		having = " HAVING " + strings.Join(s.havings, "")
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

	sq := fmt.Sprintf("SELECT%s %s FROM %s%s%s%s%s%s%s%s%s%s", keyword, column, s.from, join, where, group, having, order, limit, offset, forUpdate, lockInShareMode)

	args := make([]interface{}, 0, 20)
	args = append(args, s.wheresArgs...)
	args = append(args, s.havingsArgs...)

	return sq, args
}

func (s *SQL) toInsert() (string, []interface{}) {
	if len(s.sets) == 0 {
		panic("Insert sets is empty!")
	}

	defer s.Reset()

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", s.from, strings.Join(s.cols, ", "), strings.Repeat(", ?", len(s.cols))[2:]), s.setsArgs
}

func (s *SQL) toReplace() (string, []interface{}) {
	if len(s.sets) == 0 {
		panic("Replace sets is empty!")
	}

	defer s.Reset()

	return fmt.Sprintf("REPLACE INTO %s (%s) VALUES (%s)", s.from, strings.Join(s.cols, ", "), strings.Repeat(", ?", len(s.cols))[2:]), s.setsArgs
}

func (s *SQL) toUpdate() (string, []interface{}) {
	if len(s.sets) == 0 {
		panic("Update sets is empty!")
	}
	if len(s.wheres) == 0 {
		panic("Update where is empty!")
	}

	defer s.Reset()

	set := strings.Join(s.sets, ", ")
	where := strings.Join(s.wheres, "")
	order := ""
	if len(s.orders) > 0 {
		order = " ORDER BY " + strings.Join(s.orders, ", ")
	}
	limit := ""
	if s.limit > -1 {
		limit = fmt.Sprintf(" LIMIT %d", s.limit)
	}

	sq := fmt.Sprintf("UPDATE %s SET %s WHERE %s%s%s", s.from, set, where, order, limit)

	args := make([]interface{}, 0, 20)
	args = append(args, s.setsArgs...)
	args = append(args, s.wheresArgs...)

	return sq, args
}

func (s *SQL) toDelete() (string, []interface{}) {
	if len(s.wheres) == 0 {
		panic("Delete wheres is empty!")
	}

	defer s.Reset()

	where := strings.Join(s.wheres, "")
	order := ""
	if len(s.orders) > 0 {
		order = " ORDER BY " + strings.Join(s.orders, ", ")
	}
	limit := ""
	if s.limit > -1 {
		limit = fmt.Sprintf(" LIMIT %d", s.limit)
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s%s%s", s.from, where, order, limit), s.wheresArgs
}

func (s *SQL) SQL() (string, []interface{}) {
	switch s.mode {
	case sqlSelect, sqlCount:
		return s.toSelect()
	case sqlInsert:
		return s.toInsert()
	case sqlReplace:
		return s.toReplace()
	case sqlUpdate:
		return s.toUpdate()
	case sqlDelete:
		return s.toDelete()
	default:
		panic("not reached")
	}
}

func (s *SQL) String() string {
	sq, args := s.SQL()
	return fmt.Sprintf("%s %v", sq, args)
}

// set method

func (s *SQL) setMode(mode int) *SQL {
	s.mode = mode
	return s
}

func (s *SQL) SetORM(orm *ORM) *SQL {
	s.orm = orm
	return s
}

// clone sql

func (s *SQL) Clone() *SQL {
	sc := new(SQL)
	sc.mode = s.mode                                         // sql mode
	sc.table = s.table                                       // table
	sc.alias = s.alias                                       // table alias
	sc.keywords = make([]string, len(s.keywords))            // keywords
	sc.columns = make([]string, len(s.columns))              // columns
	sc.from = s.from                                         // from
	sc.cols = make([]string, len(s.cols))                    // cols
	sc.sets = make([]string, len(s.sets))                    // sets
	sc.setsArgs = make([]interface{}, len(s.setsArgs))       // sets args
	sc.joins = make([]string, len(s.joins))                  // joins
	sc.wheres = make([]string, len(s.wheres))                // where
	sc.wheresArgs = make([]interface{}, len(s.wheresArgs))   // where args
	sc.groups = make([]string, len(s.groups))                // group
	sc.havings = make([]string, len(s.havings))              // having
	sc.havingsArgs = make([]interface{}, len(s.havingsArgs)) // having args
	sc.orders = make([]string, len(s.orders))                // order
	sc.limit = s.limit                                       // limit
	sc.offset = s.offset                                     // offset
	sc.forUpdate = s.forUpdate                               // read lock
	sc.lockInShareMode = s.lockInShareMode                   // write lock
	sc.orm = s.orm                                           // ORM
	copy(sc.keywords, s.keywords)
	copy(sc.columns, s.columns)
	copy(sc.cols, s.cols)
	copy(sc.sets, s.sets)
	copy(sc.setsArgs, s.setsArgs)
	copy(sc.joins, s.joins)
	copy(sc.wheres, s.wheres)
	copy(sc.wheresArgs, s.wheresArgs)
	copy(sc.groups, s.groups)
	copy(sc.havings, s.havings)
	copy(sc.havingsArgs, s.havingsArgs)
	copy(sc.orders, s.orders)
	return sc
}

func (s *SQL) NewSelect() *SQL {
	return s.Clone().setMode(sqlSelect)
}

func (s *SQL) NewCount() *SQL {
	sc := s.Clone().setMode(sqlCount)
	sc.keywords = sc.keywords[0:0]
	sc.columns = []string{"count(*) AS count"}
	sc.cols = sc.cols[0:0]
	sc.sets = sc.sets[0:0]
	sc.setsArgs = sc.setsArgs[0:0]
	sc.orders = sc.orders[0:0]
	sc.limit = -1
	sc.offset = -1
	sc.forUpdate = ""
	sc.lockInShareMode = ""
	return sc
}

func (s *SQL) NewInsert() *SQL {
	return s.Clone().setMode(sqlInsert)
}

func (s *SQL) NewReplace() *SQL {
	return s.Clone().setMode(sqlReplace)
}

func (s *SQL) NewUpdate() *SQL {
	return s.Clone().setMode(sqlUpdate)
}

func (s *SQL) NewDelete() *SQL {
	return s.Clone().setMode(sqlDelete)
}
