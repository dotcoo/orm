// Copyright 2015 The dotcoo zhao. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"fmt"
	"strings"
)

const (
	sqlAs  = " AS "
	sqlAnd = " AND "
	sqlOr  = " OR "
)

type SQL struct {
	orm             *ORM          // ORM
	table           string        // table
	keywords        string        // keywords
	columns         string        // columns
	from            string        // from
	joins           string        // joins
	wheres          string        // where
	wheresArgs      []interface{} // where args
	groups          string        // group
	havings         string        // having
	havingsArgs     []interface{} // having args
	orders          string        // order
	limit           int           // limit
	offset          int           // offset
	forUpdate       string        // write lock
	lockInShareMode string        // read lock
	cols            string        // cols
	sets            string        // sets args
	setsArgs        []interface{} // sets args
}

func (s *SQL) Reset() *SQL {
	s.keywords = ""
	s.columns = ""
	s.joins = ""
	s.wheres = ""
	s.wheresArgs = s.wheresArgs[0:0]
	s.groups = ""
	s.havings = ""
	s.havingsArgs = s.havingsArgs[0:0]
	s.orders = ""
	s.limit = 0
	s.offset = 0
	s.forUpdate = ""
	s.lockInShareMode = ""
	s.cols = ""
	s.sets = ""
	s.setsArgs = s.setsArgs[0:0]
	return s
}

// sql syntax

func (s *SQL) Keywords(keywords ...string) *SQL {
	s.keywords += " " + strings.Join(keywords, " ")
	return s
}

func (s *SQL) CalcFoundRows() *SQL {
	return s.Keywords("SQL_CALC_FOUND_ROWS")
}

func (s *SQL) Columns(columns ...string) *SQL {
	if len(columns) == 0 {
		return s
	}
	s.columns += ", " + strings.Join(columns, ", ")
	return s
}

func (s *SQL) From(table string) *SQL {
	if s.orm != nil {
		table = s.orm.sqlFrom(s, table)
	}
	s.table = table
	s.from = strings.Replace(s.table, sqlAs, "` AS `", -1)
	return s
}

func (s *SQL) Set(col string, val interface{}) *SQL {
	s.cols += ", `" + col + "`"
	s.sets += ", `" + col + "` = ?"
	s.setsArgs = append(s.setsArgs, val)
	return s
}

func (s *SQL) Join(table, cond string) *SQL {
	if s.orm != nil {
		table, cond = s.orm.sqlJoin(s, table, cond)
	}
	s.joins += " LEFT JOIN `" + strings.Replace(table, sqlAs, "` AS `", 1) + "` ON " + cond
	return s
}

func (s *SQL) Where(where string, args ...interface{}) *SQL {
	if !strings.HasPrefix(where, sqlAnd) && !strings.HasPrefix(where, sqlOr) {
		where = sqlAnd + where
	}
	s.wheres += where
	s.wheresArgs = append(s.wheresArgs, args...)
	return s
}

func (s *SQL) WhereIn(where string, args ...interface{}) *SQL {
	return s.Where(strings.Replace(where, "?", strings.Repeat(", ?", len(args))[2:], 1), args...)
}

func (s *SQL) Group(groups ...string) *SQL {
	s.groups += ", " + strings.Join(groups, ", ")
	return s
}

func (s *SQL) Having(having string, args ...interface{}) *SQL {
	if !strings.HasPrefix(having, sqlAnd) && !strings.HasPrefix(having, sqlOr) {
		having = sqlAnd + having
	}
	s.havings += having
	s.havingsArgs = append(s.havingsArgs, args...)
	return s
}

func (s *SQL) Order(orders ...string) *SQL {
	s.orders += ", " + strings.Join(orders, ", ")
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
	s.sets += ", `" + col + "` = `" + col + "` + ?"
	s.setsArgs = append(s.setsArgs, val)
	return s
}

func (s *SQL) Incr(col string, val int) *SQL {
	s.sets += ", `" + col + "` = last_insert_id(`" + col + "` + ?)"
	s.setsArgs = append(s.setsArgs, val)
	return s
}

// build sql

func (s *SQL) ToSelect() (string, []interface{}) {
	column := " *"
	if s.columns != "" {
		column = s.columns[1:]
	}
	where := ""
	if s.wheres != "" {
		where = " WHERE " + s.wheres[5:]
	}
	group := ""
	if s.groups != "" {
		group = " GROUP BY " + s.groups[2:]
	}
	having := ""
	if s.havings != "" {
		having = " HAVING " + s.havings[5:]
	}
	order := ""
	if s.orders != "" {
		order = " ORDER BY " + s.orders[2:]
	}
	limit := ""
	if s.limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", s.limit)
	}
	offset := ""
	if s.offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", s.offset)
	}
	sq := "SELECT" + s.keywords + column + " FROM `" + s.from + "`" + s.joins + where + group + having + order + limit + offset + s.forUpdate + s.lockInShareMode

	args := make([]interface{}, 0, len(s.wheresArgs)+len(s.havingsArgs))
	args = append(args, s.wheresArgs...)
	args = append(args, s.havingsArgs...)

	return sq, args
}

func (s *SQL) ToInsert() (string, []interface{}) {
	return "INSERT INTO `" + s.from + "` (" + s.cols[2:] + ") VALUES (" + strings.Repeat(", ?", len(s.setsArgs))[2:] + ")", s.setsArgs
}

func (s *SQL) ToReplace() (string, []interface{}) {
	return "REPLACE INTO `" + s.from + "` (" + s.cols[2:] + ") VALUES (" + strings.Repeat(", ?", len(s.setsArgs))[2:] + ")", s.setsArgs
}

func (s *SQL) ToUpdate() (string, []interface{}) {
	where := ""
	if s.wheres != "" {
		where = " WHERE " + s.wheres[5:]
	} else {
		panic("Update where is empty!")
	}
	order := ""
	if s.orders != "" {
		order = " ORDER BY " + s.orders
	}
	limit := ""
	if s.limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", s.limit)
	}
	sq := "UPDATE `" + s.from + "` SET " + s.sets[2:] + where + order + limit

	args := make([]interface{}, 0, len(s.setsArgs)+len(s.wheresArgs))
	args = append(args, s.setsArgs...)
	args = append(args, s.wheresArgs...)

	return sq, args
}

func (s *SQL) ToDelete() (string, []interface{}) {
	where := ""
	if s.wheres != "" {
		where = " WHERE " + s.wheres[5:]
	} else {
		panic("Delete where is empty!")
	}
	order := ""
	if s.orders != "" {
		order = " ORDER BY " + s.orders
	}
	limit := ""
	if s.limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", s.limit)
	}

	return "DELETE FROM `" + s.from + "`" + where + order + limit, s.wheresArgs
}

// count

func (s *SQL) NewCount() *SQL {
	c := new(SQL)
	c.orm = s.orm
	c.table = s.table
	c.columns = ", count(*) AS count"
	c.from = s.from
	c.joins = s.joins
	c.wheres = s.wheres
	c.wheresArgs = s.wheresArgs
	c.groups = s.groups
	c.havings = s.havings
	c.havingsArgs = s.havingsArgs
	// c.orders = s.orders
	return c
}
