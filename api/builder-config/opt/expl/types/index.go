package types

import (
	"fmt"
)

type Index struct {
	sqlColumn    string
	sqlAscending bool
	value        uint
	prefix       string
}

func NewHeadIndex(n uint) Index {
	return Index{
		sqlColumn:    "head_index",
		sqlAscending: true,
		value:        n,
		prefix:       "",
	}
}

func NewTailIndex(n uint) Index {
	return Index{
		sqlColumn:    "tail_index",
		sqlAscending: false,
		value:        n,
		prefix:       "-",
	}
}

func NewPermanentIndex(n uint) Index {
	return Index{
		sqlColumn:    "permanent_index",
		sqlAscending: true,
		value:        n,
		prefix:       "p",
	}
}

func (i Index) SQLConditionStartingWith() (sqlCondition string, params []any) {
	var sqlOp string
	if i.sqlAscending {
		sqlOp = ">="
	} else {
		sqlOp = "<="
	}
	return i.sqlCondition(sqlOp)
}

func (i Index) SQLConditionMatching() (sqlCondition string, params []any) {
	return i.sqlCondition("=")
}

func (i Index) SQLConditionEndingWith() (sqlCondition string, params []any) {
	var sqlOp string
	if i.sqlAscending {
		sqlOp = "<="
	} else {
		sqlOp = ">="
	}
	return i.sqlCondition(sqlOp)
}

func (i Index) sqlCondition(sqlOp string) (sqlCondition string, params []any) {
	return fmt.Sprintf("%s %s ?", i.sqlColumn, sqlOp), []any{i.value}
}

func (i Index) String() string {
	return fmt.Sprintf("%s%d", i.prefix, i.value)
}
