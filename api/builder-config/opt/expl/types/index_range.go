package types

import (
	"fmt"
	"regexp"
)

type IndexRange struct {
	from Index
	to   Index
}

func NewIndexRange(from Index, to Index) IndexRange {
	return IndexRange{
		from: from,
		to:   to,
	}
}

func NewIndexRangeFromString(s string) (IndexRange, error) {
	sep := regexp.MustCompile("^\\pZ*(?P<From>\\PZ+?)(:(?P<To>\\PZ+?))?\\pZ*$")
	match := sep.FindStringSubmatch(s)
	if match == nil {
		return IndexRange{}, fmt.Errorf("invalid index range: %s", s)
	}

	fromStr := match[sep.SubexpIndex("From")]
	from, err := NewIndexFromString(fromStr)
	if err != nil {
		return IndexRange{}, err
	}

	toStr := match[sep.SubexpIndex("To")]
	if toStr == "" {
		return NewIndexRange(from, from), nil
	}
	to, err := NewIndexFromString(toStr)
	if err != nil {
		return IndexRange{}, err
	}

	return NewIndexRange(from, to), nil
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
