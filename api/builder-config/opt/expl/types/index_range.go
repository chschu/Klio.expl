package types

func NewIndexRange(from Index, to Index) IndexRange {
	return IndexRange{
		from: from,
		to:   to,
	}
}

type IndexRange struct {
	from Index
	to   Index
}

func (ir IndexRange) SQLCondition() (sqlCondition string, params []any) {
	if ir.from == ir.to {
		return ir.from.SQLConditionMatching()
	}
	fromSql, fromParams := ir.from.SQLConditionStartingWith()
	toSql, toParams := ir.to.SQLConditionEndingWith()
	return "(" + fromSql + ") AND (" + toSql + ")", append(fromParams, toParams...)
}

func (ir IndexRange) String() string {
	if ir.from == ir.to {
		return ir.from.String()
	}
	return ir.from.String() + ":" + ir.to.String()
}
