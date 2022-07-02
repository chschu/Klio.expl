package types

import "strings"

type IndexSpec []IndexRange

func (is IndexSpec) SqlCondition() (sqlCondition string, params []any) {
	sb := strings.Builder{}
	sb.WriteString("false")
	for _, ir := range is {
		indexRangeSql, indexRangeParams := ir.SqlCondition()
		sb.WriteString(" OR (")
		sb.WriteString(indexRangeSql)
		sb.WriteRune(')')
		params = append(params, indexRangeParams...)
	}
	return sb.String(), params
}

func IndexSpecAll() IndexSpec {
	return []IndexRange{{
		From: HeadIndex(1),
		To:   TailIndex(1),
	}}
}

func IndexSpecSingle(index Index) IndexSpec {
	return []IndexRange{{
		From: index,
		To:   index,
	}}
}
