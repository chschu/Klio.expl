package types

import (
	"strings"
)

type IndexSpec struct {
	ranges []IndexRange
}

func NewIndexSpec(ranges ...IndexRange) IndexSpec {
	return IndexSpec{
		ranges: ranges,
	}
}

func IndexSpecAll() IndexSpec {
	return NewIndexSpec(NewIndexRange(NewHeadIndex(1), NewTailIndex(1)))
}

func IndexSpecSingle(index Index) IndexSpec {
	return NewIndexSpec(NewIndexRange(index, index))
}

func (is IndexSpec) SQLCondition() (sqlCondition string, params []any) {
	sb := strings.Builder{}
	sb.WriteString("false")
	for _, ir := range is.ranges {
		indexRangeSql, indexRangeParams := ir.SQLCondition()
		sb.WriteString(" OR (")
		sb.WriteString(indexRangeSql)
		sb.WriteRune(')')
		params = append(params, indexRangeParams...)
	}
	return sb.String(), params
}

func (is IndexSpec) String() string {
	sb := strings.Builder{}
	sb.WriteString("")
	for _, ir := range is.ranges {
		sb.WriteString(ir.String())
		sb.WriteRune(' ')
	}
	return strings.TrimRight(sb.String(), " ")
}
