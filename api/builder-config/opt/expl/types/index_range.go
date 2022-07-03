package types

import "fmt"

type IndexRange interface {
	fmt.Stringer
	From() Index
	To() Index
	SqlCondition() (sqlCondition string, params []any)
}

func NewIndexRange(from Index, to Index) IndexRange {
	return &indexRange{
		from: from,
		to:   to,
	}
}

type indexRange struct {
	from Index
	to   Index
}

func (ir *indexRange) From() Index {
	return ir.from
}

func (ir *indexRange) To() Index {
	return ir.to
}

func (ir *indexRange) SqlCondition() (sqlCondition string, params []any) {
	fromSql, fromParams := ir.from.SqlCondition(IndexStartingWith)
	toSql, toParams := ir.to.SqlCondition(IndexEndingWith)
	return "(" + fromSql + ") AND (" + toSql + ")", append(fromParams, toParams...)
}

func (ir *indexRange) String() string {
	if ir.from == ir.to {
		return ir.from.String()
	}
	return ir.from.String() + ":" + ir.to.String()
}
