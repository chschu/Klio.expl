package types

import (
	"fmt"
	"regexp"
	"strconv"
)

type Index struct {
	sqlColumn    string
	sqlAscending bool
	value        uint
	prefix       string
}

var indexRegexp = regexp.MustCompile("^\\pZ*(?P<Prefix>|-|p)(?P<N>[1-9]\\d*)\\pZ*$")
var indexSubexpIndexPrefix = indexRegexp.SubexpIndex("Prefix")
var indexSubexpIndexN = indexRegexp.SubexpIndex("N")

func NewIndexFromString(s string) (Index, error) {
	match := indexRegexp.FindStringSubmatch(s)
	if match == nil {
		return Index{}, fmt.Errorf("invalid index: %s", s)
	}
	prefix := match[indexSubexpIndexPrefix]
	nStr := match[indexSubexpIndexN]

	n64, err := strconv.ParseUint(nStr, 10, 32)
	if err != nil {
		return Index{}, err
	}
	n := uint(n64)

	switch prefix {
	case "":
		return NewHeadIndex(n), nil
	case "-":
		return NewTailIndex(n), nil
	case "p":
		return NewPermanentIndex(n), nil
	}

	return Index{}, fmt.Errorf("invalid index: %s", s)
}

func NewHeadIndex(n uint) Index {
	return Index{
		sqlColumn:    "head_index",
		sqlAscending: true,
		value:        n,
		prefix:       "",
	}
}

func NewTailIndex(n uint) Index {
	return Index{
		sqlColumn:    "tail_index",
		sqlAscending: false,
		value:        n,
		prefix:       "-",
	}
}

func NewPermanentIndex(n uint) Index {
	return Index{
		sqlColumn:    "permanent_index",
		sqlAscending: true,
		value:        n,
		prefix:       "p",
	}
}

func (i Index) SQLConditionStartingWith() (sqlCondition string, params []any) {
	var sqlOp string
	if i.sqlAscending {
		sqlOp = ">="
	} else {
		sqlOp = "<="
	}
	return i.sqlCondition(sqlOp)
}

func (i Index) SQLConditionMatching() (sqlCondition string, params []any) {
	return i.sqlCondition("=")
}

func (i Index) SQLConditionEndingWith() (sqlCondition string, params []any) {
	var sqlOp string
	if i.sqlAscending {
		sqlOp = "<="
	} else {
		sqlOp = ">="
	}
	return i.sqlCondition(sqlOp)
}

func (i Index) sqlCondition(sqlOp string) (sqlCondition string, params []any) {
	return fmt.Sprintf("%s %s ?", i.sqlColumn, sqlOp), []any{i.value}
}

func (i Index) String() string {
	return fmt.Sprintf("%s%d", i.prefix, i.value)
}
