package types

import (
	"fmt"
	"strings"
)

type IndexSpec interface {
	fmt.Stringer
	SqlCondition() (sqlCondition string, params []any)
}

func NewIndexSpec(r ...IndexRange) IndexSpec {
	return indexSpec(r)
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
	params = []any{}
	for _, ir := range is {
		indexRangeSql, indexRangeParams := ir.SqlCondition()
		sb.WriteString(" OR (")
		sb.WriteString(indexRangeSql)
		sb.WriteRune(')')
		params = append(params, indexRangeParams...)
	}
	return sb.String(), params
}

func (is indexSpec) String() string {
	sb := strings.Builder{}
	sb.WriteString("")
	for _, ir := range is {
		sb.WriteString(ir.String())
		sb.WriteRune(' ')
	}
	return strings.TrimRight(sb.String(), " ")
}
