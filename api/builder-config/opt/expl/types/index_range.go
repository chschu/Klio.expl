package types

type IndexRange interface {
	From() Index
	To() Index
	SqlCondition() (sqlCondition string, params []any)
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
	fromSql, fromParams := ir.from.SqlCondition(IndexFrom)
	toSql, toParams := ir.to.SqlCondition(IndexTo)
	return "(" + fromSql + ") AND (" + toSql + ")", append(fromParams, toParams...)
}
