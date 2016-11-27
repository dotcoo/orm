#!/bin/sh

set -eu

# test SQL
cat > sql_tmp_orm.go << EOF
package orm
type ORM struct {}
func (o *ORM) sqlFrom(s *SQL, table string) string {
	return table
}
func (o *ORM) sqlJoin(s *SQL, table, cond string) (string, string) {
	return table, cond
}
type ModelInfo struct{}
EOF
go test sql_tmp_orm.go sql.go sql_test.go
rm sql_tmp_orm.go

# test ModelInfo
go test modelinfo.go modelinfo_test.go

# test ORM
go test sql.go sql_test.go modelinfo.go modelinfo_test.go orm.go orm_test.go

# test ORM safe
go test sql.go sql_test.go modelinfo.go modelinfo_test.go orm.go orm_test.go orm_safe.go

# test SQL ORM
go test sql.go sql_test.go modelinfo.go modelinfo_test.go orm.go orm_test.go orm_safe.go sql_orm.go

# test orm func
go test sql.go sql_test.go modelinfo.go modelinfo_test.go orm.go orm_test.go orm_safe.go sql_orm.go func.go
