package types

import (
	"regexp"
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

var indexSpecWhitespaceRegexp = regexp.MustCompile("\\pZ+")

func NewIndexSpecFromString(s string) (IndexSpec, error) {
	var irs []IndexRange
	irStrs := indexSpecWhitespaceRegexp.Split(s, -1)
	for _, irStr := range irStrs {
		if len(irStr) > 0 {
			ir, err := NewIndexRangeFromString(irStr)
			if err != nil {
				return IndexSpec{}, err
			}
			irs = append(irs, ir)
		}
	}
	return NewIndexSpec(irs...), nil
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
