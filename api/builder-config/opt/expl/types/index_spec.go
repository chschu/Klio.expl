package types

import (
	"strings"
)

type IndexSpec interface {
	SqlCondition() (sqlCondition string, params []any)
}

func IndexSpecAll() IndexSpec {
	return indexSpec([]IndexRange{&indexRange{
		from: HeadIndex(1),
		to:   TailIndex(1),
	}})
}

func IndexSpecSingle(index Index) IndexSpec {
	return indexSpec([]IndexRange{&indexRange{
		from: index,
		to:   index,
	}})
}

type indexSpec []IndexRange

func (is indexSpec) SqlCondition() (sqlCondition string, params []any) {
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
