package types

type IndexRange struct {
	From Index
	To   Index
}

func (ir *IndexRange) SqlCondition() (sqlCondition string, params []any) {
	fromSql, fromParams := ir.From.SqlCondition(IndexFrom)
	toSql, toParams := ir.To.SqlCondition(IndexTo)
	return "(" + fromSql + ") AND (" + toSql + ")", append(fromParams, toParams...)
}
